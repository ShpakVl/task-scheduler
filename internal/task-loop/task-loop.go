package task_loop

import (
	"errors"
	"fmt"
	"sync"
	task_package "task-planner/internal/task"

	"github.com/k0kubun/pp"
)

type event string

type Job struct {
	Task         task_package.Task
	Event        event
	processingCb func(task *task_package.Task)
}

type JobsLoop struct {
	wg                sync.WaitGroup
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
	tasksInprogress := 0
	postponedJobs := make([]Job, 0, MAX_TASKS_IN_PROGRESS)

	for job := range j.tasksSchedulerCh {
		fmt.Println("Job: ", job.Task.ID, " Event: ", job.Event)
		switch job.Event {
		case EVENT_ADD:
			if tasksInprogress < MAX_TASKS_IN_PROGRESS {
				tasksInprogress++
				pp.Println("ADDED==> ", job.Task.ID)

				go j.startJobProcessing(&job.Task, job.processingCb)

				continue
			}
			j.addToPostponed(&postponedJobs, job)

		case EVENT_DONE:
			tasksInprogress--

			if tasksInprogress == 0 && j.isShutdownPlanned {
				close(j.tasksSchedulerCh)
				j.loopSchedulerCh <- EVENT_STOP
			}

			if nextJob, err := j.getTaskPostponed(&postponedJobs); err == nil && !j.isShutdownPlanned {
				fmt.Print("Postponed task: ", nextJob.Task.ID, "\n")
				j.tasksSchedulerCh <- nextJob
			}

		case EVENT_STOP:
			j.isShutdownPlanned = true
		}

	}

}

func (l *JobsLoop) loopScheduler() {
	for event := range l.loopSchedulerCh {
		switch event {
		case EVENT_STOP:
			l.Stop()
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
	cb(task)
	j.tasksSchedulerCh <- Job{Task: *task, Event: EVENT_DONE}
}

func (j *JobsLoop) Wait() {
	j.wg.Wait()
}

func (j *JobsLoop) Stop() {
	j.wg.Done()
}

func NewTaskLoop() *JobsLoop {
	return &JobsLoop{
		wg:                sync.WaitGroup{},
		tasksSchedulerCh:  make(chan Job, 5),
		isShutdownPlanned: false,
	}
}

const MAX_TASKS_IN_PROGRESS = 3

const EVENT_DONE event = "done"
const EVENT_ADD event = "add"
const EVENT_STOP event = "stop"
