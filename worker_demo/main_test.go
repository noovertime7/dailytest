package worker_demo

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestWorker_Concurrent(t *testing.T) {
	w := NewWorker()

	tasksMap := []string{
		"task1", "task2", "task3",
	}

	var wg sync.WaitGroup
	concurrentTasks := 10

	for i := 0; i < concurrentTasks; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for _, task := range tasksMap {
				w.AddWorkerTask(task)

				err := w.Run(task, 2*time.Second, func() {
					fmt.Println(task, "running")
				})
				if err != nil {
					t.Errorf("Error running task %s: %v", task, err)
				}

				time.Sleep(2 * time.Second)

				w.Stop(task)
			}
		}()
	}

	wg.Wait()

}
