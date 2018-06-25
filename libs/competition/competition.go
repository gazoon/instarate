package competition

import (
	"context"
	"instarate/libs/competition/config"
	"instarate/libs/instagram"
	"path"
	"time"

	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/gazoon/go-utils"
	"github.com/gazoon/go-utils/logging"
	"github.com/pkg/errors"
)

const (
	CelebritiesCompetition = "celebrities"
	ModelsCompetition      = "models"
	RegularCompetition     = "regular"
	GlobalCompetition      = "global"
)

const (
	httpTimeout = time.Second * 3

	celebrityFollowersThreshold = 500000
	modelFollowersThreshold     = 10000

	nextPairGetAttempts = 10
)

var (
	AlreadyVotedErr          = errors.New("already voted")
	BadPhotoLinkErr          = errors.New("photo link doesn't contain a valid media code")
	BadProfileLinkErr        = errors.New("profile link doesn't contain a valid username")
	NotPhotoMediaErr         = errors.New("media is not a photo")
	GetNextPairNoAttemptsErr = errors.New("out of attempts to get next pair")
)

type CompetitorNotFound struct {
	Username string
}

func (self *CompetitorNotFound) Error() string {
	return fmt.Sprintf("competitor %s doesn't exist", self.Username)
}

type InstCompetitor struct {
	Username string
	*InstProfile
	*competitor
}

func (self InstCompetitor) String() string {
	return utils.ObjToString(&self)
}

type Competition struct {
	*logging.LoggerMixin
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
	return &Competition{
		competitors: competitors, profiles: profiles, photosStorage: photosStorage,
		voters: voters, LoggerMixin: logging.NewLoggerMixin("competition", nil)}, nil
}

func (self *Competition) GetPhotoUrl(competitor *InstCompetitor) string {
	return self.photosStorage.buildUrl(competitor.PhotoPath)
}

func (self *Competition) GetPosition(ctx context.Context, competitor *InstCompetitor) (int, error) {
	num, err := self.competitors.getNumberWithHigherRating(ctx, competitor.CompetitionCode, competitor.Rating)
	if err != nil {
		return 0, err
	}
	return num + 1, nil
}

func (self *Competition) Add(ctx context.Context, photoLink string) (*InstProfile, error) {
	logger := self.GetLogger(ctx)
	logger.WithField("photo_link", photoLink).Info("Add instagram competitor by photo link")
	mediaCode, err := instagram.ExtractMediaCode(photoLink)
	if err != nil {
		logger.WithFields(log.Fields{"photo_link": photoLink, "error": err}).
			Warn("Can't extract media code from the photo link")
		return nil, BadPhotoLinkErr
	}
	mediaInfo, err := instagram.GetMediaInfo(ctx, mediaCode)
	if err != nil {
		if err == instagram.MediaForbidden {
			logger.WithField("media_code", mediaCode).Warn("Media not found")
		}
		return nil, err
	}
	if !mediaInfo.IsPhoto {
		logger.WithField("media_code", mediaCode).Warn("Media is not a photo")
		return nil, NotPhotoMediaErr
	}
	followers, err := instagram.GetFollowersNumber(ctx, mediaInfo.Owner)
	if err != nil {
		return nil, err
	}
	profile := newProfile(mediaInfo.Owner, mediaCode, followers)
	logger.WithField("profile", profile).Info("Add new instagram profile")
	err = self.profiles.create(ctx, profile)
	if err == ProfileExistsErr {
		logger.WithField("username", mediaInfo.Owner).Info("Instagram profile already exists")
		return profile, ProfileExistsErr
	}
	if err != nil {
		return nil, err
	}
	_, err = self.photosStorage.upload(ctx, profile.PhotoPath, mediaInfo.Url)
	if err != nil {
		return nil, err
	}
	competitions := choseCompetition(followers)
	logger.WithField("competitions", competitions).Info("Add new profile to suitable competitions")
	for _, competitionCode := range competitions {
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
	return nil, nil, GetNextPairNoAttemptsErr
}

func (self *Competition) GetCompetitor(ctx context.Context, competitionCode, profileLink string) (*InstCompetitor, error) {
	username, err := instagram.ExtractUsername(profileLink)
	if err != nil {
		logger := self.GetLogger(ctx)
		logger.WithFields(log.Fields{"profile_link": profileLink, "error": err}).
			Warn("Can't extract username from the profile link")
		return nil, BadProfileLinkErr
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

func (self *Competition) Remove(ctx context.Context, usernames ...string) error {
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
	logger := self.GetLogger(ctx)
	logger.WithFields(log.Fields{
		"competition":     competitionCode,
		"voters_group_id": votersGroupId,
		"voter_id":        voterId,
		"winner":          winnerUsername,
		"loser":           loserUsername,
	}).Info("Save user vote")
	ok, err := self.voters.tryVote(ctx, competitionCode, votersGroupId, voterId, winnerUsername, loserUsername)
	if err != nil {
		return nil, nil, err
	}
	if !ok {
		logger.Info("User already voter")
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
	return &InstCompetitor{InstProfile: profile, competitor: competitor, Username: profile.Username}
}

func choseCompetition(followersNumber int) []string {
	var competitionByFollowers string
	if followersNumber < modelFollowersThreshold {
		competitionByFollowers = RegularCompetition
	} else if followersNumber < celebrityFollowersThreshold {
		competitionByFollowers = ModelsCompetition
	} else {
		competitionByFollowers = CelebritiesCompetition
	}
	return []string{GlobalCompetition, competitionByFollowers}
}
