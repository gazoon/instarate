package main

import (
	"instarate/libs/competition"
)

func main() {
	c, err := competition.New()
	if err != nil {
		panic(err)
	}
	c.Test()
}
