package main

import (
	"log"
	"time"
	"wait_use_demo/clock"
	"wait_use_demo/wait"
)

func main() {
	stopCh := make(chan struct{})
	realClock := &clock.RealClock{}
	backoffManager := wait.NewExponentialBackoffManager(3*time.Second, 2*time.Minute, 5*time.Minute, 2.0, 1.0, realClock)

	wait.BackoffUntil(func() {
		if err := work(stopCh); err != nil {
			log.Println("error 发生")
		}
	}, backoffManager, false, stopCh)

}

func work(chan struct{}) error {
	log.Println("我被调用了")
	time.Sleep(2 * time.Second)
	return nil
}
