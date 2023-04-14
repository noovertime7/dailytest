package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"strings"
	"time"
)

func main() {
	// 模拟打印error日志 用于日志错误检测demo
	var (
		traceID = ""
		spanId  = ""
	)
	r := gin.Default()
	r.GET("/", func(ctx *gin.Context) {
		log.Println("请求头", ctx.Request.Header)
		allStr := ctx.Request.Header.Get("Traceparent")
		allStrList := strings.Split(allStr, "-")
		if len(allStrList) == 4 {
			traceID = allStrList[1]
			spanId = allStrList[2]
		}
		log.Printf(`{"level":"error","traceID":"%s","spanId":"%s","message":"error模拟日志"}`, traceID, spanId)
		ctx.Writer.WriteString("ok")
	})

	go func() {
		for {
			log.Printf(`{"level":"info","traceID":"%s","spanId":"%s","message":"info模拟日志"}`, traceID, spanId)
			time.Sleep(time.Second)
		}
	}()

	r.Run(":8080")
}
