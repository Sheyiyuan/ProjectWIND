package core

import (
	"ProjectWIND/LOG"
	"ProjectWIND/wba"
)

type CmdListInfo map[string]wba.Cmd

type AppInfo struct {
	CmdMap map[string]wba.Cmd
	AppKey wba.AppKey
}

//func (app AppInfo) Get() AppInfo {
//	return app
//}
//
//func (app *AppInfo) Run(cmd string, args []string, msg wba.MessageEventInfo) error {
//	_, ok := app.CmdMap[cmd]
//	if !ok {
//		return errors.New("cmd not found")
//	}
//	app.CmdMap[cmd].Solve(args, msg)
//	return nil
//}

func (app *AppInfo) GetCmd() map[string]wba.Cmd {
	return app.CmdMap
}

func NewCmd(name string, help string, solve func(args []string, msg wba.MessageEventInfo), appKey wba.AppKey) wba.Cmd {
	return wba.Cmd{
		Name:   name,
		Desc:   help,
		Solve:  solve,
		AppKey: appKey,
	}
}

var AppKeyOfCore = wba.AppKey{
	Name:     "core",
	Level:    0,
	Version:  "1.0.0",
	Selector: "core",
	Option:   "core",
}

var AppCore = AppInfo{
	AppKey: AppKeyOfCore,
	CmdMap: CmdListInfo{
		"bot": NewCmd(
			"bot",
			"显示WIND版本信息",
			func(args []string, msg wba.MessageEventInfo) {
				ProtocolApi.SendMsg(msg, "WIND 0.1.0", false)
				LOG.Info("发送核心版本信息:(至：%v-%v:%v-%v)", msg.MessageType, msg.GroupId, msg.UserId, msg.Sender.Nickname)
			},
			AppKeyOfCore,
		),
		"help": NewCmd(
			"help",
			"显示帮助信息",
			func(args []string, msg wba.MessageEventInfo) {
				ProtocolApi.SendMsg(msg, "帮助信息", false)
				LOG.Info("发送帮助信息:(至：%v-%v:%v-%v)", msg.MessageType, msg.GroupId, msg.UserId, msg.Sender.Nickname)
			},
			AppKeyOfCore,
		),
	},
}
