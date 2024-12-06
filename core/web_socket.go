package core

import (
	"ProjectWIND/LOG"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
)

var gProtocolAddr string

// WebSocketHandler 接收WebSocket连接处的消息并处理
func WebSocketHandler(protocolAddr string) error {
	// 保存全局变量
	gProtocolAddr = protocolAddr
	// 解析连接URL
	u, err := url.Parse(protocolAddr)
	if err != nil {
		LOG.ERROR("Parse URL error: %v", err)
		return err
	}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		LOG.ERROR("Dial error: %v", err)
		return err
	}
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			LOG.ERROR("Close error: %v", err)
		}
	}(conn)

	LOG.INFO("WebSocket connection to %v established.", u.String())

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
	default:
		{
			// 打印接收到的消息
			LOG.WARN("Received unknown event: %s", message)
		}
	}
}

// wsSendMessage 向WebSocket服务器发送消息并返回发送状态
func wsAPI(body []byte) error {
	// 解析连接URL
	u, err := url.Parse(fmt.Sprintf("%v/api", gProtocolAddr))
	if err != nil {
		return fmt.Errorf("无效的URL: %v", err)
	}

	// 建立连接
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return fmt.Errorf("连接失败: %v", err)
	}
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			LOG.ERROR("Close error: %v", err)
		}
	}(conn)

	// 发送请求
	err = conn.WriteMessage(websocket.TextMessage, body)
	if err != nil {
		return fmt.Errorf("请求发送失败: %v", err)
	}

	return nil
}

func httpAPI(method, action string, body []byte) (int, []byte, error) {
	urlStr := fmt.Sprintf("%v/api/%v", gProtocolAddr, action)
	resp, err := http.Post(urlStr, "application/json", bytes.NewReader(body))
	if err != nil {
		return 0, nil, fmt.Errorf("请求失败: %v", err)
	}
	defer func(resp *http.Response) {
		err := resp.Body.Close()
		if err != nil {
			LOG.ERROR("Close error: %v", err)
		}
	}(resp)

	if resp.StatusCode != http.StatusOK {
		return 0, nil, fmt.Errorf("请求失败: %v", resp.Status)
	}

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, fmt.Errorf("读取响应失败: %v", err)
	}
	return resp.StatusCode, body, nil
}
