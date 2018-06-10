package sender

import (
	"context"
	"instarate/scheduler/tasks"

	"github.com/gazoon/go-utils/consumer"
	"github.com/gazoon/go-utils/logging"
)

type TasksPipeline struct {
	*logging.LoggerMixin
	getTask  func(context.Context) (*tasks.Task, error)
	sendTask func(context.Context, *tasks.Task)
}

func NewTasksPipeline(getTask func(context.Context) (*tasks.Task, error),
	sendTask func(context.Context, *tasks.Task)) *TasksPipeline {

	return &TasksPipeline{
		getTask:     getTask,
		sendTask:    sendTask,
		LoggerMixin: logging.NewLoggerMixin("tasks_pipe", nil),
	}
}

func (self *TasksPipeline) Fetch(ctx context.Context) consumer.Process {
	defer func() {
		if r := recover(); r != nil {
			self.LogError(ctx, r)
		}
	}()
	task, err := self.getTask(ctx)
	if err != nil {
		self.LogError(ctx, err)
		return nil
	}
	if task == nil {
		return nil
	}
	return func() {
		self.sendTask(ctx, task)
	}
}
