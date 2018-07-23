package main

import (
	log "github.com/Sirupsen/logrus"
	"instarate/scheduler/services"
)

func main() {
	publisher := services.InitTaskPublisher()
	log.Info("Create tasks publisher indexes")
	if err := publisher.CreateIndexes(); err != nil {
		panic(err)
	}

	reader := services.InitTaskReader()
	log.Info("Create tasks reader indexes")
	if err := reader.CreateIndexes(); err != nil {
		panic(err)
	}

	log.Info("Done!")
}
