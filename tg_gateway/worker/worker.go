package worker

import (
	"context"
	"github.com/gazoon/go-utils"
	"github.com/gazoon/go-utils/logging"
	"github.com/gazoon/go-utils/queue"
	"gopkg.in/telegram-bot-api.v4"
	"instarate/tg_gateway/config"
)

type Worker struct {
	*logging.LoggerMixin
	queue *queue.MongoWriter
}

func New(config *config.ConfigSchema) (*Worker, error) {
	mongoWriter, err := queue.NewMongoWriter(config.MongoQueue)
	if err != nil {
		return nil, err
	}
	logger := logging.NewLoggerMixin("worker", nil)
	return &Worker{queue: mongoWriter, LoggerMixin: logger}, nil
}

func (self *Worker) ProcessUpdate(ctx context.Context, update *tgbotapi.Update) error {
	logger := self.GetLogger(ctx)
	logger.WithField("update", utils.ObjToString(update)).Info("Update received")
	chatId, message := buildQueueMessage(update)
	if message == nil {
		logger.Info("Unsupported update, skip")
		return nil
	}
	return self.queue.Put(ctx, chatId, message)
}

func buildQueueMessage(update *tgbotapi.Update) (int, map[string]interface{}) {
	if update.Message != nil {
		message := update.Message
		messageConverted := convertMessage(message)
		var messageType string
		if message.Text != "" {
			messageType = "text"
			messageConverted["text"] = message.Text
			if message.ReplyToMessage != nil {
				replyToConverted := convertUser(message.ReplyToMessage.From)
				messageConverted["reply_to_user"] = replyToConverted
			} else {
				messageConverted["reply_to_user"] = nil
			}
		} else if message.NewChatMembers != nil {
			messageType = "new_chat_users"
			newMembersRaw := *message.NewChatMembers
			newMembers := make([]map[string]interface{}, len(newMembersRaw))
			for i := range newMembersRaw {
				newMembers[i] = convertUser(&newMembersRaw[i])
			}
			messageConverted["new_users"] = newMembers
		} else {
			return 0, nil
		}
		return int(message.Chat.ID), map[string]interface{}{
			"type": messageType,
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
		"chat_id":       int(message.Chat.ID),
		"message_id":    message.MessageID,
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
