package main

import (
	"fmt"
	"github.com/gorhill/cronexpr"
	"strings"
	"time"
)

// 生成 Cron 表达式
func generateCronExpression(backupPeriod string, daysOfWeek []string, dayOfMonth int, timeOfDay time.Time, repeatDuration int, repeatFrequency string) string {
	cronExpression := ""

	// 设置时间字段
	hour := timeOfDay.Hour()
	minute := timeOfDay.Minute()

	// 构建时间字段
	cronTime := fmt.Sprintf("%d %d", minute, hour)

	// 根据备份周期生成 Cron 表达式
	switch backupPeriod {
	case "每天":
		cronExpression = fmt.Sprintf("%s * * ?", cronTime)
	case "每周":
		if len(daysOfWeek) == 0 {
			// 如果没有选择任何天，默认为周一到周五
			cronExpression = fmt.Sprintf("%s ? * 2-6", cronTime)
		} else {
			// 将选中的星期转换为Cron格式
			weekDays := make([]string, len(daysOfWeek))
			for i, day := range daysOfWeek {
				weekDays[i] = convertDayToCron(day)
			}
			cronExpression = fmt.Sprintf("%s ? * %s", cronTime, strings.Join(weekDays, ","))
		}
	case "每月":
		if dayOfMonth == 0 {
			// 如果没有选择任何日期，默认为每月1日
			dayOfMonth = 1
		}
		cronExpression = fmt.Sprintf("%s %d * ?", cronTime, dayOfMonth)
	}

	// 考虑重复发起的频率和持续时间
	//if repeatDuration > 0 && repeatFrequency != "" {
	//	durationUnit := convertFrequencyToCron(repeatFrequency)
	//	cronExpression += fmt.Sprintf(" /%d %s", repeatDuration, durationUnit)
	//}

	return cronExpression
}

// 将表单中的星期转换为Cron格式
func convertDayToCron(day string) string {
	switch day {
	case "周一":
		return "2"
	case "周二":
		return "3"
	case "周三":
		return "4"
	case "周四":
		return "5"
	case "周五":
		return "6"
	case "周六":
		return "7"
	case "周日":
		return "1"
	default:
		return ""
	}
}

// 将重复频率转换为Cron格式
func convertFrequencyToCron(frequency string) string {
	switch frequency {
	case "小时":
		return "H"
	case "天":
		return "D"
	case "周":
		return "W"
	case "月":
		return "M"
	default:
		return ""
	}
}

func main() {
	times := cronexpr.MustParse(" * * * * ?").NextN(time.Now(), 5)

	for _, t := range times {
		fmt.Println(t.String())
	}

}
