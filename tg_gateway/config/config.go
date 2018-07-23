package config

import (
	"github.com/gazoon/go-utils"
)

var (
	Config *ConfigSchema
)

type ConfigSchema struct {
	utils.RootConfig `yaml:",inline"`
	MongoQueue       *utils.MongoDBSettings `yaml:"mongo_queue"`
	KnownBots        map[string]string      `yaml:"known_bots"`
	BotToken         string                 `yaml:"bot_token"`
	PublicUrl        string                 `yaml:"public_url"`
	Sentry           *utils.SentrySettings  `yaml:"sentry"`
}

func (self *ConfigSchema) String() string {
	return utils.ObjToString(self)
}

func init() {
	confDir := utils.GetCurrentFileDir()
	conf := &ConfigSchema{}
	err := utils.ParseConfig(confDir, conf)
	if err != nil {
		panic(err)
	}
	Config = conf
}
