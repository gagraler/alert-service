package handler

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gagraler/alert-service/internel/biz/models"
	"github.com/gagraler/alert-service/internel/utils"
)

/**
 * @author: x.gallagher.anderson@gmail.com
 * @time: 2024/1/11 22:22
 * @file: alert manager.go
 * @description: alert manager
 */

// var key = flag.String("key", "", "lark bot sign key")

var (
	key = os.Getenv("LARK_BOT_SIGN_KEY")
)

// JudgeAlertType 根据告警类型，调用不同的处理函数
func JudgeAlertType(notification models.Notification) (*models.LarkRequest, error) {
	alertType := notification.CommonLabels.AlertType

	switch alertType {
	case "container":
		return ContainerTransformHandler(notification)
	case "host":
		return HostTransformHandler(notification)
	case "middleware":
		return MiddleWareTransformHandler(notification)
	default:
		return nil, fmt.Errorf("unsupported alert type: %s", alertType)
	}
}

// buildCommonContent 公共部分的内容
func buildCommonContent(notification models.Notification, builder *strings.Builder) {
	builder.WriteString(fmt.Sprintf("**告警级别:** %s\n", notification.CommonLabels.Level))
	builder.WriteString(fmt.Sprintf("**环境:** %s\n", notification.CommonLabels.Env))
}

// buildAlertContent 告警部分的内容
func buildAlertContent(alert models.Alert, builder *strings.Builder) {
	builder.WriteString(fmt.Sprintf("**告警规则:** %s\n", alert.Labels["alertname"]))
	builder.WriteString(fmt.Sprintf("**摘要:** %s\n**详情:** %s\n", alert.Annotations.Summary, alert.Annotations.Description))
	builder.WriteString(fmt.Sprintf("**开始时间:** %s\n", utils.UTCTranLocal(alert.StartsAt)))
	builder.WriteString("\n")
}

// ContainerTransformHandler 容器告警
func ContainerTransformHandler(notification models.Notification) (*models.LarkRequest, error) {
	var builder strings.Builder
	buildCommonContent(notification, &builder)

	// 每条告警逐个获取，拼接到一起
	for _, alert := range notification.Alerts {
		builder.WriteString(fmt.Sprintf("**命名空间/Pod:** %s:%s\n", alert.Labels["namespace"], alert.Labels["pod"]))
		builder.WriteString(fmt.Sprintf("**实例:** %s\n", alert.Labels["instance"]))
		builder.WriteString(fmt.Sprintf("**PromQL:** %s\n", notification.CommonLabels.PromQL))
		buildAlertContent(alert, &builder)
		builder.WriteString("\n")
	}

	// 构造出飞书机器人所需的数据结构
	return AlertFiringTransformHandle(builder, notification.Alerts[0].Labels["alertname"]), nil
}

// HostTransformHandler 主机告警
func HostTransformHandler(notification models.Notification) (*models.LarkRequest, error) {
	var builder strings.Builder
	buildCommonContent(notification, &builder)

	// 每条告警逐个获取，拼接到一起
	for _, alert := range notification.Alerts {
		builder.WriteString(fmt.Sprintf("**主机名称:** %s\n", alert.Labels["hostname"]))
		builder.WriteString(fmt.Sprintf("**实例:** %s\n", alert.Labels["instance"]))
		buildAlertContent(alert, &builder)
		builder.WriteString("\n")
	}

	// 构造出飞书机器人所需的数据结构
	return AlertFiringTransformHandle(builder, notification.Alerts[0].Labels["alertname"]), nil
}

// MiddleWareTransformHandler 中间件告警
func MiddleWareTransformHandler(notification models.Notification) (*models.LarkRequest, error) {
	var builder strings.Builder
	buildCommonContent(notification, &builder)

	// 每条告警逐个获取，拼接到一起
	for _, alert := range notification.Alerts {
		builder.WriteString(fmt.Sprintf("**实例:** %s\n", alert.Labels["instance"]))
		buildAlertContent(alert, &builder)
		builder.WriteString("\n")
	}

	// 构造出飞书机器人所需的数据结构
	return AlertFiringTransformHandle(builder, notification.Alerts[0].Labels["alertname"]), nil
}

func AlertFiringTransformHandle(builder strings.Builder, alertName string) *models.LarkRequest {

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
						Content: builder.String(),
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
		builder strings.Builder
	)

	buildCommonContent(notification, &builder)
	builder.WriteString(fmt.Sprintf("**告警状态:** %s\n", notification.Status))

	// 每条告警逐个获取，拼接到一起
	for _, alert := range notification.Alerts {
		builder.WriteString(fmt.Sprintf("**实例:** %s\n", alert.Labels["instance"]))
		builder.WriteString(fmt.Sprintf("**告警规则:** %s\n", alert.Labels["alertname"]))
		builder.WriteString(fmt.Sprintf("**开始时间:** %s\n", utils.UTCTranLocal(alert.StartsAt)))
		builder.WriteString(fmt.Sprintf("**恢复时间:** %s\n", utils.UTCTranLocal(alert.EndsAt)))
		builder.WriteString(fmt.Sprintf("**持续时间:** %s\n", utils.ConvertDurationToReadable(alert.EndsAt.Sub(alert.StartsAt))))
		builder.WriteString("\n")
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
						Content: builder.String(),
						Tag:     "lark_md",
					},
				},
			},
		},
	}

	return resolvedReq, nil
}
