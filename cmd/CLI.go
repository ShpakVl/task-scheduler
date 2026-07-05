package main

import (
	"fmt"
	"task-planner/internal/task"
	task_manager "task-planner/internal/task-manager"
	user_commands "task-planner/internal/user-commands"
)

func initCLI(taskManager *task_manager.TaskManager) {

	Command := user_commands.NewUserCommands()

	Command.Register("add", func(description []string) error {
		taskManager.AddTask(*task.NewTask("11", 100, 2))
		fmt.Print("Task added")
		return nil
	})
	//Command.Register("getAll", func(description []string) error {
	//	pp.Println(taskManager.GetAllTasks())
	//
	//	return nil
	//})

	Command.Register("stop", func(description []string) error {
		taskManager.StopTaskManager()

		return nil
	})

	Command.Init()

}
