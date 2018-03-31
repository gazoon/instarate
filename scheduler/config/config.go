package config

import (
	"github.com/gazoon/go-utils"
)

type Config struct {
	utils.RootConfig `yaml:",inline"`
	MongoQueue       *utils.MongoDBSettings `yaml:"mongo_queue"`
	MongoTasks       *utils.MongoDBSettings `yaml:"mongo_tasks"`
	TasksConsumer    *struct {
		FetchDelay int `yaml:"fetch_delay"`
	} `yaml:"tasks_consumer"`
}

func (self *Config) String() string {
	return utils.ObjToString(self)
}
