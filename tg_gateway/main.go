package main

import (
	"path"

	"instarate/tg_gateway/config"
	"instarate/tg_gateway/webhook"
	"instarate/tg_gateway/worker"

	"github.com/gazoon/go-utils"
	"github.com/gazoon/go-utils/logging"
)

func main() {
	conf := &config.Config{}
	configPath := path.Join(utils.GetCurrentFileDir(), "config")
	err := utils.ParseConfig(configPath, conf)
	if err != nil {
		panic(err)
	}
	logging.PatchStdLog(conf.LogLevel, conf.ServiceName)
	updatesWorker, err := worker.New(conf)
	if err != nil {
		panic(err)
	}
	webhookServer := webhook.New(conf.Port, conf.KnownBots, updatesWorker)
	webhookServer.Run()
	utils.WaitingForShutdown()
	webhookServer.Stop()
}
