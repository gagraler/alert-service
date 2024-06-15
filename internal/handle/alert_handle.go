package handle

import (
	"bytes"
	"embed"
	_ "embed"
	"fmt"
	"github.com/gagraler/alert-service/internal/model"
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
	AlertStatus string
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
func (a *AlertTemplate) BuildingAlertTemplate(notification model.Notification) ([]*model.LarkRequest, error) {

	var req []*model.LarkRequest

	for _, v := range notification.Alerts {
		a.AlertName = v.Labels["alertname"]
		a.AlertStatus = v.Status
		a.AlertLevel = v.Labels["severity"]
		a.Env = v.Labels["env"]
		a.NameSpace = v.Labels["namespace"]
		a.Pod = v.Labels["pod"]
		a.Job = v.Labels["job"]
		a.PromQL = v.Labels["expr"]
		a.Summary = v.Annotations.Summary
		a.Description = v.Annotations.Description
		a.StartsAt = v.StartsAt

		tmpl, err := template.ParseFS(tmplFS, "tmpl/alert_notification.tmpl")
		if err != nil {
			return nil, fmt.Errorf("failed to parse template: %v", err)
		}

		var configBuffer bytes.Buffer

		// writer buffer
		if err := tmpl.ExecuteTemplate(&configBuffer, "AlertTemplate", a); err != nil {
			return nil, fmt.Errorf("failed to render template: %v", err)
		}

		handle := AlertHandle(configBuffer, a.AlertName)
		if v.Status == "firing" {
			handle.Card.Header.Template = "red"
		} else {
			handle.Card.Header.Template = "green"
		}

		req = append(req, handle)
	}

	if len(req) == 0 {
		return nil, fmt.Errorf("no alerts found in the notification")
	}

	return req, nil
}

// AlertHandle 告警处理
func AlertHandle(buffer bytes.Buffer, alertName string) *model.LarkRequest {

	var (
		ts = time.Now().Unix()
	)
	sign, err := util.GenSign(key, ts)
	if err != nil {
		return nil
	}

	alertReq := &model.LarkRequest{
		TimeStamp: strconv.FormatInt(time.Now().Unix(), 10),
		Sign:      sign,
		MsgType:   "interactive",
		Card: model.Card{
			Header: model.Header{
				Title: model.Title{
					Tag:     "plain_text",
					Content: alertName,
				},
				Template: "red",
			},
			Elements: []model.Elements{
				{
					Tag: "div",
					Text: model.Text{
						Content: buffer.String(),
						Tag:     "lark_md",
					},
				},
			},
		},
	}

	return alertReq
}
