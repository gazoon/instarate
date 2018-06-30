package config

import (
	"github.com/gazoon/go-utils"
)

type Config struct {
	utils.RootConfig `yaml:",inline"`
	MongoQueue       *utils.MongoDBSettings `yaml:"mongo_queue"`
	KnownBots        map[string]string      `yaml:"known_bots"`
	PublicUrl        string                 `yaml:"public_url"`
	Sentry           *utils.SentrySettings  `yaml:"sentry"`
}

func (self *Config) String() string {
	return utils.ObjToString(self)
}
