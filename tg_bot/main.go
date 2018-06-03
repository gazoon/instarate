package main

import (
	"instarate/tg_bot/config"
	"instarate/tg_bot/messenger"
	"path"

	"context"
	"github.com/gazoon/bot_libs/queue"
	"github.com/gazoon/go-utils"
	"github.com/gazoon/go-utils/consumer"
	"github.com/gazoon/go-utils/localization"
	"instarate/tg_bot/chats"
	"instarate/tg_bot/core"
)

func main() {
	rootDir := utils.GetCurrentFileDir()
	conf := &config.Config{}
	configPath := path.Join(rootDir, "config")
	err := utils.ParseConfig(configPath, conf)
	if err != nil {
		panic(err)
	}
	localesDir := path.Join(rootDir, "locales")
	locales, err := localization.NewManager(localesDir)
	if err != nil {
		panic(err)
	}
	tg, err := messenger.NewTelegram(conf.Telegram.Token)
	chatsStorage, err := chats.NewMongoStorage(conf.MongoChats)
	if err != nil {
		panic(err)
	}
	bot := core.NewBot(chatsStorage, tg, locales)
	incomingQueue, err := queue.NewMongoReader(conf.MongoQueue)
	if err != nil {
		panic(err)
	}
	incomingQueue.GetNext()
	getTask := func(ctx context.Context) interface{} { return taskStorage.GetTask(ctx) }
	worker := consumer.New(getTask, bot.OnMessage, conf.QueueConsumer.FetchDelay)
	worker.Run()
	utils.WaitingForShutdown()
	worker.Stop()
}
