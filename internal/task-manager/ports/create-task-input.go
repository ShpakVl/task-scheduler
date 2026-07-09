package ports

import "task-planner/internal/task"

type CreateTaskInput struct {
	ID          int
	Description string
	status      task.Status
	progress    uint
}
