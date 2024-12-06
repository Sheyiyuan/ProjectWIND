package core

import (
	"ProjectWIND/LOG"
	"ProjectWIND/typed"
	"encoding/json"
	"fmt"
	"strings"
)

func HandleMessage(msgJson []byte) {
	var msg typed.MessageEventInfo
	err := json.Unmarshal(msgJson, &msg)
	if err != nil {
		LOG.ERROR("Unmarshalling message: %v", err)
	}
	// 处理消息
	LOG.INFO("收到消息:(来自：%v-%v:%v-%v)%v", msg.MessageType, msg.GroupId, msg.UserId, msg.Sender.Nickname, msg.RawMessage)
	//如果消息文本内容为bot，发送框架信息。
	cmd, args := CmdSplit(msg)
	_, ok := CmdList[cmd]
	if ok {
		err = CmdList[cmd].Run(cmd, args, msg)
		if err != nil {
			LOG.ERROR("消息发送失败: %v", err)
		}
	}
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

func CmdSplit(msg typed.MessageEventInfo) (string, []string) {
	text := msg.RawMessage
	if strings.HasPrefix(text, fmt.Sprintf("[CQ:at,qq=%d]", msg.SelfId)) {
		text = strings.TrimPrefix(text, fmt.Sprintf("[CQ:at,qq=%d]", msg.SelfId))
	} else {
		if statusCheck(msg) {
			return "", []string{}
		}
	}
	//检查有无application.CmdList中的命令前缀
	for _, prefix := range cmdPrefix {
		if strings.HasPrefix(text, prefix) {
			text = strings.TrimPrefix(text, prefix)
			for cmd := range CmdList {
				if strings.HasPrefix(text, cmd) {
					text = strings.TrimPrefix(text, cmd)
					return cmd, strings.Split(text, " ")
				}
			}
		}
	}
	return "", []string{}
}

func statusCheck(msg typed.MessageEventInfo) bool {
	//TODO: 检查当前组群工作状态
	return false
}

var cmdPrefix = []string{"/", "!", "／", "！", ".", "。"}
