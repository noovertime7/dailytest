package main

import (
	"fmt"
	"time"
)

func convertToCronExpression(timeStr string) (string, error) {
	layout := "2006-01-02 15:04:05"
	targetTime, err := time.Parse(layout, timeStr)
	if err != nil {
		return "", err
	}

	// 生成 Cron 表达式，省略了秒部分
	cronExpression := fmt.Sprintf("%d %d %d %d *", targetTime.Minute(), targetTime.Hour(), targetTime.Day(), targetTime.Month())

	return cronExpression, nil
}

func main() {
	timeStr := "2024-01-04 00:11:24"

	cronExpression, err := convertToCronExpression(timeStr)
	if err != nil {
		fmt.Println("转换为Cron表达式失败:", err)
		return
	}

	fmt.Println("Cron 表达式:", cronExpression)
}
