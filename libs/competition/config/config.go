package config

import (
	"github.com/gazoon/go-utils"
)

var (
	Config *ConfigSchema
)

type ConfigSchema struct {
	MongoCompetitors *utils.MongoDBSettings `yaml:"mongo_competitors"`
	MongoProfiles    *utils.MongoDBSettings `yaml:"mongo_profiles"`
	MongoVoters      *utils.MongoDBSettings `yaml:"mongo_voters"`
	GoogleStorage    *struct {
		BucketName string `yaml:"bucket_name"`
	} `yaml:"google_storage"`
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
