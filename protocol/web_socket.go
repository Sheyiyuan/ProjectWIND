package protocol

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

const (
	severURL = "your_url"
)

func WebSocketHandler() (*websocket.Conn, error) {
	u, err := url.Parse(severURL)
	if err != nil {
		log.Println("Parse URL error:", err)
		return nil, err
	}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Println("[ERROR] Dial error:", err)
		return nil, err
	}
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			log.Println("[ERROR] Close error:", err)
		}
	}(conn)

	log.Println("[INFO] New connection established.")

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("[ERROR] Read error:", err)
			break
		}

		// 将接收到的消息交给另一个函数处理
		processMessage(messageType, message)
	}
	return conn, nil
}

func processMessage(messageType int, message []byte) {
	if messageType != websocket.TextMessage {
		log.Println("[INFO] Invalid message type:", messageType)
		return
	}
	//message json解析
	var messageMap map[string]interface{}
	err := json.Unmarshal(message, &messageMap)
	if err != nil {
		log.Println("[ERROR] Unmarshal error:", err)
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
			log.Printf("[WARN] Received message: %s", message)
		}
	}
}

// wsSendMessage 向WebSocket服务器发送消息并返回发送状态
func wsSendMessage(message []byte) (bool, error) {
	// 解析连接URL
	u, err := url.Parse(fmt.Sprintf("%v/api", severURL))
	if err != nil {
		return false, fmt.Errorf("无效的URL: %v", err)
	}

	// 建立连接
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return false, fmt.Errorf("连接失败: %v", err)
	}
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			log.Println("[ERROR] Close error:", err)
		}
	}(conn)

	// 发送消息
	err = conn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		return false, fmt.Errorf("发送消息失败: %v", err)
	}

	return true, nil
}
