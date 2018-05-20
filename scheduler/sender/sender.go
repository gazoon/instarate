package sender

import (
	"instarate/scheduler/tasks"

	"context"
	"github.com/gazoon/bot_libs/queue"
	"github.com/gazoon/go-utils"
	"github.com/gazoon/go-utils/logging"
)

type Sender struct {
	*logging.LoggerMixin
	queueWriter *queue.MongoWriter
	queueName   string
}

func New(mongoSettings *utils.MongoDBSettings) (*Sender, error) {
	writer, err := queue.NewMongoWriter(mongoSettings)
	if err != nil {
		return nil, err
	}
	logger := logging.NewLoggerMixin("sender", nil)
	return &Sender{queueWriter: writer, LoggerMixin: logger, queueName: mongoSettings.Collection}, nil
}

func (self *Sender) SendTask(ctx context.Context, data interface{}) {
	task := data.(*tasks.Task)
	ctx = utils.FillContext(ctx)
	logger := self.GetLogger(ctx)
	logger.Debugf("Send task to the queue: %s", task)
	self.queueWriter.Put(ctx, self.queueName, task.ChatId, task)
}
