package task_loop

import (
	"context"
	"errors"
	"fmt"
	"sync"
	task_package "task-planner/internal/task"
	"task-planner/internal/task-manager/ports"

	"github.com/k0kubun/pp"
)

type event string

type Job struct {
	Task         task_package.Task
	Event        event
	processingCb ports.ProcessingCb
	stopCb       ports.StopCb
}

type JobsLoop struct {
	wg                      sync.WaitGroup
	ctx                     context.Context
	mu                      sync.Mutex
	once                    *sync.Once
	isShutdownPlanned       bool
	tasksSchedulerCh        chan Job
	loopSchedulerCh         chan event
	activeTasksCloseActions map[int]func(taskId int) error
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
			pp.Println(*tasksInProgress, "Tasks in progress")
			ctx, cancel := context.WithCancel(j.ctx)

			go j.startJobProcessing(ctx, &job.Task, job.processingCb)

			j.activeTasksCloseActions[job.Task.ID] = func(id int) error {
				cancel()
				err := job.stopCb(id)

				if err != nil {
					pp.Println("err:", err)
				}

				return err
			}
			return
		}
		j.addToPostponed(postponedJobs, job)

	case EVENT_DONE:
		*tasksInProgress--

		if j.shouldCloseLoop(*tasksInProgress) {
			j.tasksSchedulerCh <- Job{Event: EVENT_STOP}
			return
		}
		delete(j.activeTasksCloseActions, job.Task.ID)

		j.tasksSchedulerCh <- Job{Event: START_NEXT_TASK}

	case START_NEXT_TASK:
		if nextJob, err := j.getTaskPostponed(postponedJobs); err == nil && !j.isShutdownPlanned {
			j.tasksSchedulerCh <- nextJob
		}

	case EVENT_CANCEL:
		{
			if cancelCb, ok := j.activeTasksCloseActions[job.Task.ID]; ok {
				*tasksInProgress--
				err := cancelCb(job.Task.ID)
				if err != nil {
					pp.Println(err)
				}
				delete(j.activeTasksCloseActions, job.Task.ID)
				j.tasksSchedulerCh <- Job{Event: START_NEXT_TASK}

			}

		}

	case EVENT_STOP:
		fmt.Print()
		close(j.tasksSchedulerCh)
	}

}
func (j *JobsLoop) shouldCloseLoop(tasksInProgress int) bool {
	return j.isShutdownPlanned && tasksInProgress == 0
}
func (j *JobsLoop) planNextTaskProcessing() {}

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

func (j *JobsLoop) CancelTask(taskId int) {
	j.tasksSchedulerCh <- Job{Task: task_package.Task{ID: taskId}, Event: EVENT_CANCEL}

}

func (j *JobsLoop) AddTask(task task_package.Task, processingCb ports.ProcessingCb, stopCb ports.StopCb) {
	j.tasksSchedulerCh <- Job{Task: task, Event: EVENT_ADD, processingCb: processingCb, stopCb: stopCb}
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

func (j *JobsLoop) startJobProcessing(ctx context.Context, task *task_package.Task, cb ports.ProcessingCb) {
	// Process task as blocking operation and publish event after processing
	err := cb(ctx, task)

	if err == nil {
		j.tasksSchedulerCh <- Job{Task: *task, Event: EVENT_DONE}
	}
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
		wg:                      sync.WaitGroup{},
		tasksSchedulerCh:        make(chan Job, 5),
		loopSchedulerCh:         make(chan event),
		once:                    &sync.Once{},
		ctx:                     context.Background(),
		activeTasksCloseActions: make(map[int]func(int) error),
		isShutdownPlanned:       false,
	}
}

const MAX_TASKS_IN_PROGRESS = 3

const EVENT_DONE event = "done"
const EVENT_ADD event = "add"
const EVENT_CANCEL event = "cancel"
const EVENT_STOP event = "stop"
const START_NEXT_TASK = "start_next_task"
