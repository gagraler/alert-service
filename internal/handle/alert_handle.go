package handle

import (
	"bytes"
	"embed"
	_ "embed"
	"fmt"
	models2 "github.com/gagraler/alert-service/internal/model"
	"os"
	"strconv"
	"text/template"
	"time"

	"github.com/gagraler/alert-service/internal/util"
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

// BuildingAlertTemplate 构建告警模板
func (a *AlertTemplate) BuildingAlertTemplate(notification models2.Notification) (*models2.LarkRequest, error) {

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
		return &models2.LarkRequest{}, fmt.Errorf("failed to parse template: %v", err)
	}

	var configBuffer bytes.Buffer

	// writer buffer
	if err := tmpl.ExecuteTemplate(&configBuffer, "AlertTemplate", a); err != nil {
		return &models2.LarkRequest{}, fmt.Errorf("failed to render template: %v", err)
	}

	handle := AlertHandle(configBuffer, a.AlertName)
	if notification.Status == "firing" {
		handle.Card.Header.Template = "red"
	} else {
		handle.Card.Header.Template = "green"
	}

	return handle, nil
}

// AlertHandle 告警处理
func AlertHandle(buffer bytes.Buffer, alertName string) *models2.LarkRequest {

	var (
		ts = time.Now().Unix()
	)
	sign, err := util.GenSign(key, ts)
	if err != nil {
		return nil
	}

	alertReq := &models2.LarkRequest{
		TimeStamp: strconv.FormatInt(time.Now().Unix(), 10),
		Sign:      sign,
		MsgType:   "interactive",
		Card: models2.Card{
			Header: models2.Header{
				Title: models2.Title{
					Tag:     "plain_text",
					Content: alertName,
				},
				Template: "red",
			},
			Elements: []models2.Elements{
				{
					Tag: "div",
					Text: models2.Text{
						Content: buffer.String(),
						Tag:     "lark_md",
					},
				},
			},
		},
	}

	return alertReq
}
