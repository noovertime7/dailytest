package worker_demo

import (
	"fmt"
	"k8s.io/apimachinery/pkg/util/wait"
	"log"
	"sync"
	"time"
)

type Worker interface {
	AddWorkerTask(name string)
	Run(name string, period time.Duration, f func()) error
	Stop(name string)
	StopAll()
}

type workers struct {
	Tasks sync.Map
}

type workerTask struct {
	stopCh chan struct{}
}

func NewWorker() Worker {
	return &workers{
		Tasks: sync.Map{},
	}
}

func (w *workers) AddWorkerTask(name string) {
	task := &workerTask{
		stopCh: make(chan struct{}),
	}
	w.Tasks.Store(name, task)
}

func (w *workers) Run(name string, period time.Duration, f func()) error {
	task, ok := w.Tasks.Load(name)
	if !ok {
		return fmt.Errorf("%s not registered", name)
	}

	taskObj := task.(*workerTask)
	go wait.Until(f, period, taskObj.stopCh)
	return nil
}

func (w *workers) Stop(name string) {
	task, ok := w.Tasks.Load(name)
	if !ok {
		return
	}

	taskObj := task.(*workerTask)
	close(taskObj.stopCh)
	w.Tasks.Delete(name)

	log.Println(name, "stop")
}

func (w *workers) StopAll() {
	w.Tasks.Range(func(key, value interface{}) bool {
		taskObj := value.(*workerTask)
		close(taskObj.stopCh)
		w.Tasks.Delete(key)
		return true
	})
}

//func main() {
//	w := NewWorker()
//
//	tasksMap := []string{
//		"task1", "tasl2", "task3",
//	}
//
//	for _, task := range tasksMap {
//		w.AddWorkerTask(task)
//
//		w.Run(task, 1*time.Second, func() {
//			fmt.Println(task, "running")
//		})
//
//		time.Sleep(2 * time.Second)
//		// 停止任务
//
//		w.Stop(task)
//	}
//
//	time.Sleep(2 * time.Hour)
//}
