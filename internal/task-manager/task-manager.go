package task_manager

//type TaskManager struct {
//	db    *storage.TaskStorage
//	wg    sync.WaitGroup
//	queue []task_package.Task
//}
//
//func (t *TaskManager) StartTaskManager() {
//	t.wg.Add(1)
//
//	go t.StartProcessingLoop()
//}
//
//func (t *TaskManager) StartProcessingLoop() {
//	for task := range t.tasksScheduler {
//
//	}
//}
//
//func (t *TaskManager) taskProcessing(task task_package.Task) {
//	taskProcessingProgress := make(chan task_package.Task)
//
//	go service.EmulateTaskProgress(taskProcessingProgress, task)
//}
//
//func (t *TaskManager) StopTaskManager() {
//
//	t.wg.Done()
//}
//
//func (t *TaskManager) AddTask(task task_package.Task) {
//
//}
//
//func NewTaskManager(db *storage.TaskStorage) *TaskManager {
//	return &TaskManager{
//		db:             db,
//		wg:             sync.WaitGroup{},
//		tasksScheduler: make(chan task_package.Task),
//	}
//}
