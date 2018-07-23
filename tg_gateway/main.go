package main

import (
	. "instarate/tg_gateway/config"
	"instarate/tg_gateway/webhook"
	"instarate/tg_gateway/worker"

	"github.com/gazoon/go-utils"
	"github.com/gazoon/go-utils/logging"
)

func main() {
	logging.PatchStdLog(Config.LogLevel, Config.ServiceName)
	err := utils.InitializeSentry(Config.Sentry.DSN)
	if err != nil {
		panic(err)
	}
	updatesWorker, err := worker.New(Config)
	if err != nil {
		panic(err)
	}
	webhookServer := webhook.New(Config.Port, Config.BotToken, updatesWorker.ProcessUpdate)
	webhookServer.Run()
	utils.WaitingForShutdown()
	webhookServer.Stop()
}
