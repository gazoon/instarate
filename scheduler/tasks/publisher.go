package tasks

import (
	"github.com/gazoon/go-utils"
	"instarate/scheduler/config"
	"path"
)

type Publisher struct {
	storage *Storage
}

func NewPublisher() (*Publisher, error) {
	conf := &config.Config{}
	configPath := path.Join(utils.GetCurrentFileDir(), "../config")
	err := utils.ParseConfig(configPath, conf)
	if err != nil {
		return nil, err
	}
	storage, err := NewStorage(conf.MongoTasks)
	if err != nil {
		return nil, err
	}
	return &Publisher{storage}, nil
}

func (self *Publisher) CreateTask(task *Task) error {
	return self.storage.CreateTask(task)
}

func (self *Publisher) CreateOrReplaceTask(task *Task) error {
	return self.storage.CreateOrReplaceTask(task)
}

func (self *Publisher) DeleteTask(chatId int, name string) error {
	return self.storage.DeleteTask(chatId, name)
}
