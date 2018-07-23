package services

import (
	. "instarate/scheduler/config"
	"instarate/scheduler/tasks"
)

func InitTaskReader() *tasks.Reader {
	taskStorage, err := tasks.NewReader(Config.MongoTasks)
	if err != nil {
		panic(err)
	}
	return taskStorage
}

func InitTaskPublisher() *tasks.Publisher {
	taskStorage, err := tasks.NewPublisher(Config.MongoTasks)
	if err != nil {
		panic(err)
	}
	return taskStorage
}
