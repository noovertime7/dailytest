package main

import (
	"github.com/sirupsen/logrus"
	"math/rand"
	"time"
)

func main() {
	// 模拟打印error日志 用于日志错误检测demo
	for {
		res := rand.Intn(100)
		if res > 50 {
			logrus.Error(res)
		} else {
			logrus.Info(res)
		}
		time.Sleep(3 * time.Second)
	}
}
