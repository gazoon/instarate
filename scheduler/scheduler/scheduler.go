package scheduler

import (
	"instarate/scheduler/services"
	"instarate/scheduler/tasks"
)

func InitScheduler() *tasks.Publisher {
	return services.InitTaskPublisher()
}
