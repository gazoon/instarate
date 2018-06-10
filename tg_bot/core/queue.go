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

type MessagesPipeline struct {
	*logging.LoggerMixin
	queue     *queue.MongoReader
	onMessage func(context.Context, *MessageEnvelope) error
}

func NewMessagesPipe(incomingQueue *queue.MongoReader,
	onMessage func(context.Context, *MessageEnvelope) error) *MessagesPipeline {

	return &MessagesPipeline{
		queue:       incomingQueue,
		onMessage:   onMessage,
		LoggerMixin: logging.NewLoggerMixin("messages_pipe", nil),
	}
}

func (self *MessagesPipeline) Fetch(ctx context.Context) consumer.Process {
	defer func() {
		if r := recover(); r != nil {
			self.LogError(ctx, r)
		}
	}()
	msg, err := self.queue.GetNext(ctx)
	if err != nil {
		self.LogError(ctx, err)
		return nil
	}
	if msg == nil {
		return nil
	}
	return func() {
		self.process(ctx, msg)
	}
}

func (self *MessagesPipeline) process(ctx context.Context, msg *queue.ReadyMessage) {
	defer func() {
		if r := recover(); r != nil {
			self.LogError(ctx, r)
		}
		defer func() {
			if r := recover(); r != nil {
				self.LogError(ctx, r)
			}
		}()
		err := self.queue.FinishProcessing(ctx, msg.ProcessingId)
		if err != nil {
			self.LogError(ctx, err)
		}
	}()
	ctx = request.NewContext(ctx, msg.RequestId)
	ctx = logging.NewContext(ctx, logging.WithRequestID(msg.RequestId))
	messageEnvelope := &MessageEnvelope{}
	err := mapstructure.Decode(msg.Payload, messageEnvelope)
	if err != nil {
		self.LogError(ctx, errors.Wrapf(err, "queue message payload parsing %v", msg.Payload))
	}
	err = self.onMessage(ctx, messageEnvelope)
	if err != nil {
		self.LogError(ctx, err)
	}
}
