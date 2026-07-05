package task

import "errors"

type Task struct {
	ID          int
	description string
	status
	progress    uint
	overallSize uint
}

func (t *Task) ChangeStatus(newStatus status) (bool, error) {
	if t.validateStatus(newStatus) {
		t.status = newStatus

		return true, nil
	}

	return false, errors.New("unsupported status")
}

func (t *Task) ChangeProgress(newProgress uint) (uint, error) {
	if t.validateProgress(newProgress) {
		t.progress = newProgress

		return t.progress, nil
	}
	return 0, errors.New("progress must be between 0 and 100")
}

func (t *Task) GetDescription() string {
	return t.description
}

func (t Task) GetProgressPercent() uint {
	return t.progress / t.overallSize
}

func NewTask(description string, overallSize uint) *Task {
	return &Task{
		description: description,
		overallSize: overallSize,
		status:      STATUS_NOT_STARTED,
	}
}
