package wba

import "strings"

type AppInfo struct {
	AppKey              AppKey
	Author              string
	Description         string
	Homepage            string
	License             string
	CmdMap              map[string]Cmd
	MessageEventHandler func(msg MessageEventInfo)
	NoticeEventHandler  func(msg NoticeEventInfo)
	RequestEventHandler func(msg RequestEventInfo)
	MetaEventHandler    func(msg MetaEventInfo)
	ScheduledTasks      map[string]ScheduledTaskInfo
	API                 map[string]interface{}
}

func (ai *AppInfo) AddCmd(name string, cmd Cmd) {
	ai.CmdMap[name] = cmd
}

func (ai *AppInfo) AddNoticeEventHandler(ScheduledTask ScheduledTaskInfo) {
	ai.ScheduledTasks[ScheduledTask.Name] = ScheduledTask
}

func (ai *AppInfo) AddScheduledTask(task ScheduledTaskInfo) {
	ai.ScheduledTasks[task.Name] = task
}

type AppInfoOption func(ei *AppInfo)

func WithDescription(description string) AppInfoOption {
	return func(ei *AppInfo) {
		ei.Description = description
	}
}

func WithWebUrl(webUrl string) AppInfoOption {
	return func(ei *AppInfo) {
		ei.Homepage = webUrl
	}
}

func WithLicense(license string) AppInfoOption {
	return func(ei *AppInfo) {
		ei.License = license
	}
}

func WithSelector(level uint8, selector string, option string) AppInfoOption {
	return func(ei *AppInfo) {
		ei.AppKey.Level = level
		ei.AppKey.Option = option
		ei.AppKey.Selector = selector
	}
}

func WithLevel(level uint8) AppInfoOption {
	if level < 1 {
		level = 1
	}
	return func(ei *AppInfo) {
		ei.AppKey.Level = level
	}
}

func toCamelCase(s string) string {
	if s == "" {
		return s
	}
	return strings.ToLower(s[:1]) + s[1:]
}

func NewApp(name string, version string, author string, opts ...AppInfoOption) AppInfo {
	Ext := AppInfo{
		AppKey: AppKey{
			Name:     name,
			Version:  version,
			Level:    2,
			Selector: name,
			Option:   name,
		},
		Author:         author,
		Description:    "A simple and easy-to-use bot framework",
		Homepage:       "https://github.com/Sheyiyuan/wind_app_model",
		License:        "MIT",
		CmdMap:         make(map[string]Cmd),
		ScheduledTasks: map[string]ScheduledTaskInfo{},
		API:            map[string]interface{}{},
	}
	for _, opt := range opts {
		opt(&Ext)
	}
	return Ext
}

//func (ai *AppInfo) NewCmd(name string, description string, solve func(args []string, msg MessageEventInfo)) Cmd {
//	return Cmd{
//		Name:   name,
//		Desc:   description,
//		Solve:  solve,
//		AppKey: ai.AppKey,
//	}
//}

func (ai *AppInfo) NewCmd() Cmd {
	return Cmd{
		Name:   "",
		Desc:   "",
		Solve:  nil,
		AppKey: ai.AppKey,
	}
}

func (ai *AppInfo) NewScheduledTask(name string, description string, cron string, task func()) ScheduledTaskInfo {
	return ScheduledTaskInfo{
		Name: name,
		Desc: description,
		Cron: cron,
		Task: task,
	}
}

type Cmd struct {
	Name   string
	Desc   string
	Solve  func(args []string, msg MessageEventInfo)
	AppKey AppKey
}

type ScheduledTaskInfo struct {
	Name string `json:"Name,omitempty"`
	Desc string `json:"desc,omitempty"`
	Task func() `json:"task,omitempty"`
	Cron string `json:"cron,omitempty"`
}

type AppKey struct {
	Name     string        `json:"name"`
	Level    Priority      `json:"level"`
	Version  string        `json:"version"`
	Selector SelectorLabel `json:"selector"`
	Option   OptionLabel   `json:"option"`
}

// Priority 是一个整数类型，用于表示命令的优先级。只能是不小于1且不大于4的整数。
type Priority = uint8

// VersionLabel 是一个字符串类型，用于表示版本标签。
type VersionLabel = string

// SelectorLabel 是一个字符串类型，用于表示选择器的标签。
type SelectorLabel = string

// OptionLabel 是一个字符串类型，用于表示选择器选项的标签。
type OptionLabel = string

// SessionLabel 是一个字符串类型，用于表示聊天会话的标签，格式为[平台:类型-ID]，如"QQ:group-1145141919810"
type SessionLabel = string

// CmdSetLabel 是一个字符串类型，用于表示命令集的标签。
type CmdSetLabel = string

// CmdLabel 是一个字符串类型，用于表示命令的标签。
type CmdLabel = string

// CmdList 是一个字符串到 wba.Cmd 的映射，用于存储命令的列表。
type CmdList = map[string]Cmd

// SessionInfo 是一个结构体，用于存储会话信息。
//
// 字段:
//   - Platform: 表示会话的平台，如"QQ"、"Telegram"等。
//   - SessionType: 表示会话的类型，如"group"、"private"等。
//   - SessionId: 表示会话的ID，如群号、私聊号等。
type SessionInfo struct {
	Platform    string
	SessionType string
	SessionId   int64
}

func (s *SessionInfo) Load(platform string, msg MessageEventInfo) SessionInfo {
	s.Platform = platform
	s.SessionType = msg.MessageType
	if s.SessionType == "group" {
		s.SessionId = msg.GroupId
	}
	if s.SessionType == "private" {
		s.SessionId = msg.UserId
	}
	return *s
}

// VersionInfo 是一个结构体，用于存储版本信息。
//
// 字段:
//   - BigVersion: 表示大版本号。
//   - SmallVersion: 表示小版本号。
//   - FixVersion: 表示修复版本号。
type VersionInfo struct {
	BigVersion   uint8
	SmallVersion uint8
	FixVersion   uint8
}
