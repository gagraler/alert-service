package models

/**
 * @author: gagral.x@gmail.com
 * @time: 2024/1/11 21:45
 * @file: lark_mutual_model.go
 * @description: 飞书交互的数据结构
 */

type Title struct {
	Tag     string `json:"tag"`
	Content string `json:"content"`
}

type Header struct {
	Title    Title  `json:"title"`
	Template string `json:"template"`
}

type Text struct {
	Content string `json:"content"`
	Tag     string `json:"tag"`
}

type Elements struct {
	Tag  string `json:"tag"`
	Text Text   `json:"text"`
}

type Card struct {
	Header   Header     `json:"header"`
	Elements []Elements `json:"elements"`
}

// LarkRequest 飞书机器人支持的数据结构
type LarkRequest struct {
	TimeStamp string `json:"timestamp"`
	Sign      string `json:"sign"`
	MsgType   string `json:"msg_type"`
	Card      Card   `json:"card"`
}

// LarkResponse 响应体相关
type LarkResponse struct {
	Code int      `json:"code"`
	Msg  string   `json:"msg"`
	Data struct{} `json:"data"`
}
