package ports

import "task-planner/internal/task"

type Storage interface {
	GetById(id int) (task.Task, error)
	GetAll() []task.Task
	Create(task task.Task) (task.Task, error)
	Update(task task.Task) (task.Task, error)
}
