package main

import (
	"context"
	log "github.com/Sirupsen/logrus"
	"instarate/libs/competition"
)

func main() {
	ctx := context.Background()
	profilesStorage := competition.InitProfilesStorage()
	competitorsStorage := competition.InitCompetitorsStorage()
	log.Info("Get all profiles from the db...")
	profiles, err := profilesStorage.GetAll(ctx)
	if err != nil {
		panic(err)
	}
	log.Info("Start allocating competitors")
	for _, profile := range profiles {
		competitionCodes := competition.ChoseCompetitions(profile.Followers)
		for _, competitionCode := range competitionCodes {
			competitor := competition.NewCompetitor(profile.Username, competitionCode)
			err = competitorsStorage.Create(ctx, competitor)
			if err == nil {
				log.WithFields(log.Fields{
					"username":    profile.Username,
					"competition": competitionCode,
					"followers":   profile.Followers}).
					Info("Competitor allocated to the competition")
				continue
			}
			if _, ok := err.(competition.CompetitorAlreadyExists); ok {
				continue
			}
			panic(err)
		}
	}
	log.Info("Done!")
}
