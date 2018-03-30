package worker

import (
	"context"
	"github.com/gazoon/bot_libs/queue"
	"github.com/gazoon/go-utils"
	"github.com/gazoon/go-utils/logging"
	"gopkg.in/telegram-bot-api.v4"
	"instarate/tg_gateway/config"
)

type Worker struct {
	*logging.LoggerMixin
	queue *queue.MongoWriter
}

func New(config *config.Config) (*Worker, error) {
	mongoWriter, err := queue.NewMongoWriter(config.MongoQueue)
	if err != nil {
		return nil, err
	}
	logger := logging.NewLoggerMixin("worker", nil)
	return &Worker{queue: mongoWriter, LoggerMixin: logger}, nil
}

func (self *Worker) ProcessUpdate(ctx context.Context, queueName string, update *tgbotapi.Update) error {
	logger := self.GetLogger(ctx)
	logger.Debugf("Receive update %s", utils.ObjToString(update))
	chatId, message := buildQueueMessage(update)
	if message == nil {
		logger.Debugf("Unsupported update, skip")
		return nil
	}
	return self.queue.Put(ctx, queueName, chatId, message)
}

func buildQueueMessage(update *tgbotapi.Update) (int, map[string]interface{}) {
	if update.Message != nil && update.Message.Text != "" {
		message := update.Message
		var replyToConverted map[string]interface{}
		if message.ReplyToMessage != nil {
			replyToConverted = convertMessage(message.ReplyToMessage)
		}
		messageConverted := convertMessage(message)
		messageConverted["reply_to"] = replyToConverted
		return int(message.Chat.ID), map[string]interface{}{
			"type": "text",
			"data": messageConverted,
		}
	} else if update.CallbackQuery != nil && update.CallbackQuery.Message != nil {
		callback := update.CallbackQuery
		return int(callback.Message.Chat.ID), map[string]interface{}{
			"type": "callback",
			"data": map[string]interface{}{
				"callback_id":   callback.ID,
				"user":          convertUser(callback.From),
				"parent_msg_id": callback.Message.MessageID,
				"payload":       callback.Data,
				"chat_id":       callback.Message.Chat.ID,
				"is_group_chat": isGroupChat(callback.Message),
			},
		}
	}
	return 0, nil

}

func isGroupChat(message *tgbotapi.Message) bool {
	return message.Chat.Type != "private"
}

func convertMessage(message *tgbotapi.Message) map[string]interface{} {
	return map[string]interface{}{
		"user":          convertUser(message.From),
		"text":          message.Text,
		"chat_id":       int(message.Chat.ID),
		"is_group_chat": isGroupChat(message),
	}
}

func convertUser(user *tgbotapi.User) map[string]interface{} {
	name := user.FirstName
	if user.LastName != "" {
		name += " " + user.LastName
	}
	return map[string]interface{}{
		"id":       user.ID,
		"name":     name,
		"username": user.UserName,
	}
}
