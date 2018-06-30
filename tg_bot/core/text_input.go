package core

import (
	"context"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"gopkg.in/telegram-bot-api.v4"
	"instarate/libs/competition"
	"instarate/libs/instagram"
	"instarate/tg_bot/messages"
	"instarate/tg_bot/models"
	"strconv"
	"strings"
	"time"
)

var (
	defaultCommandsLanguage = models.EnLanguage
)

type TextMessageHandler func(context.Context, *models.Chat, *messages.TextMessage) error

type TextCommand struct {
	Name    string
	Handler TextMessageHandler
}

var (
	noTextCommand = &TextCommand{}
)

var (
	leftVoteCmd               = "left_vote"
	rightVoteCmd              = "right_vote"
	nextTopGirlCmd            = "next_top_girl"
	showTopCmd                = "show_top"
	girlProfileCmd            = "girl_profile"
	startCmd                  = "start"
	addGirlCmd                = "add_girl"
	helpCmd                   = "help"
	chatSettingsCmd           = "chat_settings"
	globalCompetitionCmd      = "global_competition"
	celebritiesCompetitionCmd = "celebrities_competition"
	modelsCompetitionCmd      = "models_competition"
	regularCompetitionCmd     = "regular_competition"
	enableNotificationsCmd    = "enable_notifications"
	disableNotificationsCmd   = "disable_notifications"
	setRussianCmd             = "set_russian"
	setEnglishCmd             = "set_english"
	setVotingTimeoutCmd       = "set_voting_timeout"
)

func (self *Bot) onText(ctx context.Context, chat *models.Chat, message *messages.TextMessage) error {
	logger := self.GetLogger(ctx)

	if !message.IsAppealToBot(self.info.Name) &&
		!message.IsReplyToBot(self.info.Username) && message.IsGroupChat {
		logger.Debug("Message doesn't apply to the bot, skip")
		return nil
	}

	command := self.detectTextCommand(chat, message)
	if command != noTextCommand {
		logger.WithField("command_name", command.Name).Info("Handling text command")
		return command.Handler(ctx, chat, message)
	}

	logger.Info("Handling regular text message")
	return self.handleRegularText(ctx, chat, message)
}

func (self *Bot) detectTextCommand(chat *models.Chat, message *messages.TextMessage) *TextCommand {
	for _, cmdInfo := range self.commandsRegistry {
		for _, lang := range []string{defaultCommandsLanguage, chat.Language} {
			commandMsgstr := self.locales.GettextD(lang, "commands", cmdInfo.Name)
			for _, commandText := range strings.Split(commandMsgstr, "|") {
				if message.TextContains(commandText) {
					return cmdInfo
				}
			}
		}
	}
	return noTextCommand
}

func (self *Bot) buildCommandsList() []*TextCommand {
	commands := []*TextCommand{
		{leftVoteCmd, self.leftVoteCmd},
		{rightVoteCmd, self.rightVoteCmd},
		{nextTopGirlCmd, self.nextGirlCmd},
		{showTopCmd, self.showTopCmd},
		{girlProfileCmd, self.girlProfileCmd},
		{startCmd, self.startCmd},
		{addGirlCmd, self.addGirlCmd},
		{helpCmd, self.helpCmd},
		{chatSettingsCmd, self.chatSettingsCmd},
		{globalCompetitionCmd, self.globalCompetitionCmd},
		{celebritiesCompetitionCmd, self.celebritiesCompetitionCmd},
		{modelsCompetitionCmd, self.modelsCompetitionCmd},
		{regularCompetitionCmd, self.regularCompetitionCmd},
		{enableNotificationsCmd, self.enableNotificationsCmd},
		{disableNotificationsCmd, self.disableNotificationsCmd},
		{setRussianCmd, self.setRussianCmd},
		{setEnglishCmd, self.setEnglishCmd},
		{setVotingTimeoutCmd, self.setVotingTimeoutCmd},
	}
	return commands
}

func (self *Bot) startCmd(ctx context.Context, chat *models.Chat, message *messages.TextMessage) error {
	return self.trySendNextGirlsPair(ctx, chat)
}

func (self *Bot) addGirlCmd(ctx context.Context, chat *models.Chat, message *messages.TextMessage) error {
	photoLink := message.GetCommandArg()
	if photoLink == "" {
		_, err := self.messenger.SendText(ctx, chat.Id, self.gettext(chat, "add_girl_no_link"))
		return err
	}
	err := self.addGirl(ctx, chat, photoLink)
	if err == competition.BadPhotoLinkErr {
		_, err := self.messenger.SendText(ctx, chat.Id, self.gettext(chat, "add_girl_no_link"))
		return err
	} else if err == instagram.MediaForbidden {
		_, err := self.messenger.SendText(ctx, chat.Id, self.gettext(chat, "add_girl_media_forbidden"))
		return err
	}
	return err
}

func (self *Bot) showTopCmd(ctx context.Context, chat *models.Chat, message *messages.TextMessage) error {
	arg := message.GetCommandArg()
	offset := 0
	if arg != "" {
		startPosition, err := strconv.Atoi(arg)
		if err != nil {
			logger := self.GetLogger(ctx).WithFields(
				log.Fields{"command_arg": arg, "parsing_error": err},
			)
			logger.Warn("Show top command expects an int arg")
		} else if startPosition > 0 {
			offset = startPosition - 1
		}
	}
	return self.sendGirlFromTop(ctx, chat, offset)
}

func (self *Bot) nextGirlCmd(ctx context.Context, chat *models.Chat, message *messages.TextMessage) error {
	err := self.sendGirlFromTop(ctx, chat, chat.CurrentTopOffset+1)
	return err
}

func (self *Bot) girlProfileCmd(ctx context.Context, chat *models.Chat, message *messages.TextMessage) error {
	girlLink := message.GetCommandArg()
	if girlLink == "" {
		_, err := self.messenger.SendText(ctx, chat.Id, self.gettext(chat, "get_girl_no_username"))
		return err
	}
	girl, err := self.competition.GetCompetitor(ctx, chat.CompetitionCode, girlLink)
	if err == competition.BadProfileLinkErr {
		_, err := self.messenger.SendText(ctx, chat.Id, self.gettext(chat, "get_girl_no_username"))
		return err
	}
	if noGirlErr, ok := err.(*competition.CompetitorNotFound); ok {
		text := self.gettext(chat, "get_girl_no_girl", noGirlErr.Username, noGirlErr.Username)
		_, err := self.messenger.SendText(ctx, chat.Id, text)
		return err
	}
	if err != nil {
		return err
	}
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

func (self *Bot) helpCmd(ctx context.Context, chat *models.Chat, message *messages.TextMessage) error {
	_, err := self.messenger.SendText(ctx, chat.Id, self.gettext(chat, "help_message"))
	return err
}

func (self *Bot) chatSettingsCmd(ctx context.Context, chat *models.Chat, message *messages.TextMessage) error {
	_, err := self.messenger.SendText(ctx, chat.Id, self.gettext(chat, "chat_settings_commands"))
	return err
}

func (self *Bot) leftVoteCmd(ctx context.Context, chat *models.Chat, message *messages.TextMessage) error {
	if chat.LastMatch == nil {
		self.GetLogger(ctx).Warn("Received left vote command, but chat doesn't have a match")
		return nil
	}
	return self.processVoteMessage(ctx, chat, message, chat.LastMatch.LeftGirlUsername, chat.LastMatch.RightGirlUsername)
}

func (self *Bot) rightVoteCmd(ctx context.Context, chat *models.Chat, message *messages.TextMessage) error {
	if chat.LastMatch == nil {
		self.GetLogger(ctx).Warn("Received right vote command, but chat doesn't have a match")
		return nil
	}
	return self.processVoteMessage(ctx, chat, message, chat.LastMatch.RightGirlUsername, chat.LastMatch.LeftGirlUsername)
}

func (self *Bot) globalCompetitionCmd(ctx context.Context, chat *models.Chat, message *messages.TextMessage) error {
	text := self.gettext(chat, "global_competition_enabled")
	if _, err := self.messenger.SendText(ctx, chat.Id, text, func(msg *tgbotapi.MessageConfig) {
		msg.ReplyMarkup = tgbotapi.ReplyKeyboardRemove{}
	}); err != nil {
		return err
	}
	chat.CompetitionCode = competition.GlobalCompetition
	return nil
}

func (self *Bot) celebritiesCompetitionCmd(ctx context.Context, chat *models.Chat, message *messages.TextMessage) error {
	text := self.gettext(chat, "celebrities_competition_enabled")
	if _, err := self.messenger.SendText(ctx, chat.Id, text, func(msg *tgbotapi.MessageConfig) {
		msg.ReplyMarkup = tgbotapi.ReplyKeyboardRemove{}
	}); err != nil {
		return err
	}
	chat.CompetitionCode = competition.CelebritiesCompetition
	return nil
}

func (self *Bot) modelsCompetitionCmd(ctx context.Context, chat *models.Chat, message *messages.TextMessage) error {
	text := self.gettext(chat, "models_competition_enabled")
	if _, err := self.messenger.SendText(ctx, chat.Id, text, func(msg *tgbotapi.MessageConfig) {
		msg.ReplyMarkup = tgbotapi.ReplyKeyboardRemove{}
	}); err != nil {
		return err
	}
	chat.CompetitionCode = competition.ModelsCompetition
	return nil
}

func (self *Bot) regularCompetitionCmd(ctx context.Context, chat *models.Chat, message *messages.TextMessage) error {
	text := self.gettext(chat, "regular_competition_enabled")
	if _, err := self.messenger.SendText(ctx, chat.Id, text, func(msg *tgbotapi.MessageConfig) {
		msg.ReplyMarkup = tgbotapi.ReplyKeyboardRemove{}
	}); err != nil {
		return err
	}
	chat.CompetitionCode = competition.RegularCompetition
	return nil
}

func (self *Bot) enableNotificationsCmd(ctx context.Context, chat *models.Chat, message *messages.TextMessage) error {
	chat.SelfActivationAllowed = true
	_, err := self.messenger.SendText(ctx, chat.Id, self.gettext(chat, "daily_notification_enabled"))
	return err
}

func (self *Bot) disableNotificationsCmd(ctx context.Context, chat *models.Chat, message *messages.TextMessage) error {
	chat.SelfActivationAllowed = false
	if _, err := self.messenger.SendText(ctx, chat.Id, self.gettext(chat, "daily_notification_disabled")); err != nil {
		return err
	}
	err := self.scheduler.DeleteTask(ctx, chat.Id, messages.DailyActivationTaskType)
	return err
}

func (self *Bot) setRussianCmd(ctx context.Context, chat *models.Chat, message *messages.TextMessage) error {
	chat.Language = models.RuLanguage
	_, err := self.messenger.SendText(ctx, chat.Id, self.gettext(chat, "switch_to_language"))
	return err
}

func (self *Bot) setEnglishCmd(ctx context.Context, chat *models.Chat, message *messages.TextMessage) error {
	chat.Language = models.EnLanguage
	_, err := self.messenger.SendText(ctx, chat.Id, self.gettext(chat, "switch_to_language"))
	return err
}

func (self *Bot) setVotingTimeoutCmd(ctx context.Context, chat *models.Chat, message *messages.TextMessage) error {
	logger := self.GetLogger(ctx)
	arg := message.GetCommandArg()
	timeoutSeconds, err := strconv.Atoi(arg)
	if err != nil {
		logger.WithFields(log.Fields{"command_arg": arg, "parsing_error": err}).
			Warn("Can't parse voting timeout arg as an int")
		_, err := self.messenger.SendText(ctx, chat.Id, self.gettext(chat, "set_voting_timeout_no_number"))
		return err
	}
	timeout := time.Second * time.Duration(timeoutSeconds)
	if models.MinVotingTimeout < timeout || timeout > sessionDuration {
		lowerBound := models.MinVotingTimeoutSeconds
		upperBound := int(sessionDuration / time.Minute)
		logger.WithField("timeout", timeout).Warn("User entered invalid timeout")
		text := self.gettext(chat, "set_voting_timeout_out_of_range", lowerBound, upperBound)
		_, err := self.messenger.SendText(ctx, chat.Id, text)
		return err
	}
	chat.VotingTimeout = timeoutSeconds
	text := self.gettext(chat, "voting_timeout_is_set", timeoutSeconds)
	_, err = self.messenger.SendText(ctx, chat.Id, text)
	return err
}

func (self *Bot) handleRegularText(ctx context.Context, chat *models.Chat, message *messages.TextMessage) error {
	photoLink := message.GetLastWord()
	if photoLink == "" {
		return self.sendDontGetYou(ctx, chat)
	}
	err := self.addGirl(ctx, chat, photoLink)
	if err == competition.BadPhotoLinkErr || err == instagram.MediaForbidden {
		return self.sendDontGetYou(ctx, chat)
	}
	return err
}

func (self *Bot) sendDontGetYou(ctx context.Context, chat *models.Chat) error {
	_, err := self.messenger.SendText(ctx, chat.Id, self.gettext(chat, "dont_get_you"))
	return err
}
