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

func (self *Sender) SendTask(ctx context.Context, task *tasks.Task) error {
	ctx = utils.FillContext(ctx)
	logger := self.GetLogger(ctx)
	logger.Debugf("Send task to the queue: %s", task)
	messageData := map[string]interface{}{
		"chat_id": task.ChatId, "do_at": task.DoAt,
	}
	for k, v := range task.Args {
		messageData[k] = v
	}
	message := map[string]interface{}{
		"type": task.Name, "data": messageData,
	}
	err := self.queueWriter.Put(ctx, self.queueName, task.ChatId, message)
	return err
}
