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
	onMessage func(context.Context, *MessageEnvelope)
}

func NewMessagesPipe(incomingQueue *queue.MongoReader,
	onMessage func(context.Context, *MessageEnvelope)) *MessagesPipeline {

	return &MessagesPipeline{
		queue:       incomingQueue,
		onMessage:   onMessage,
		LoggerMixin: logging.NewLoggerMixin("messages_pipe", nil),
	}
}

func (self *MessagesPipeline) Fetch(ctx context.Context) consumer.Process {
	msg, err := self.queue.GetNext(ctx)
	if err != nil {
		logger := self.GetLogger(ctx)
		logger.WithError(err).Error("Cant obtain message from the incoming queue")
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
			err := r.(error)
			logger := self.GetLogger(ctx)
			logger.WithError(err).Error("Can't process queue message")
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
	self.onMessage(ctx, messageEnvelope)
}
