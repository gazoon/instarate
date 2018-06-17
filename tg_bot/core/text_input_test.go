package core

import (
	"github.com/gazoon/go-utils"
	"github.com/gazoon/go-utils/localization"
	"instarate/tg_bot/messages"
	"instarate/tg_bot/models"
	"path"
	"testing"
)

func TestCommands(t *testing.T) {
	rootDir := utils.GetCurrentFileDir()
	localesDir := path.Join(rootDir, "../locales")
	locales, err := localization.NewManager(localesDir)
	if err != nil {
		panic(err)
	}
	bot := NewBot(nil, nil, nil, locales, &utils.BotInfo{Name: "bot", Username: "bot"})

	enChat := &models.Chat{Language: "en"}
	enTable := []struct {
		Text        string
		CommandName string
	}{
		{"left", leftVoteCmd},
		{"right", rightVoteCmd},
		{"start", startCmd},
		{"next", nextTopGirlCmd},
		{"show the best", showTopCmd},
		{"girl profile girl_username", girlProfileCmd},
		{"info about girl_username", girlProfileCmd},
		{"addGirl image_link", addGirlCmd},
		{"what can you do?", helpCmd},
		{"tell about yourself", helpCmd},
		{"tell about", noTextCommand.Name},
		{"settings", chatSettingsCmd},
		{"/enableNotifications@botname", enableNotificationsCmd},
		{"setVotingTimeout@botname 10", setVotingTimeoutCmd},
		{"here is no command", noTextCommand.Name},
	}
	for i, pair := range enTable {
		msg := &messages.TextMessage{Text: pair.Text}
		command := bot.detectTextCommand(enChat, msg)
		if command.Name != pair.CommandName {
			t.Errorf("English. Text command detected incorrectly, got: %s, expected: %s, pair: %d", command.Name, pair.CommandName, i+1)
		}
	}

	ruChat := &models.Chat{Language: "ru"}
	ruTable := []struct {
		Text        string
		CommandName string
	}{
		{"левая", leftVoteCmd},
		{"правая", rightVoteCmd},
		{"начать", startCmd},
		{"дальше", nextTopGirlCmd},
		{"кто в топе", showTopCmd},
		{"покажи лучших", showTopCmd},
		{"информация о girl_username", girlProfileCmd},
		{"покажи статистику о girl_username", girlProfileCmd},
		{"добавь девушку image_link", addGirlCmd},
		{"что ты можешь?", helpCmd},
		{"помощь", helpCmd},
		{"настройки", chatSettingsCmd},
		{"/enableNotifications@botname", enableNotificationsCmd},
		{"setVotingTimeout@botname 10", setVotingTimeoutCmd},
		{"текст без команд", noTextCommand.Name},
	}
	for i, pair := range ruTable {
		msg := &messages.TextMessage{Text: pair.Text}
		command := bot.detectTextCommand(ruChat, msg)
		if command.Name != pair.CommandName {
			t.Errorf("Russian chat. Text command detected incorrectly, got: %s, expected: %s, pair: %d", command.Name, pair.CommandName, i+1)
		}
	}
}
