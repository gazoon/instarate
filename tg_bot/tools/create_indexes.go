package main

import (
	log "github.com/Sirupsen/logrus"
	"instarate/tg_bot/services"
)

func main() {
	chats := services.InitChatsStorage()
	log.Info("Create chat storage indexes")
	if err := chats.CreateIndexes(); err != nil {
		panic(err)
	}

	cache := services.InitCache()
	log.Info("Create cache indexes")
	if err := cache.CreateIndexes(); err != nil {
		panic(err)
	}

	queueReader := services.InitQueueReader()
	log.Info("Create queue reader indexes")
	if err := queueReader.CreateIndexes(); err != nil {
		panic(err)
	}

	queueWriter := services.InitQueueWriter()
	log.Info("Create queue writer indexes")
	if err := queueWriter.CreateIndexes(); err != nil {
		panic(err)
	}

	log.Info("Done!")
}
