package handler

import (
	"fmt"
	"github.com/keington/alertService/internel/biz/models"
	"github.com/keington/alertService/pkg/database"
	"log/slog"
	"time"
)

/**
 * @author: x.gallagher.anderson@gmail.com
 * @time: 2024/1/14 20:41
 * @file: alert_message_persistence_handle.go
 * @description:
 */

func PersistenceHandle(notification models.Notification) {
	var alertList models.AlertList

	// 构建插入数据对象
	alertList.Env = notification.CommonLabels.Env
	alertList.AlertName = notification.Alerts[0].Labels["alertname"]
	alertList.AlertLevel = notification.CommonLabels.Level
	alertList.InstanceName = notification.Alerts[0].Labels["instance"]
	alertList.Description = notification.Alerts[0].Annotations.Description
	alertList.Promql = notification.CommonLabels.PromQL
	alertList.StartTime = notification.Alerts[0].StartsAt
	alertList.EndTime = notification.Alerts[0].EndsAt
	durationTime := notification.Alerts[0].EndsAt.Sub(notification.Alerts[0].StartsAt)
	alertList.DurationTime = convertDurationToReadable(durationTime)

	if err := database.DB.Create(&alertList).Error; err != nil {
		slog.Error("Failed to insert data into database", err.Error())
		return
	}
}

// convertDurationToReadable 将Duration转换为可读的格式
func convertDurationToReadable(duration time.Duration) string {

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
		return "0秒"
	}
	return result
}
