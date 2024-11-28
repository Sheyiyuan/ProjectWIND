package protocol

import (
	"ProjectWIND/typed"
	"encoding/json"
	"log"
)

func HandleMessage(msgJson []byte) {
	var msg typed.MessageEventInfo
	err := json.Unmarshal(msgJson, &msg)
	if err != nil {
		log.Println("[ERROR] unmarshalling message: ", err)
	}
	// 处理消息
	log.Printf("[INFO] 收到消息:(来自：%v-%v:%v-%v)%v", msg.MessageType, msg.GroupId, msg.UserId, msg.Sender.Nickname, msg.RawMessage)
	//一个简单的测试
	if msg.RawMessage == "wind test" {
		log.Println("[INFO] 收到wind test")
		switch msg.MessageType {
		case "group":
			{
				_, err := SendMessage(msg.MessageType, "wind test success", msg.GroupId, false)
				if err != nil {
					log.Println("[ERROR] send message: ", err)
				}
				break
			}
		case "private":
			{
				_, err := SendMessage(msg.MessageType, "wind test success", msg.UserId, false)
				if err != nil {
					log.Println("[ERROR] send message: ", err)
				}
				break
			}
		default:
			{
				log.Println("[ERROR] 不支持的消息类型")
				break
			}
		}
	}
}

func HandleNotice(msgJson []byte) {
	var notice typed.NoticeEventInfo
	err := json.Unmarshal(msgJson, &notice)
	if err != nil {
		log.Println("[ERROR] unmarshalling notice: ", err)
	}
	// 处理通知
}

func HandleRequest(msgJson []byte) {
	var request typed.NoticeEventInfo
	err := json.Unmarshal(msgJson, &request)
	if err != nil {
		log.Println("[ERROR] unmarshalling request: ", err)
	}
	// 处理请求
}

func HandleMetaEvent(msgJson []byte) {
	var meta typed.NoticeEventInfo
	err := json.Unmarshal(msgJson, &meta)
	if err != nil {
		log.Println("[ERROR]  unmarshalling meta: ", err)
	}
	// 处理元事件
}
