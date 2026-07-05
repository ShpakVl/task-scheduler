package storage

import (
	"errors"
	TaskModel "task-planner/internal/task"
)

type TaskStorage struct {
	tasks         map[int]*TaskModel.Task
	lastCreatedID int
}

func (t *TaskStorage) GetById(id int) (TaskModel.Task, error) {
	savedTask, ok := t.tasks[id]
	var err error

	if !ok {
		err = errors.New("task with this ID was not found")
	}

	return *savedTask, err
}

func (t *TaskStorage) Create(task TaskModel.Task) (TaskModel.Task, error) {
	currentTaskId := t.lastCreatedID + 1

	if !t.isTaskAlreadyExists(currentTaskId) {
		// Reassign ID of the task param because currently it is filled with default value (0)
		task.ID = currentTaskId

		t.tasks[currentTaskId] = &task
		t.increaseLastCreatedId()

		return *t.tasks[currentTaskId], nil
	}

	return TaskModel.Task{}, errors.New("task with this ID already exists")
}

func (t *TaskStorage) Update(task TaskModel.Task) (TaskModel.Task, error) {
	if t.isTaskAlreadyExists(task.ID) {
		t.tasks[task.ID] = &task

		return *t.tasks[task.ID], nil
	}
	return TaskModel.Task{}, errors.New("cannot update -> task with this ID was not found")

}

func (t *TaskStorage) increaseLastCreatedId() {
	t.lastCreatedID++
}

func NewTaskStorage() TaskStorage {
	return TaskStorage{
		tasks: make(map[int]*TaskModel.Task),
	}
}
