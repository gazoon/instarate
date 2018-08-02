package main

import (
	"context"
	log "github.com/Sirupsen/logrus"
	"instarate/libs/competition"
	"strings"
)

func main() {
	profiles := competition.InitProfilesStorage()
	ctx := context.Background()
	profilesList, err := profiles.GetAll(ctx)
	if err != nil {
		panic(err)
	}
	for _, p := range profilesList {
		if !strings.HasPrefix(p.PhotoPath, competition.ProfilePhotosFolder) {
			log.Infof("Add prefix to %s", p.PhotoPath)
			p.PhotoPath = competition.ProfilePhotosFolder + p.PhotoPath
			if err := profiles.Save(ctx, p); err != nil {
				panic(err)
			}
		}
	}
	log.Info("Done!")
}
