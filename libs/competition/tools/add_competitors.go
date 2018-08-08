package main

import (
	"context"
	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"instarate/libs/competition"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		panic("You must provide path to the file with links.")
	}
	filePath := os.Args[1]
	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(errors.Wrap(err, "can't read links file"))
	}
	links := strings.Split(string(fileContent), "\n")
	api := competition.InitCompetition()
	ctx := context.Background()
	outFile, err := os.Create("error_links.txt")
	if err != nil {
		panic(err)
	}
	defer outFile.Close()
	for _, link := range links {
		_, err := api.Add(ctx, link)
		if err == nil {
			continue
		}
		if err == competition.ProfileExistsErr {
			log.WithField("link", link).Warn("Competitor already exists.")
			continue
		}
		log.WithFields(log.Fields{"link": link, "reason": err}).
			Error("Competitor wasn't added. Skip.")
		if _, err := outFile.WriteString(link + "\n"); err != nil {
			panic(err)
		}

	}
	log.Info("Done!")
}
