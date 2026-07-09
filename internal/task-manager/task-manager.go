package task_manager

import (
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

	t.processor.AddTask(createdTask, t.taskProcessingCallback)

	return nil
}

func (t *TaskManager) taskProcessingCallback(taskToProcess *task.Task) {
	/*
		ch -->
				emulate(ch, task)
						---> ch <- task
		val := <-ch
				val ---> monitoring
				val ---> saveToDB

	*/
	progressCh := make(chan task.Task)

	go service.EmulateTaskProgress(progressCh, taskToProcess)

	for processedTask := range progressCh {
		pp.Println("[ID]", processedTask.ID, " [%]", processedTask.GetProgress())

		if processedTask.GetStatus() == task.STATUS_DONE {
			pp.Println("Task is done")

			t.db.Update(processedTask)

			return
		}
	}
}

func (t *TaskManager) GetAllTasks() []task.Task {
	return t.db.GetAll()
}

func NewTaskManager(db ports.TaskRepository, processor ports.TaskProcessor) *TaskManager {
	return &TaskManager{
		db:        db,
		processor: processor,
	}
}
