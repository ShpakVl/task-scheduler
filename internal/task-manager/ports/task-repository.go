package ports

import (
	task_package "task-planner/internal/task"
)

type TaskRepository interface {
	GetById(id int) (task_package.Task, error)
	GetAll() []task_package.Task
	Create(task CreateTaskInput) (task_package.Task, error)
	Update(task task_package.Task) (task_package.Task, error)
}

type TaskProcessor interface {
	Start()
	Wait()
	Stop()

	AddTask(task task_package.Task, processingCb func(task *task_package.Task))
}
