package competition

import (
	"context"
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
	photosStorage, err := newGoogleStorage(conf.GoogleStorage.BucketName)
	if err != nil {
		return nil, err
	}
	return &Competition{competitors: competitors, photosStorage: photosStorage}, nil
}

func (self *Competition) Test() {
	s, err := self.photosStorage.upload(context.Background(), "test22.png", "https://scontent-frt3-1.cdninstagram.com/vp/663e956f9bfe4beff327b8283723e205/5B54901A/t51.2885-15/e35/29739290_160601121289134_5239513480079343616_n.jpg")
	if err != nil {
		panic(err)
	}
	fmt.Println(s)
}
