package core

import (
	"context"
	"fmt"
	"github.com/gazoon/go-utils/localization"
	"github.com/gazoon/go-utils/logging"
	"github.com/gazoon/go-utils/request"
	"github.com/pkg/errors"
	"gopkg.in/telegram-bot-api.v4"
	"instarate/tg_bot/chats"
	"instarate/tg_bot/messages"
	"instarate/tg_bot/messenger"
	"instarate/tg_bot/models"
)

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

func (self *Bot) OnMessage(ctx context.Context, data interface{}) {
	messageEnvelope, ok := data.(map[string]interface{})
	if !ok {
		self.Logger.Errorf("Message envelope must be a map, got: %v", data)
		return
	}
	message, err := messages.Instantiate(messageEnvelope)
	if err != nil {
		self.Logger.Error(err)
		return
	}
	ctx = initializeContext(ctx, messageEnvelope, message)
	defer func() {
		if r := recover(); r != nil {
			err := errors.Errorf("panic recovered: %v", r)
			self.handleError(ctx, message, err)
		}
	}()
	err = self.processMessage(ctx, message)
	if err != nil {
		self.handleError(ctx, message, err)
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
	logger.Error(err)
	if _, ok := message.(messages.UserMessage); !ok {
		return
	}
	// user expects reaction, so we should show him a error
	chatId := message.GetChatId()
	errorText := self.gettext(models.DefaultLang, "unknown_error")
	_, err = self.messenger.SendText(ctx, chatId, errorText, func(msg *tgbotapi.MessageConfig) {
		msg.ReplyMarkup = tgbotapi.ReplyKeyboardRemove{}
	})
	if err != nil {
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

func initializeContext(ctx context.Context, messageEnvelope map[string]interface{},
	message messages.Message) context.Context {

	requestId, _ := messageEnvelope["request_id"].(string)
	if requestId == "" {
		requestId = request.NewRequestId()
	}
	ctx = request.NewContext(ctx, requestId)
	logger := logging.WithRequestID(requestId).WithField("chat_id", message.GetChatId())

	ctx = logging.NewContext(ctx, logger)
	return ctx
}
