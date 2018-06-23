package core

import (
	"context"
	"instarate/libs/competition"
	"instarate/libs/instagram"
	"instarate/tg_bot/messages"
	"instarate/tg_bot/models"
	"strings"
)

var (
	defaultCommandsLanguage = "en"
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
	return nil
}

func (self *Bot) nextGirlCmd(ctx context.Context, chat *models.Chat, message *messages.TextMessage) error {
	return nil
}

func (self *Bot) girlProfileCmd(ctx context.Context, chat *models.Chat, message *messages.TextMessage) error {
	return nil
}

func (self *Bot) helpCmd(ctx context.Context, chat *models.Chat, message *messages.TextMessage) error {
	return nil
}

func (self *Bot) chatSettingsCmd(ctx context.Context, chat *models.Chat, message *messages.TextMessage) error {
	return nil
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
	return nil
}

func (self *Bot) celebritiesCompetitionCmd(ctx context.Context, chat *models.Chat, message *messages.TextMessage) error {
	return nil
}

func (self *Bot) modelsCompetitionCmd(ctx context.Context, chat *models.Chat, message *messages.TextMessage) error {
	return nil
}

func (self *Bot) regularCompetitionCmd(ctx context.Context, chat *models.Chat, message *messages.TextMessage) error {
	return nil
}

func (self *Bot) enableNotificationsCmd(ctx context.Context, chat *models.Chat, message *messages.TextMessage) error {
	return nil
}

func (self *Bot) disableNotificationsCmd(ctx context.Context, chat *models.Chat, message *messages.TextMessage) error {
	return nil
}

func (self *Bot) setRussianCmd(ctx context.Context, chat *models.Chat, message *messages.TextMessage) error {
	return nil
}

func (self *Bot) setEnglishCmd(ctx context.Context, chat *models.Chat, message *messages.TextMessage) error {
	return nil
}

func (self *Bot) setVotingTimeoutCmd(ctx context.Context, chat *models.Chat, message *messages.TextMessage) error {
	return nil
}

func (self *Bot) handleRegularText(ctx context.Context, chat *models.Chat, message *messages.TextMessage) error {
	return nil
}
