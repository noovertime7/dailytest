package main

import (
	"github.com/gin-gonic/gin"
	"github.com/noovertime7/dailytest/websocket-manager/manager"
	"log"
	"os"
	"os/signal"
)

func main() {
	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	r := gin.Default()
	gin.SetMode(gin.DebugMode)
	r.GET("/ws/:name", func(ctx *gin.Context) {
		manager.WebSocketConnect(ctx)
	})
	if err := r.Run(":9091"); err != nil {
		log.Fatalln(err)
	}
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
}
