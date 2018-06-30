package main

import (
	"path"

	"instarate/scheduler/config"

	"instarate/scheduler/sender"
	"instarate/scheduler/tasks"

	"github.com/gazoon/go-utils"
	"github.com/gazoon/go-utils/consumer"
	"github.com/gazoon/go-utils/logging"
)

var logger = logging.WithPackage("main")

func main() {

	conf := &config.Config{}
	configPath := path.Join(utils.GetCurrentFileDir(), "config")
	err := utils.ParseConfig(configPath, conf)
	if err != nil {
		panic(err)
	}
	logging.PatchStdLog(conf.LogLevel, conf.ServiceName)
	taskSender, err := sender.New(conf.MongoQueue)
	if err != nil {
		panic(err)
	}
	taskStorage, err := tasks.NewStorage(conf.MongoTasks)
	if err != nil {
		panic(err)
	}
	tasksPipe := sender.NewTasksPipeline(taskStorage.GetAndRemoveTask, taskSender.SendTask)
	worker := consumer.New(tasksPipe.Fetch, conf.TasksConsumer.FetchDelay)
	worker.Run()
	logger.Info("Scheduler successfully started")
	utils.WaitingForShutdown()
	worker.Stop()
}
