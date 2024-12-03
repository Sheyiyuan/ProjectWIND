package protocol

import (
	"ProjectWIND/typed"
	"encoding/json"
	"errors"
)

/*
关于API的说明：

1.所有API请求按照OneBot11标准，使用JSON格式进行数据交换。api命名为由原文档中蛇形命名法改为双驼峰命名法。

2.无响应的API请求使用ws协议处理，有响应的API请求使用http协议处理。

3.wind会从配置文件中读取API请求的url，请确保正确填写。
*/

//1.无响应API,使用ws协议处理

func SendMsg(msg typed.MessageEventInfo, message string, autoEscape bool) error {
	// 构建发送消息的JSON数据
	var messageData typed.APIRequestInfo

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
			return errors.New("invalid type")
		}
	}
	messageData.Params.Message = message
	messageData.Params.AutoEscape = autoEscape
	messageJson, err := json.Marshal(messageData)
	if err != nil {
		return err
	}
	// 发送消息
	err = wsAPI(messageJson)
	return err
}

func SendPrivateMsg(msg typed.MessageEventInfo, message string, autoEscape bool) error {
	// 构建发送消息的JSON数据
	var messageData typed.APIRequestInfo
	messageData.Action = "send_private_msg"
	messageData.Params.UserId = msg.UserId
	messageData.Params.Message = message
	messageData.Params.AutoEscape = autoEscape
	messageJson, err := json.Marshal(messageData)
	if err != nil {
		return err
	}
	// 发送消息
	err = wsAPI(messageJson)
	return err
}

func SendGroupMsg(msg typed.MessageEventInfo, message string, autoEscape bool) error {
	// 构建发送消息的JSON数据
	var messageData typed.APIRequestInfo
	messageData.Action = "send_group_msg"
	messageData.Params.GroupId = msg.GroupId
	messageData.Params.Message = message
	messageData.Params.AutoEscape = autoEscape
	messageJson, err := json.Marshal(messageData)
	if err != nil {
		return err
	}
	// 发送消息
	err = wsAPI(messageJson)
	return err
}

func DeleteMsg(msg typed.MessageEventInfo, msgId int64) error {
	// 构建删除消息的JSON数据
	var messageData typed.APIRequestInfo
	messageData.Action = "delete_msg"
	messageData.Params.MessageId = msg.MessageId
	messageJson, err := json.Marshal(messageData)
	if err != nil {
		return err
	}
	err = wsAPI(messageJson)
	return err
}

func sendLike(userId int64, times int) error {
	// 构建发送赞的JSON数据
	var messageData typed.APIRequestInfo
	messageData.Action = "send_like"
	messageData.Params.UserId = userId
	messageData.Params.Times = times
	messageJson, err := json.Marshal(messageData)
	if err != nil {
		return err
	}
	err = wsAPI(messageJson)
	return nil
}

func setGroupKick(groupId int64, userId int64, rejectAddRequest bool) error {
	var messageData typed.APIRequestInfo
	messageData.Action = "set_group_kick"
	messageData.Params.GroupId = groupId
	messageData.Params.UserId = userId
	messageData.Params.RejectAddRequest = rejectAddRequest
	messageJson, err := json.Marshal(messageData)
	if err != nil {
		return err
	}
	err = wsAPI(messageJson)
	return nil
}

func setGroupBan(groupId int64, userId int64, duration int32) error {
	var messageData typed.APIRequestInfo
	messageData.Action = "set_group_ban"
	messageData.Params.GroupId = groupId
	messageData.Params.UserId = userId
	messageData.Params.Duration = duration
	messageJson, err := json.Marshal(messageData)
	if err != nil {
		return err
	}
	err = wsAPI(messageJson)
	return nil
}
func setGroupAnonymousBan(groupId int64, flag string, duration int32) error {
	var messageData typed.APIRequestInfo
	messageData.Action = "set_group_anonymous_ban"
	messageData.Params.GroupId = groupId
	messageData.Params.Flag = flag
	messageData.Params.Duration = duration
	messageJson, err := json.Marshal(messageData)
	if err != nil {
		return err
	}
	err = wsAPI(messageJson)
	return nil
}

func setGroupWholeBan(groupId int64, enable bool) error {
	var messageData typed.APIRequestInfo
	messageData.Action = "set_group_whole_ban"
	messageData.Params.GroupId = groupId
	messageData.Params.Enable = enable
	messageJson, err := json.Marshal(messageData)
	if err != nil {
		return err
	}
	err = wsAPI(messageJson)
	return nil
}

func setGroupAdmin(groupId int64, userId int64, enable bool) error {
	var messageData typed.APIRequestInfo
	messageData.Action = "set_group_admin"
	messageData.Params.GroupId = groupId
	messageData.Params.UserId = userId
	messageData.Params.Enable = enable
	messageJson, err := json.Marshal(messageData)
	if err != nil {
		return err
	}
	err = wsAPI(messageJson)
	return nil
}

func setGroupAnonymous(groupId int64, enable bool) error {
	var messageData typed.APIRequestInfo
	messageData.Action = "set_group_anonymous"
	messageData.Params.GroupId = groupId
	messageData.Params.Enable = enable
	messageJson, err := json.Marshal(messageData)
	if err != nil {
		return err
	}
	err = wsAPI(messageJson)
	return nil
}

func setGroupCard(groupId int64, userId int64, card string) error {
	var messageData typed.APIRequestInfo
	messageData.Action = "set_group_card"
	messageData.Params.GroupId = groupId
	messageData.Params.UserId = userId
	messageData.Params.Card = card
	messageJson, err := json.Marshal(messageData)
	if err != nil {
		return err
	}
	err = wsAPI(messageJson)
	return nil
}

func setGroupName(groupId int64, groupName string) error {
	var messageData typed.APIRequestInfo
	messageData.Action = "set_group_name"
	messageData.Params.GroupId = groupId
	messageData.Params.GroupName = groupName
	messageJson, err := json.Marshal(messageData)
	if err != nil {
		return err
	}
	err = wsAPI(messageJson)
	return nil
}

func setGroupLeave(groupId int64, isDismiss bool) error {
	var messageData typed.APIRequestInfo
	messageData.Action = "set_group_leave"
	messageData.Params.GroupId = groupId
	messageData.Params.IsDismiss = isDismiss
	messageJson, err := json.Marshal(messageData)
	if err != nil {
		return err
	}
	err = wsAPI(messageJson)
	return nil
}

func setGroupSpecialTitle(groupId int64, userId int64, specialTitle string, duration int32) error {
	var messageData typed.APIRequestInfo
	messageData.Action = "set_group_special_title"
	messageData.Params.GroupId = groupId
	messageData.Params.UserId = userId
	messageData.Params.SpecialTitle = specialTitle
	messageData.Params.Duration = duration
	messageJson, err := json.Marshal(messageData)
	if err != nil {
		return err
	}
	err = wsAPI(messageJson)
	return nil
}

func setFriendAddRequest(flag string, approve bool, remark string) error {
	var messageData typed.APIRequestInfo
	messageData.Action = "set_friend_add_request"
	messageData.Params.Flag = flag
	messageData.Params.Approve = approve
	messageData.Params.Remark = remark
	messageJson, err := json.Marshal(messageData)
	if err != nil {
		return err
	}
	err = wsAPI(messageJson)
	return nil
}

// 2.有响应API，使用http协议处理
func GetMsg(messageId int32) (typed.MessageEventInfo, error) {
	// 构建获取消息的JSON数据
	var requestData typed.ParamsInfo
	var msg typed.MessageEventInfo
	action := "get_msg"
	requestData.MessageId = messageId
	body, err := json.Marshal(requestData)
	if err != nil {
		return typed.MessageEventInfo{}, err
	}
	// 发送请求
	_, response, err := httpAPI("POST", action, body)
	if err != nil {
		return typed.MessageEventInfo{}, err
	}
	// 解析响应
	err = json.Unmarshal(response, &msg)
	if err != nil {
		return typed.MessageEventInfo{}, err
	}
	return msg, nil
}
