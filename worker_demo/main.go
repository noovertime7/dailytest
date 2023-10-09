package main

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"k8s.io/apimachinery/pkg/util/wait"
	"sync"
	"time"
)

type Manager struct {
	Tasks map[string]chan struct{}
	lock  *sync.Mutex
}

func (m *Manager) Run(name string, period time.Duration, f func()) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	stopCh, has := m.Tasks[name]
	if !has {
		return fmt.Errorf("%s not register", name)
	}
	go wait.Until(f, period, stopCh)
	return nil
	// 启动任务
	//go func() {
	//	for {
	//		select {
	//		case <-m.Tasks[name]:
	//			fmt.Println(name, "stopped")
	//			return
	//		default:
	//
	//			// 执行任务逻辑
	//		}
	//	}
	//}()
}

func (m *Manager) Stop(name string) {
	m.lock.Lock()
	defer m.lock.Unlock()
	log.Info("发送stop信号")
	m.Tasks[name] <- struct{}{}
	delete(m.Tasks, name)
}

func main() {
	m := Manager{Tasks: map[string]chan struct{}{
		"task1": make(chan struct{}),
	}, lock: &sync.Mutex{}}

	m.Run("task1", 1*time.Second, func() {
		fmt.Println("task1 run")
	})

	time.Sleep(2 * time.Second)
	// 停止任务

	m.Stop("task1")
	time.Sleep(2 * time.Hour)
}
