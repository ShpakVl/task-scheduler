package ports

import "task-planner/internal/task"

type CreateTaskInput struct {
	Description string
	status      task.Status
	progress    uint
}
