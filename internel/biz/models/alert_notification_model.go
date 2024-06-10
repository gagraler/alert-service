package models

import "time"

/**
 * @author: gagral.x@gmail.com
 * @time: 2024/1/11 21:43
 * @file: alert_notification_model.go
 * @description: 告警通知的数据结构
 */

type Annotations struct {
	Summary     string `json:"summary"`
	Description string `json:"description"`
}

type Alert struct {
	Status      string            `json:"status"`
	Labels      map[string]string `json:"labels"`
	Annotations Annotations       `json:"annotations"`
	StartsAt    time.Time         `json:"startsAt"`
	EndsAt      time.Time         `json:"endsAt"`
}

type CommonLabels struct {
	Env       string `json:"env"`
	Level     string `json:"level"`
	PromQL    string `json:"promql"`
	AlertType string `json:"alertType"`
}

type Notification struct {
	Receiver          string            `json:"receiver"`
	Status            string            `json:"status"`
	Alerts            []Alert           `json:"alerts"`
	GroupLabels       map[string]string `json:"groupLabels"`
	CommonLabels      CommonLabels      `json:"commonLabels"`
	CommonAnnotations map[string]string `json:"commonAnnotations"`
	ExternalURL       string            `json:"externalURL"`
	Version           string            `json:"version"`
	GroupKey          string            `json:"groupKey"`
}
