package config

import (
	"github.com/gazoon/go-utils"
)

type Config struct {
	MongoCompetitors *utils.MongoDBSettings `yaml:"mongo_competitors"`
	MongoProfiles    *utils.MongoDBSettings `yaml:"mongo_profiles"`
	MongoVoters      *utils.MongoDBSettings `yaml:"mongo_voters"`
	GoogleStorage    *struct {
		BucketName string `yaml:"bucket_name"`
	} `yaml:"google_storage"`
}

func (self *Config) String() string {
	return utils.ObjToString(self)
}
