package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
)

var errorNum int

func init() {
	flag.IntVar(&errorNum, "errorNum", 10, "set errorNum")
}

func TestGin(ctx *gin.Context) {
	ctx.String(200, "test")
}

func main() {
	flag.Parse()

	r := gin.Default()

	r.GET("/test", TestGin)

	//srv := &http.Server{
	//	Addr:    ":38081",
	//	Handler: r,
	//}
	if err := r.Run(":8081"); err != nil {
		log.Fatal(err)
	}

	//go func() {
	//	log.Info("start...")
	//	if err := srv.ListenAndServe; err != nil {
	//		log.Error(err)
	//	}
	//}()

	//for {
	//	log.WithFields(log.Fields{
	//		"level": log.GetLevel(),
	//	}).Info("测试info日志")
	//	if errorNum > 5 {
	//		log.WithFields(log.Fields{"level": log.GetLevel()}).Error("测试ERROR日志 报错了")
	//		errorNum = 0
	//	}
	//	errorNum++
	//	time.Sleep(1 * time.Second)
	//}
}
func init() {
	file, err := os.OpenFile("testError.log", os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}

	// 日志作为JSON而不是默认的ASCII格式器.
	log.SetFormatter(&log.JSONFormatter{})

	// 输出到标准输出,可以是任何io.Writer
	//log.SetOutput(os.Stdout)
	writers := []io.Writer{
		file,
		os.Stdout}
	//同时写文件和屏幕
	fileAndStdoutWriter := io.MultiWriter(writers...)
	log.SetOutput(fileAndStdoutWriter)
	// 只记录xx级别或以上的日志
	log.SetLevel(log.DebugLevel)
}
