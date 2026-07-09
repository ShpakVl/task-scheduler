package main

import (
	"strconv"
	"task-planner/internal/storage"
	task_loop "task-planner/internal/task-loop"
	task_manager "task-planner/internal/task-manager"
	"task-planner/internal/task-manager/ports"
	"time"

	"github.com/k0kubun/pp"
)

func main() {
	DB := storage.NewTaskStorage()
	Processor := task_loop.NewTaskLoop()
	TaskManager := task_manager.NewTaskManager(DB, Processor)
	TaskManager.StartTaskManager()

	t := time.Now()
	for i := 1; i <= 30; i++ {
		TaskManager.AddTask(ports.CreateTaskInput{
			ID:          i,
			Description: "ID=" + strconv.Itoa(i),
		})
	}

	go func() {
		time.Sleep(time.Second * 5)

		pp.Println("DB--> ", TaskManager.GetAllTasks())

		TaskManager.StopTaskManager()
	}()

	//!!!BLOCKING TASK!!!///
	TaskManager.WaitTaskManager()

	pp.Println(time.Since(t).Seconds())
}
