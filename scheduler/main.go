package main

import (
	. "instarate/scheduler/config"

	"github.com/gazoon/go-utils"
	"github.com/gazoon/go-utils/consumer"
	"github.com/gazoon/go-utils/logging"
	"instarate/scheduler/sender"
	"instarate/scheduler/services"
)

var logger = logging.WithPackage("main")

func main() {

	logging.PatchStdLog(Config.LogLevel, Config.ServiceName)
	err := utils.InitializeSentry(Config.Sentry.DSN)
	if err != nil {
		panic(err)
	}
	taskSender, err := sender.New(Config.MongoQueue)
	if err != nil {
		panic(err)
	}
	taskStorage := services.InitTaskReader()
	tasksPipe := sender.NewTasksPipeline(taskStorage.GetAndRemoveTask, taskSender.SendTask)
	worker := consumer.New(tasksPipe.Fetch, Config.TasksConsumer.FetchDelay)
	worker.Run()
	logger.Info("Scheduler successfully started")
	utils.WaitingForShutdown()
	worker.Stop()
}
