package main

// main.go

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var exitCh = make(chan int)

func HelloWeb(c *gin.Context) {
	c.String(http.StatusOK, "Hello, Go\n")
}

func main() {

	r := gin.Default()
	r.GET("/hello", func(c *gin.Context) {
		HelloWeb(c)
	})

	s := &http.Server{
		Addr:    ":8091",
		Handler: r,
	}

	go func() {
		if err := s.ListenAndServe(); err != nil {
			log.Printf("Listen err: %s\n", err)
		}
	}()

	go gracefulExit(s)

	<-exitCh

	log.Printf("Server exit")
}

func gracefulExit(srv *http.Server) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM, os.Kill)

	sig := <-signalChan
	log.Printf("catch signal, %+v", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second) // 4秒后退出
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Printf("server exiting")
	close(exitCh)
}
