package main

import (
	"runtime"
	"strconv"
	"task-planner/internal/storage"
	task_package "task-planner/internal/task"
	task_manager "task-planner/internal/task-manager"
)

func main() {
	runtime.GOMAXPROCS(2)

	var TasksDB storage.TaskStorage = storage.NewTaskStorage()
	var Runner *task_manager.Runner = task_manager.NewRunner()
	var Queue *task_manager.Queue = task_manager.NewQueue(Runner.TaskProcessing)
	var TaskManager *task_manager.TaskManager = task_manager.NewTaskManager(&TasksDB, Queue, Runner)

	TaskManager.StartTaskManager()

	for i := 0; i < 40; i++ {
		TaskManager.CreateTask(task_package.NewTask("INDEX "+strconv.Itoa(i), 100, i))
	}

	go initCLI(TaskManager)

	TaskManager.WaitTaskManager()
}
