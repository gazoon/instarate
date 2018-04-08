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
	photosStorage *googleStorage
	profiles      *profilesStorage
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
	return &Competition{competitors: competitors, profiles: profiles, photosStorage: photosStorage}, nil
}

func (self *Competition) Test() {
	p, err := self.profiles.get("deleks_lina")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", p)
	ps, err := self.profiles.getMultiple([]string{"anyuta_rai", "galina_dub"})
	if err != nil {
		panic(err)
	}
	for _, p := range ps {
		fmt.Printf("%+v\n", p)
	}
	err = self.profiles.delete([]string{"anyuta_rai"})
	if err != nil {
		panic(err)
	}
}
