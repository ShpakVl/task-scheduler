package repository

import "task-planner/internal/task"

type TaskRepository interface {
	GetById(id int) (task.Task, error)
	GetAll() []task.Task
	Create(task task.Task) (task.Task, error)
	Update(task task.Task) (task.Task, error)
}
