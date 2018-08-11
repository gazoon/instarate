package core

import (
	"context"
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
	randomGirlCmd             = "random_girl"
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

	if !message.IsAppealToBot(self.botInfo.Name) &&
		!message.IsReplyToBot(self.botInfo.Username) && message.IsGroupChat {
		logger.Info("Message doesn't apply to the bot, skip")
		return nil
	}

	command := self.detectTextCommand(chat, message)
	if command != noTextCommand {
		logger.WithField("command_name", command.Name).Info("Handling text command")
		return command.Handler(ctx, chat, message)
	}

	logger.Info("Handling regular text message")
	return self.handleFreeInput(ctx, chat, message)
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
		{randomGirlCmd, self.randomGirlCmd},
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
		_, err := self.messenger.SendText(ctx, chat.Id, self.gettext(chat, "add_girl_bad_link"))
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
	chat.ResetLastMatch()
	chat.CurrentTopOffset = offset
	return self.sendGirlFromTop(ctx, chat)
}

func (self *Bot) nextGirlCmd(ctx context.Context, chat *models.Chat, message *messages.TextMessage) error {
	chat.CurrentTopOffset += 1
	err := self.sendGirlFromTop(ctx, chat)
	return err
}

func (self *Bot) girlProfileCmd(ctx context.Context, chat *models.Chat, message *messages.TextMessage) error {
	girlLink := message.GetCommandArg()
	girl, err := self.competition.GetCompetitor(ctx, chat.CompetitionCode, girlLink)
	if err == competition.BadProfileLinkErr {
		_, err := self.messenger.SendText(ctx, chat.Id, self.gettext(chat, "get_girl_no_username"))
		return err
	}
	if noGirlErr, ok := err.(*competition.CompetitorNotFound); ok {
		text := self.gettext(chat, "get_girl_no_girl", noGirlErr.Username, noGirlErr.Username)
		_, err := self.messenger.SendMarkdown(ctx, chat.Id, text)
		return err
	}
	if err != nil {
		return err
	}
	err = self.sendGirlProfile(ctx, chat, girl)
	return err
}

func (self *Bot) randomGirlCmd(ctx context.Context, chat *models.Chat, message *messages.TextMessage) error {
	girl, err := self.competition.GetRandomCompetitor(ctx, chat.CompetitionCode)
	if err != nil {
		return err
	}
	err = self.sendGirlProfile(ctx, chat, girl)
	return err
}

func (self *Bot) helpCmd(ctx context.Context, chat *models.Chat, message *messages.TextMessage) error {
	return self.sendHelpText(ctx, chat)
}

func (self *Bot) chatSettingsCmd(ctx context.Context, chat *models.Chat, message *messages.TextMessage) error {
	_, err := self.messenger.SendText(ctx, chat.Id, self.gettext(chat, "chat_settings_commands"))
	return err
}

func (self *Bot) leftVoteCmd(ctx context.Context, chat *models.Chat, message *messages.TextMessage) error {
	return self.processVoteMessage(ctx, chat, message, leftVote)
}

func (self *Bot) rightVoteCmd(ctx context.Context, chat *models.Chat, message *messages.TextMessage) error {
	return self.processVoteMessage(ctx, chat, message, rightVote)
}

func (self *Bot) globalCompetitionCmd(ctx context.Context, chat *models.Chat, message *messages.TextMessage) error {
	return self.changeCompetition(ctx, chat, competition.GlobalCompetition)
}

func (self *Bot) celebritiesCompetitionCmd(ctx context.Context, chat *models.Chat, message *messages.TextMessage) error {
	return self.changeCompetition(ctx, chat, competition.CelebritiesCompetition)
}

func (self *Bot) modelsCompetitionCmd(ctx context.Context, chat *models.Chat, message *messages.TextMessage) error {
	return self.changeCompetition(ctx, chat, competition.ModelsCompetition)
}

func (self *Bot) regularCompetitionCmd(ctx context.Context, chat *models.Chat, message *messages.TextMessage) error {
	return self.changeCompetition(ctx, chat, competition.RegularCompetition)
}

func (self *Bot) changeCompetition(ctx context.Context, chat *models.Chat, newCompetition string) error {
	var msgid string
	if newCompetition == competition.RegularCompetition {
		msgid = "regular_competition_enabled"
	} else if newCompetition == competition.ModelsCompetition {
		msgid = "models_competition_enabled"
	} else if newCompetition == competition.CelebritiesCompetition {
		msgid = "celebrities_competition_enabled"
	} else {
		msgid = "global_competition_enabled"
	}
	text := self.gettext(chat, msgid)
	if _, err := self.messenger.SendText(ctx, chat.Id, text, func(msg *tgbotapi.MessageConfig) {
		msg.ReplyMarkup = tgbotapi.ReplyKeyboardRemove{RemoveKeyboard: true}
	}); err != nil {
		return err
	}
	chat.CompetitionCode = newCompetition
	chat.ResetLastMatch()
	chat.ResetTopOffset()
	return nil

}

func (self *Bot) enableNotificationsCmd(ctx context.Context, chat *models.Chat, message *messages.TextMessage) error {
	if _, err := self.messenger.SendText(ctx, chat.Id, self.gettext(chat, "daily_notification_enabled")); err != nil {
		return err
	}
	if !chat.SelfActivationAllowed {
		chat.SelfActivationAllowed = true
		err := self.scheduleDailyNotification(ctx, chat)
		return err
	}
	return nil
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
	if timeout < models.DefaultTimeout || timeout > sessionDuration {
		lowerBound := models.DefaultTimeoutSeconds
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

func (self *Bot) handleFreeInput(ctx context.Context, chat *models.Chat, message *messages.TextMessage) error {
	link := message.GetLastWord()
	if link == "" {
		return self.sendDontGetYou(ctx, chat)
	}

	var err error

	girl, err := self.competition.GetCompetitor(ctx, chat.CompetitionCode, link)
	if _, ok := err.(*competition.CompetitorNotFound); !ok && err != competition.BadProfileLinkErr {
		return err
	}
	if err == nil {
		err = self.sendGirlProfile(ctx, chat, girl)
		return err
	}

	err = self.addGirl(ctx, chat, link)
	if err != competition.BadPhotoLinkErr && err != instagram.MediaForbidden {
		return err
	}

	err = self.sendDontGetYou(ctx, chat)
	return err
}

func (self *Bot) sendDontGetYou(ctx context.Context, chat *models.Chat) error {
	_, err := self.messenger.SendText(ctx, chat.Id, self.gettext(chat, "dont_get_you"))
	return err
}
func (self *Bot) sendHelpText(ctx context.Context, chat *models.Chat) error {
	_, err := self.messenger.SendText(ctx, chat.Id, self.gettext(chat, "help_message"))
	return err
}
