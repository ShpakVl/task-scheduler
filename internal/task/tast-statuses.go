package task

type status string

const STATUS_NOT_STARTED status = "not_started"
const STATUS_IN_PROGRESS status = "in_progress"
const STATUS_DONE status = "done"
const STATUS_FAILED status = "failed"

var SUPPORTED_STATUSES = map[status]bool{
	STATUS_NOT_STARTED: true,
	STATUS_IN_PROGRESS: true,
	STATUS_DONE:        true,
	STATUS_FAILED:      true,
}
