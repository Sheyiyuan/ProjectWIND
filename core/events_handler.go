package core

import (
	"ProjectWIND/LOG"
	"ProjectWIND/database"
	"ProjectWIND/wba"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

func HandleMessage(msgJson []byte) {
	var msg wba.MessageEventInfo
	err := json.Unmarshal(msgJson, &msg)
	if err != nil {
		LOG.Error("消息事件反序列化失败: %v", err)
	}
	// 处理消息
	LOG.Info("收到消息:(来自：%v-%v:%v-%v)%v", msg.MessageType, msg.GroupId, msg.UserId, msg.Sender.Nickname, msg.RawMessage)
	isCmd, ok, err := CmdHandle(msg)
	if err != nil {
		LOG.Error("命令处理失败: %v", err)
	}
	if !isCmd && ok {
		// 处理消息
		// TODO: 处理消息
	}
}

func HandleNotice(msgJson []byte) {
	var notice wba.NoticeEventInfo
	err := json.Unmarshal(msgJson, &notice)
	if err != nil {
		LOG.Error("通知事件反序列化失败: %v", err)
	}
	// TODO: 处理通知
}

func HandleRequest(msgJson []byte) {
	var request wba.NoticeEventInfo
	err := json.Unmarshal(msgJson, &request)
	if err != nil {
		LOG.Error("请求事件反序列化失败: %v", err)
	}
	// TODO: 处理请求
}

func HandleMetaEvent(msgJson []byte) {
	var meta wba.NoticeEventInfo
	err := json.Unmarshal(msgJson, &meta)
	if err != nil {
		LOG.Error("元事件反序列化失败: %v", err)
	}
	// TODO: 处理元事件
}

func CmdHandle(msg wba.MessageEventInfo) (isCmd bool, ok bool, err error) {
	// 获取消息的原始文本
	text := msg.GetText()
	// 初始化会话工作状态
	session := wba.SessionInfo{}
	session = session.Load("QQ", msg)
	// 从数据库加载配置
	var key string
	if session.SessionType == "group" {
		key = "group"
	} else {
		key = "user"
	}
	conf, ok := database.MasterGet("GCAS_Config", key, strconv.FormatInt(session.SessionId, 10), "GCAS_Config")
	if !ok {
		conf = "{}"
	}
	if err = GlobalCmdAgentSelector.LoadConfig(session, conf.(string)); err != nil {
		LOG.Error("加载配置失败: %v", err)
		return false, false, err
	}

	// 检查命令前缀
	for _, prefix := range cmdPrefix {
		if strings.HasPrefix(text, prefix) {
			// 解析命令
			text = strings.TrimPrefix(text, prefix)
			cmdSlice := strings.Split(text, " ")
			if text == "" {
				return false, true, nil
			}
			cmd := cmdSlice[0]
			args := cmdSlice[1:]

			result, count := GlobalCmdAgentSelector.FindCmd(cmd, true)
			if count == 0 {
				LOG.Debug("未找到命令: %s", cmd)
				return false, true, nil
			}

			// 执行找到的第一个匹配命令
			cmdFunc := result[0]
			if err := safeExecuteCmd(cmdFunc, args, msg); err != nil {
				LOG.Error("执行命令失败: %v", err)
				return true, false, err
			}
			return true, true, nil
		}
	}
	return false, true, nil
}

func safeExecuteCmd(cmdFunc wba.Cmd, args []string, msg wba.MessageEventInfo) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	cmdFunc.Solve(args, msg)
	return nil
}

var cmdPrefix = []string{"/", "!", "／", "！", ".", "。"}
