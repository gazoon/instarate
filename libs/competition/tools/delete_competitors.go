package main

import (
	"context"
	"instarate/libs/competition"
	"os"
)

func main() {
	c := competition.InitCompetition()
	if len(os.Args) < 2 {
		panic("You must provide list of usernames to remove")
	}
	err := c.Remove(context.Background(), os.Args[1:]...)
	if err != nil {
		panic(err)
	}
}
