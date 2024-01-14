package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/keington/alertService/internel/biz/handler"
	"github.com/keington/alertService/internel/biz/models"
	"net/http"
)

/**
 * @author: x.gallagher.anderson@gmail.com
 * @time: 2024/1/14 18:45
 * @file: alert_message_persistence_controller.go
 * @description: 告警持久化
 */

func AlertMessagePersistenceController(c *gin.Context) {
	var notification models.Notification

	if err := c.BindJSON(&notification); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to parse JSON payload",
		})
	}

	// 持久化
	handler.PersistenceHandle(notification)
}
