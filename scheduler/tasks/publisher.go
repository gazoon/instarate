package tasks

import (
	"context"
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

func (self *Publisher) CreateTask(ctx context.Context, task *Task) error {
	return self.storage.CreateTask(ctx, task)
}

func (self *Publisher) CreateOrReplaceTask(ctx context.Context, task *Task) error {
	return self.storage.CreateOrReplaceTask(ctx, task)
}

func (self *Publisher) DeleteTask(ctx context.Context, chatId int, name string) error {
	return self.storage.DeleteTask(ctx, chatId, name)
}
