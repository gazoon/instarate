package config

import (
	"github.com/gazoon/go-utils"
	"path"
)

var (
	Config  *ConfigSchema
	RootDir string
)

type ConfigSchema struct {
	utils.RootConfig `yaml:",inline"`
	MongoCache       *utils.MongoDBSettings `yaml:"mongo_cache"`
	MongoChats       *utils.MongoDBSettings `yaml:"mongo_chats"`
	MongoQueue       *utils.MongoDBSettings `yaml:"mongo_queue"`
	Telegram         *struct {
		Token string `yaml:"token"`
	} `yaml:"telegram"`
	Bot           *utils.BotInfo        `yaml:"bot"`
	Sentry        *utils.SentrySettings `yaml:"sentry"`
	QueueConsumer *struct {
		FetchDelay int `yaml:"fetch_delay"`
	} `yaml:"queue_consumer"`
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
	RootDir = path.Join(confDir, "./..")
}
