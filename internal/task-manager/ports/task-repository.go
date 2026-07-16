package ports

import (
	"context"
	task_package "task-planner/internal/task"
)

type TaskRepository interface {
	GetById(id int) (task_package.Task, error)
	GetAll() []task_package.Task
	Create(task CreateTaskInput) (task_package.Task, error)
	Update(task task_package.Task) (task_package.Task, error)
}

type ProcessingCb func(ctx context.Context, task *task_package.Task) error
type StopCb func(taskId int) error

type TaskProcessor interface {
	Start()
	Wait()
	Stop()

	CancelTask(taskId int)
	AddTask(
		task task_package.Task,
		processingCb ProcessingCb,
		stopCb StopCb,
	)
}
