package core

import (
	"context"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/gazoon/go-utils/logging"
	"github.com/pkg/errors"
	"instarate/tg_bot/messages"
)

func initializeContext(ctx context.Context, message messages.Message) context.Context {
	logger := log.WithField("chat_id", message.GetChatId())
	ctx = logging.NewContext(ctx, logger)
	return ctx
}

func instantiateMessage(messageEnvelope *MessageEnvelope) (messages.Message, error) {
	messageType := messageEnvelope.Type
	messageData := messageEnvelope.Data
	var message messages.Message
	var err error
	switch messageType {
	case messages.TextType:
		message, err = messages.TextMessageFromData(messageData)
	case messages.CallbackType:
		message, err = messages.CallbackFromData(messageData)
	case messages.NextPairTaskType:
		message, err = messages.NextPairTaskFromData(messageData)
	case messages.DailyActivationTaskType:
		message, err = messages.DailyActivationTaskFromData(messageData)
	default:
		return nil, errors.Errorf("unknown message type: %s", messageType)
	}
	if err != nil {
		return nil, errors.Wrapf(err, "can't build message type of %s", messageType)
	}
	return message, nil
}

func buildVotersGroupId(chatId int) string {
	return fmt.Sprintf("tg_chat:%d", chatId)
}

func buildVoterId(user *messages.User) string {
	return fmt.Sprintf("tg_user:%d", user.Id)
}
