package task

func (t Task) validateStatus(newStatus Status) bool {

	return SUPPORTED_STATUSES[newStatus]
}
func (t Task) validateProgress(newProgress uint) bool {
	return newProgress > 100 || newProgress < 0
}
