package main

import (
	"path"

	"fmt"
	"github.com/gazoon/go-utils"
	"instarate/tg_gateway/config"
	"instarate/tg_gateway/webhook"
	"instarate/tg_gateway/worker"
)

func main() {
	conf := &config.Config{}
	configPath := path.Join(utils.GetCurrentFileDir(), "config")
	err := utils.ParseConfig(configPath, conf)
	if err != nil {
		panic(err)
	}
	fmt.Println(conf)
	updatesWorker, err := worker.New(conf)
	if err != nil {
		panic(err)
	}
	webhookServer := webhook.New(conf.Port, conf.KnownBots, updatesWorker)
	webhookServer.Run()
	utils.WaitingForShutdown()
	webhookServer.Stop()
}
