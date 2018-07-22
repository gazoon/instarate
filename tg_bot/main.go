package main

import (
	. "instarate/tg_bot/config"
	"instarate/tg_bot/core"

	"github.com/gazoon/go-utils"
	"github.com/gazoon/go-utils/consumer"
	"github.com/gazoon/go-utils/logging"
	"instarate/tg_bot/services"
)

var logger = logging.WithPackage("main")

func main() {
	logging.PatchStdLog(Config.LogLevel, Config.ServiceName)
	err := utils.InitializeSentry(Config.Sentry.DSN)
	if err != nil {
		panic(err)
	}
	bot := services.InitBot()
	incomingQueue := services.InitIncomingQueue()
	messagesPipe := core.NewMessagesPipe(incomingQueue, bot.OnMessage)
	worker := consumer.New(messagesPipe.Fetch, Config.QueueConsumer.FetchDelay)
	worker.Run()
	logger.Info("Bot successfully started")
	utils.WaitingForShutdown()
	worker.Stop()
}
