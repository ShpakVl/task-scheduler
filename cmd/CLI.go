package main

import (
	"strconv"
	task_manager "task-planner/internal/task-manager"
	"task-planner/internal/task-manager/ports"
	user_commands "task-planner/internal/user-commands"

	"github.com/k0kubun/pp"
)

func initCLI(taskManager *task_manager.TaskManager) {

	Command := user_commands.NewUserCommands()

	Command.Register(COMMAND_ADD, func(description []string) error {
		err := taskManager.AddTask(ports.CreateTaskInput{Description: description[0]})

		return err
	})
	Command.Register(COMMAND_GET_ALL, func(description []string) error {
		_, err := pp.Println(taskManager.GetAllTasks())

		return err
	})

	Command.Register(COMMAND_STOP, func(description []string) error {
		taskManager.StopTaskManager()

		return nil
	})

	Command.Register(COMMAND_CANCEL, func(description []string) error {
		id, _ := strconv.ParseInt(description[0], 0, 32)
		taskManager.CancelTask(int(id))

		return nil
	})

	Command.Init()
}

const COMMAND_ADD = "add"
const COMMAND_GET_ALL = "get_all"
const COMMAND_CANCEL = "cancel" //cancel task
const COMMAND_STOP = "stop"     //stop task task manager + event loop after completing current tasks
