package models

/**
 * @author: x.gallagher.anderson@gmail.com
 * @time: 2024/1/11 21:45
 * @file: lark_mutual_model.go
 * @description: 飞书交互的数据结构
 */

type Content struct {
	Text string `json:"text"`
}

// LarkRequest 飞书机器人支持的POST数据结构
type LarkRequest struct {
	MsgType string  `json:"msg_type"`
	Content Content `json:"content"`
}

type Data struct {
}

// LarkResponse 响应体相关
type LarkResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data Data   `json:"data"`
}
