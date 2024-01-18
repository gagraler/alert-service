package handler

import (
	"github.com/keington/alertService/internel/biz/models"
	"github.com/keington/alertService/internel/utils"
	"github.com/keington/alertService/pkg/database"
	"log/slog"
)

/**
 * @author: x.gallagher.anderson@gmail.com
 * @time: 2024/1/14 20:41
 * @file: persistence_handle.go
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
	alertList.Level = notification.CommonLabels.Level
	alertList.Status = notification.Status
	alertList.StartTime = notification.Alerts[0].StartsAt
	alertList.EndTime = notification.Alerts[0].EndsAt
	durationTime := notification.Alerts[0].EndsAt.Sub(notification.Alerts[0].StartsAt)
	alertList.DurationTime = utils.ConvertDurationToReadable(durationTime)

	if err := database.DB.Create(&alertList).Error; err != nil {
		slog.Error("Failed to insert data into database", err.Error())
		return
	}
}
