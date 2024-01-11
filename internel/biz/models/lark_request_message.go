package models

/**
 * @author: x.gallagher.anderson@gmail.com
 * @time: 2024/1/11 21:45
 * @file: lark_request_message.go
 * @description: lark_request_message
 */

// LarkMessageRequest 飞书机器人支持的POST数据结构
type LarkMessageRequest struct {
	MsgType string  `json:"msg_type"`
	Content Content `json:"content"`
}

type Content struct {
	Text string `json:"text"`
}

// LarkResponse 响应体相关
type LarkResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data Data   `json:"data"`
}

type Data struct {
}
