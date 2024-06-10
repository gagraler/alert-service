package controller

import (
	"bytes"
	"github.com/gagraler/alert-service/pkg/logger"
	"io"
	"net/http"
	"os"

	"github.com/bytedance/sonic"
	"github.com/gagraler/alert-service/internel/biz/handler"
	"github.com/gagraler/alert-service/internel/biz/models"
	"github.com/gin-gonic/gin"
)

/**
 * @author: gagral.x@gmail.com
 * @time: 2024/1/11 22:24
 * @file: alert_message_webhook_controller.go
 * @description: lark_webhook_router
 */

var (
	hookUrl = os.Getenv("LARK_BOT_URL")
	log     = logger.SugaredLogger()
)

// AlertMessageWebhookController 路由
func AlertMessageWebhookController(c *gin.Context) {

	var notification models.Notification

	err := c.ShouldBindJSON(&notification)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	log.Debugf("received AlertManager alarm: %s", notification)

	switch notification.Alerts[0].Status {
	case "resolved":
		handleResolvedAlert(c, notification)
	case "firing":
		handleFiringAlert(c, notification)
	default:
		log.Info("unknown alert status, skip sending message to lark server")
		c.JSON(http.StatusOK, gin.H{
			"message": "unknown alert status, skip sending message to lark server",
		})
	}
}

// handleResolvedAlert 处理告警恢复的情况
func handleResolvedAlert(c *gin.Context, notification models.Notification) {
	larkReq, err := handler.AlertResolvedTransformHandle(notification)
	if err != nil {
		log.Errorf("failed to transform alertManager notification: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	log.Info("alert has been resolved, skip sending message to lark server")

	sendMessageToLarkServer(c, larkReq, notification)
}

// handleFiringAlert 处理告警触发的情况
func handleFiringAlert(c *gin.Context, notification models.Notification) {

	req := new(handler.AlertTemplate)
	log.Infof("%s the alert status is: %s", notification.GroupLabels["alertname"], notification.Status)
	larkReq, err := req.RenderingAlertTemplate(notification)
	if err != nil {
		// Handle the error
		log.Error("failed to transform alertManager notification: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	log.Infof("%s the alert is firing and starts sending messages to the lark server", notification.GroupLabels["alertname"])
	log.Infof("alert startAt: %v", notification.Alerts[0].StartsAt)
	sendMessageToLarkServer(c, larkReq, notification)
}

// sendMessageToLarkServer 发送消息到飞书机器人
func sendMessageToLarkServer(c *gin.Context, larkRequest *models.LarkRequest, notification models.Notification) {

	bytesData, _ := sonic.Marshal(larkRequest)
	req, _ := http.NewRequest("POST", hookUrl, bytes.NewReader(bytesData))
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error("request to lark server failed: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Error("failed to close response body: ", err)
		}
	}(res.Body)

	// 持久化
	go handler.PersistenceHandle(notification)

	body, _ := io.ReadAll(res.Body)
	var larkRes models.LarkResponse
	err = sonic.Unmarshal(body, &larkRes)
	if err != nil {
		log.Error("failed to obtain response from lark server: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    larkRes.Code,
		"message": larkRes.Msg,
		"data":    larkRes.Data,
	})

	log.Infof("%s alert request to lark server result, code: %v, message: %s", notification.GroupLabels["alertname"], larkRes.Code, larkRes.Msg)
}
