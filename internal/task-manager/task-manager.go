package task_manager

import (
	"context"
	"task-planner/internal/task"
	"task-planner/internal/task-manager/ports"
	"task-planner/internal/task-manager/service"

	"github.com/k0kubun/pp"
)

type TaskManager struct {
	db        ports.TaskRepository
	processor ports.TaskProcessor
}

func (t *TaskManager) StartTaskManager() {
	t.processor.Start()
}

func (t *TaskManager) StopTaskManager() {
	t.processor.Stop()
}

func (t *TaskManager) WaitTaskManager() {
	t.processor.Wait()
}

func (t *TaskManager) AddTask(task ports.CreateTaskInput) error {
	createdTask, err := t.db.Create(task)
	if err != nil {
		return err
	}

	t.processor.AddTask(createdTask, t.taskProcessingCallback, t.taskStopCallback)

	return nil
}
func (t *TaskManager) taskStopCallback(stoppedTaskId int) error {
	lastTaskState, err := t.GetTaskById(stoppedTaskId)
	if err != nil {
		return err
	}
	err = lastTaskState.SetStatus(task.STATUS_CANCELLED)

	if err != nil {
		return err
	}

	t.taskProcessMonitoring(lastTaskState)
	pp.Println("Task is cancelled ", lastTaskState)
	_, err = t.db.Update(lastTaskState)

	return err
}

func (t *TaskManager) taskProcessingCallback(ctx context.Context, taskToProcess *task.Task) error {
	progressCh := make(chan service.EmulateEvent)

	taskToProcess.SetStatus(task.STATUS_IN_PROGRESS)

	go service.EmulateTaskProgress(ctx, progressCh, taskToProcess)

	for progressEvent := range progressCh {
		if progressEvent.Err != nil {
			return progressEvent.Err
		}
		t.taskProcessMonitoring(*progressEvent.Task)
		t.db.Update(*progressEvent.Task)

		if progressEvent.Task.GetStatus() == task.STATUS_DONE {
			pp.Println("Task is done ", progressEvent.Task.ID)
		}

	}

	return nil
}

func (t *TaskManager) taskProcessMonitoring(processedTask task.Task) {
	pp.Println("[ID]", processedTask.ID, "staus-> ", processedTask.GetStatus(), " [%]", processedTask.GetProgress())
}

func (t *TaskManager) GetTaskById(id int) (task.Task, error) {
	return t.db.GetById(id)
}
func (t *TaskManager) GetAllTasks() []task.Task {
	return t.db.GetAll()
}

func (t *TaskManager) CancelTask(taskId int) {
	t.processor.CancelTask(taskId)
}

func NewTaskManager(db ports.TaskRepository, processor ports.TaskProcessor) *TaskManager {
	return &TaskManager{
		db:        db,
		processor: processor,
	}
}
