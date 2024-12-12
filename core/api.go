package core

import (
	"ProjectWIND/LOG"
	"ProjectWIND/wba"
	"crypto/rand"
	"encoding/json"
	"fmt"
)

type apiInfo struct{}

/*
关于API的说明：

1.所有API请求按照OneBot11标准，使用JSON格式进行数据交换。api命名为由原文档中蛇形命名法改为双驼峰命名法。

2.无响应的API请求使用ws协议处理，有响应的API需添加echo字段。

3.wind会从配置文件中读取API请求的url，请确保正确填写。
*/

//1.无响应API,使用ws协议处理

// SendMsg 发送消息(自动判断消息类型)
func (a *apiInfo) SendMsg(msg wba.MessageEventInfo, message string, autoEscape bool) {
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
			LOG.ERROR("发送消息(SendMsg)时，消息类型错误: %v", messageType)
		}
	}
	messageData.Params.Message = message
	messageData.Params.AutoEscape = autoEscape
	messageJson, err := json.Marshal(messageData)
	if err != nil {
		LOG.ERROR("发送消息时，构建JSON数据失败: %v", err)
		return
	}
	// 发送消息
	err = wsAPI(messageJson)
	if err != nil {
		LOG.ERROR("发送消息时，发送失败: %v", err)
		return
	}
	LOG.INFO("发送消息(SendMsg)(至：%v-%v:%v-%v):%v", msg.MessageType, msg.GroupId, msg.UserId, msg.Sender.Nickname, message)
	return
}

// SendPrivateMsg 发送私聊消息
func (a *apiInfo) SendPrivateMsg(msg wba.MessageEventInfo, message string, autoEscape bool) {
	// 构建发送消息的JSON数据
	var messageData wba.APIRequestInfo
	messageData.Action = "send_private_msg"
	messageData.Params.UserId = msg.UserId
	messageData.Params.Message = message
	messageData.Params.AutoEscape = autoEscape
	messageJson, err := json.Marshal(messageData)
	if err != nil {
		LOG.ERROR("发送消息(SendPrivateMsg)时，构建JSON数据失败: %v", err)
		return
	}
	// 发送消息
	err = wsAPI(messageJson)
	if err != nil {
		LOG.ERROR("发送消息(SendPrivateMsg)时，发送失败: %v", err)
		return
	}
	LOG.INFO("发送消息(SendPrivateMsg)(至：%v-%v:%v-%v):%v", msg.MessageType, msg.GroupId, msg.UserId, msg.Sender.Nickname, message)
	return
}

// SendGroupMsg 发送群消息
func (a *apiInfo) SendGroupMsg(msg wba.MessageEventInfo, message string, autoEscape bool) {
	// 构建发送消息的JSON数据
	var messageData wba.APIRequestInfo
	messageData.Action = "send_group_msg"
	messageData.Params.GroupId = msg.GroupId
	messageData.Params.Message = message
	messageData.Params.AutoEscape = autoEscape
	messageJson, err := json.Marshal(messageData)
	if err != nil {
		LOG.ERROR("发送消息(SendGroupMsg)时，构建JSON数据失败: %v", err)
		return
	}
	// 发送消息
	err = wsAPI(messageJson)
	if err != nil {
		LOG.ERROR("发送消息(SendGroupMsg)时，发送失败: %v", err)
		return
	}
	LOG.INFO("发送消息(SendGroupMsg)(至：%v-%v:%v-%v):%v", msg.MessageType, msg.GroupId, msg.UserId, msg.Sender.Nickname, message)
	return
}

// DeleteMsg 撤回消息
func (a *apiInfo) DeleteMsg(msg wba.MessageEventInfo) {
	// 构建删除消息的JSON数据
	var messageData wba.APIRequestInfo
	messageData.Action = "delete_msg"
	messageData.Params.MessageId = msg.MessageId
	messageJson, err := json.Marshal(messageData)
	if err != nil {
		LOG.ERROR("撤回消息(DeleteMsg)时，构建JSON数据失败: %v", err)
		return
	}
	err = wsAPI(messageJson)
	if err != nil {
		LOG.ERROR("撤回消息(DeleteMsg)时，发送失败: %v", err)
		return
	}
	LOG.INFO("撤回消息(DeleteMsg):[id:%v]%v", msg.MessageId, msg.RawMessage)
	return
}

// SendLike 发送赞
func (a *apiInfo) SendLike(userId int64, times int) {
	// 构建发送赞的JSON数据
	var messageData wba.APIRequestInfo
	messageData.Action = "send_like"
	messageData.Params.UserId = userId
	messageData.Params.Times = times
	messageJson, err := json.Marshal(messageData)
	if err != nil {
		LOG.ERROR("发送赞(SendLike)时，构建JSON数据失败: %v", err)
		return
	}
	err = wsAPI(messageJson)
	if err != nil {
		LOG.ERROR("发送赞(SendLike)时，发送失败: %v", err)
		return
	}
	LOG.INFO("发送赞(SendLike)(至：%v):%v", userId, times)
	return
}

// SetGroupKick 将指定用户移出群聊(需要群主或管理员权限)
func (a *apiInfo) SetGroupKick(groupId int64, userId int64, rejectAddRequest bool) {
	var messageData wba.APIRequestInfo
	messageData.Action = "set_group_kick"
	messageData.Params.GroupId = groupId
	messageData.Params.UserId = userId
	messageData.Params.RejectAddRequest = rejectAddRequest
	messageJson, err := json.Marshal(messageData)
	if err != nil {
		LOG.ERROR("移出群聊(SetGroupKick)时，构建JSON数据失败: %v", err)
		return
	}
	err = wsAPI(messageJson)
	if err != nil {
		LOG.ERROR("移出群聊(SetGroupKick)时，发送失败: %v", err)
		return
	}
	LOG.INFO("移出群聊(SetGroupKick)(从：%v-%v):%v", groupId, userId, rejectAddRequest)
	return
}

// SetGroupBan 将指定用户禁言(需要群主或管理员权限)
func (a *apiInfo) SetGroupBan(groupId int64, userId int64, duration int32) {
	var messageData wba.APIRequestInfo
	messageData.Action = "set_group_ban"
	messageData.Params.GroupId = groupId
	messageData.Params.UserId = userId
	messageData.Params.Duration = duration
	messageJson, err := json.Marshal(messageData)
	if err != nil {
		LOG.ERROR("禁言群成员(SetGroupBan)时，构建JSON数据失败: %v", err)
		return
	}
	err = wsAPI(messageJson)
	if err != nil {
		LOG.ERROR("禁言群成员(SetGroupBan)时，执行失败: %v", err)
		return
	}
	LOG.INFO("禁言群成员(SetGroupBan)(在：%v-%v):%v", groupId, userId, duration)
	return
}

// SetGroupWholeBan 设置全员禁言(需要群主或管理员权限)
func (a *apiInfo) SetGroupWholeBan(groupId int64, enable bool) {
	var messageData wba.APIRequestInfo
	messageData.Action = "set_group_whole_ban"
	messageData.Params.GroupId = groupId
	messageData.Params.Enable = enable
	messageJson, err := json.Marshal(messageData)
	if err != nil {
		LOG.ERROR("设置全员禁言(SetGroupWholeBan)时，构建JSON数据失败: %v", err)
		return
	}
	err = wsAPI(messageJson)
	if err != nil {
		LOG.ERROR("设置全员禁言(SetGroupWholeBan)时，执行失败: %v", err)
		return
	}
	LOG.INFO("设置全员禁言(SetGroupWholeBan)(在：%v):%v", groupId, enable)
	return
}

// SetGroupAdmin 设置群管理员(需要群主权限)
func (a *apiInfo) SetGroupAdmin(groupId int64, userId int64, enable bool) {
	var messageData wba.APIRequestInfo
	messageData.Action = "set_group_admin"
	messageData.Params.GroupId = groupId
	messageData.Params.UserId = userId
	messageData.Params.Enable = enable
	messageJson, err := json.Marshal(messageData)
	if err != nil {
		LOG.ERROR("设置群管理员(SetGroupAdmin)时，构建JSON数据失败: %v", err)
		return
	}
	err = wsAPI(messageJson)
	if err != nil {
		LOG.ERROR("设置群管理员(SetGroupAdmin)时，执行失败: %v", err)
		return
	}
	LOG.INFO("设置群管理员(SetGroupAdmin)(在：%v-%v):%v", groupId, userId, enable)
	return
}

// SetGroupCard 设置群名片(需要群主或管理员权限)
func (a *apiInfo) SetGroupCard(groupId int64, userId int64, card string) {
	var messageData wba.APIRequestInfo
	messageData.Action = "set_group_card"
	messageData.Params.GroupId = groupId
	messageData.Params.UserId = userId
	messageData.Params.Card = card
	messageJson, err := json.Marshal(messageData)
	if err != nil {
		LOG.ERROR("设置群名片(SetGroupCard)时，构建JSON数据失败: %v", err)
		return
	}
	err = wsAPI(messageJson)
	if err != nil {
		LOG.ERROR("设置群名片(SetGroupCard)时，执行失败: %v", err)
		return
	}
	LOG.INFO("设置群名片(SetGroupCard)(在：%v-%v):%v", groupId, userId, card)
	return
}

// SetGroupName 设置群名称(可能需要群主或管理员权限)
func (a *apiInfo) SetGroupName(groupId int64, groupName string) {
	var messageData wba.APIRequestInfo
	messageData.Action = "set_group_name"
	messageData.Params.GroupId = groupId
	messageData.Params.GroupName = groupName
	messageJson, err := json.Marshal(messageData)
	if err != nil {
		LOG.ERROR("设置群名称(SetGroupName)时，构建JSON数据失败: %v", err)
		return
	}
	err = wsAPI(messageJson)
	if err != nil {
		LOG.ERROR("设置群名称(SetGroupName)时，执行失败: %v", err)
		return
	}
	LOG.INFO("设置群名称(SetGroupName)(在：%v):%v", groupId, groupName)
	return
}

// SetGroupLeave 退出群聊
func (a *apiInfo) SetGroupLeave(groupId int64, isDismiss bool) {
	var messageData wba.APIRequestInfo
	messageData.Action = "set_group_leave"
	messageData.Params.GroupId = groupId
	messageData.Params.IsDismiss = isDismiss
	messageJson, err := json.Marshal(messageData)
	if err != nil {
		LOG.ERROR("退出群聊(SetGroupLeave)时，构建JSON数据失败: %v", err)
		return
	}
	err = wsAPI(messageJson)
	if err != nil {
		LOG.ERROR("退出群聊(SetGroupLeave)时，执行失败: %v", err)
		return
	}
	LOG.INFO("退出群聊(SetGroupLeave)(在：%v):%v", groupId, isDismiss)
	return
}

// SetGroupSpecialTitle 设置群专属头衔(需要群主权限)
func (a *apiInfo) SetGroupSpecialTitle(groupId int64, userId int64, specialTitle string, duration int32) {
	var messageData wba.APIRequestInfo
	messageData.Action = "set_group_special_title"
	messageData.Params.GroupId = groupId
	messageData.Params.UserId = userId
	messageData.Params.SpecialTitle = specialTitle
	messageData.Params.Duration = duration
	messageJson, err := json.Marshal(messageData)
	if err != nil {
		LOG.ERROR("设置群特殊头衔(SetGroupSpecialTitle)时，构建JSON数据失败: %v", err)
		return
	}
	err = wsAPI(messageJson)
	if err != nil {
		LOG.ERROR("设置群特殊头衔(SetGroupSpecialTitle)时，执行失败: %v", err)
		return
	}
	LOG.INFO("设置群特殊头衔(SetGroupSpecialTitle)(在：%v-%v):%v-%v", groupId, userId, specialTitle, duration)
	return
}

// SetFriendAddRequest 处理加好友请求
func (a *apiInfo) SetFriendAddRequest(flag string, approve bool, remark string) {
	var messageData wba.APIRequestInfo
	messageData.Action = "set_friend_add_request"
	messageData.Params.Flag = flag
	messageData.Params.Approve = approve
	messageData.Params.Remark = remark
	messageJson, err := json.Marshal(messageData)
	if err != nil {
		LOG.ERROR("处理加好友请求(SetFriendAddRequest)时，构建JSON数据失败: %v", err)
		return
	}
	err = wsAPI(messageJson)
	if err != nil {
		LOG.ERROR("处理加好友请求(SetFriendAddRequest)时，执行失败: %v", err)
		return
	}
	LOG.INFO("处理加好友请求(SetFriendAddRequest)(在：%v):%v-%v-%v", flag, approve, remark)
	return
}

// SetGroupAddRequest 处理加群请求/邀请
func (a *apiInfo) SetGroupAddRequest(flag string, subType string, approve bool, reason string) {
	var messageData wba.APIRequestInfo
	messageData.Action = "set_group_add_request"
	messageData.Params.Flag = flag
	messageData.Params.SubType = subType
	messageData.Params.Approve = approve
	messageData.Params.Reason = reason
	messageJson, err := json.Marshal(messageData)
	if err != nil {
		LOG.ERROR("处理加群请求/邀请(SetGroupAddRequest)时，构建JSON数据失败: %v", err)
		return
	}
	err = wsAPI(messageJson)
	if err != nil {
		LOG.ERROR("处理加群请求/邀请(SetGroupAddRequest)时，执行失败: %v", err)
		return
	}
	LOG.INFO("处理加群请求/邀请(SetGroupAddRequest)(在：%v-%v-%v):%v", flag, subType, approve, reason)
	return
}

// 2.有响应API，需添加echo字段

func (a *apiInfo) GetLoginInfo(flag string, approve bool) {

}

var AppApi apiInfo

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
