package protocol

import (
	"ProjectWIND/LOG"
	"ProjectWIND/typed"
	"encoding/json"
)

func HandleMessage(msgJson []byte) {
	var msg typed.MessageEventInfo
	err := json.Unmarshal(msgJson, &msg)
	if err != nil {
		LOG.FATAL("Unmarshalling message: %v", err)
	}
	// 处理消息
	LOG.INFO("收到消息:(来自：%v-%v:%v-%v)%v", msg.MessageType, msg.GroupId, msg.UserId, msg.Sender.Nickname, msg.RawMessage)
	// TODO: 处理消息内容

}

func HandleNotice(msgJson []byte) {
	var notice typed.NoticeEventInfo
	err := json.Unmarshal(msgJson, &notice)
	if err != nil {
		LOG.ERROR("Unmarshalling notice: %v", err)
	}
	// TODO: 处理通知
}

func HandleRequest(msgJson []byte) {
	var request typed.NoticeEventInfo
	err := json.Unmarshal(msgJson, &request)
	if err != nil {
		LOG.ERROR("Unmarshalling request: %v", err)
	}
	// TODO: 处理请求
}

func HandleMetaEvent(msgJson []byte) {
	var meta typed.NoticeEventInfo
	err := json.Unmarshal(msgJson, &meta)
	if err != nil {
		LOG.ERROR("Unmarshalling meta: %v", err)
	}
	// TODO: 处理元事件
}
