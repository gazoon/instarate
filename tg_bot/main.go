package main

import (
	"fmt"
	"github.com/gazoon/go-utils"
	"github.com/gazoon/go-utils/localization"
	"gopkg.in/telegram-bot-api.v4"
	"instarate/tg_bot/messenger"
	"os"
	"path"
)

func main() {
	localesDir := path.Join(utils.GetCurrentFileDir(), "locales")
	lm, err := localization.NewManager(localesDir)
	if err != nil {
		panic(err)
	}
	println(lm.GettextD("ru", "messages", "propose_to_vote"))
	m, err := messenger.NewTelegram("480997285:AAEwT3739sBnTz0RSqhEz8TNh4wvJUuqn20")
	if err != nil {
		panic(err)
	}
	f, err := os.Open("tmp.jpg")
	if err != nil {
		panic(err)
	}
	msgid, fileId, err := m.SendBinaryPhoto(nil, 231193206, f, func(opts *tgbotapi.PhotoConfig) {
		opts.Caption = "balalblalblblal"
		opts.ReplyMarkup = tgbotapi.ReplyKeyboardMarkup{
			OneTimeKeyboard: true,
			Keyboard:        [][]tgbotapi.KeyboardButton{{{Text: "teeee"}}},
		}
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(msgid, fileId)
}
