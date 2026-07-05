package storage

func (t *TaskStorage) isTaskAlreadyExists(id int) bool {
	_, ok := t.tasks[id]

	return ok
}
