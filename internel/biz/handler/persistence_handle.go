package handler

import (
	"log/slog"

	"github.com/gagraler/alert-service/internel/biz/models"
	"github.com/gagraler/alert-service/internel/biz/models/entity"
	"github.com/gagraler/alert-service/internel/utils"
	"github.com/gagraler/alert-service/pkg/database"
)

/**
 * @author: x.gallagher.anderson@gmail.com
 * @time: 2024/1/14 20:41
 * @file: persistence_handle.go
 * @description:
 */

func PersistenceHandle(notification models.Notification) {
	var alertList entity.AlertList

	// 构建插入数据对象
	alertList.Env = notification.CommonLabels.Env
	alertList.AlertName = notification.Alerts[0].Labels["alertname"]
	alertList.AlertLevel = notification.CommonLabels.Level
	alertList.InstanceName = notification.Alerts[0].Labels["instance"]
	alertList.Description = notification.Alerts[0].Annotations.Description
	alertList.Promql = notification.CommonLabels.PromQL
	alertList.Status = notification.Status
	alertList.StartTime = notification.Alerts[0].StartsAt
	alertList.EndTime = notification.Alerts[0].EndsAt
	durationTime := notification.Alerts[0].EndsAt.Sub(notification.Alerts[0].StartsAt)
	alertList.DurationTime = utils.ConvertDurationToReadable(durationTime)

	if err := database.DB.Create(&alertList).Error; err != nil {
		slog.Error("Failed to insert data into database", err)
		return
	}
}
