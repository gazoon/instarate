package core

import (
	"context"
	"github.com/gazoon/bot_libs/queue"
	"github.com/gazoon/go-utils/consumer"
	"github.com/gazoon/go-utils/logging"
	"github.com/gazoon/go-utils/request"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

type MessagesPipe struct {
	*logging.LoggerMixin
	queue *queue.MongoReader
	bot   *Bot
}

func NewMessagesPipe(incomingQueue *queue.MongoReader, bot *Bot) *MessagesPipe {
	return &MessagesPipe{
		queue:       incomingQueue,
		bot:         bot,
		LoggerMixin: logging.NewLoggerMixin("messages_pipe", nil)}
}

func (self *MessagesPipe) Fetch(ctx context.Context) consumer.Process {
	msg, err := self.queue.GetNext()
	if err != nil {
		self.Logger.Errorf("Cant obtain message from the incoming queue: %s", err)
		return nil
	}
	if msg == nil {
		return nil
	}
	return func() {
		self.process(ctx, msg)
	}
}

func (self *MessagesPipe) process(ctx context.Context, msg *queue.ReadyMessage) {
	defer func() {
		if r := recover(); r != nil {
			self.GetLogger(ctx).Errorf("Can't process queue message; error: %v", r)
		}
		err := self.queue.FinishProcessing(ctx, msg.ProcessingId)
		if err != nil {
			self.GetLogger(ctx).Errorf(
				"Can't finish queue message processing; error: %s, processing_id: %s",
				err, msg.ProcessingId,
			)
		}
	}()
	ctx = request.NewContext(ctx, msg.RequestId)
	ctx = logging.NewContext(ctx, logging.WithRequestID(msg.RequestId))
	messageEnvelope := &MessageEnvelope{}
	err := mapstructure.Decode(msg.Payload, messageEnvelope)
	if err != nil {
		panic(errors.Wrapf(err, "queue message payload parsing %v", msg.Payload))
	}
	self.bot.OnMessage(ctx, messageEnvelope)
}
