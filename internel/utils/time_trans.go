package utils

import (
	"fmt"
	"time"
)

/**
 * @author: x.gallagher.anderson@gmail.com
 * @time: 2024/1/15 21:41
 * @file: time_trans.go
 * @description: time convert
 */

// ConvertDurationToReadable 将Duration转换为可读的格式
func ConvertDurationToReadable(duration time.Duration) string {

	days := int(duration.Hours() / 24)
	hours := int(duration.Hours()) % 24
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60
	milliseconds := duration.Milliseconds() % 1000

	var result string

	if days > 0 {
		result += fmt.Sprintf("%d天", days)
	}
	if hours > 0 {
		result += fmt.Sprintf("%d小时", hours)
	}
	if minutes > 0 {
		result += fmt.Sprintf("%d分钟", minutes)
	}
	if seconds > 0 {
		result += fmt.Sprintf("%d秒", seconds)
	}
	if milliseconds > 0 {
		result += fmt.Sprintf("%d毫秒", milliseconds)
	}
	if result == "" {
		return "-"
	}
	return result
}

func UTCTranLocal(utcTime time.Time) string {

	var layout = "2006-01-02 15:04:05"
	return utcTime.Local().Format(layout)
}
