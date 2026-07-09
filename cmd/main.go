package main

import (
	"strconv"
	task_package "task-planner/internal/task"
	task_loop "task-planner/internal/task-loop"
	"task-planner/internal/task-manager/service"
	"time"
)

func main() {
	//DB := storage.NewTaskStorage()
	//TaskManager := task_manager.NewTaskManager(&DB)

	TaskLoop := task_loop.NewTaskLoop()
	TaskLoop.Start()

	for i := 1; i <= 10; i++ {
		TaskLoop.AddTask(*task_package.NewTask("INDEX "+strconv.Itoa(i), 100, i), func(task *task_package.Task) {
			service.EmulateTaskProgress(task)
		})
	}

	//go initCLI(TaskManager)
	go func() {
		time.Sleep(time.Second * 40)
		TaskLoop.Stop()

	}()
	TaskLoop.Wait()
	//time.Sleep(time.Minute * 1)

}
