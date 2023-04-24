package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM, os.Kill)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-signalChan:
				log.Printf("exit")
				return
			default:
			}

			log.Printf("default waiting")
			time.Sleep(10 * time.Second)
			select {
			case <-signalChan:
				log.Printf("work finish exit")
				return
			}
		}
	}()
	wg.Wait()
	//sig := <-signalChan
	//log.Printf("catch signal, %+v", sig)

}
