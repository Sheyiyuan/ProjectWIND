package core

import (
	"ProjectWIND/LOG"
	"ProjectWIND/wba"
	"encoding/json"
	"fmt"
	"strings"
)

func HandleMessage(msgJson []byte) {
	var msg wba.MessageEventInfo
	err := json.Unmarshal(msgJson, &msg)
	if err != nil {
		LOG.ERROR("消息事件反序列化失败: %v", err)
	}
	// 处理消息
	LOG.INFO("收到消息:(来自：%v-%v:%v-%v)%v", msg.MessageType, msg.GroupId, msg.UserId, msg.Sender.Nickname, msg.RawMessage)
	//如果消息文本内容为bot，发送框架信息。
	cmd, args := CmdSplit(msg)
	_, ok := CmdMap[cmd]
	if ok {
		LOG.DEBUG("执行命令：%v %v", cmd, args)
		CmdMap[cmd].SOLVE(args, msg)
	}
	// TODO: 处理消息内容
}

func HandleNotice(msgJson []byte) {
	var notice wba.NoticeEventInfo
	err := json.Unmarshal(msgJson, &notice)
	if err != nil {
		LOG.ERROR("通知事件反序列化失败: %v", err)
	}
	// TODO: 处理通知
}

func HandleRequest(msgJson []byte) {
	var request wba.NoticeEventInfo
	err := json.Unmarshal(msgJson, &request)
	if err != nil {
		LOG.ERROR("请求事件反序列化失败: %v", err)
	}
	// TODO: 处理请求
}

func HandleMetaEvent(msgJson []byte) {
	var meta wba.NoticeEventInfo
	err := json.Unmarshal(msgJson, &meta)
	if err != nil {
		LOG.ERROR("元事件反序列化失败: %v", err)
	}
	// TODO: 处理元事件
}

func CmdSplit(msg wba.MessageEventInfo) (string, []string) {
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
			for cmd := range CmdMap {
				if strings.HasPrefix(text, cmd) {
					text = strings.TrimPrefix(text, cmd)
					text = strings.TrimPrefix(text, " ")
					return cmd, strings.Split(text, " ")
				}
			}
		}
	}
	return "", []string{}
}

func statusCheck(msg wba.MessageEventInfo) bool {
	//TODO: 检查当前组群工作状态
	return false
}

var cmdPrefix = []string{"/", "!", "／", "！", ".", "。"}
