package core

import (
	"ProjectWIND/LOG"
	"ProjectWIND/wba"
	"errors"
)

type CmdListInfo map[string]wba.Cmd

type AppInfo struct {
	CmdMap map[string]wba.Cmd
}

func (app AppInfo) Get() AppInfo {
	return app
}

func (app *AppInfo) Run(cmd string, args []string, msg wba.MessageEventInfo) error {
	_, ok := app.CmdMap[cmd]
	if !ok {
		return errors.New("cmd not found")
	}
	app.CmdMap[cmd].SOLVE(args, msg)
	return nil
}

func (app *AppInfo) Init(Api wba.WindAPI) error {
	return nil
}

func (app *AppInfo) GetCmd() map[string]wba.Cmd {
	return app.CmdMap
}

func NewCmd(name string, help string, solve func(args []string, msg wba.MessageEventInfo)) wba.Cmd {
	return wba.Cmd{
		NAME:  name,
		DESC:  help,
		SOLVE: solve,
	}
}

var AppCore = AppInfo{
	CmdMap: CmdListInfo{
		"bot": NewCmd(
			"bot",
			"显示WIND版本信息",
			func(args []string, msg wba.MessageEventInfo) {
				AppApi.SendMsg(msg, "WIND 0.1.0", false)
				LOG.INFO("发送核心版本信息:(至：%v-%v:%v-%v)", msg.MessageType, msg.GroupId, msg.UserId, msg.Sender.Nickname)
			},
		),
	},
}
