package main

import (
	"instarate/tg_bot/config"
	"instarate/tg_bot/messenger"
	"path"

	"instarate/tg_bot/chats"
	"instarate/tg_bot/core"

	"github.com/gazoon/go-utils"
	"github.com/gazoon/go-utils/consumer"
	"github.com/gazoon/go-utils/localization"
	"github.com/gazoon/go-utils/logging"
	"github.com/gazoon/go-utils/queue"
	"instarate/libs/competition"
	"instarate/scheduler/tasks"
	"instarate/tg_bot/cache"
)

var logger = logging.WithPackage("main")

func main() {
	rootDir := utils.GetCurrentFileDir()
	conf := &config.Config{}
	configPath := path.Join(rootDir, "config")
	err := utils.ParseConfig(configPath, conf)
	if err != nil {
		panic(err)
	}
	logging.PatchStdLog(conf.LogLevel, conf.ServiceName)
	localesDir := path.Join(rootDir, "locales")
	locales, err := localization.NewManager(localesDir)
	if err != nil {
		panic(err)
	}
	tg, err := messenger.NewTelegram(conf.Telegram.Token)
	if err != nil {
		panic(err)
	}
	chatsStorage, err := chats.NewMongoStorage(conf.MongoChats)
	if err != nil {
		panic(err)
	}
	cacheStorage, err := cache.NewMongo(conf.MongoCache)
	if err != nil {
		panic(err)
	}
	scheduler, err := tasks.NewPublisher()
	if err != nil {
		panic(err)
	}
	competitionAPI, err := competition.New()
	if err != nil {
		panic(err)
	}
	bot := core.NewBot(competitionAPI, chatsStorage, cacheStorage, tg, scheduler, locales, conf.Bot)
	incomingQueue, err := queue.NewMongoReader(conf.MongoQueue)
	if err != nil {
		panic(err)
	}
	messagesPipe := core.NewMessagesPipe(incomingQueue, bot.OnMessage)
	worker := consumer.New(messagesPipe.Fetch, conf.QueueConsumer.FetchDelay)
	worker.Run()
	logger.Info("Bot successfully started")
	utils.WaitingForShutdown()
	worker.Stop()
}
