package core

import (
	"context"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/gazoon/go-utils"
	"github.com/gazoon/go-utils/localization"
	"github.com/gazoon/go-utils/logging"
	"github.com/pkg/errors"
	"gopkg.in/telegram-bot-api.v4"
	"instarate/libs/competition"
	"instarate/scheduler/tasks"
	"instarate/tg_bot/cache"
	"instarate/tg_bot/chats"
	"instarate/tg_bot/concatenation"
	"instarate/tg_bot/messages"
	"instarate/tg_bot/messenger"
	"instarate/tg_bot/models"
	"reflect"
	"strings"
	"time"
)

const (
	leftVote  = "LEFT_VOTE"
	rightVote = "RIGHT_VOTE"
)

var (
	sessionDuration = 20 * time.Minute
)

type MessageEnvelope struct {
	Type string      `mapstructure:"type"`
	Data interface{} `mapstructure:"data"`
}

type Bot struct {
	*logging.LoggerMixin
	chats            *chats.MongoStorage
	messenger        *messenger.Telegram
	scheduler        *tasks.Publisher
	locales          *localization.Manager
	commandsRegistry []*TextCommand
	botInfo          *utils.BotInfo
	competition      *competition.Competition
	cache            *cache.Mongo
}

func NewBot(competitionAPI *competition.Competition, chatsStorage *chats.MongoStorage,
	cacheStorage *cache.Mongo, telegramMessenger *messenger.Telegram,
	scheduler *tasks.Publisher, locales *localization.Manager, info *utils.BotInfo) *Bot {

	b := &Bot{
		chats:       chatsStorage,
		messenger:   telegramMessenger,
		scheduler:   scheduler,
		locales:     locales,
		LoggerMixin: logging.NewLoggerMixin("bot", nil),
		botInfo:     info,
		competition: competitionAPI,
		cache:       cacheStorage,
	}
	b.commandsRegistry = b.buildCommandsList()
	return b
}

func (self *Bot) OnMessage(ctx context.Context, messageEnvelope *MessageEnvelope) error {
	message, err := instantiateMessage(messageEnvelope)
	if err != nil {
		return err
	}
	defer func() {
		if r := recover(); r != nil {
			self.sendError(ctx, message)
			panic(r)
		}
	}()
	ctx = initializeContext(ctx, message)
	logger := self.GetLogger(ctx)
	logger.WithFields(log.Fields{
		"message":      message,
		"message_type": messageEnvelope.Type,
	}).Info("Process message")
	err = self.processMessage(ctx, message)
	if err != nil {
		self.sendError(ctx, message)
	}
	return err
}

func (self *Bot) processMessage(ctx context.Context, message messages.Message) error {
	chat, err := self.getOrCreateChat(ctx, message)
	if err != nil {
		return err
	}
	if err = self.dispatchMessage(ctx, chat, message); err != nil {
		if resetErr := self.chats.ResetState(ctx, chat.Id); resetErr != nil {
			self.LogError(ctx, resetErr)
		}
		return err
	}
	if err = self.chats.Save(ctx, chat); err != nil {
		return err
	}
	return nil
}

func (self *Bot) sendError(ctx context.Context, message messages.Message) {
	if _, ok := message.(messages.UserMessage); !ok {
		logger := self.GetLogger(ctx).WithField("message_type", reflect.TypeOf(message))
		logger.Info(
			"Error on not a user message, no need to send the error message to the chat",
		)
		return
	}
	// user expects reaction, so we should show him a error
	chatId := message.GetChatId()
	errorText := self.locales.Gettext(models.DefaultLang, "unknown_error")
	_, err := self.messenger.SendText(ctx, chatId, errorText, func(msg *tgbotapi.MessageConfig) {
		msg.ReplyMarkup = tgbotapi.ReplyKeyboardRemove{RemoveKeyboard: true}
	})
	if err != nil {
		self.LogError(ctx, errors.Wrap(err, "during error message sending, another occurred"))
	}
}

func (self *Bot) dispatchMessage(ctx context.Context, chat *models.Chat, message messages.Message) error {
	switch actualMessage := message.(type) {
	case *messages.TextMessage:
		return self.onText(ctx, chat, actualMessage)
	case *messages.NewChatUsers:
		return self.onNewChatUsers(ctx, chat, actualMessage)
	case *messages.Callback:
		return self.onCallback(ctx, chat, actualMessage)
	case *messages.NextPairTask:
		return self.onNextPairTask(ctx, chat, actualMessage)
	case *messages.DailyActivationTask:
		return self.onDailyActivationTask(ctx, chat, actualMessage)
	case *messages.CancelKeyboardTask:
		return self.onCancelKeyboardTask(ctx, chat, actualMessage)
	default:
		return errors.Errorf("can't dispatch message %T", actualMessage)
	}
	return nil
}
func (self *Bot) onNewChatUsers(ctx context.Context, chat *models.Chat, message *messages.NewChatUsers) error {
	if message.IsBotAdded(self.botInfo.Username) {
		return self.sendHelpText(ctx, chat)
	}
	self.GetLogger(ctx).WithField("new_members", message.NewUsers).
		Info("New chat members don't contain the bot, skip")
	return nil
}

func (self *Bot) onCallback(ctx context.Context, chat *models.Chat, message *messages.Callback) error {
	self.GetLogger(ctx).WithField("callback", message).Info("Callback received")
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

func (self *Bot) trySendNextGirlsPair(ctx context.Context, chat *models.Chat) error {
	if chat.LastMatch == nil {
		return self.sendNextGirlsPair(ctx, chat)
	}
	timeToShow := chat.LastMatch.ShownAt.Add(time.Duration(chat.GetVotingTimeout()) * time.Second)
	if timeToShow.Before(utils.UTCNow()) {
		return self.sendNextGirlsPair(ctx, chat)
	}
	task := tasks.NewTask(
		messages.NextPairTaskType, chat.Id, timeToShow,
		map[string]interface{}{"last_match_message_id": chat.LastMatch.MessageId},
	)
	self.GetLogger(ctx).WithField("time_to_show", timeToShow).Info("Schedule send next pair task")
	return self.scheduler.CreateOrReplaceTask(ctx, task)
}

func (self *Bot) sendNextGirlsPair(ctx context.Context, chat *models.Chat) error {
	votersGroupId := buildVotersGroupId(chat.Id)
	girl1, girl2, err := self.competition.GetNextPair(ctx, chat.CompetitionCode, votersGroupId)
	if err == competition.GetNextPairNoAttemptsErr {
		self.LogErrorWithFields(ctx, err, map[string]interface{}{
			"chat_id": chat.Id, "competition": chat.CompetitionCode})
		text := self.gettext(chat, "no_more_girls_in_competition")
		_, err = self.messenger.SendText(ctx, chat.Id, text, func(msg *tgbotapi.MessageConfig) {
			msg.ReplyMarkup = tgbotapi.ReplyKeyboardRemove{RemoveKeyboard: true}
		})
		chat.ResetLastMatch()
		return err
	}
	leftGirl, rightGirl := girl1, girl2
	var tgFileId string
	tgFileId, err = self.getMatchPhoto(ctx, girl1.PhotoPath, girl2.PhotoPath)
	if err != nil {
		return err
	}
	if tgFileId == "" {
		tgFileId, err = self.getMatchPhoto(ctx, girl2.PhotoPath, girl1.PhotoPath)
		if err != nil {
			return err
		}
		leftGirl, rightGirl = girl2, girl1
	}
	captionText := fmt.Sprintf("%s vs %s", leftGirl.GetProfileLink(), rightGirl.GetProfileLink())
	logger := self.GetLogger(ctx)
	logger.WithField("caption_text", captionText).Info("Send pair match")
	keyboard := tgbotapi.ReplyKeyboardMarkup{
		ResizeKeyboard: true, Keyboard: [][]tgbotapi.KeyboardButton{
			{tgbotapi.KeyboardButton{Text: "Left"}, tgbotapi.KeyboardButton{Text: "Right"}},
			{tgbotapi.KeyboardButton{Text: "Next pair"}},
		}}
	var messageId int
	if tgFileId != "" {
		logger.WithField("tg_file_id", tgFileId).Info("Use cached match photo")
		messageId, _, err = self.messenger.SendPhoto(ctx, chat.Id, tgFileId, func(settings *tgbotapi.PhotoConfig) {
			settings.Caption = captionText
			settings.ReplyMarkup = keyboard
		})
		if err != nil {
			return err
		}
	} else {
		leftPhotoUrl := self.competition.GetPhotoUrl(leftGirl)
		rightPhotoUrl := self.competition.GetPhotoUrl(rightGirl)
		logger.WithFields(log.Fields{"left_url": leftPhotoUrl, "right_url": rightPhotoUrl}).Info("Concatenate photos")
		imageBody, err := concatenation.Concatenate(ctx, leftPhotoUrl, rightPhotoUrl)
		if err != nil {
			return err
		}
		messageId, tgFileId, err = self.messenger.SendBinaryPhoto(ctx, chat.Id, imageBody, func(settings *tgbotapi.PhotoConfig) {
			settings.Caption = captionText
			settings.ReplyMarkup = keyboard
		})
		if err != nil {
			return err
		}
		if err = self.setMatchPhoto(ctx, leftGirl.PhotoPath, rightGirl.PhotoPath, tgFileId); err != nil {
			return err
		}
	}
	if err := self.scheduleDailyNotification(ctx, chat); err != nil {
		return err
	}
	if err := self.scheduleCancelKeyboard(ctx, chat); err != nil {
		return err
	}
	chat.LastMatch = models.NewMatch(messageId, leftGirl.Username, rightGirl.Username)
	return nil
}

func (self *Bot) processVoteMessage(ctx context.Context, chat *models.Chat,
	message messages.UserMessage, voteSide string) error {
	if chat.LastMatch == nil {
		self.GetLogger(ctx).Infof("Received %s vote command, but chat doesn't have a match", voteSide)
		return self.sendDontGetYou(ctx, chat)
	}
	var winner, loser string
	if voteSide == leftVote {
		winner, loser = chat.LastMatch.LeftGirlUsername, chat.LastMatch.RightGirlUsername
	} else {
		winner, loser = chat.LastMatch.RightGirlUsername, chat.LastMatch.LeftGirlUsername
	}
	votersGroupId := buildVotersGroupId(message.GetChatId())
	voterId := buildVoterId(message.GetUser())
	_, _, err := self.competition.Vote(ctx, chat.CompetitionCode, votersGroupId, voterId, winner, loser)
	if err == competition.AlreadyVotedErr {
		return nil
	}
	if _, ok := err.(*competition.CompetitorNotFound); ok {
		return errors.Wrap(err, "winner or loser is no longer in the competition")
	}
	if err != nil {
		return err
	}
	if err := self.scheduleCancelKeyboard(ctx, chat); err != nil {
		return err
	}
	return nil
}

func (self *Bot) scheduleDailyNotification(ctx context.Context, chat *models.Chat) error {
	logger := self.GetLogger(ctx)
	if !chat.SelfActivationAllowed {
		logger.Info("Self activation disabled for the chat")
		return nil
	}
	activationTime := utils.UTCNow().Add(24*time.Hour - sessionDuration)
	task := tasks.NewTaskWithoutArgs(messages.DailyActivationTaskType, chat.Id, activationTime)
	logger.WithField("activation_time", activationTime).Info("Schedule next day activation")
	return self.scheduler.CreateOrReplaceTask(ctx, task)
}

func (self *Bot) scheduleCancelKeyboard(ctx context.Context, chat *models.Chat) error {
	logger := self.GetLogger(ctx)
	activationTime := utils.UTCNow().Add(models.CancelKeyboardAfter)
	task := tasks.NewTaskWithoutArgs(messages.CancelKeyboardTaskType, chat.Id, activationTime)
	logger.WithField("activation_time", activationTime).Info("Schedule cancel keyboard")
	return self.scheduler.CreateOrReplaceTask(ctx, task)
}

func (self *Bot) sendGirlFromTop(ctx context.Context, chat *models.Chat) error {
	offset := chat.CurrentTopOffset
	amount := 2
	girls, err := self.competition.GetTop(ctx, chat.CompetitionCode, amount, offset)
	if err != nil {
		return err
	}
	if len(girls) == 0 {
		totalNumber, err := self.competition.GetCompetitorsNumber(ctx, chat.CompetitionCode)
		if err != nil {
			return err
		}
		logger := self.GetLogger(ctx).WithFields(log.Fields{"offset": offset, "total_number": totalNumber})
		logger.Warn("Girl offset exceeds the total number of girls")
		_, err = self.messenger.SendText(ctx, chat.Id, self.gettext(chat, "no_more_girls_in_top"), func(msg *tgbotapi.MessageConfig) {
			msg.ReplyMarkup = tgbotapi.ReplyKeyboardRemove{RemoveKeyboard: true}
		})
		chat.ResetTopOffset()
		return err
	}
	var keyboard interface{} = tgbotapi.ReplyKeyboardMarkup{
		ResizeKeyboard: true,
		Keyboard: [][]tgbotapi.KeyboardButton{
			{tgbotapi.KeyboardButton{Text: "Next girl"}},
		}}
	if len(girls) == 1 {
		keyboard = tgbotapi.ReplyKeyboardRemove{RemoveKeyboard: true}
		chat.ResetTopOffset()
	}
	girl := girls[0]
	caption := self.getPlaceInCompetitionText(chat, offset+1) + girl.GetProfileLink()
	err = self.sendSingleGirlPhoto(ctx, chat, girl, func(settings *tgbotapi.PhotoConfig) {
		settings.ReplyMarkup = keyboard
		settings.Caption = caption
	})
	if err != nil {
		return err
	}
	err = self.scheduleCancelKeyboard(ctx, chat)
	return err
}

func (self *Bot) addGirl(ctx context.Context, chat *models.Chat, photoLink string) error {
	if !strings.HasPrefix(photoLink, "https://instagram.com/p/") &&
		!strings.HasPrefix(photoLink, "https://www.instagram.com/p/") {
		return competition.BadPhotoLinkErr
	}
	profile, err := self.competition.Add(ctx, photoLink)
	if err == competition.NotPhotoMediaErr {
		_, err = self.messenger.SendText(ctx, chat.Id, self.gettext(chat, "add_girl_not_photo"))
		return err
	}
	if err != nil && err != competition.ProfileExistsErr {
		return err
	}

	var msgid string
	if err == competition.ProfileExistsErr {
		msgid = "girl_already_added"
	} else {
		msgid = "girl_successfully_added"
	}
	text := self.gettext(chat, msgid, profile.Username, profile.GetProfileLink())
	_, err = self.messenger.SendMarkdown(ctx, chat.Id, text)
	return err
}

func (self *Bot) sendGirlProfile(ctx context.Context, chat *models.Chat, girl *competition.InstCompetitor) error {
	titleText := fmt.Sprintf("[%s](%s)", girl.Username, girl.GetProfileLink())
	if _, err := self.messenger.SendMarkdown(ctx, chat.Id, titleText); err != nil {
		return err
	}
	if err := self.sendSingleGirlPhoto(ctx, chat, girl); err != nil {
		return err
	}
	place, err := self.competition.GetPosition(ctx, girl)
	if err != nil {
		return err
	}
	profileText := self.gettext(chat, "girl_statistics", place, girl.Wins, girl.Loses)
	if _, err := self.messenger.SendText(ctx, chat.Id, profileText); err != nil {
		return err
	}
	return nil
}

func (self *Bot) getPlaceInCompetitionText(chat *models.Chat, place int) string {
	var vars []interface{}
	if chat.Language == models.EnLanguage {
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
	} else {
		vars = []interface{}{place}
	}
	return self.gettext(chat, "place_in_competition", vars...)
}

func (self *Bot) sendSingleGirlPhoto(ctx context.Context, chat *models.Chat,
	girl *competition.InstCompetitor, opts ...func(settings *tgbotapi.PhotoConfig)) error {
	tgFileId, ok, err := self.cache.Get(ctx, girl.PhotoPath)
	if err != nil {
		return err
	}
	isNew := false
	photoUri := tgFileId
	if !ok {
		photoUri = self.competition.GetPhotoUrl(girl)
		isNew = true
	}
	_, tgFileId, err = self.messenger.SendPhoto(ctx, chat.Id, photoUri, opts...)
	if err != nil {
		return err
	}
	if isNew {
		if err := self.cache.Set(ctx, girl.PhotoPath, tgFileId); err != nil {
			return err
		}
	}
	return nil
}

func (self *Bot) gettext(chat *models.Chat, msgid string, vars ...interface{}) string {
	return self.locales.Gettext(chat.Language, msgid, vars...)
}
