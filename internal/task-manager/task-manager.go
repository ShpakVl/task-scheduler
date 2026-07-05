package task_manager

import (
	"errors"
	"sync"
	task_package "task-planner/internal/task"
	"task-planner/internal/task-manager/repository"

	"github.com/k0kubun/pp"
)

type TaskManager struct {
	tasksDB repository.TaskRepository
	wg      *sync.WaitGroup
	Queue
	Runner
}

func (t TaskManager) GetAllTasks() []task_package.Task {
	return t.tasksDB.GetAll()
}

func (t *TaskManager) CreateTask(taskToCreate *task_package.Task) (task_package.Task, error) {
	createdTask, err := t.tasksDB.Create(*taskToCreate)
	t.Queue.AddTask(createdTask)

	if err != nil {
		return task_package.Task{}, errors.New("failed to create task")
	}

	return createdTask, nil
}

func (t *TaskManager) StartTaskManager() {
	t.wg.Add(1)
	go t.Queue.StartQueueProcessingLoop()
}

func (t *TaskManager) WaitTaskManager() {
	t.wg.Wait()
}

func (t *TaskManager) StopTaskManager() {
	t.Queue.StopQueueProcessingLoop()
	pp.Println("AllTasks==> ", t.GetAllTasks())

	t.wg.Done()
}

func NewTaskManager(tasksDB repository.TaskRepository, queue *Queue, runner *Runner) *TaskManager {
	return &TaskManager{
		wg:      &sync.WaitGroup{},
		tasksDB: tasksDB,
		Queue:   *queue,
		Runner:  *runner,
	}
}
