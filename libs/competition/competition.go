package competition

import (
	"context"
	"fmt"
	"instarate/libs/competition/config"
	"instarate/libs/instagram"
	"path"
	"time"

	"github.com/gazoon/go-utils"
	"github.com/pkg/errors"
)

const (
	CelebritiesCompetition = "celebrities"
	ModelsCompetition      = "models"
	NormalCompetition      = "normal"
	GlobalCompetition      = "global"
)

const (
	httpTimeout = time.Second * 3

	celebrityFollowersThreshold = 500000
	modelFollowersThreshold     = 10000

	nextPairGetAttempts = 10
)

var (
	AlreadyVotedErr = errors.New("already voted")
)

type InstCompetitor struct {
	*InstProfile
	*competitor
}

type Competition struct {
	competitors   *competitorsStorage
	profiles      *profilesStorage
	voters        *votersStorage
	photosStorage *googleStorage
}

func New() (*Competition, error) {

	conf := &config.Config{}
	configPath := path.Join(utils.GetCurrentFileDir(), "config")
	err := utils.ParseConfig(configPath, &conf)
	if err != nil {
		return nil, err
	}
	competitors, err := newCompetitorsStorage(conf.MongoCompetitors)
	if err != nil {
		return nil, err
	}
	profiles, err := newProfilesStorage(conf.MongoProfiles)
	if err != nil {
		return nil, err
	}
	photosStorage, err := newGoogleStorage(conf.GoogleStorage.BucketName)
	if err != nil {
		return nil, err
	}
	voters, err := newVotersStorage(conf.MongoVoters)
	if err != nil {
		return nil, err
	}
	return &Competition{competitors: competitors, profiles: profiles, photosStorage: photosStorage, voters: voters}, nil
}

func (self *Competition) Add(ctx context.Context, photoUrl string) (*InstProfile, error) {
	mediaCode, err := instagram.ExtractMediaCode(photoUrl)
	if err != nil {
		return nil, err
	}
	mediaInfo, err := instagram.GetMediaInfo(ctx, mediaCode)
	if err != nil {
		return nil, err
	}
	if !mediaInfo.IsPhoto {
		return nil, errors.Errorf("media %s is not a photo", mediaCode)
	}
	followers, err := instagram.GetFollowersNumber(ctx, mediaInfo.Owner)
	if err != nil {
		return nil, err
	}
	profile := createProfile(mediaInfo.Owner, mediaCode, followers)
	err = self.profiles.save(ctx, profile)
	if err != nil {
		return nil, err
	}
	_, err = self.photosStorage.upload(ctx, profile.PhotoPath, mediaInfo.Url)
	if err != nil {
		return nil, err
	}
	for _, competitionCode := range choseCompetition(followers) {
		compttr := createCompetitor(mediaInfo.Owner, competitionCode)
		err = self.competitors.create(ctx, compttr)
		if err != nil {
			return nil, err
		}
	}
	return profile, nil
}

func (self *Competition) GetCompetitorsNumber(ctx context.Context, competitionCode string) (int, error) {
	return self.competitors.getCompetitorsNumber(ctx, competitionCode)
}

func (self *Competition) GetNextPair(ctx context.Context, competitionCode, votersGroupId string) (*InstCompetitor, *InstCompetitor, error) {
	for i := 0; i < nextPairGetAttempts; i++ {
		competitor1, competitor2, err := self.competitors.getRandomPair(ctx, competitionCode)
		if err != nil {
			return nil, nil, err
		}
		haveSeenPair, err := self.voters.haveSeenPair(ctx, competitionCode, votersGroupId, competitor1.Username, competitor2.Username)
		if err != nil {
			return nil, nil, err
		}
		if haveSeenPair {
			continue
		}
		return self.convertPairToInstCompetitors(ctx, competitor1, competitor2)
	}
	return nil, nil, errors.Errorf("out of attempts to get next pair in %s for %s", competitionCode, votersGroupId)
}

func (self *Competition) GetCompetitor(ctx context.Context, competitionCode, username string) (*InstCompetitor, error) {
	username, err := instagram.ExtractUsername(username)
	if err != nil {
		return nil, err
	}
	compttr, err := self.competitors.get(ctx, competitionCode, username)
	if err != nil {
		return nil, err
	}
	profile, err := self.profiles.get(ctx, username)
	if err != nil {
		return nil, err
	}
	return combineProfileAndCompetitor(profile, compttr), nil
}

func (self *Competition) Remove(ctx context.Context, usernames []string) error {
	var err error
	for i := range usernames {
		usernames[i], err = instagram.ExtractUsername(usernames[i])
		if err != nil {
			return err
		}
	}
	err = self.competitors.delete(ctx, usernames)
	if err != nil {
		return err
	}
	err = self.profiles.delete(ctx, usernames)
	if err != nil {
		return err
	}
	return nil
}

func (self *Competition) GetTop(ctx context.Context, competitionCode string, number, offset int) ([]*InstCompetitor, error) {
	competitors, err := self.competitors.getTop(ctx, competitionCode, number, offset)
	if err != nil {
		return nil, err
	}
	return self.convertToInstCompetitors(ctx, competitors...)
}

func (self *Competition) Vote(ctx context.Context, competitionCode, votersGroupId, voterId, winnerUsername, loserUsername string) (*InstCompetitor, *InstCompetitor, error) {
	ok, err := self.voters.tryVote(ctx, competitionCode, votersGroupId, voterId, winnerUsername, loserUsername)
	if err != nil {
		return nil, nil, err
	}
	if !ok {
		return nil, nil, AlreadyVotedErr
	}

	winner, err := self.competitors.get(ctx, competitionCode, winnerUsername)
	if err != nil {
		return nil, nil, err
	}
	loser, err := self.competitors.get(ctx, competitionCode, loserUsername)
	if err != nil {
		return nil, nil, err
	}

	winner.Rating, loser.Rating = recalculateEloRating(winner.Rating, loser.Rating)
	winner.Wins += 1
	winner.Matches += 1
	loser.Loses += 1
	loser.Matches += 1

	self.competitors.update(ctx, winner)
	if err != nil {
		return nil, nil, err
	}
	self.competitors.update(ctx, loser)
	if err != nil {
		return nil, nil, err
	}

	return self.convertPairToInstCompetitors(ctx, winner, loser)
}

func (self *Competition) convertPairToInstCompetitors(ctx context.Context, competitor1, competitor2 *competitor) (*InstCompetitor, *InstCompetitor, error) {
	competitorsPair, err := self.convertToInstCompetitors(ctx, competitor1, competitor2)
	if err != nil {
		return nil, nil, err
	}
	if len(competitorsPair) != 2 {
		return nil, nil, errors.Errorf("expected two competitors, got: %v", competitorsPair)
	}
	return competitorsPair[0], competitorsPair[1], nil
}

func (self *Competition) convertToInstCompetitors(ctx context.Context, competitors ...*competitor) ([]*InstCompetitor, error) {
	usernames := make([]string, len(competitors))
	for i := range usernames {
		usernames[i] = competitors[i].Username
	}
	profiles, err := self.profiles.getMultiple(ctx, usernames)
	if err != nil {
		return nil, err
	}
	profilesMapping := make(map[string]*InstProfile, len(profiles))
	for _, profile := range profiles {
		profilesMapping[profile.Username] = profile
	}
	if len(competitors) != len(profilesMapping) {
		return nil, errors.New("number of profiles is not equal to number of competitors")
	}
	result := make([]*InstCompetitor, len(competitors))
	for i, competitor := range competitors {
		profile, ok := profilesMapping[competitor.Username]
		if !ok {
			return nil, errors.Errorf("cant find profile for competitor %s", competitor.Username)
		}
		result[i] = combineProfileAndCompetitor(profile, competitor)
	}
	return result, nil
}

func combineProfileAndCompetitor(profile *InstProfile, competitor *competitor) *InstCompetitor {
	return &InstCompetitor{InstProfile: profile, competitor: competitor}
}

func choseCompetition(followersNumber int) []string {
	var competitionByFollowers string
	if followersNumber < modelFollowersThreshold {
		competitionByFollowers = NormalCompetition
	} else if followersNumber < celebrityFollowersThreshold {
		competitionByFollowers = ModelsCompetition
	} else {
		competitionByFollowers = CelebritiesCompetition
	}
	return []string{GlobalCompetition, competitionByFollowers}
}

func (self *Competition) Test() {
	ok, err := self.voters.tryVote(nil, "global", "tt", "22", "1", "2")
	if err != nil {
		panic(err)
	}
	fmt.Println(ok)
	seen, err := self.voters.haveSeenPair(nil, "global", "tt", "2", "1")
	if err != nil {
		panic(err)
	}
	fmt.Println(seen)
}
