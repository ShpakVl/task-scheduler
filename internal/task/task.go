package task

import "errors"

type Task struct {
	ID          int
	Description string
	status      Status
	progress    uint
}

func (t *Task) GetStatus() Status {
	return t.status
}

func (t *Task) SetStatus(newStatus Status) error {
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

func NewTask(description string, id int) *Task {
	return &Task{
		Description: description,
		status:      STATUS_NOT_STARTED,
		ID:          id,
	}
}
