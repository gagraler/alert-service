package message

import (
	"bytes"
	"github.com/bytedance/sonic"
	handler2 "github.com/gagraler/alert-service/internal/handle"
	models2 "github.com/gagraler/alert-service/internal/model"
	"github.com/gagraler/alert-service/pkg/logger"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
)

/**
 * @author: gagral.x@gmail.com
 * @time: 2024/6/14 20:47
 * @file: lark.go
 * @description: 发送消息到飞书机器人
 */

var hookUrl = os.Getenv("LARK_BOT_URL")
var log = logger.SugaredLogger()

// SendMessageToLarkServer 发送消息到飞书机器人
func SendMessageToLarkServer(c *gin.Context, larkRequest *models2.LarkRequest, notification models2.Notification) {

	bytesData, _ := sonic.Marshal(larkRequest)
	req, _ := http.NewRequest(http.MethodPost, hookUrl, bytes.NewReader(bytesData))
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
	go handler2.PersistenceHandle(notification)

	body, _ := io.ReadAll(res.Body)
	var larkRes models2.LarkResponse
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

	log.Infof("%s alert request to lark server result, code: %v, internal: %s", notification.GroupLabels["alertname"], larkRes.Code, larkRes.Msg)
}
