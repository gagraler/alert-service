package entity

import "time"

/**
 * @author: x.gallagher.anderson@gmail.com
 * @time: 2024/1/14 20:28
 * @file: alert_list_entity.go
 * @description:
 */

type AlertList struct {
	ID           int64     `json:"id" gorm:"id"`
	Env          string    `json:"env" gorm:"env"`
	AlertName    string    `json:"alert_name" gorm:"alert_name"`
	AlertLevel   string    `json:"alert_level" gorm:"alert_level"`
	InstanceName string    `json:"instance_name" gorm:"instance_name"`
	Description  string    `json:"description" gorm:"description"`
	Promql       string    `json:"promql" gorm:"promql"`
	Status       string    `json:"status" gorm:"status"`
	StartTime    time.Time `json:"start_time" gorm:"start_time"`
	EndTime      time.Time `json:"end_time" gorm:"end_time"`
	DurationTime string    `json:"duration_time" gorm:"duration_time"`
}

// TableName 表名称
func (*AlertList) TableName() string {
	return "alert_list"
}
