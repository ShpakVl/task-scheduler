package main

import (
	"strconv"
	task_package "task-planner/internal/task"
	task_loop "task-planner/internal/task-loop"
	"task-planner/internal/task-manager/service"
	"time"

	"github.com/k0kubun/pp"
)

func main() {
	//DB := storage.NewTaskStorage()
	//TaskManager := task_manager.NewTaskManager(&DB)

	TaskLoop := task_loop.NewTaskLoop()
	TaskLoop.Start()
	t := time.Now()
	for i := 1; i <= 3; i++ {
		TaskLoop.AddTask(*task_package.NewTask("INDEX "+strconv.Itoa(i), 100, i), func(task *task_package.Task) {
			service.EmulateTaskProgress(task)
		})
	}
	go func() {
		time.Sleep(time.Millisecond * 2)
		TaskLoop.Stop()
	}()

	//go initCLI(TaskManager)

	TaskLoop.Wait()
	pp.Println(time.Since(t).Seconds())
}
