package service

import (
	"task-planner/internal/task"
	"time"
)

func EmulateTaskProgress(ch chan<- task.Task, taskToProcess task.Task) {
	for {
		if taskToProcess.GetProgress() == 100 {

			taskToProcess.SetStatus(task.STATUS_DONE)

			ch <- taskToProcess

			close(ch)

			return

		}
		time.Sleep(time.Nanosecond * 1)
		taskToProcess.SetProgress(taskToProcess.GetProgress() + 20)

		ch <- taskToProcess
	}
}
