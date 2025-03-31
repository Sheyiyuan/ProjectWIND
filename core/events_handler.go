package core

import (
	"ProjectWIND/LOG"
	"ProjectWIND/database"
	"ProjectWIND/typed"
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
	CmdHandle(msg)
	fmt.Printf("%#v\n", msg.Message)
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

func CmdSplit(msg wba.MessageEventInfo) (string, []string) {
	//text := msg.RawMessage
	//if strings.HasPrefix(text, fmt.Sprintf("[CQ:at,qq=%d]", msg.SelfId)) {
	//	text = strings.TrimPrefix(text, fmt.Sprintf("[CQ:at,qq=%d]", msg.SelfId))
	//} else {
	//	if !statusCheck(msg) {
	//		return "", []string{}
	//	}
	//}
	////检查有无application.CmdMap中的命令前缀
	//for _, prefix := range cmdPrefix {
	//	if strings.HasPrefix(text, prefix) {
	//		text = strings.TrimPrefix(text, prefix)
	//		for cmdList := range CmdMap {
	//			for cmd := range CmdMap[cmdList] {
	//				if strings.HasPrefix(text, cmd) {
	//					text = strings.TrimPrefix(text, cmd)
	//					text = strings.TrimPrefix(text, " ")
	//					return cmd, strings.Split(text, " ")
	//				}
	//			}
	//		}
	//	}
	//}
	return "", []string{}
}

func CmdHandle(msg wba.MessageEventInfo) {
	// 获取消息的原始文本
	text := msg.GetText()
	// 初始化会话工作状态
	WorkStatus := typed.SessionWorkSpace{}

	if msg.MessageType == "group" {
		// 如果消息类型是群消息
		// 从数据库中获取工作状态信息
		WorkStatusJson, ok := database.MasterGet("WorkStatus", "global", strconv.FormatInt(msg.SelfId, 10), strconv.FormatInt(msg.GroupId, 10))
		// 尝试将获取到的工作状态信息反序列化为 WorkStatus 结构体
		err := json.Unmarshal([]byte(WorkStatusJson.(string)), &WorkStatus)
		// 如果获取失败或反序列化失败
		if !ok || err != nil {
			// 初始化一个新的工作状态
			WorkStatus = typed.SessionWorkSpace{
				SessionId: fmt.Sprintf("%d-%d", msg.SelfId, msg.GroupId),
				Rule:      "coc",
				Enable:    true,
				AppEnable: make(map[wba.AppKey]bool),
				CmdEnable: make(map[string]bool),
				WorkLevel: 0,
			}
			// 将新的工作状态序列化为 JSON 字符串
			WorkStatusJson, err := json.Marshal(WorkStatus)
			if err != nil {
				// 如果序列化失败，记录错误信息
				LOG.Error("命令处理过程中，WorkStatusJson序列化失败: %v", err)
			}
			// 将序列化后的工作状态保存到数据库中
			database.MasterSet("WorkStatus", "global", strconv.FormatInt(msg.SelfId, 10), strconv.FormatInt(msg.GroupId, 10), string(WorkStatusJson))
		}
	}

	// 调用 CmdSplit 函数将消息文本拆分为命令和参数
	cmd, args := CmdSplit(msg)

	// 检查是否找到了有效的命令
	if cmd != "" {
		// 遍历命令映射表 CmdMap
		for _, cmdList := range CmdMap {
			for cmdKey, Cmd := range cmdList {
				// 如果找到了匹配的命令
				if cmdKey == cmd {
					// 检查会话工作状态
					if WorkStatus.SessionId != "" {
						if !WorkStatus.Enable {
							// 如果会话未启用，忽略该命令
							LOG.Debug("忽略指令:%s,当前会话中bot未启用", cmd)
							continue
						}
						// 检查 APP 或命令是否启用，以及工作级别和规则是否匹配
						if !WorkStatus.AppEnable[Cmd.AppKey] || !WorkStatus.CmdEnable[cmd] || WorkStatus.WorkLevel > Cmd.AppKey.Level || (WorkStatus.Rule != Cmd.AppKey.Rule && WorkStatus.Rule != "none") {
							// 如果不满足条件，忽略该命令
							LOG.Debug("忽略指令:%s,当前会话中APP或命令未启用", cmd)
							continue
						}
					}
					// 执行命令
					Cmd.Solve(args, msg)
					// 记录执行的命令和参数
					LOG.Info("执行命令：%v %v", cmd, args)
				}
			}
		}
	} else {
		// 如果未找到有效的命令，记录未找到命令的信息
		LOG.Info("未找到命令：%v", strings.Split(text, " ")[0])
	}
}

var cmdPrefix = []string{"/", "!", "／", "！", ".", "。"}
