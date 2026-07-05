package task_manager

import (
	"errors"
	"task-planner/internal/task"
	"task-planner/internal/task-manager/repository"

	"github.com/k0kubun/pp"
)

type TaskManager struct {
	tasksDB      repository.TaskRepository
	currentTasks map[int]chan task.Task
}

func (t *TaskManager) CreateTask(task *task.Task) (chan task.Task, error) {
	//ctx := context.Background()

	createdTask, err := t.tasksDB.Create(*task)
	pp.Println(createdTask)
	if err != nil {
		return nil, errors.New("failed to create task")
	}
	//var chTaskProgress chan int = make(chan int)
	//
	//go service.EmulateTaskProgress(ctx, chTaskProgress, createdTask)
	//
	//select {
	//case <-chTaskProgress:
	//	{
	//		t.tasksDB.Update(createdTask)
	//	}
	//}

	return nil, nil
}

func NewTaskManager(tasksDB repository.TaskRepository) *TaskManager {
	return &TaskManager{
		currentTasks: make(map[int]chan task.Task),
		tasksDB:      tasksDB,
	}
}
