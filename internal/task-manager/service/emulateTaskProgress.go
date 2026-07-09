package service

import (
	"task-planner/internal/task"
	"time"

	"github.com/k0kubun/pp"
)

func EmulateTaskProgress(taskToProcess *task.Task) {
	for {
		if taskToProcess.GetProgress() == 100 {

			taskToProcess.SetStatus(task.STATUS_DONE)

			//ch <- taskToProcess

			//close(ch)

			return

		}
		time.Sleep(time.Millisecond * 30)

		taskToProcess.SetProgress(taskToProcess.GetProgress() + 50)

		pp.Println(taskToProcess)

		//ch <- taskToProcess
	}
}
