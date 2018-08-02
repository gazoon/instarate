package main

import (
	"fmt"
	"instarate/libs/instagram"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		panic("You must provide url to the resource.")
	}
	url := os.Args[1]
	jsonData, err := instagram.RequestResourceData(url)
	if err != nil {
		panic(err)
	}
	fmt.Println(jsonData)
}
