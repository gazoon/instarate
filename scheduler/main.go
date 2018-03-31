package main

import (
	"path"

	"instarate/scheduler/config"

	"instarate/scheduler/sender"
	"instarate/scheduler/tasks"

	"fmt"
	"github.com/gazoon/go-utils"
	"github.com/gazoon/go-utils/consumer"
	"github.com/gazoon/go-utils/logging"
)

func main() {
	conf := &config.Config{}
	configPath := path.Join(utils.GetCurrentFileDir(), "config")
	err := utils.ParseConfig(configPath, conf)
	if err != nil {
		panic(err)
	}
	fmt.Println(conf)
	logging.PatchStdLog(conf.LogLevel, conf.ServiceName)
	taskSender, err := sender.New(conf.MongoQueue)
	if err != nil {
		panic(err)
	}
	taskStorage, err := tasks.NewStorage(conf.MongoTasks)
	if err != nil {
		panic(err)
	}
	worker := consumer.New(taskStorage.GetTask, taskSender.SendTask, "", conf.TasksConsumer.FetchDelay)
	worker.Run()
	utils.WaitingForShutdown()
	worker.Stop()
}
