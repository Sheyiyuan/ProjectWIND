package core

import (
	"ProjectWIND/LOG"
	"ProjectWIND/wba"
	"crypto/rand"
	"fmt"
)

type protocolAPI struct{}

//一、Protocol模块

/*
关于Protocol模块的说明

1.所有API请求按照OneBot11标准，使用JSON格式进行数据交换。api命名为由原文档中蛇形命名法改为双驼峰命名法。

2.无响应的API请求使用ws协议处理，有响应的API需添加echo字段。

3.wind会从配置文件中读取API请求的url，请确保正确填写。
*/

//1.无响应API,使用ws协议处理

// UnsafelySendMsg 发送消息(自动判断消息类型)
//
// 注意：该API不安全，除非你知道你在做什么。
//
// 参数：
//
//   - messageType: 消息类型，可选值为 "private" 和 "group
//
//   - groupId: 群号，当messageType为"group"时必填
//
//   - userId: 用户ID，当messageType为"private"时必填
//
//   - message: 要发送的消息内容，类型为字符串
//
//   - autoEscape: 是否自动转义（解析消息内容中的CQ码），可选值为 true 和 false
func (p *protocolAPI) UnsafelySendMsg(messageType string, groupId int64, userId int64, message string, autoEscape bool) {
	// 构建发送消息的JSON数据
	var messageData wba.APIRequestInfo
	messageData.Action = "send_msg"
	switch messageType {
	case "private":
		{
			messageData.Params.UserId = userId
			break
		}
	case "group":
		{
			messageData.Params.GroupId = groupId
			break
		}
	default:
		{
			LOG.Error("发送消息(UnsafelySendMsg)时，消息类型错误: %v", messageType)
		}
	}
	messageData.Params.Message = message
	messageData.Params.AutoEscape = autoEscape
	// 发送消息
	_, err := wsAPI(messageData)
	if err != nil {
		LOG.Error("发送消息时，发送失败: %v", err)
		return
	}
	LOG.Info("发送消息(UnsafelySendMsg)(至：%v-%v:%v):%v", messageType, groupId, userId, message)
	return
}

// UnsafelySendPrivateMsg 发送私聊消息
//
// 注意：该API不安全，除非你知道你在做什么。
//
// 参数：
//
//   - userId: 用户ID
//
//   - message: 要发送的消息内容，类型为字符串
//
//   - autoEscape: 是否自动转义（解析消息内容中的CQ码），可选值为 true 和 false
func (p *protocolAPI) UnsafelySendPrivateMsg(userId int64, message string, autoEscape bool) {
	// 构建发送消息的JSON数据
	var messageData wba.APIRequestInfo
	messageData.Action = "send_private_msg"
	messageData.Params.UserId = userId
	messageData.Params.Message = message
	messageData.Params.AutoEscape = autoEscape
	// 发送消息
	_, err := wsAPI(messageData)
	if err != nil {
		LOG.Error("发送私聊消息(UnsafelySendPrivateMsg)时，发送失败: %v", err)
		return
	}
	LOG.Info("发送私聊消息(UnsafelySendPrivateMsg)(至：%v):%v", userId, message)
	return
}

// UnsafelySendGroupMsg 发送群消息
//
// 注意：该API不安全，除非你知道你在做什么。
//
// 参数：
//
//   - groupId: 群号
//
//   - message: 要发送的消息内容，类型为字符串
//
//   - autoEscape: 是否自动转义（解析消息内容中的CQ码），可选值为 true 和 false
func (p *protocolAPI) UnsafelySendGroupMsg(groupId int64, message string, autoEscape bool) {
	// 构建发送消息的JSON数据
	var messageData wba.APIRequestInfo
	messageData.Action = "send_group_msg"
	messageData.Params.GroupId = groupId
	messageData.Params.Message = message
	messageData.Params.AutoEscape = autoEscape
	// 发送消息
	_, err := wsAPI(messageData)
	if err != nil {
		LOG.Error("发送群消息(UnsafelySendGroupMsg)时，发送失败: %v", err)
		return
	}
	LOG.Info("发送群消息(UnsafelySendGroupMsg)(至：%v):%v", groupId, message)
	return
}

// SendMsg 回复消息(自动判断消息类型)
//
// 参数：
//
//   - msg: 消息事件信息
//
//   - message: 要发送的消息内容，类型为字符串
//
//   - autoEscape: 是否自动转义（解析消息内容中的CQ码），可选值为 true 和 false
func (p *protocolAPI) SendMsg(msg wba.MessageEventInfo, message string, autoEscape bool) {
	// 构建发送消息的JSON数据
	var messageData wba.APIRequestInfo

	messageType := msg.MessageType

	messageData.Action = "send_msg"
	switch messageType {
	case "private":
		{
			messageData.Params.UserId = msg.UserId
			break
		}
	case "group":
		{
			messageData.Params.GroupId = msg.GroupId
			break
		}
	default:
		{
			LOG.Error("回复消息(SendMsg)时，消息类型错误: %v", messageType)
		}
	}
	messageData.Params.Message = message
	messageData.Params.AutoEscape = autoEscape
	// 发送消息
	_, err := wsAPI(messageData)
	if err != nil {
		LOG.Error("回复消息时，发送失败: %v", err)
		return
	}
	LOG.Info("回复消息(SendMsg)(至：%v-%v:%v-%v):%v", msg.MessageType, msg.GroupId, msg.UserId, msg.Sender.Nickname, message)
	return
}

// SendPrivateMsg 回复私聊消息
//
// 参数：
//
//   - msg: 原始消息事件信息
//
//   - message: 要发送的消息内容，类型为字符串
//
//   - autoEscape: 是否自动转义，可选值为 true 和 false
func (p *protocolAPI) SendPrivateMsg(msg wba.MessageEventInfo, message string, autoEscape bool) {
	// 构建发送消息的JSON数据
	var messageData wba.APIRequestInfo
	messageData.Action = "send_private_msg"
	messageData.Params.UserId = msg.UserId
	messageData.Params.Message = message
	messageData.Params.AutoEscape = autoEscape
	// 发送消息
	_, err := wsAPI(messageData)
	if err != nil {
		LOG.Error("回复消息(SendPrivateMsg)时，发送失败: %v", err)
		return
	}
	LOG.Info("回复消息(SendPrivateMsg)(至：%v-%v:%v-%v):%v", msg.MessageType, msg.GroupId, msg.UserId, msg.Sender.Nickname, message)
	return
}

// SendGroupMsg 回复群消息
//
// 参数：
//
//   - msg: 原始消息事件信息
//
//   - message: 要发送的消息内容，类型为字符串
//
//   - autoEscape: 是否自动转义，可选值为 true 和 false
func (p *protocolAPI) SendGroupMsg(msg wba.MessageEventInfo, message string, autoEscape bool) {
	// 构建发送消息的JSON数据
	var messageData wba.APIRequestInfo
	messageData.Action = "send_group_msg"
	messageData.Params.GroupId = msg.GroupId
	messageData.Params.Message = message
	messageData.Params.AutoEscape = autoEscape
	// 发送消息
	_, err := wsAPI(messageData)
	if err != nil {
		LOG.Error("回复消息(SendGroupMsg)时，发送失败: %v", err)
		return
	}
	LOG.Info("回复消息(SendGroupMsg)(至：%v-%v:%v-%v):%v", msg.MessageType, msg.GroupId, msg.UserId, msg.Sender.Nickname, message)
	return
}

// UnsafelyDeleteMsg 撤回消息
//
// 注意：该API不安全，除非你知道你在做什么。
//
// 参数：
//
//   - messageId: 要撤回的消息ID
func (p *protocolAPI) UnsafelyDeleteMsg(messageId int32) {
	// 构建删除消息的JSON数据
	var messageData wba.APIRequestInfo
	messageData.Action = "delete_msg"
	messageData.Params.MessageId = messageId
	_, err := wsAPI(messageData)
	if err != nil {
		LOG.Error("撤回消息(UnsafeDeleteMsg)时，发送失败: %v", err)
		return
	}
	LOG.Info("撤回消息(UnsafeDeleteMsg):[id:%v]", messageId)
	return
}

// DeleteMsg 撤回消息
//
// 参数：
//
//   - msg: 原始消息事件信息
func (p *protocolAPI) DeleteMsg(msg wba.MessageEventInfo) {
	// 构建删除消息的JSON数据
	var messageData wba.APIRequestInfo
	messageData.Action = "delete_msg"
	messageData.Params.MessageId = msg.MessageId
	_, err := wsAPI(messageData)
	if err != nil {
		LOG.Error("撤回消息(DeleteMsg)时，发送失败: %v", err)
		return
	}
	LOG.Info("撤回消息(DeleteMsg):[id:%v]%v", msg.MessageId, msg.RawMessage)
	return
}

// SendLike 发送赞
//
// 参数：
//
//   - userId: 要赞的用户ID
//
//   - times: 赞的次数
func (p *protocolAPI) SendLike(userId int64, times int) {
	// 构建发送赞的JSON数据
	var messageData wba.APIRequestInfo
	messageData.Action = "send_like"
	messageData.Params.UserId = userId
	messageData.Params.Times = times
	_, err := wsAPI(messageData)
	if err != nil {
		LOG.Error("发送赞(SendLike)时，发送失败: %v", err)
		return
	}
	LOG.Info("发送赞(SendLike)(至：%v):%v", userId, times)
	return
}

// SetGroupKick 将指定用户移出群聊(需要群主或管理员权限)
//
// 参数：
//
//   - groupId: 群号
//
//   - userId: 用户ID
//
//   - rejectAddRequest: 是否拒绝该用户的后续加群请求，可选值为 true 和 false
func (p *protocolAPI) SetGroupKick(groupId int64, userId int64, rejectAddRequest bool) {
	var messageData wba.APIRequestInfo
	messageData.Action = "set_group_kick"
	messageData.Params.GroupId = groupId
	messageData.Params.UserId = userId
	messageData.Params.RejectAddRequest = rejectAddRequest
	_, err := wsAPI(messageData)
	if err != nil {
		LOG.Error("移出群聊(SetGroupKick)时，发送失败: %v", err)
		return
	}
	LOG.Info("移出群聊(SetGroupKick)(从：%v-%v):%v", groupId, userId, rejectAddRequest)
	return
}

// SetGroupBan 将指定用户禁言(需要群主或管理员权限)
//
// 参数：
//
//   - groupId: 群号
//
//   - userId: 用户ID
//
//   - duration: 禁言时长，单位为秒，0表示取消禁言
func (p *protocolAPI) SetGroupBan(groupId int64, userId int64, duration int32) {
	var messageData wba.APIRequestInfo
	messageData.Action = "set_group_ban"
	messageData.Params.GroupId = groupId
	messageData.Params.UserId = userId
	messageData.Params.Duration = duration
	_, err := wsAPI(messageData)
	if err != nil {
		LOG.Error("禁言群成员(SetGroupBan)时，执行失败: %v", err)
		return
	}
	LOG.Info("禁言群成员(SetGroupBan)(在：%v-%v):%v", groupId, userId, duration)
	return
}

// SetGroupWholeBan 设置全员禁言(需要群主或管理员权限)
//
// 参数：
//
//   - groupId: 群号
//
//   - enable: 是否启用全员禁言，可选值为 true 和 false
func (p *protocolAPI) SetGroupWholeBan(groupId int64, enable bool) {
	var messageData wba.APIRequestInfo
	messageData.Action = "set_group_whole_ban"
	messageData.Params.GroupId = groupId
	messageData.Params.Enable = enable
	_, err := wsAPI(messageData)
	if err != nil {
		LOG.Error("设置全员禁言(SetGroupWholeBan)时，执行失败: %v", err)
		return
	}
	LOG.Info("设置全员禁言(SetGroupWholeBan)(在：%v):%v", groupId, enable)
	return
}

// SetGroupAdmin 设置群管理员(需要群主权限)
//
// 参数：
//
//   - groupId: 群号
//
//   - userId: 用户ID
//
//   - enable: 是否设置为管理员，可选值为 true 和 false
func (p *protocolAPI) SetGroupAdmin(groupId int64, userId int64, enable bool) {
	var messageData wba.APIRequestInfo
	messageData.Action = "set_group_admin"
	messageData.Params.GroupId = groupId
	messageData.Params.UserId = userId
	messageData.Params.Enable = enable
	_, err := wsAPI(messageData)
	if err != nil {
		LOG.Error("设置群管理员(SetGroupAdmin)时，执行失败: %v", err)
		return
	}
	LOG.Info("设置群管理员(SetGroupAdmin)(在：%v-%v):%v", groupId, userId, enable)
	return
}

// SetGroupCard 设置群名片（可能需要群主或管理员权限）
//
// 参数：
//
//   - groupId: 群号
//
//   - userId: 用户ID
//
//   - card: 新的群名片
func (p *protocolAPI) SetGroupCard(groupId int64, userId int64, card string) {
	var messageData wba.APIRequestInfo
	messageData.Action = "set_group_card"
	messageData.Params.GroupId = groupId
	messageData.Params.UserId = userId
	messageData.Params.Card = card
	_, err := wsAPI(messageData)
	if err != nil {
		LOG.Error("设置群名片(SetGroupCard)时，执行失败: %v", err)
		return
	}
	LOG.Info("设置群名片(SetGroupCard)(在：%v-%v):%v", groupId, userId, card)
	return
}

// SetGroupName 设置群名称(可能需要群主或管理员权限)
//
// 参数：
//
//   - groupId: 群号
//
//   - groupName: 新的群名称
func (p *protocolAPI) SetGroupName(groupId int64, groupName string) {
	var messageData wba.APIRequestInfo
	messageData.Action = "set_group_name"
	messageData.Params.GroupId = groupId
	messageData.Params.GroupName = groupName
	_, err := wsAPI(messageData)
	if err != nil {
		LOG.Error("设置群名称(SetGroupName)时，执行失败: %v", err)
		return
	}
	LOG.Info("设置群名称(SetGroupName)(在：%v):%v", groupId, groupName)
	return
}

// SetGroupLeave 退出群聊
//
// 参数：
//
//   - groupId: 群号
//
//   - isDismiss: 是否解散群聊，仅当退出的是群主时有效，可选值为 true 和 false
func (p *protocolAPI) SetGroupLeave(groupId int64, isDismiss bool) {
	var messageData wba.APIRequestInfo
	messageData.Action = "set_group_leave"
	messageData.Params.GroupId = groupId
	messageData.Params.IsDismiss = isDismiss
	_, err := wsAPI(messageData)
	if err != nil {
		LOG.Error("退出群聊(SetGroupLeave)时，执行失败: %v", err)
		return
	}
	LOG.Info("退出群聊(SetGroupLeave)(在：%v):%v", groupId, isDismiss)
	return
}

// SetGroupSpecialTitle 设置群专属头衔(需要群主权限)
//
// 参数：
//
//   - groupId: 群号
//
//   - userId: 用户ID
//
//   - specialTitle: 新的专属头衔
func (p *protocolAPI) SetGroupSpecialTitle(groupId int64, userId int64, specialTitle string) {
	var messageData wba.APIRequestInfo
	messageData.Action = "set_group_special_title"
	messageData.Params.GroupId = groupId
	messageData.Params.UserId = userId
	messageData.Params.SpecialTitle = specialTitle
	messageData.Params.Duration = -1
	_, err := wsAPI(messageData)
	if err != nil {
		LOG.Error("设置群特殊头衔(SetGroupSpecialTitle)时，执行失败: %v", err)
		return
	}
	LOG.Info("设置群特殊头衔(SetGroupSpecialTitle)(在：%v-%v):%v-%v", groupId, userId, specialTitle)
	return
}

// SetFriendAddRequest 处理加好友请求
//
// 参数：
//
//   - flag: 请求标识，由上报的事件中获得
//
//   - approve: 是否同意请求，可选值为 true 和 false
//
//   - remark: 设置好友的备注信息
func (p *protocolAPI) SetFriendAddRequest(flag string, approve bool, remark string) {
	var messageData wba.APIRequestInfo
	messageData.Action = "set_friend_add_request"
	messageData.Params.Flag = flag
	messageData.Params.Approve = approve
	messageData.Params.Remark = remark
	_, err := wsAPI(messageData)
	if err != nil {
		LOG.Error("处理加好友请求(SetFriendAddRequest)时，执行失败: %v", err)
		return
	}
	LOG.Info("处理加好友请求(SetFriendAddRequest)(在：%v):%v-%v-%v", flag, approve, remark)
	return
}

// SetGroupAddRequest 处理加群请求/邀请
//
// 参数：
//
//   - flag: 请求标识，由上报的事件中获得
//
//   - subType: 子类型，可能是 "invite" 或 "add", 由上报的事件中获得
//
//   - approve: 是否同意请求，可选值为 true 和 false
//
//   - reason: 拒绝请求的原因，仅当 approve 为 false 时有效
func (p *protocolAPI) SetGroupAddRequest(flag string, subType string, approve bool, reason string) {
	var messageData wba.APIRequestInfo
	messageData.Action = "set_group_add_request"
	messageData.Params.Flag = flag
	messageData.Params.SubType = subType
	messageData.Params.Approve = approve
	messageData.Params.Reason = reason
	_, err := wsAPI(messageData)
	if err != nil {
		LOG.Error("处理加群请求/邀请(SetGroupAddRequest)时，执行失败: %v", err)
		return
	}
	LOG.Info("处理加群请求/邀请(SetGroupAddRequest)(在：%v-%v-%v):%v", flag, subType, approve, reason)
	return
}

// SetRestart 重启
func (p *protocolAPI) SetRestart(delay int32) {
	var messageData wba.APIRequestInfo
	messageData.Action = "set_restart"
	messageData.Params.Delay = delay
	_, err := wsAPI(messageData)
	if err != nil {
		LOG.Error("设置重启(SetRestart)时，执行失败: %v", err)
		return
	}
	LOG.Info("设置重启(SetRestart):%v", delay)
	return
}

// CleanCache 清理缓存
func (p *protocolAPI) CleanCache() {
	var messageData wba.APIRequestInfo
	messageData.Action = "clean_cache"
	_, err := wsAPI(messageData)
	if err != nil {
		LOG.Error("清理缓存(CleanCache)时，执行失败: %v", err)
		return
	}
	LOG.Info("清理缓存(CleanCache)")
	return
}

// 2.有响应API，需添加echo字段，统一返回响应结构体

// GetLoginInfo 获取登录信息
func (p *protocolAPI) GetLoginInfo() (Response wba.APIResponseInfo) {
	LOG.Info("获取登录信息(GetLoginInfo)")
	var messageData wba.APIRequestInfo
	var err error
	messageData.Action = "get_login_info"
	messageData.Echo, err = GenerateUUID()
	if err != nil {
		LOG.Error("获取登录信息(GetLoginInfo)时，生成UUID失败: %v", err)
		return wba.APIResponseInfo{}
	}
	Response, err = wsAPI(messageData)
	if err != nil {
		LOG.Error("获取登录信息(GetLoginInfo)时，执行失败: %v", err)
		return wba.APIResponseInfo{}
	}
	return Response
}

// GetVersionInfo 获取协议信息
func (p *protocolAPI) GetVersionInfo() (Response wba.APIResponseInfo) {
	LOG.Info("获取协议信息(GetVersionInfo)")
	var messageData wba.APIRequestInfo
	var err error
	messageData.Action = "get_version_info"
	messageData.Echo, err = GenerateUUID()
	if err != nil {
		LOG.Error("获取协议信息(GetVersionInfo)时，生成UUID失败: %v", err)
		return wba.APIResponseInfo{}
	}
	Response, err = wsAPI(messageData)
	if err != nil {
		LOG.Error("获取登录信息(GetVersionInfo)时，执行失败: %v", err)
		return wba.APIResponseInfo{}
	}
	return Response
}

// GetMsg 获取消息
func (p *protocolAPI) GetMsg(messageId int32) (Response wba.APIResponseInfo) {
	LOG.Info("获取消息(GetMsg)")
	var messageData wba.APIRequestInfo
	var err error
	messageData.Action = "get_msg"
	messageData.Params.MessageId = messageId
	messageData.Echo, err = GenerateUUID()
	if err != nil {
		LOG.Error("获取消息(GetMsg)时，生成UUID失败: %v", err)
		return wba.APIResponseInfo{}
	}
	Response, err = wsAPI(messageData)
	if err != nil {
		LOG.Error("获取消息(GetMsg)时，执行失败: %v", err)
		return wba.APIResponseInfo{}
	}
	return Response
}

// GetForwardMsg 获取合并转发消息
func (p *protocolAPI) GetForwardMsg(id string) (Response wba.APIResponseInfo) {
	LOG.Info("获取合并转发消息(GetForwardMsg)")
	var messageData wba.APIRequestInfo
	var err error
	messageData.Action = "get_forward_msg"
	messageData.Params.Id = id
	messageData.Echo, err = GenerateUUID()
	if err != nil {
		LOG.Error("获取合并转发消息(GetForwardMsg)时，生成UUID失败: %v", err)
		return wba.APIResponseInfo{}
	}
	Response, err = wsAPI(messageData)
	if err != nil {
		LOG.Error("获取合并转发消息(GetForwardMsg)时，执行失败: %v", err)
		return wba.APIResponseInfo{}
	}
	return Response
}

// GetStrangerInfo 获取陌生人信息
func (p *protocolAPI) GetStrangerInfo(userId int64, noCache bool) (Response wba.APIResponseInfo) {
	LOG.Info("获取陌生人信息(GetStrangerInfo)")
	var messageData wba.APIRequestInfo
	var err error
	messageData.Action = "get_stranger_info"
	messageData.Params.UserId = userId
	messageData.Params.NoCache = noCache
	messageData.Echo, err = GenerateUUID()
	if err != nil {
		LOG.Error("获取陌生人信息(GetStrangerInfo)时，生成UUID失败: %v", err)
		return wba.APIResponseInfo{}
	}
	Response, err = wsAPI(messageData)
	if err != nil {
		LOG.Error("获取陌生人信息(GetStrangerInfo)时，执行失败: %v", err)
		return wba.APIResponseInfo{}
	}
	return Response
}

// GetFriendList 获取好友列表
func (p *protocolAPI) GetFriendList() (Response wba.APIResponseInfo) {
	LOG.Info("获取好友列表(GetFriendList)")
	var messageData wba.APIRequestInfo
	var err error
	messageData.Action = "get_friend_list"
	messageData.Echo, err = GenerateUUID()
	if err != nil {
		LOG.Error("获取好友列表(GetFriendList)时，生成UUID失败: %v", err)
		return wba.APIResponseInfo{}
	}
	Response, err = wsAPI(messageData)
	if err != nil {
		LOG.Error("获取好友列表(GetFriendList)时，执行失败: %v", err)
		return wba.APIResponseInfo{}
	}
	return Response
}

// GetGroupList 获取群列表
func (p *protocolAPI) GetGroupList() (Response wba.APIResponseInfo) {
	LOG.Info("获取群列表(GetGroupList)")
	var messageData wba.APIRequestInfo
	var err error
	messageData.Action = "get_group_list"
	messageData.Echo, err = GenerateUUID()
	if err != nil {
		LOG.Error("获取群列表(GetGroupList)时，生成UUID失败: %v", err)
		return wba.APIResponseInfo{}
	}
	Response, err = wsAPI(messageData)
	if err != nil {
		LOG.Error("获取群列表(GetGroupList)时，执行失败: %v", err)
		return wba.APIResponseInfo{}
	}
	return Response
}

// GetGroupInfo 获取群信息
func (p *protocolAPI) GetGroupInfo(groupId int64, noCache bool) (Response wba.APIResponseInfo) {
	LOG.Info("获取群信息(GetGroupInfo)")
	var messageData wba.APIRequestInfo
	var err error
	messageData.Action = "get_group_info"
	messageData.Params.GroupId = groupId
	messageData.Params.NoCache = noCache
	messageData.Echo, err = GenerateUUID()
	if err != nil {
		LOG.Error("获取群信息(GetGroupInfo)时，生成UUID失败: %v", err)
		return wba.APIResponseInfo{}
	}
	Response, err = wsAPI(messageData)
	if err != nil {
		LOG.Error("获取群信息(GetGroupInfo)时，执行失败: %v", err)
		return wba.APIResponseInfo{}
	}
	return Response
}

// GetGroupMemberInfo 获取群成员信息
func (p *protocolAPI) GetGroupMemberInfo(groupId int64, userId int64, noCache bool) (Response wba.APIResponseInfo) {
	LOG.Info("获取群成员信息(GetGroupMemberInfo)")
	var messageData wba.APIRequestInfo
	var err error
	messageData.Action = "get_group_member_info"
	messageData.Params.GroupId = groupId
	messageData.Params.UserId = userId
	messageData.Params.NoCache = noCache
	messageData.Echo, err = GenerateUUID()
	if err != nil {
		LOG.Error("获取群成员信息(GetGroupMemberInfo)时，生成UUID失败: %v", err)
		return wba.APIResponseInfo{}
	}
	Response, err = wsAPI(messageData)
	if err != nil {
		LOG.Error("获取群成员信息(GetGroupMemberInfo)时，执行失败: %v", err)
		return wba.APIResponseInfo{}
	}
	return Response
}

// GetGroupMemberList 获取群成员列表
func (p *protocolAPI) GetGroupMemberList(groupId int64) (Response wba.APIResponseInfo) {
	LOG.Info("获取群成员列表(GetGroupMemberList)")
	var messageData wba.APIRequestInfo
	var err error
	messageData.Action = "get_group_member_list"
	messageData.Params.GroupId = groupId
	messageData.Echo, err = GenerateUUID()
	if err != nil {
		LOG.Error("获取群成员列表(GetGroupMemberList)时，生成UUID失败: %v", err)
		return wba.APIResponseInfo{}
	}
	Response, err = wsAPI(messageData)
	if err != nil {
		LOG.Error("获取群成员列表(GetGroupMemberList)时，执行失败: %v", err)
		return wba.APIResponseInfo{}
	}
	return Response
}

// GetGroupHonorInfo 获取群荣誉信息
func (p *protocolAPI) GetGroupHonorInfo(groupId int64, Type string) (Response wba.APIResponseInfo) {
	LOG.Info("获取群荣誉信息(GetGroupHonorInfo)")
	var messageData wba.APIRequestInfo
	var err error
	messageData.Action = "get_group_honor_info"
	messageData.Params.GroupId = groupId
	messageData.Params.Type = Type
	messageData.Echo, err = GenerateUUID()
	if err != nil {
		LOG.Error("获取群荣誉信息(GetGroupHonorInfo)时，生成UUID失败: %v", err)
		return wba.APIResponseInfo{}
	}
	Response, err = wsAPI(messageData)
	if err != nil {
		LOG.Error("获取群荣誉信息(GetGroupHonorInfo)时，执行失败: %v", err)
		return wba.APIResponseInfo{}
	}
	return Response
}

// GetCookies 获取Cookies
func (p *protocolAPI) GetCookies(domain string) (Response wba.APIResponseInfo) {
	LOG.Info("获取Cookies(GetCookies)")
	var messageData wba.APIRequestInfo
	var err error
	messageData.Action = "get_cookies"
	messageData.Params.Domain = domain
	messageData.Echo, err = GenerateUUID()
	if err != nil {
		LOG.Error("获取Cookies(GetCookies)时，生成UUID失败: %v", err)
		return wba.APIResponseInfo{}
	}
	Response, err = wsAPI(messageData)
	if err != nil {
		LOG.Error("获取Cookies(GetCookies)时，执行失败: %v", err)
		return wba.APIResponseInfo{}
	}
	return Response
}

// GetCSRFToken 获取CSRF Token
func (p *protocolAPI) GetCSRFToken() (Response wba.APIResponseInfo) {
	LOG.Info("获取CSRF Token(GetCSRFToken)")
	var messageData wba.APIRequestInfo
	var err error
	messageData.Action = "get_csrf_token"
	messageData.Echo, err = GenerateUUID()
	if err != nil {
		LOG.Error("获取CSRF Token(GetCSRFToken)时，生成UUID失败: %v", err)
		return wba.APIResponseInfo{}
	}
	Response, err = wsAPI(messageData)
	if err != nil {
		LOG.Error("获取CSRF Token(GetCSRFToken)时，执行失败: %v", err)
		return wba.APIResponseInfo{}
	}
	return Response
}

// GetCredentials 获取登录令牌
func (p *protocolAPI) GetCredentials(domain string) (Response wba.APIResponseInfo) {
	LOG.Info("获取登录令牌(GetCredentials)")
	var messageData wba.APIRequestInfo
	var err error
	messageData.Action = "get_credentials"
	messageData.Params.Domain = domain
	messageData.Echo, err = GenerateUUID()
	if err != nil {
		LOG.Error("获取登录令牌(GetCredentials)时，生成UUID失败: %v", err)
		return wba.APIResponseInfo{}
	}
	Response, err = wsAPI(messageData)
	if err != nil {
		LOG.Error("获取登录令牌(GetCredentials)时，执行失败: %v", err)
		return wba.APIResponseInfo{}
	}
	return Response
}

// GetRecord 获取语音
func (p *protocolAPI) GetRecord(file string, outFormat string) (Response wba.APIResponseInfo) {
	LOG.Info("获取语音(GetRecord)")
	var messageData wba.APIRequestInfo
	var err error
	messageData.Action = "get_record"
	messageData.Params.File = file
	messageData.Params.OutFormat = outFormat
	messageData.Echo, err = GenerateUUID()
	if err != nil {
		LOG.Error("获取语音(GetRecord)时，生成UUID失败: %v", err)
		return wba.APIResponseInfo{}
	}
	Response, err = wsAPI(messageData)
	if err != nil {
		LOG.Error("获取语音(GetRecord)时，执行失败: %v", err)
		return wba.APIResponseInfo{}
	}
	return Response
}

// GetImage 获取图片
func (p *protocolAPI) GetImage(file string) (Response wba.APIResponseInfo) {
	LOG.Info("获取图片(GetImage)")
	var messageData wba.APIRequestInfo
	var err error
	messageData.Action = "get_image"
	messageData.Params.File = file
	messageData.Echo, err = GenerateUUID()
	if err != nil {
		LOG.Error("获取图片(GetImage)时，生成UUID失败: %v", err)
		return wba.APIResponseInfo{}
	}
	Response, err = wsAPI(messageData)
	if err != nil {
		LOG.Error("获取图片(GetImage)时，执行失败: %v", err)
		return wba.APIResponseInfo{}
	}
	return Response
}

// CanSendImage 检查是否可以发送图片
func (p *protocolAPI) CanSendImage() (Response wba.APIResponseInfo) {
	LOG.Info("检查是否可以发送图片(CanSendImage)")
	var messageData wba.APIRequestInfo
	var err error
	messageData.Action = "can_send_image"
	messageData.Echo, err = GenerateUUID()
	if err != nil {
		LOG.Error("检查是否可以发送图片(CanSendImage)时，生成UUID失败: %v", err)
		return wba.APIResponseInfo{}
	}
	Response, err = wsAPI(messageData)
	if err != nil {
		LOG.Error("检查是否可以发送图片(CanSendImage)时，执行失败: %v", err)
		return wba.APIResponseInfo{}
	}
	return Response
}

// CanSendRecord 检查是否可以发送语音
func (p *protocolAPI) CanSendRecord() (Response wba.APIResponseInfo) {
	LOG.Info("检查是否可以发送语音(CanSendRecord)")
	var messageData wba.APIRequestInfo
	var err error
	messageData.Action = "can_send_record"
	messageData.Echo, err = GenerateUUID()
	if err != nil {
		LOG.Error("检查是否可以发送语音(CanSendRecord)时，生成UUID失败: %v", err)
		return wba.APIResponseInfo{}
	}
	Response, err = wsAPI(messageData)
	if err != nil {
		LOG.Error("检查是否可以发送语音(CanSendRecord)时，执行失败: %v", err)
		return wba.APIResponseInfo{}
	}
	return Response
}

// GetStatus 获取状态
func (p *protocolAPI) GetStatus() (Response wba.APIResponseInfo) {
	LOG.Info("获取状态(GetStatus)")
	var messageData wba.APIRequestInfo
	var err error
	messageData.Action = "get_status"
	messageData.Echo, err = GenerateUUID()
	if err != nil {
		LOG.Error("获取状态(GetStatus)时，生成UUID失败: %v", err)
		return wba.APIResponseInfo{}
	}
	Response, err = wsAPI(messageData)
	if err != nil {
		LOG.Error("获取状态(GetStatus)时，执行失败: %v", err)
		return wba.APIResponseInfo{}
	}
	return Response
}

// 文件管理模块
//TODO: 文件管理模块待实现

//终端连接模块
//TODO: 终端模块待实现

//核心信息调用模块

var ProtocolApi protocolAPI

func GenerateUUID() (string, error) {
	uuid := make([]byte, 16)
	_, err := rand.Read(uuid)
	if err != nil {
		return "", err
	}

	// 设置UUID版本号（版本4），将第6字节的高4位设置为0100
	uuid[6] = (uuid[6] & 0x0F) | 0x40
	// 设置UUID变体（RFC 4122规范定义的变体），将第8字节的高4位设置为10
	uuid[8] = (uuid[8] & 0x3F) | 0x80

	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}
