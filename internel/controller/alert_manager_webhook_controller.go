package controller

import (
	"bytes"
	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
	"github.com/keington/alart-service/internel/biz/handler"
	"github.com/keington/alart-service/internel/biz/models"
	"io"
	"log/slog"
	"net/http"
	"time"
)

/**
 * @author: x.gallagher.anderson@gmail.com
 * @time: 2024/1/11 22:24
 * @file: alert_manager_webhook_controller.go
 * @description: lark_webhook_router
 */

const LarkRobotURL = "https://open.larksuite.com/open-apis/bot/v2/hook/27562c31-1810-4c08-b2ef-344ad2b99648"

// AlertManagerWebhookController 飞书机器人的路由
func AlertManagerWebhookController(c *gin.Context) {
	var notification models.Notification

	err := c.ShouldBindJSON(&notification)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	slog.Info("received AlertManager alarm: ", notification)

	// 根据alert manager的请求构造飞书消息的请求数据结构
	larkRequest, err := handler.TransformHandler(notification)
	if err != nil {
		slog.Error("[ERROR] failed to transform alertManager notification: ", err)
	}

	// 向飞书服务器发送POST请求
	bytesData, _ := sonic.Marshal(larkRequest)
	req, _ := http.NewRequest("POST", LarkRobotURL, bytes.NewReader(bytesData))
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	// 飞书服务器可能通信失败
	if err != nil {
		slog.Error("[ERROR] request to lark server failed: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
		}
	}(res.Body)
	body, _ := io.ReadAll(res.Body)

	var larkResponse models.LarkResponse
	err = sonic.Unmarshal(body, &larkResponse)
	// 飞书服务器返回的包可能有问题
	if err != nil {
		slog.Error("[ERROR] failed to obtain response from lark server: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	slog.Info("successfully sent message to lark server")
	timeStamp := time.Now().Local().Format("2006-01-02 15:04:05")
	c.JSON(http.StatusOK, gin.H{
		"message":   "successful receive alert notification message!",
		"timeStamp": timeStamp,
	})
}
