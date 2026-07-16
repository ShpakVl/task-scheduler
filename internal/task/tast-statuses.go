package task

type Status string

const STATUS_NOT_STARTED Status = "not_started"
const STATUS_IN_PROGRESS Status = "in_progress"
const STATUS_DONE Status = "done"
const STATUS_FAILED Status = "failed"
const STATUS_CANCELLED Status = "cancelled"

var SUPPORTED_STATUSES = map[Status]bool{
	STATUS_NOT_STARTED: true,
	STATUS_IN_PROGRESS: true,
	STATUS_DONE:        true,
	STATUS_FAILED:      true,
	STATUS_CANCELLED:   true,
}
