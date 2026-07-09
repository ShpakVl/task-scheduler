package task_loop

import (
	"errors"
	"sync"
	task_package "task-planner/internal/task"
)

type event string

type Job struct {
	Task         task_package.Task
	Event        event
	processingCb func(task *task_package.Task)
}

type JobsLoop struct {
	wg                sync.WaitGroup
	once              *sync.Once
	isShutdownPlanned bool
	tasksSchedulerCh  chan Job
	loopSchedulerCh   chan event
}

func (j *JobsLoop) Start() {
	j.wg.Add(1) //To not finish main thread

	// Blocking task due to RANGE inside
	go j.startProcessingLoop()
}

func (j *JobsLoop) startProcessingLoop() {
	defer j.wg.Done()

	tasksInProgress := 0
	postponedJobs := make([]Job, 0, MAX_TASKS_IN_PROGRESS)

	for j.loopSchedulerCh != nil || j.tasksSchedulerCh != nil {
		select {
		case job, ok := <-j.tasksSchedulerCh:
			if !ok {
				j.tasksSchedulerCh = nil
				continue
			}
			j.handleJobEvent(job, &postponedJobs, &tasksInProgress)

		case loopEvent, ok := <-j.loopSchedulerCh:
			if !ok {
				j.loopSchedulerCh = nil
				continue
			}
			j.handleLoopSchedulerEvent(loopEvent, tasksInProgress)
		}

	}

}

func (j *JobsLoop) handleJobEvent(job Job, postponedJobs *[]Job, tasksInProgress *int) {
	switch job.Event {
	case EVENT_ADD:
		if *tasksInProgress < MAX_TASKS_IN_PROGRESS && !j.isShutdownPlanned {
			*tasksInProgress++
			go j.startJobProcessing(&job.Task, job.processingCb)

			return
		}
		j.addToPostponed(postponedJobs, job)

	case EVENT_DONE:
		*tasksInProgress--
		if *tasksInProgress == 0 && j.isShutdownPlanned {
			close(j.tasksSchedulerCh)

			return
		}
		if nextJob, err := j.getTaskPostponed(postponedJobs); err == nil && !j.isShutdownPlanned {
			j.tasksSchedulerCh <- nextJob
		}
	}
}

func (j *JobsLoop) handleLoopSchedulerEvent(loopEvent event, tasksInProgress int) {
	switch loopEvent {
	case EVENT_STOP:
		j.isShutdownPlanned = true
		close(j.loopSchedulerCh)
		if tasksInProgress == 0 && j.isShutdownPlanned {
			close(j.tasksSchedulerCh)
		}
	}
}

func (j *JobsLoop) AddTask(task task_package.Task, processingCb func(task *task_package.Task)) {
	j.tasksSchedulerCh <- Job{Task: task, Event: EVENT_ADD, processingCb: processingCb}
}

func (j *JobsLoop) addToPostponed(postponedJobs *[]Job, job Job) {
	*postponedJobs = append(*postponedJobs, job)
}

func (j *JobsLoop) getTaskPostponed(postponedJob *[]Job) (Job, error) {
	if len(*postponedJob) > 0 {
		nexJob := (*postponedJob)[0]

		*postponedJob = (*postponedJob)[1:]

		return nexJob, nil
	}

	return Job{}, errors.New("no more tasks to process")
}

func (j *JobsLoop) startJobProcessing(task *task_package.Task, cb func(task *task_package.Task)) {
	// Process task as blocking operation and publish event after processing
	cb(task)

	j.tasksSchedulerCh <- Job{Task: *task, Event: EVENT_DONE}
}

func (j *JobsLoop) Wait() {
	j.wg.Wait()
}

func (j *JobsLoop) Stop() {
	j.once.Do(func() {
		j.loopSchedulerCh <- EVENT_STOP
	})

}

func NewTaskLoop() *JobsLoop {
	return &JobsLoop{
		wg:                sync.WaitGroup{},
		tasksSchedulerCh:  make(chan Job, 5),
		loopSchedulerCh:   make(chan event),
		once:              &sync.Once{},
		isShutdownPlanned: false,
	}
}

const MAX_TASKS_IN_PROGRESS = 3

const EVENT_DONE event = "done"
const EVENT_ADD event = "add"
const EVENT_STOP event = "stop"
