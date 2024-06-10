package handler

import (
	"bytes"
	"embed"
	_ "embed"
	"fmt"
	"os"
	"strconv"
	"text/template"
	"time"

	"github.com/gagraler/alert-service/internel/biz/models"
	"github.com/gagraler/alert-service/internel/utils"
)

/**
 * @author: gagral.x@gmail.com
 * @time: 2024/1/11 22:22
 * @file: alert manager.go
 * @description: alert manager
 */

//go:embed tmpl/alert_notification.tmpl
var tmplFS embed.FS

var (
	key = os.Getenv("LARK_BOT_SIGN_KEY")
)

type AlertTemplate struct {
	AlertName   string
	AlertLevel  string
	Env         string
	NameSpace   string
	Job         string
	Pod         string
	PromQL      string
	Summary     string
	Description string
	StartsAt    time.Time
}

func (a *AlertTemplate) RenderingAlertTemplate(notification models.Notification) (*models.LarkRequest, error) {

	a.AlertName = notification.Alerts[0].Labels["alertname"]
	a.AlertLevel = notification.Alerts[0].Labels["severity"]
	a.Env = notification.CommonLabels.Env
	a.NameSpace = notification.Alerts[0].Labels["namespace"]
	a.Pod = notification.Alerts[0].Labels["pod"]
	a.Job = notification.Alerts[0].Labels["job"]
	a.PromQL = notification.CommonLabels.PromQL
	a.Summary = notification.Alerts[0].Annotations.Summary
	a.Description = notification.Alerts[0].Annotations.Description
	a.StartsAt = notification.Alerts[0].StartsAt

	tmpl, err := template.ParseFS(tmplFS, "tmpl/alert_notification.tmpl")
	if err != nil {
		return &models.LarkRequest{}, fmt.Errorf("failed to parse template: %v", err)
	}

	var configBuffer bytes.Buffer

	// writer buffer
	if err := tmpl.ExecuteTemplate(&configBuffer, "AlertTemplate", a); err != nil {
		return &models.LarkRequest{}, fmt.Errorf("failed to render template: %v", err)
	}

	return AlertFiringTransformHandle(configBuffer, a.AlertName), nil
}

func AlertFiringTransformHandle(buffer bytes.Buffer, alertName string) *models.LarkRequest {

	var (
		ts = time.Now().Unix()
	)
	sign, err := utils.GenSign(key, ts)
	if err != nil {
		return nil
	}

	firingReq := &models.LarkRequest{
		TimeStamp: strconv.FormatInt(time.Now().Unix(), 10),
		Sign:      sign,
		MsgType:   "interactive",
		Card: models.Card{
			Header: models.Header{
				Title: models.Title{
					Tag:     "plain_text",
					Content: alertName,
				},
				Template: "red",
			},
			Elements: []models.Elements{
				{
					Tag: "div",
					Text: models.Text{
						Content: buffer.String(),
						Tag:     "lark_md",
					},
				},
			},
		},
	}

	return firingReq
}

// AlertResolvedTransformHandle 告警恢复
func AlertResolvedTransformHandle(notification models.Notification) (*models.LarkRequest, error) {
	var (
		alert AlertTemplate
	)

	alertTemplate, err := alert.RenderingAlertTemplate(notification)
	if err != nil {
		return nil, err
	}

	var (
		ts = time.Now().Unix()
	)
	sign, err := utils.GenSign(key, ts)
	if err != nil {
		return nil, err
	}

	resolvedReq := &models.LarkRequest{
		TimeStamp: strconv.FormatInt(time.Now().Unix(), 10),
		Sign:      sign,
		MsgType:   "interactive",
		Card: models.Card{
			Header: models.Header{
				Title: models.Title{
					Tag:     "plain_text",
					Content: notification.Alerts[0].Labels["alertname"],
				},
				Template: "green",
			},
			Elements: []models.Elements{
				{
					Tag: "div",
					Text: models.Text{
						Content: fmt.Sprintf("%v", alertTemplate),
						Tag:     "lark_md",
					},
				},
			},
		},
	}

	return resolvedReq, nil
}
