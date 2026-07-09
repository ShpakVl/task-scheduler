package service

import (
	"task-planner/internal/task"
	"time"
)

func EmulateTaskProgress(progressCh chan<- task.Task, taskToProcess *task.Task) {
	for {
		if taskToProcess.GetProgress() == 100 {
			defer close(progressCh)

			taskToProcess.SetStatus(task.STATUS_DONE)
			progressCh <- *taskToProcess
			return

		}
		time.Sleep(time.Millisecond * 30)

		taskToProcess.SetProgress(taskToProcess.GetProgress() + 50)
		progressCh <- *taskToProcess
	}
}
