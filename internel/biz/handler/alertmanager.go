package handler

import (
	"bytes"
	"fmt"
	"github.com/keington/alart-service/internel/biz/models"
)

/**
 * @author: x.gallagher.anderson@gmail.com
 * @time: 2024/1/11 22:22
 * @file: alert manager.go
 * @description: alert manager
 */

// TransformToLarkHandler 根据alertManager的对象，创建出飞书消息的对象
func TransformToLarkHandler(notification models.Notification) (larkRequest *models.LarkMessageRequest, err error) {
	var buffer bytes.Buffer

	// 先拿到分组情况
	buffer.WriteString(fmt.Sprintf("通知组%s，状态[%s]\n告警项\n\n", notification.GroupKey, notification.Status))

	// 每条告警逐个获取，拼接到一起
	for _, alert := range notification.Alerts {
		buffer.WriteString(fmt.Sprintf("摘要：%s\n详情：%s\n", alert.Annotations["summary"], alert.Annotations["description"]))
		buffer.WriteString(fmt.Sprintf("开始时间: %s\n\n", alert.StartsAt.Format("15:04:05")))
	}

	// 构造出飞书机器人所需的数据结构
	larkRequest = &models.LarkMessageRequest{
		MsgType: "text",
		Content: models.Content{
			Text: buffer.String(),
		},
	}

	return larkRequest, nil
}
