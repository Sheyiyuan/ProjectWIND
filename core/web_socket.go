package core

import (
	"ProjectWIND/LOG"
	"ProjectWIND/wba"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"net/url"
)

var gProtocolAddr string
var gToken string

// WebSocketHandler 接收WebSocket连接处的消息并处理
func WebSocketHandler(protocolAddr string, token string) error {
	// 保存全局变量
	gProtocolAddr = protocolAddr
	gToken = token
	// 解析连接URL
	u, err := url.Parse(protocolAddr)
	if err != nil {
		LOG.ERROR("Parse URL error: %v", err)
		return err
	}

	// 创建一个带有Authorization头的HTTP请求
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		LOG.FATAL("创建请求出错:%v", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	// 配置WebSocket连接升级器
	dialer := websocket.DefaultDialer
	// 使用升级器建立WebSocket连接
	conn, _, err := dialer.Dial(req.URL.String(), req.Header)
	if err != nil {
		LOG.FATAL("建立WebSocket连接出错:%v", err)
	}
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			LOG.ERROR("Close error: %v", err)
		}
	}(conn)
	LOG.INFO("WebSocket connection to %v established.", u.String())
	logInfo := AppApi.GetLoginInfo()
	LOG.INFO("连接到账号: %v（%v）", logInfo.Data.Nickname, logInfo.Data.UserId)

	// 定义通道,缓存消息和消息类型，防止消息处理阻塞
	messageChan := make(chan []byte, 32)
	messageTypeChan := make(chan int, 32)
	for {
		// 接收消息并放入通道
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			LOG.ERROR("ReadMessage error: %v", err)
			return err
		}
		messageChan <- message
		messageTypeChan <- messageType

		// 启动一个新的goroutine来处理消息
		go func() {
			defer func() {
				// 处理完成后从通道中移除消息
				<-messageChan
				<-messageTypeChan
			}()
			processMessage(messageType, message)
		}()
	}
}

// processMessage 处理接收到的消息
func processMessage(messageType int, message []byte) {
	if messageType != websocket.TextMessage {
		LOG.ERROR("Invalid message type: %v", messageType)
		return
	}
	//message json解析
	var messageMap map[string]interface{}
	err := json.Unmarshal(message, &messageMap)
	if err != nil {
		LOG.ERROR("Unmarshal error: %v", err)
		return
	}
	// 处理接收到的消息
	messageTypeStr := messageMap["post_type"]
	switch messageTypeStr {
	case "message":
		{
			// 处理message消息
			HandleMessage(message)
			return
		}
	case "notice":
		{
			// 处理notice消息
			HandleNotice(message)
			return
		}
	case "request":
		{
			// 处理request消息
			HandleRequest(message)
			return
		}
	case "meta_event":
		{
			// 处理meta_event消息
			HandleMetaEvent(message)
			return
		}
	}
}

// wsSendMessage 向WebSocket服务器发送消息并返回发送状态
func wsAPI(body wba.APIRequestInfo) (Response wba.APIResponseInfo, err error) {
	// 序列化请求体
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return wba.APIResponseInfo{}, err
	}
	// 解析连接URL
	u, err := url.Parse(gProtocolAddr)
	if err != nil {
		LOG.ERROR("Parse URL error: %v", err)
		return wba.APIResponseInfo{}, err
	}
	// 创建一个带有Authorization头的HTTP请求
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		LOG.FATAL("创建请求出错:%v", err)
	}
	req.Header.Set("Authorization", "Bearer "+gToken)
	// 配置WebSocket连接升级器
	dialer := websocket.DefaultDialer
	// 使用升级器建立WebSocket连接
	conn, _, err := dialer.Dial(req.URL.String(), req.Header)
	if err != nil {
		LOG.FATAL("建立WebSocket连接出错:%v", err)
	}
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			LOG.ERROR("Close error: %v", err)
		}
	}(conn)
	err = conn.WriteMessage(websocket.TextMessage, bodyBytes)
	if err != nil {
		return wba.APIResponseInfo{}, fmt.Errorf("请求发送失败: %v", err)
	}
	if body.Action == "get_group_list" || body.Action == "get_member_list" {
		// 处理get_group_list和get_member_list请求,直接返回
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				return wba.APIResponseInfo{}, fmt.Errorf("响应接收失败: %v", err)
			}
			var Response wba.APIResponseInfo
			err = json.Unmarshal(message, &Response)
			if err != nil {
				return wba.APIResponseInfo{}, fmt.Errorf("unmarshal error: %v", err)
			}
			if Response.Echo == body.Echo {
				return Response, nil
			}
		}
	}
	//检查是否含有echo字段
	if body.Echo != "" {
		// 接收响应消息,直到收到echo字段一致的消息
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				return wba.APIResponseInfo{}, fmt.Errorf("响应接收失败: %v", err)
			}
			var Response wba.APIResponseInfo
			err = json.Unmarshal(message, &Response)
			if err != nil {
				return wba.APIResponseInfo{}, fmt.Errorf("unmarshal error: %v", err)
			}
			if Response.Echo == body.Echo {
				return Response, nil
			}
		}
	}
	return wba.APIResponseInfo{}, nil
}
