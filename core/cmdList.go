package core

import (
	"ProjectWIND/LOG"
	"ProjectWIND/typed"
)

var CmdList = map[string]typed.ExtInfo{
	"bot":  ExtCore,
	"help": ExtCore,
}

var ExtCore = typed.ExtInfo{
	Run: func(cmd string, args []string, msg typed.MessageEventInfo) error {
		if cmd == "help" {
			err := SendMsg(msg, "假装有帮助信息", false)
			LOG.INFO("发送核心帮助信息:(至：%v-%v:%v-%v)", msg.MessageType, msg.GroupId, msg.UserId, msg.Sender.Nickname)
			if err != nil {
				return err
			}
		}
		if cmd == "bot" {
			err := SendMsg(msg, "WIND 0.1.0", false)
			LOG.INFO("发送核心版本信息:(至：%v-%v:%v-%v)", msg.MessageType, msg.GroupId, msg.UserId, msg.Sender.Nickname)
			if err != nil {
				return err
			}
		}
		return nil
	},
}
