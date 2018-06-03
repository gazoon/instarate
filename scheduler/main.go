package main

import (
	"path"

	"instarate/scheduler/config"

	"instarate/scheduler/sender"
	"instarate/scheduler/tasks"

	"context"
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
	logging.PatchStdLog(conf.LogLevel, conf.ServiceName)
	taskSender, err := sender.New(conf.MongoQueue)
	if err != nil {
		panic(err)
	}
	taskStorage, err := tasks.NewStorage(conf.MongoTasks)
	if err != nil {
		panic(err)
	}

	// consumer expects a function that returns interface{}
	// but .GetTasks() returns *Task
	// so we need an adopter function
	getTask := func(ctx context.Context) interface{} { return taskStorage.GetTask(ctx) }
	worker := consumer.New(getTask, taskSender.SendTask, conf.TasksConsumer.FetchDelay)
	worker.Run()
	utils.WaitingForShutdown()
	worker.Stop()
}
