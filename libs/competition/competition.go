package competition

import (
	"fmt"
	"github.com/gazoon/go-utils"
	"instarate/libs/competition/config"
	"path"
	"time"
)

const (
	httpTimeout = time.Second * 3
)

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

func (self *Competition) Test() {
	ok, err := self.voters.tryVote("global", "tt", "22", "1", "2")
	if err != nil {
		panic(err)
	}
	fmt.Println(ok)
	seen, err := self.voters.haveSeenPair("global", "tt", "2", "1")
	if err != nil {
		panic(err)
	}
	fmt.Println(seen)
}
