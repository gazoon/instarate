package main

import (
	log "github.com/Sirupsen/logrus"
	"instarate/libs/competition"
)

func main() {
	competitors := competition.InitCompetitorsStorage()
	log.Info("Create competitors storage indexes")
	if err := competitors.CreateIndexes(); err != nil {
		panic(err)
	}

	profiles := competition.InitProfilesStorage()
	log.Info("Create profiles storage indexes")
	if err := profiles.CreateIndexes(); err != nil {
		panic(err)
	}

	voters := competition.InitVotersStorage()
	log.Info("Create voters storage indexes")
	if err := voters.CreateIndexes(); err != nil {
		panic(err)
	}

	log.Info("Done!")
}
