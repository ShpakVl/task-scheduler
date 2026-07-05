package task

import "errors"

type Task struct {
	ID          int
	Description string
	status      status
	progress    uint
	OverallSize uint
}

func (t *Task) GetStatus() status {
	return t.status
}

func (t *Task) SetStatus(newStatus status) error {
	if t.validateStatus(newStatus) {
		t.status = newStatus
		return nil
	}

	return errors.New("invalid status")
}

func (t *Task) GetProgress() uint {
	return t.progress
}
func (t *Task) SetProgress(newProgress uint) error {
	if t.progress > 100 || t.progress < 0 {
		return errors.New("invalid progress")
	}
	t.progress = newProgress

	return nil
}

func NewTask(description string, overallSize uint, id int) *Task {
	return &Task{
		Description: description,
		OverallSize: overallSize,
		status:      STATUS_NOT_STARTED,
		ID:          id,
	}
}
