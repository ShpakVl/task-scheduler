package task_manager

import (
	task_package "task-planner/internal/task"
	"task-planner/internal/task-manager/service"

	"github.com/k0kubun/pp"
)

type Runner struct{}

func (r *Runner) TaskProcessing(ch chan task_package.Task, task task_package.Task) {
	taskProgressCh := make(chan task_package.Task)

	go service.EmulateTaskProgress(taskProgressCh, task)
	for processedTask := range taskProgressCh {
		ch <- processedTask
		r.TaskProcessingMonitoring(processedTask)
	}

	defer close(ch)
}

func (r *Runner) TaskProcessingMonitoring(processedTask task_package.Task) {
	pp.Println("ID:", processedTask.ID, "Status:", processedTask.GetStatus(), "Progress:", processedTask.GetProgress())
}

func NewRunner() *Runner {
	return &Runner{}
}
