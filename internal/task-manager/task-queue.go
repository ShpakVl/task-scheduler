package task_manager

import (
	"errors"
	"fmt"
	"slices"
	"sync"
	task_package "task-planner/internal/task"

	"github.com/k0kubun/pp"
)

type Queue struct {
	processingCb    func(ch chan task_package.Task, task task_package.Task)
	taskSchedulerCh chan task_package.Task
	mu              sync.RWMutex
	tasksInProcess  int
	postponed       []task_package.Task
}

func (q *Queue) StartTask(
	processingCb func(ch chan task_package.Task, task task_package.Task),
	taskToProcess task_package.Task,
) {
	taskCh := make(chan task_package.Task)
	go processingCb(taskCh, taskToProcess)

	for processingTask := range taskCh {
		if processingTask.GetStatus() == task_package.STATUS_DONE {
			q.OnTaskProcessed()
		}
	}
}

func (q *Queue) OnTaskProcessed() {
	q.mu.Lock()
	q.tasksInProcess--

	nextTask, ok := q.getTaskToProcess()

	if ok {
		q.taskSchedulerCh <- nextTask
		q.tasksInProcess++
	}
	q.mu.Unlock()
}

func (q *Queue) AddTask(taskToAdd task_package.Task) {
	q.mu.Lock()

	if isQueueCapacityAvailable(q.tasksInProcess) {
		q.tasksInProcess++
		q.taskSchedulerCh <- taskToAdd

		q.mu.Unlock()

		return
	}

	q.addTaskToPostponed(taskToAdd)

	q.mu.Unlock()
}

func (q *Queue) getTaskToProcess() (task_package.Task, bool) {
	if len(q.postponed) > 0 {
		pp.Println("BEF=>", q.postponed)

		taskToProcess := q.postponed[0]
		q.removeTaskFromPostponed(taskToProcess)
		pp.Println("AFTER-==> ", q.postponed)

		return taskToProcess, true
	}

	return task_package.Task{}, false
}

func (q *Queue) addTaskToPostponed(postponedTask task_package.Task) {
	addToQueue(&q.postponed, postponedTask)
}

func (q *Queue) removeTaskFromPostponed(postponedTask task_package.Task) {
	err := removeFromQueue(&q.postponed, postponedTask)

	if err != nil {
		fmt.Println(err)
	}
}

func (q *Queue) StartQueueProcessingLoop() {
	for taskToProcess := range q.taskSchedulerCh {
		go q.StartTask(q.processingCb, taskToProcess)
	}
}

func (q *Queue) StopQueueProcessingLoop() {
	close(q.taskSchedulerCh)
}

func NewQueue(taskProcessingCb func(ch chan task_package.Task, task task_package.Task)) *Queue {
	return &Queue{
		processingCb:    taskProcessingCb,
		taskSchedulerCh: make(chan task_package.Task, MAX_QUEUE_CAPACITY),
		postponed:       make([]task_package.Task, 0, 10),
		mu:              sync.RWMutex{},
	}
}

const MAX_QUEUE_CAPACITY = 3

func addToQueue(queueToModify *[]task_package.Task, taskToAdd task_package.Task) {
	*queueToModify = append(*queueToModify, taskToAdd)
}

func removeFromQueue(queueToModify *[]task_package.Task, taskToRemove task_package.Task) error {
	idx := getQueueTaskIdx(*queueToModify, taskToRemove.ID)
	if idx == -1 {
		return errors.New("postponed task with this ID was not found")
	}

	*queueToModify = append((*queueToModify)[:idx], (*queueToModify)[idx+1:]...)

	return nil
}

func getQueueTaskIdx(queue []task_package.Task, taskID int) int {
	return slices.IndexFunc(queue, func(task task_package.Task) bool {
		return task.ID == taskID
	})

}

func isQueueCapacityAvailable(currentlyInProgressCount int) bool {
	return currentlyInProgressCount < MAX_QUEUE_CAPACITY
}
