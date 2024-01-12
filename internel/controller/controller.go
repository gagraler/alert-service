package controller

import "github.com/gin-gonic/gin"

/**
 * @author: x.gallagher.anderson@gmail.com
 * @time: 2024/1/12 15:35
 * @file: controller.go
 * @description: 路由
 */

// InitializeController 初始化
func InitializeController(r *gin.Engine) {

	apiPath := r.Group("/api/v1/alertService")

	controller(apiPath)
}

func controller(r *gin.RouterGroup) {

	r.POST("/alertMessage/hook", AlertManagerWebhookController)
}
