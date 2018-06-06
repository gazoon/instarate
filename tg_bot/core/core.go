package core

import (
	"context"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/gazoon/go-utils/localization"
	"github.com/gazoon/go-utils/logging"
	"github.com/pkg/errors"
	"gopkg.in/telegram-bot-api.v4"
	"instarate/tg_bot/chats"
	"instarate/tg_bot/messages"
	"instarate/tg_bot/messenger"
	"instarate/tg_bot/models"
)

type MessageEnvelope struct {
	Type string      `mapstructure:"type"`
	Data interface{} `mapstructure:"data"`
}

type Bot struct {
	*logging.LoggerMixin
	chats     *chats.MongoStorage
	messenger *messenger.Telegram
	locales   *localization.Manager
}

func NewBot(chatsStorage *chats.MongoStorage, telegramMessenger *messenger.Telegram,
	locales *localization.Manager) *Bot {

	return &Bot{
		chats:       chatsStorage,
		messenger:   telegramMessenger,
		locales:     locales,
		LoggerMixin: logging.NewLoggerMixin("bot", nil),
	}
}

func (self *Bot) OnMessage(ctx context.Context, messageEnvelope *MessageEnvelope) {
	message, err := instantiateMessage(messageEnvelope)
	if err != nil {
		self.GetLogger(ctx).Error(err)
		return
	}
	defer func() {
		r := recover()
		if r == nil {
			return
		}
		err := r.(error)
		self.handleError(ctx, message, err)
	}()
	ctx = initializeContext(ctx, message)
	err = self.processMessage(ctx, message)
	if err != nil {
		panic(err)
	}
}

func (self *Bot) processMessage(ctx context.Context, message messages.Message) error {
	chat, err := self.getOrCreateChat(ctx, message)
	if err != nil {
		return err
	}
	if err = self.dispatchMessage(ctx, chat, message); err != nil {
		return err
	}
	if err = self.chats.Save(ctx, chat); err != nil {
		return err
	}
	return nil
}

func (self *Bot) handleError(ctx context.Context, message messages.Message, err error) {
	logger := self.GetLogger(ctx)
	logger.WithError(err).WithField("message", message).Error("Message processing failed")
	if _, ok := message.(messages.UserMessage); !ok {
		return
	}
	// user expects reaction, so we should show him a error
	chatId := message.GetChatId()
	errorText := self.gettext(models.DefaultLang, "unknown_error")
	if _, err := self.messenger.SendText(ctx, chatId, errorText, func(msg *tgbotapi.MessageConfig) {
		msg.ReplyMarkup = tgbotapi.ReplyKeyboardRemove{}
	}); err != nil {
		logger.Errorf("Can't send error message: %s", err)
	}
}

func (self *Bot) dispatchMessage(ctx context.Context, chat *models.Chat, message messages.Message) error {
	switch actualMessage := message.(type) {
	case *messages.TextMessage:
		self.onText(ctx, chat, actualMessage)
	case *messages.Callback:
		self.onCallback(ctx, chat, actualMessage)
	case *messages.NextPairTask:
		self.onNextPairTask(ctx, chat, actualMessage)
	case *messages.DailyActivationTask:
		self.onDailyActivationTask(ctx, chat, actualMessage)
	default:
		return errors.Errorf("can't dispatch message %T", actualMessage)
	}
	return nil
}

func (self *Bot) onText(ctx context.Context, chat *models.Chat, message *messages.TextMessage) error {
	self.GetLogger(ctx).Info("On text")
	text := self.gettext("ru", "propose_to_vote")
	if _, err := self.messenger.SendText(ctx, chat.Id, text); err != nil {
		return err
	}
	return nil
}

func (self *Bot) onCallback(ctx context.Context, chat *models.Chat, message *messages.Callback) error {

	return nil
}

func (self *Bot) onNextPairTask(ctx context.Context, chat *models.Chat, message *messages.NextPairTask) error {

	return nil
}

func (self *Bot) onDailyActivationTask(ctx context.Context, chat *models.Chat, message *messages.DailyActivationTask) error {

	return nil
}

func (self *Bot) getOrCreateChat(ctx context.Context, message messages.Message) (*models.Chat, error) {
	chatId := message.GetChatId()
	chat, err := self.chats.Get(ctx, chatId)
	if err != nil {
		return nil, err
	}
	if chat != nil {
		return chat, nil
	}
	userMessage, ok := message.(messages.UserMessage)
	if !ok {
		return nil, errors.Errorf(
			"chat doesn't exist, can't create a new one with message type %T",
			message)
	}
	membersNum, err := self.messenger.GetMembersNum(ctx, chatId)
	chat = models.NewChat(chatId, membersNum, userMessage.GetIsGroupChat())
	return chat, nil
}

func (self *Bot) gettext(lang, msgid string, vars ...interface{}) string {
	if lang == "en" && msgid == "place_in_competition" {
		place := vars[0].(int)
		var format string
		switch place % 10 {
		case 1:
			format = "%dst"
		case 2:
			format = "%dnd"
		case 3:
			format = "%drd"
		default:
			format = "%dth"
		}
		vars = []interface{}{fmt.Sprintf(format, place)}
	}
	return self.locales.Gettext(lang, msgid, vars...)
}

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