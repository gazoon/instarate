package main

import (
	"fmt"
	"instarate/scheduler/tasks"
)

func main() {
	d := map[string]interface{}{
		"chat_id": 1,
		"name":    "true",
		"args":    map[string]interface{}{"foo": "bar"},
		"do_at":   9898,
	}
	c, err := tasks.TaskFromData(d)
	if err != nil {
		panic(err)
	}
	fmt.Println(c)
}
