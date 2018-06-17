package config

import (
	"github.com/gazoon/go-utils"
)

type Config struct {
	utils.RootConfig `yaml:",inline"`
	MongoCache       *utils.MongoDBSettings `yaml:"mongo_cache"`
	MongoChats       *utils.MongoDBSettings `yaml:"mongo_chats"`
	MongoQueue       *utils.MongoDBSettings `yaml:"mongo_queue"`
	Telegram         *struct {
		Token string `yaml:"token"`
	} `yaml:"telegram"`
	Bot *utils.BotInfo `yaml:"bot"`

	QueueConsumer *struct {
		FetchDelay int `yaml:"fetch_delay"`
	} `yaml:"queue_consumer"`
}

func (self *Config) String() string {
	return utils.ObjToString(self)
}
