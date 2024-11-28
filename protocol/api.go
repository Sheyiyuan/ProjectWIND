package protocol

import (
	"ProjectWIND/typed"
	"encoding/json"
	"errors"
)

func init() {

}

func SendMessage(messageType string, message string, targetId int64, autoEscape bool) (bool, error) {
	// 构建发送消息的JSON数据
	var messageData typed.APIRequest
	messageData.Action = "send_msg"
	switch messageType {
	case "private":
		messageData.Params.UserId = targetId
		break
	case "group":
		messageData.Params.GroupId = targetId
		break
	default:
		return false, errors.New("invalid type")
	}
	messageData.Params.Message = message
	messageData.Params.AutoEscape = autoEscape
	messageJson, err := json.Marshal(messageData)
	if err != nil {
		return false, err
	}
	// 发送消息
	_, err = wsSendMessage(messageJson)
	if err != nil {
		return false, err
	}
	return true, nil
}
