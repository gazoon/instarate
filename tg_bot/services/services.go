package services

import (
	"github.com/gazoon/go-utils/localization"
	"github.com/gazoon/go-utils/queue"
	"instarate/libs/competition"
	"instarate/scheduler/scheduler"
	"instarate/tg_bot/cache"
	"instarate/tg_bot/chats"
	. "instarate/tg_bot/config"
	"instarate/tg_bot/core"
	"instarate/tg_bot/messenger"
	"path"
)

func InitLocalization() *localization.Manager {
	localesDir := path.Join(RootDir, "locales")
	locales, err := localization.NewManager(localesDir)
	if err != nil {
		panic(err)
	}
	return locales
}

func InitTelegramMessenger() *messenger.Telegram {
	tg, err := messenger.NewTelegram(Config.Telegram.Token)
	if err != nil {
		panic(err)
	}
	return tg
}

func InitChatsStorage() *chats.MongoStorage {
	chatsStorage, err := chats.NewMongoStorage(Config.MongoChats)
	if err != nil {
		panic(err)
	}
	return chatsStorage
}

func InitCache() *cache.Mongo {
	cacheStorage, err := cache.NewMongo(Config.MongoCache)
	if err != nil {
		panic(err)
	}
	return cacheStorage
}

func InitBot() *core.Bot {
	competitionAPI := competition.InitCompetition()
	chatsStorage := InitChatsStorage()
	cacheStorage := InitCache()
	tg := InitTelegramMessenger()
	schedulerAPI := scheduler.InitScheduler()
	locales := InitLocalization()
	bot := core.NewBot(competitionAPI, chatsStorage, cacheStorage, tg, schedulerAPI, locales, Config.Bot)
	return bot
}

func InitQueueReader() *queue.MongoReader {
	queueReader, err := queue.NewMongoReader(Config.MongoQueue)
	if err != nil {
		panic(err)
	}
	return queueReader
}

func InitQueueWriter() *queue.MongoWriter {
	queueWriter, err := queue.NewMongoWriter(Config.MongoQueue)
	if err != nil {
		panic(err)
	}
	return queueWriter
}
