package service

import (
	"context"
	"task-planner/internal/task"
	"time"
)

type EmulateEvent struct {
	Err  error
	Task *task.Task
}

func EmulateTaskProgress(ctx context.Context, progressCh chan<- EmulateEvent, taskToProcess *task.Task) {
	defer close(progressCh)

	for {
		select {
		case <-ctx.Done():
			progressCh <- EmulateEvent{Err: ctx.Err(), Task: taskToProcess}

			return

		default:
			time.Sleep(time.Second * 4)
			taskToProcess.SetProgress(taskToProcess.GetProgress() + 20)

			if taskToProcess.GetProgress() == 100 {
				taskToProcess.SetStatus(task.STATUS_DONE)
				progressCh <- EmulateEvent{Err: ctx.Err(), Task: taskToProcess}

				return
			}

			progressCh <- EmulateEvent{Err: ctx.Err(), Task: taskToProcess}
		}
	}
}
