package controller

import (
	"github.com/bytedance/sonic"
	"github.com/gagraler/alert-service/internal/handle"
	"github.com/gagraler/alert-service/internal/message"
	"github.com/gagraler/alert-service/internal/model"
	"github.com/gagraler/alert-service/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

/**
 * @author: gagral.x@gmail.com
 * @time: 2024/1/11 22:24
 * @file: alert_message_webhook_controller.go
 * @description: lark_webhook_router
 */

var log = logger.SugaredLogger()

// AlertMessageWebhookController 路由
func AlertMessageWebhookController(c *gin.Context) {

	var notification model.Notification

	err := c.ShouldBindJSON(&notification)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	notificationJSON, err := sonic.MarshalString(notification)
	if err != nil {
		log.Errorf("Failed to marshal notification to JSON: %v", err)
	} else {
		log.Debugf("Received AlertManager alarm: %s", string(notificationJSON))
	}

	//log.Debugf("received AlertManager alarm: %s", c.Params)

	req := new(handle.AlertTemplate)
	log.Infof("%s the alert status is: %s", notification.GroupLabels["alertname"], notification.Status)
	larkReqs, err := req.BuildingAlertTemplate(notification)
	if err != nil {
		// Handle the error
		log.Error("failed to transform alertManager notification: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	for _, larkReq := range larkReqs {
		log.Infof("%s the alert is firing and starts sending messages to the lark server", notification.GroupLabels["alertname"])
		log.Infof("alert startAt: %v", notification.Alerts[0].StartsAt)
		message.SendMessageToLarkServer(c, larkReq, notification)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Alerts processed successfully",
	})
}
