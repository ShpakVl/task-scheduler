package main

import (
	"fmt"
	"task-planner/internal/storage"
	"task-planner/internal/task"
	task_manager "task-planner/internal/task-manager"
	user_commands "task-planner/internal/user-commands"
)

func main() {
	var TasksDB storage.TaskStorage = storage.NewTaskStorage()
	var TaskManager *task_manager.TaskManager = task_manager.NewTaskManager(&TasksDB)

	Command := user_commands.NewUserCommands()

	Command.Register("add", func(description []string) error {
		_, err := TaskManager.CreateTask(task.NewTask(description[0], 100))
		fmt.Print("Task added", err)
		return err
	})
	Command.Init()

}
