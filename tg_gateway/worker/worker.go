package worker

import (
	"context"
	"github.com/gazoon/bot_libs/queue"
	"instarate/tg_gateway/config"
)

type Worker struct {
	config *config.Config
	queue  *queue.MongoWriter
}

func New(config *config.Config) (*Worker, error) {
	mongoWriter, err := queue.NewMongoWriter(config.MongoQueue)
	if err != nil {
		return nil, err
	}
	return &Worker{config, mongoWriter}, nil
}

func (self *Worker) ProcessUpdate(ctx context.Context, update interface{}) error {
	return nil
}
