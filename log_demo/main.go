package main

import (
	"flag"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"time"
)

var errorNum int

func init() {
	flag.IntVar(&errorNum, "errorNum", 10, "set errorNum")
}

func main() {
	flag.Parse()
	for {
		log.WithFields(log.Fields{
			"level": log.GetLevel(),
		}).Info("测试info日志")
		if errorNum > 5 {
			log.WithFields(log.Fields{
				"level": log.GetLevel(),
			}).Error("测试ERROR日志 报错了")
			errorNum = 0
		}
		errorNum++
		time.Sleep(1 * time.Second)
	}
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
