package core

import (
	"ProjectWIND/LOG"
	"ProjectWIND/wba"
	"crypto/rand"
	"fmt"
)

type apiInfo struct{}

//一、Protocol模块

/*
关于Protocol模块的说明

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
	// 发送消息
	_, err := wsAPI(messageData)
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
	// 发送消息
	_, err := wsAPI(messageData)
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
	// 发送消息
	_, err := wsAPI(messageData)
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
	_, err := wsAPI(messageData)
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
	_, err := wsAPI(messageData)
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
	_, err := wsAPI(messageData)
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
	_, err := wsAPI(messageData)
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
	_, err := wsAPI(messageData)
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
	_, err := wsAPI(messageData)
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
	_, err := wsAPI(messageData)
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
	_, err := wsAPI(messageData)
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
	_, err := wsAPI(messageData)
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
	_, err := wsAPI(messageData)
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
	_, err := wsAPI(messageData)
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
	_, err := wsAPI(messageData)
	if err != nil {
		LOG.ERROR("处理加群请求/邀请(SetGroupAddRequest)时，执行失败: %v", err)
		return
	}
	LOG.INFO("处理加群请求/邀请(SetGroupAddRequest)(在：%v-%v-%v):%v", flag, subType, approve, reason)
	return
}

// SetRestart 重启
func (a *apiInfo) SetRestart(delay int32) {
	var messageData wba.APIRequestInfo
	messageData.Action = "set_restart"
	messageData.Params.Delay = delay
	_, err := wsAPI(messageData)
	if err != nil {
		LOG.ERROR("设置重启(SetRestart)时，执行失败: %v", err)
		return
	}
	LOG.INFO("设置重启(SetRestart):%v", delay)
	return
}

// CleanCache 清理缓存
func (a *apiInfo) CleanCache() {
	var messageData wba.APIRequestInfo
	messageData.Action = "clean_cache"
	_, err := wsAPI(messageData)
	if err != nil {
		LOG.ERROR("清理缓存(CleanCache)时，执行失败: %v", err)
		return
	}
	LOG.INFO("清理缓存(CleanCache)")
	return
}

// 2.有响应API，需添加echo字段

// GetLoginInfo 获取登录信息
func (a *apiInfo) GetLoginInfo() (Response wba.APIResponseInfo) {
	LOG.INFO("获取登录信息(GetLoginInfo)")
	var messageData wba.APIRequestInfo
	var err error
	messageData.Action = "get_login_info"
	messageData.Echo, err = GenerateUUID()
	if err != nil {
		LOG.ERROR("获取登录信息(GetLoginInfo)时，生成UUID失败: %v", err)
		return wba.APIResponseInfo{}
	}
	Response, err = wsAPI(messageData)
	if err != nil {
		LOG.ERROR("获取登录信息(GetLoginInfo)时，执行失败: %v", err)
		return wba.APIResponseInfo{}
	}
	return Response
}

// GetVersionInfo 获取协议信息
func (a *apiInfo) GetVersionInfo() (Response wba.APIResponseInfo) {
	LOG.INFO("获取协议信息(GetVersionInfo)")
	var messageData wba.APIRequestInfo
	var err error
	messageData.Action = "get_version_info"
	messageData.Echo, err = GenerateUUID()
	if err != nil {
		LOG.ERROR("获取协议信息(GetVersionInfo)时，生成UUID失败: %v", err)
		return wba.APIResponseInfo{}
	}
	Response, err = wsAPI(messageData)
	if err != nil {
		LOG.ERROR("获取登录信息(GetVersionInfo)时，执行失败: %v", err)
		return wba.APIResponseInfo{}
	}
	return Response
}

// GetMsg 获取消息
func (a *apiInfo) GetMsg(messageId int32) (Response wba.APIResponseInfo) {
	LOG.INFO("获取消息(GetMsg)")
	var messageData wba.APIRequestInfo
	var err error
	messageData.Action = "get_msg"
	messageData.Params.MessageId = messageId
	messageData.Echo, err = GenerateUUID()
	if err != nil {
		LOG.ERROR("获取消息(GetMsg)时，生成UUID失败: %v", err)
		return wba.APIResponseInfo{}
	}
	Response, err = wsAPI(messageData)
	if err != nil {
		LOG.ERROR("获取消息(GetMsg)时，执行失败: %v", err)
		return wba.APIResponseInfo{}
	}
	return Response
}

// GetForwardMsg 获取合并转发消息
func (a *apiInfo) GetForwardMsg(id string) (Response wba.APIResponseInfo) {
	LOG.INFO("获取合并转发消息(GetForwardMsg)")
	var messageData wba.APIRequestInfo
	var err error
	messageData.Action = "get_forward_msg"
	messageData.Params.Id = id
	messageData.Echo, err = GenerateUUID()
	if err != nil {
		LOG.ERROR("获取合并转发消息(GetForwardMsg)时，生成UUID失败: %v", err)
		return wba.APIResponseInfo{}
	}
	Response, err = wsAPI(messageData)
	if err != nil {
		LOG.ERROR("获取合并转发消息(GetForwardMsg)时，执行失败: %v", err)
		return wba.APIResponseInfo{}
	}
	return Response
}

// GetStrangerInfo 获取陌生人信息
func (a *apiInfo) GetStrangerInfo(userId int64, noCache bool) (Response wba.APIResponseInfo) {
	LOG.INFO("获取陌生人信息(GetStrangerInfo)")
	var messageData wba.APIRequestInfo
	var err error
	messageData.Action = "get_stranger_info"
	messageData.Params.UserId = userId
	messageData.Params.NoCache = noCache
	messageData.Echo, err = GenerateUUID()
	if err != nil {
		LOG.ERROR("获取陌生人信息(GetStrangerInfo)时，生成UUID失败: %v", err)
		return wba.APIResponseInfo{}
	}
	Response, err = wsAPI(messageData)
	if err != nil {
		LOG.ERROR("获取陌生人信息(GetStrangerInfo)时，执行失败: %v", err)
		return wba.APIResponseInfo{}
	}
	return Response
}

// GetFriendList 获取好友列表
func (a *apiInfo) GetFriendList() (Response wba.APIResponseInfo) {
	LOG.INFO("获取好友列表(GetFriendList)")
	var messageData wba.APIRequestInfo
	var err error
	messageData.Action = "get_friend_list"
	messageData.Echo, err = GenerateUUID()
	if err != nil {
		LOG.ERROR("获取好友列表(GetFriendList)时，生成UUID失败: %v", err)
		return wba.APIResponseInfo{}
	}
	Response, err = wsAPI(messageData)
	if err != nil {
		LOG.ERROR("获取好友列表(GetFriendList)时，执行失败: %v", err)
		return wba.APIResponseInfo{}
	}
	return Response
}

// GetGroupList 获取群列表
func (a *apiInfo) GetGroupList() (Response wba.APIResponseInfo) {
	LOG.INFO("获取群列表(GetGroupList)")
	var messageData wba.APIRequestInfo
	var err error
	messageData.Action = "get_group_list"
	messageData.Echo, err = GenerateUUID()
	if err != nil {
		LOG.ERROR("获取群列表(GetGroupList)时，生成UUID失败: %v", err)
		return wba.APIResponseInfo{}
	}
	Response, err = wsAPI(messageData)
	if err != nil {
		LOG.ERROR("获取群列表(GetGroupList)时，执行失败: %v", err)
		return wba.APIResponseInfo{}
	}
	return Response
}

// GetGroupInfo 获取群信息
func (a *apiInfo) GetGroupInfo(groupId int64, noCache bool) (Response wba.APIResponseInfo) {
	LOG.INFO("获取群信息(GetGroupInfo)")
	var messageData wba.APIRequestInfo
	var err error
	messageData.Action = "get_group_info"
	messageData.Params.GroupId = groupId
	messageData.Params.NoCache = noCache
	messageData.Echo, err = GenerateUUID()
	if err != nil {
		LOG.ERROR("获取群信息(GetGroupInfo)时，生成UUID失败: %v", err)
		return wba.APIResponseInfo{}
	}
	Response, err = wsAPI(messageData)
	if err != nil {
		LOG.ERROR("获取群信息(GetGroupInfo)时，执行失败: %v", err)
		return wba.APIResponseInfo{}
	}
	return Response
}

// GetGroupMemberInfo 获取群成员信息
func (a *apiInfo) GetGroupMemberInfo(groupId int64, userId int64, noCache bool) (Response wba.APIResponseInfo) {
	LOG.INFO("获取群成员信息(GetGroupMemberInfo)")
	var messageData wba.APIRequestInfo
	var err error
	messageData.Action = "get_group_member_info"
	messageData.Params.GroupId = groupId
	messageData.Params.UserId = userId
	messageData.Params.NoCache = noCache
	messageData.Echo, err = GenerateUUID()
	if err != nil {
		LOG.ERROR("获取群成员信息(GetGroupMemberInfo)时，生成UUID失败: %v", err)
		return wba.APIResponseInfo{}
	}
	Response, err = wsAPI(messageData)
	if err != nil {
		LOG.ERROR("获取群成员信息(GetGroupMemberInfo)时，执行失败: %v", err)
		return wba.APIResponseInfo{}
	}
	return Response
}

// GetGroupMemberList 获取群成员列表
func (a *apiInfo) GetGroupMemberList(groupId int64) (Response wba.APIResponseInfo) {
	LOG.INFO("获取群成员列表(GetGroupMemberList)")
	var messageData wba.APIRequestInfo
	var err error
	messageData.Action = "get_group_member_list"
	messageData.Params.GroupId = groupId
	messageData.Echo, err = GenerateUUID()
	if err != nil {
		LOG.ERROR("获取群成员列表(GetGroupMemberList)时，生成UUID失败: %v", err)
		return wba.APIResponseInfo{}
	}
	Response, err = wsAPI(messageData)
	if err != nil {
		LOG.ERROR("获取群成员列表(GetGroupMemberList)时，执行失败: %v", err)
		return wba.APIResponseInfo{}
	}
	return Response
}

// GetGroupHonorInfo 获取群荣誉信息
func (a *apiInfo) GetGroupHonorInfo(groupId int64, Type string) (Response wba.APIResponseInfo) {
	LOG.INFO("获取群荣誉信息(GetGroupHonorInfo)")
	var messageData wba.APIRequestInfo
	var err error
	messageData.Action = "get_group_honor_info"
	messageData.Params.GroupId = groupId
	messageData.Params.Type = Type
	messageData.Echo, err = GenerateUUID()
	if err != nil {
		LOG.ERROR("获取群荣誉信息(GetGroupHonorInfo)时，生成UUID失败: %v", err)
		return wba.APIResponseInfo{}
	}
	Response, err = wsAPI(messageData)
	if err != nil {
		LOG.ERROR("获取群荣誉信息(GetGroupHonorInfo)时，执行失败: %v", err)
		return wba.APIResponseInfo{}
	}
	return Response
}

// GetCookies 获取Cookies
func (a *apiInfo) GetCookies(domain string) (Response wba.APIResponseInfo) {
	LOG.INFO("获取Cookies(GetCookies)")
	var messageData wba.APIRequestInfo
	var err error
	messageData.Action = "get_cookies"
	messageData.Params.Domain = domain
	messageData.Echo, err = GenerateUUID()
	if err != nil {
		LOG.ERROR("获取Cookies(GetCookies)时，生成UUID失败: %v", err)
		return wba.APIResponseInfo{}
	}
	Response, err = wsAPI(messageData)
	if err != nil {
		LOG.ERROR("获取Cookies(GetCookies)时，执行失败: %v", err)
		return wba.APIResponseInfo{}
	}
	return Response
}

// GetCSRFToken 获取CSRF Token
func (a *apiInfo) GetCSRFToken() (Response wba.APIResponseInfo) {
	LOG.INFO("获取CSRF Token(GetCSRFToken)")
	var messageData wba.APIRequestInfo
	var err error
	messageData.Action = "get_csrf_token"
	messageData.Echo, err = GenerateUUID()
	if err != nil {
		LOG.ERROR("获取CSRF Token(GetCSRFToken)时，生成UUID失败: %v", err)
		return wba.APIResponseInfo{}
	}
	Response, err = wsAPI(messageData)
	if err != nil {
		LOG.ERROR("获取CSRF Token(GetCSRFToken)时，执行失败: %v", err)
		return wba.APIResponseInfo{}
	}
	return Response
}

// GetCredentials 获取登录令牌
func (a *apiInfo) GetCredentials(domain string) (Response wba.APIResponseInfo) {
	LOG.INFO("获取登录令牌(GetCredentials)")
	var messageData wba.APIRequestInfo
	var err error
	messageData.Action = "get_credentials"
	messageData.Params.Domain = domain
	messageData.Echo, err = GenerateUUID()
	if err != nil {
		LOG.ERROR("获取登录令牌(GetCredentials)时，生成UUID失败: %v", err)
		return wba.APIResponseInfo{}
	}
	Response, err = wsAPI(messageData)
	if err != nil {
		LOG.ERROR("获取登录令牌(GetCredentials)时，执行失败: %v", err)
		return wba.APIResponseInfo{}
	}
	return Response
}

// GetRecord 获取语音
func (a *apiInfo) GetRecord(file string, outFormat string) (Response wba.APIResponseInfo) {
	LOG.INFO("获取语音(GetRecord)")
	var messageData wba.APIRequestInfo
	var err error
	messageData.Action = "get_record"
	messageData.Params.File = file
	messageData.Params.OutFormat = outFormat
	messageData.Echo, err = GenerateUUID()
	if err != nil {
		LOG.ERROR("获取语音(GetRecord)时，生成UUID失败: %v", err)
		return wba.APIResponseInfo{}
	}
	Response, err = wsAPI(messageData)
	if err != nil {
		LOG.ERROR("获取语音(GetRecord)时，执行失败: %v", err)
		return wba.APIResponseInfo{}
	}
	return Response
}

// GetImage 获取图片
func (a *apiInfo) GetImage(file string) (Response wba.APIResponseInfo) {
	LOG.INFO("获取图片(GetImage)")
	var messageData wba.APIRequestInfo
	var err error
	messageData.Action = "get_image"
	messageData.Params.File = file
	messageData.Echo, err = GenerateUUID()
	if err != nil {
		LOG.ERROR("获取图片(GetImage)时，生成UUID失败: %v", err)
		return wba.APIResponseInfo{}
	}
	Response, err = wsAPI(messageData)
	if err != nil {
		LOG.ERROR("获取图片(GetImage)时，执行失败: %v", err)
		return wba.APIResponseInfo{}
	}
	return Response
}

// CanSendImage 检查是否可以发送图片
func (a *apiInfo) CanSendImage() (Response wba.APIResponseInfo) {
	LOG.INFO("检查是否可以发送图片(CanSendImage)")
	var messageData wba.APIRequestInfo
	var err error
	messageData.Action = "can_send_image"
	messageData.Echo, err = GenerateUUID()
	if err != nil {
		LOG.ERROR("检查是否可以发送图片(CanSendImage)时，生成UUID失败: %v", err)
		return wba.APIResponseInfo{}
	}
	Response, err = wsAPI(messageData)
	if err != nil {
		LOG.ERROR("检查是否可以发送图片(CanSendImage)时，执行失败: %v", err)
		return wba.APIResponseInfo{}
	}
	return Response
}

// CanSendRecord 检查是否可以发送语音
func (a *apiInfo) CanSendRecord() (Response wba.APIResponseInfo) {
	LOG.INFO("检查是否可以发送语音(CanSendRecord)")
	var messageData wba.APIRequestInfo
	var err error
	messageData.Action = "can_send_record"
	messageData.Echo, err = GenerateUUID()
	if err != nil {
		LOG.ERROR("检查是否可以发送语音(CanSendRecord)时，生成UUID失败: %v", err)
		return wba.APIResponseInfo{}
	}
	Response, err = wsAPI(messageData)
	if err != nil {
		LOG.ERROR("检查是否可以发送语音(CanSendRecord)时，执行失败: %v", err)
		return wba.APIResponseInfo{}
	}
	return Response
}

// GetStatus 获取状态
func (a *apiInfo) GetStatus() (Response wba.APIResponseInfo) {
	LOG.INFO("获取状态(GetStatus)")
	var messageData wba.APIRequestInfo
	var err error
	messageData.Action = "get_status"
	messageData.Echo, err = GenerateUUID()
	if err != nil {
		LOG.ERROR("获取状态(GetStatus)时，生成UUID失败: %v", err)
		return wba.APIResponseInfo{}
	}
	Response, err = wsAPI(messageData)
	if err != nil {
		LOG.ERROR("获取状态(GetStatus)时，执行失败: %v", err)
		return wba.APIResponseInfo{}
	}
	return Response
}

//二、LOG模块

/*
关于LOG模块的说明

1.日志模块使用go-logging库，日志级别分为DEBUG、INFO、WARN、ERROR。

2.日志模块提供LogWith方法，可以自定义日志级别，调用级别为DEBUG时，会打印输出调用者的文件名、函数名、行号。

3.日志模块提供Log方法，默认日志级别为INFO。
*/

func (a *apiInfo) LogWith(level string, content string, args ...interface{}) {
	switch level {
	case "DEBUG":
		LOG.DEBUG(content, args...)
		return
	case "WARN":
		LOG.WARN(content, args...)
		return
	case "ERROR":
		LOG.ERROR(content, args...)
		return
	default:
		LOG.INFO(content, args...)
		return
	}
}

func (a *apiInfo) Log(content string, args ...interface{}) {
	LOG.INFO(content, args...)
}

//database模块
//TODO: 数据库模块待实现

// 文件管理模块
//TODO: 文件管理模块待实现

//终端连接模块
//TODO: 终端模块待实现

//核心信息调用模块

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
