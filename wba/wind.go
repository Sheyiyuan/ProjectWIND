package wba

import (
	"fmt"
)

type APP interface {
	Get() AppInfo
	Init(api WindStandardProtocolAPI) error
}

// WindStandardProtocolAPI Wind标准协议API,提供了onebot11标准中的API接口。
type WindStandardProtocolAPI interface {
	// UnsafelySendMsg [不安全][需要master权限]按照QQ号或群号发送消息到指定的群组或用户。
	// 参数:
	// - msgType: 消息类型，如 "group" 或 "private"。
	// - groupId: 群组 ID，若为私聊则为 0。
	// - userId: 用户 ID，若为群聊则为 0。
	// - message: 要发送的消息内容。
	// - autoEscape: 消息内容是否作为纯文本发送（即不解析 CQ 码）。
	UnsafelySendMsg(msgType string, groupId int64, userId int64, message string, autoEscape bool)

	// UnsafelySendPrivateMsg [不安全][需要master权限]按照QQ号发送私聊消息给指定用户。
	// 参数:
	// - userId: 用户 ID。
	// - message: 要发送的消息内容。
	// - autoEscape: 消息内容是否作为纯文本发送（即不解析 CQ 码）。
	UnsafelySendPrivateMsg(userId int64, message string, autoEscape bool)

	// UnsafelySendGroupMsg [不安全][需要master权限]按照群号发送群聊消息到指定群组。
	// 参数:
	// - groupId: 群组 ID。
	// - message: 要发送的消息内容。
	// - autoEscape: 消息内容是否作为纯文本发送（即不解析 CQ 码）。
	UnsafelySendGroupMsg(groupId int64, message string, autoEscape bool)

	// SendMsg 按照提供的消息事件信息回复指定消息。
	// 参数:
	// - msg: 消息事件信息。
	// - message: 要回复的消息内容。
	// - autoEscape: 是否自动转义消息内容。
	SendMsg(msg MessageEventInfo, message string, autoEscape bool)

	// SendPrivateMsg 按照提供的消息事件信息回复私聊消息。
	// 参数:
	// - msg: 消息事件信息。
	// - message: 要回复的消息内容。
	// - autoEscape: 消息内容是否作为纯文本发送（即不解析 CQ 码）。
	SendPrivateMsg(msg MessageEventInfo, message string, autoEscape bool)

	// SendGroupMsg 按照提供的消息事件信息回复群聊消息。
	// 参数:
	// - msg: 消息事件信息。
	// - message: 要回复的消息内容。
	// - autoEscape: 消息内容是否作为纯文本发送（即不解析 CQ 码）。
	SendGroupMsg(msg MessageEventInfo, message string, autoEscape bool)

	// UnsafelyDeleteMsg [不安全][需要master权限]撤回指定消息。
	// 参数:
	// - msgId: 消息 ID。
	UnsafelyDeleteMsg(msgId int32)

	// DeleteMsg 按照提供的消息事件信息撤回指定的消息。
	// 参数:
	// - msg: 消息事件信息。
	DeleteMsg(msg MessageEventInfo)

	// SendLike 给指定用户发送点赞，支持指定点赞次数。
	// 参数:
	// - userId: 用户 ID。
	// - times: 赞的次数，每个好友每天最多 10 次。
	SendLike(userId int64, times int)

	// SetGroupKick 将指定用户踢出指定群组，支持拒绝该用户再次加入。
	// 参数:
	// - groupId: 群组 ID。
	// - userId: 用户 ID。
	// - rejectAddRequest: 是否拒绝该用户再次加入。
	SetGroupKick(groupId int64, userId int64, rejectAddRequest bool)

	// SetGroupBan 禁言指定用户在指定群组的发言，支持指定禁言时长。
	// 参数:
	// - groupId: 群组 ID。
	// - userId: 用户 ID。
	// - duration: 禁言时长，单位为秒，0 表示取消禁言。
	SetGroupBan(groupId int64, userId int64, duration int32)

	// SetGroupWholeBan 开启或关闭指定群组的全员禁言。
	// 参数:
	// - groupId: 群组 ID。
	// - enable: 是否开启全员禁言。
	SetGroupWholeBan(groupId int64, enable bool)

	// SetGroupAdmin 设置或取消指定用户在指定群组的管理员权限。
	// 参数:
	// - groupId: 群组 ID。
	// - userId: 用户 ID。
	// - enable: 是否设置为管理员。
	SetGroupAdmin(groupId int64, userId int64, enable bool)

	// SetGroupLeave 退出指定群组。
	// 参数:
	// - groupId: 群组 ID。
	// - isDismiss: 是否解散，如果登录号是群主，则仅在此项为 true 时能够解散。
	SetGroupLeave(groupId int64, isDismiss bool)

	// SetGroupCard 设置指定用户在指定群组的群名片。
	// 参数:
	// - groupId: 群组 ID。
	// - userId: 用户 ID。
	// - card: 群名片内容。
	SetGroupCard(groupId int64, userId int64, card string)

	// SetGroupName 设置指定群组的名称。
	// 参数:
	// - groupId: 群组 ID。
	// - groupName: 群组名称。
	SetGroupName(groupId int64, groupName string)

	// SetGroupSpecialTitle 设置指定用户在指定群组的专属头衔，支持指定头衔有效期。
	// 参数:
	// - groupId: 群组 ID。
	// - userId: 用户 ID。
	// - specialTitle: 专属头衔内容。
	SetGroupSpecialTitle(groupId int64, userId int64, specialTitle string)

	// SetFriendAddRequest 处理好友添加请求，支持同意或拒绝并添加备注。
	// 参数:
	// - flag: 请求标识(需从上报的数据中获得)。
	// - approve: 是否同意请求。
	// - remark: 备注信息（仅在同意时有效）。
	SetFriendAddRequest(flag string, approve bool, remark string)

	// SetGroupAddRequest 处理群组添加请求，支持同意或拒绝并添加理由。
	// 参数:
	// - flag: 请求标识。
	// - subType: 请求子类型。
	// - approve: 是否同意请求。
	// - reason: 拒绝理由。
	SetGroupAddRequest(flag string, subType string, approve bool, reason string)

	// GetLoginInfo 获取当前登录的账号信息。
	// 返回: API 响应信息。
	GetLoginInfo() APIResponseInfo

	// GetVersionInfo 获取当前程序的版本信息。
	// 返回: API 响应信息。
	GetVersionInfo() APIResponseInfo

	// GetMsg 获取指定消息 ID 的消息内容。
	// 参数:
	// - msgId: 消息 ID。
	// 返回: API 响应信息。
	GetMsg(msgId int32) APIResponseInfo

	// GetForwardMsg 获取指定转发消息 ID 的消息内容。
	// 参数:
	// - msgId: 转发消息 ID。
	// 返回: API 响应信息。
	GetForwardMsg(msgId string) APIResponseInfo

	// GetGroupList 获取当前登录账号加入的所有群组列表。
	// 返回: API 响应信息。
	GetGroupList() APIResponseInfo

	// GetGroupMemberList 获取指定群组的所有成员列表。
	// 参数:
	// - groupId: 群组 ID。
	// 返回: API 响应信息。
	GetGroupMemberList(groupId int64) APIResponseInfo

	// GetGroupMemberInfo 获取指定群组中指定用户的详细信息，支持缓存控制。
	// 参数:
	// - groupId: 群组 ID。
	// - userId: 用户 ID。
	// - noCache: 是否不使用缓存。
	// 返回: API 响应信息。
	GetGroupMemberInfo(groupId int64, userId int64, noCache bool) APIResponseInfo

	// GetFriendList 获取当前登录账号的所有好友列表。
	// 返回: API 响应信息。
	GetFriendList() APIResponseInfo

	// GetStrangerInfo 获取指定用户的详细信息，支持缓存控制。
	// 参数:
	// - userId: 用户 ID。
	// - noCache: 是否不使用缓存。
	// 返回: API 响应信息。
	GetStrangerInfo(userId int64, noCache bool) APIResponseInfo

	// GetGroupInfo 获取指定群组的详细信息，支持缓存控制。
	// 参数:
	// - groupId: 群组 ID。
	// - noCache: 是否不使用缓存。
	// 返回: API 响应信息。
	GetGroupInfo(groupId int64, noCache bool) APIResponseInfo

	// GetGroupHonorInfo 获取指定群组的荣誉信息。
	// 参数:
	// - groupId: 群组 ID。
	// - Type: 荣誉类型。
	// 返回: API 响应信息。
	GetGroupHonorInfo(groupId int64, Type string) APIResponseInfo

	// GetStatus 获取当前程序的运行状态信息。
	// 返回: API 响应信息。
	GetStatus() APIResponseInfo

	//	// GetCookies 获取指定域名的 Cookies 信息。
	// 参数:
	// - domain: 域名。
	// 返回: API 响应信息。
	GetCookies(domain string) APIResponseInfo

	// GetCSRFToken 获取 CSRF 令牌。
	// 返回: API 响应信息。
	GetCSRFToken() APIResponseInfo

	// GetCredentials 获取指定域名的凭证信息。
	// 参数:
	// - domain: 域名。
	// 返回: API 响应信息。
	GetCredentials(domain string) APIResponseInfo

	// GetImage 获取指定文件的图片信息。
	// 参数:
	// - file: 文件路径或标识符。
	// 返回: API 响应信息。
	GetImage(file string) APIResponseInfo

	// GetRecord 获取指定文件的语音记录信息，支持指定输出格式。
	// 参数:
	// - file: 文件路径或标识符。
	// - outFormat: 输出格式。
	// 返回: API 响应信息。
	GetRecord(file string, outFormat string) APIResponseInfo

	// CanSendImage 检查是否可以发送图片。
	// 返回: API 响应信息。
	CanSendImage() APIResponseInfo

	// CanSendRecord 检查是否可以发送语音记录。
	// 返回: API 响应信息。
	CanSendRecord() APIResponseInfo

	// SetRestart 设置程序在指定延迟后重启。
	// 参数:
	// - delay: 延迟时间，单位为秒。
	SetRestart(delay int32)

	// CleanCache 清理程序的缓存。
	CleanCache()

	// LogWith 使用指定日志级别记录日志，支持可变参数占位符。
	// 参数:
	// - level: 日志级别: "trace", "debug", "info", "notice", "warn", "error"。
	// - log: 日志内容。
	// - args: 可变参数，用于格式化日志内容。
	LogWith(level string, log string, args ...interface{})

	// Log 记录日志，级别为 "info"，支持可变参数占位符。
	// 参数:
	// - log: 日志内容。
	// - args: 可变参数，用于格式化日志内容。
	Log(log string, args ...interface{})
}

type Database interface {
	VarSet(app AppInfo, datamap string, unit string, id string, key string, value string)
	VarGet(app AppInfo, datamap string, unit string, id string, key string) (string, bool)
	GetIntConfig(app AppInfo, datamap string, key string) (int64, bool)
	GetStringConfig(app AppInfo, datamap string, key string) (string, bool)
	GetFloatConfig(app AppInfo, datamap string, key string) (float64, bool)
	GetIntSliceConfig(app AppInfo, datamap string, key string) ([]int64, bool)
	GetStringSliceConfig(app AppInfo, datamap string, key string) ([]string, bool)
}

type AppInfo struct {
	Name                string
	Version             string
	Author              string
	Description         string
	Namespace           string
	Homepage            string
	License             string
	AppType             string
	Rule                string
	CmdMap              map[string]Cmd
	MessageEventHandler func(msg MessageEventInfo)
	NoticeEventHandler  func(msg NoticeEventInfo)
	RequestEventHandler func(msg RequestEventInfo)
	MetaEventHandler    func(msg MetaEventInfo)
	ScheduledTasks      map[string]ScheduledTaskInfo
	API                 map[string]interface{}
}

func (ai AppInfo) Get() AppInfo {
	return ai
}

func (ai *AppInfo) Init(api WindStandardProtocolAPI) error {
	WSP = api
	return nil
}

func (ai *AppInfo) AddCmd(name string, cmd Cmd) {
	ai.CmdMap[name] = cmd
}

func (ai *AppInfo) AddNoticeEventHandler(ScheduledTask ScheduledTaskInfo) {
	ai.ScheduledTasks[ScheduledTask.Name] = ScheduledTask
}

type AppInfoOption func(ei *AppInfo)

func WithName(name string) AppInfoOption {
	return func(ei *AppInfo) {
		ei.Name = name
	}
}

func WithVersion(version string) AppInfoOption {
	return func(ei *AppInfo) {
		ei.Version = version
	}
}

func WithAuthor(author string) AppInfoOption {
	return func(ei *AppInfo) {
		ei.Author = author
	}
}

func WithDescription(description string) AppInfoOption {
	return func(ei *AppInfo) {
		ei.Description = description
	}
}

func WithNamespace(namespace string) AppInfoOption {
	return func(ei *AppInfo) {
		ei.Namespace = namespace
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

func WithAppType(appType string) AppInfoOption {
	return func(ei *AppInfo) {
		ei.AppType = appType
	}
}

func WithRule(rule string) AppInfoOption {
	return func(ei *AppInfo) {
		ei.Rule = fmt.Sprintf("rule_%s", rule)
	}
}

func NewApp(opts ...AppInfoOption) AppInfo {
	Ext := AppInfo{
		Name:        "WSP",
		Version:     "v1.0.0",
		Author:      "WSP",
		Description: "A simple and easy-to-use bot framework",
		Namespace:   "PUBLIC",
		Homepage:    "https://github.com/Sheyiyuan/wind_app_model",
		License:     "MIT",
		AppType:     "fun",
		Rule:        "none",
		CmdMap:      make(map[string]Cmd),
	}
	for _, opt := range opts {
		opt(&Ext)
	}
	return Ext
}

func NewCmd(name string, description string, solve func(args []string, msg MessageEventInfo)) Cmd {
	return Cmd{
		NAME:  name,
		DESC:  description,
		SOLVE: solve,
	}
}

func NewScheduledTask(name string, description string, cron string, task func()) ScheduledTaskInfo {
	return ScheduledTaskInfo{
		Name: name,
		Desc: description,
		Cron: cron,
		Task: task,
	}
}

type Cmd struct {
	NAME  string
	DESC  string
	SOLVE func(args []string, msg MessageEventInfo)
}

type MessageEventInfo struct {
	Time        int64         `json:"time,omitempty"`
	SelfId      int64         `json:"self_id,omitempty"`
	PostType    string        `json:"post_type,omitempty"`
	MessageType string        `json:"message_type,omitempty"`
	SubType     string        `json:"sub_type,omitempty"`
	MessageId   int32         `json:"message_id,omitempty"`
	GroupId     int64         `json:"group_id,omitempty"`
	UserId      int64         `json:"user_id,omitempty"`
	Anonymous   AnonymousInfo `json:"anonymous"`
	Message     []MessageInfo `json:"message,omitempty"`
	RawMessage  string        `json:"raw_message,omitempty"`
	Font        int32         `json:"font,omitempty"`
	Sender      SenderInfo    `json:"sender"`
}

type NoticeEventInfo struct {
	Time       int64    `json:"time,omitempty"`
	SelfId     int64    `json:"self_id,omitempty"`
	PostType   string   `json:"post_type,omitempty"`
	NoticeType string   `json:"notice_type,omitempty"`
	GroupId    int64    `json:"group_id,omitempty"`
	UserId     int64    `json:"user_id,omitempty"`
	File       FileInfo `json:"file,omitempty"`
	SubType    string   `json:"sub_type,omitempty"`
	OperatorId int64    `json:"operator_id,omitempty"`
	Duration   int64    `json:"duration,omitempty"`
	MessageId  int64    `json:"message,omitempty"`
	TargetId   int64    `json:"target_id,omitempty"`
	HonorType  string   `json:"honor_type,omitempty"`
}

type RequestEventInfo struct {
	Time        int64  `json:"time,omitempty"`
	SelfId      int64  `json:"self_id,omitempty"`
	PostType    string `json:"post_type,omitempty"`
	RequestType string `json:"request_type,omitempty"`
	SubType     string `json:"sub_type,omitempty"`
	UserId      int64  `json:"user_id,omitempty"`
	Comment     string `json:"comment,omitempty"`
	Flag        string `json:"flag,omitempty"`
	GroupId     int64  `json:"group_id,omitempty"`
}

type MetaEventInfo struct {
	Time          int64  `json:"time,omitempty"`
	SelfId        int64  `json:"self_id,omitempty"`
	PostType      string `json:"post_type,omitempty"`
	MetaEventType string `json:"meta_event_type,omitempty"`
	SubType       string `json:"sub_type,omitempty"`
	Status        string `json:"status,omitempty"`
	Interval      int64  `json:"interval,omitempty"`
}

type FileInfo struct {
	Id    string `json:"id,omitempty"`
	Name  string `json:"Name,omitempty"`
	Size  int64  `json:"size,omitempty"`
	Busid int64  `json:"bucket,omitempty"`
}

type SenderInfo struct {
	UserId   int64  `json:"user_id,omitempty"`
	Nickname string `json:"nickname,omitempty"`
	Card     string `json:"card,omitempty"`
	Sex      string `json:"sex,omitempty"`
	Age      int32  `json:"age,omitempty"`
	Area     string `json:"area,omitempty"`
	Level    string `json:"level,omitempty"`
	Role     string `json:"role,omitempty"`
	Title    string `json:"title,omitempty"`
}

type AnonymousInfo struct {
	Id   string `json:"id,omitempty"`
	Name string `json:"Name,omitempty"`
	Flag string `json:"flag,omitempty"`
}

type MessageInfo struct {
	Type string          `json:"type,omitempty"`
	Data MessageDataInfo `json:"data"`
}

type MessageDataInfo struct {
	Type    string `json:"type,omitempty"`
	Text    string `json:"text,omitempty"`
	Id      string `json:"id,omitempty"`
	File    string `json:"file,omitempty"`
	Url     string `json:"url,omitempty"`
	Magic   string `json:"magic,omitempty"`
	Qq      string `json:"qq,omitempty"`
	Title   string `json:"title,omitempty"`
	Content string `json:"content,omitempty"`
	Image   string `json:"image,omitempty"`
	Audio   string `json:"audio,omitempty"`
	Lat     string `json:"lat,omitempty"`
	Lon     string `json:"lon,omitempty"`
	Data    string `json:"data,omitempty"`
}

type ParamsInfo struct {
	Message          string `json:"message,omitempty"`
	UserId           int64  `json:"user_id,omitempty"`
	GroupId          int64  `json:"group_id,omitempty"`
	AutoEscape       bool   `json:"auto_escape,omitempty"`
	MessageId        int32  `json:"message_id,omitempty"`
	Id               string `json:"id,omitempty"`
	RejectAddRequest bool   `json:"reject_add_request,omitempty"`
	Duration         int32  `json:"duration,omitempty"`
	Enable           bool   `json:"enable,omitempty"`
	Card             string `json:"card,omitempty"`
	GroupName        string `json:"group_name,omitempty"`
	IsDismiss        bool   `json:"is_dismiss,omitempty"`
	SpecialTitle     string `json:"special_title,omitempty"`
	Flag             string `json:"flag,omitempty"`
	Approve          bool   `json:"approve,omitempty"`
	Remark           string `json:"remark,omitempty"`
	Type             string `json:"type,omitempty"`
	SubType          string `json:"sub_type,omitempty"`
	Reason           string `json:"reason,omitempty"`
	NoCache          bool   `json:"no_cache,omitempty"`
	File             string `json:"file,omitempty"`
	Times            int    `json:"times,omitempty"`
	Domain           string `json:"domain,omitempty"`
	OutFormat        string `json:"out_format,omitempty"`
	Delay            int32  `json:"delay,omitempty"`
}

type APIRequestInfo struct {
	Action string     `json:"action,omitempty"`
	Params ParamsInfo `json:"params"`
	Echo   string     `json:"echo,omitempty"`
}

type APIResponseInfo struct {
	Status  string           `json:"status,omitempty"`
	Retcode int64            `json:"retcode,omitempty"`
	Data    ResponseDataInfo `json:"data,omitempty"`
	Echo    string           `json:"echo,omitempty"`
}

type APIResponseListInfo struct {
	Status  string             `json:"status,omitempty"`
	Retcode int64              `json:"retcode,omitempty"`
	Data    []ResponseDataInfo `json:"data,omitempty"`
	Echo    string             `json:"echo,omitempty"`
}

type ResponseDataInfo struct {
	UserId           int64                  `json:"user_id,omitempty"`
	Nickname         string                 `json:"nickname,omitempty"`
	Sex              string                 `json:"sex,omitempty"`
	Age              int32                  `json:"age,omitempty"`
	Remark           string                 `json:"remark,omitempty"`
	GroupId          int64                  `json:"group_id,omitempty"`
	GroupName        string                 `json:"group_name,omitempty"`
	MemberCount      int32                  `json:"member_count,omitempty"`
	MaxMemberCount   int32                  `json:"max_member_count,omitempty"`
	Card             string                 `json:"card,omitempty"`
	Area             string                 `json:"area,omitempty"`
	JoinTime         int32                  `json:"join_time,omitempty"`
	LastSentTime     int32                  `json:"last_sent_time,omitempty"`
	Level            string                 `json:"level,omitempty"`
	Role             string                 `json:"role,omitempty"`
	Unfriendly       bool                   `json:"unfriendly,omitempty"`
	Title            string                 `json:"title,omitempty"`
	TitleExpireTime  int32                  `json:"title_expire_time,omitempty"`
	CardChangeable   bool                   `json:"card_changeable,omitempty"`
	CurrentTalkative CurrentTalkativeInfo   `json:"current_talkative,omitempty"`
	TalkativeList    []CurrentTalkativeInfo `json:"talkative_list,omitempty"`
	PerformerList    []HonorInfo            `json:"performer_list,omitempty"`
	LegendList       []HonorInfo            `json:"legend_list,omitempty"`
	StrongNewbieList []HonorInfo            `json:"strong_newbie_list,omitempty"`
	EmoticonList     []HonorInfo            `json:"emoticon_list,omitempty"`
	Cookies          string                 `json:"cookies,omitempty"`
	Token            string                 `json:"token,omitempty"`
	CsrfToken        string                 `json:"csrf_token,omitempty"`
	File             string                 `json:"file,omitempty"`
	OutFormat        string                 `json:"out_format,omitempty"`
	Yes              bool                   `json:"yes,omitempty"`
	Online           bool                   `json:"online,omitempty"`
	Good             bool                   `json:"good,omitempty"`
	AppName          string                 `json:"app_name,omitempty"`
	AppVersion       string                 `json:"app_version,omitempty"`
	ProtocolVersion  string                 `json:"protocol_version,omitempty"`
	Time             int64                  `json:"time,omitempty"`
	MessageType      string                 `json:"message_type,omitempty"`
	MessageId        int32                  `json:"message_id,omitempty"`
	RealId           int32                  `json:"real_id,omitempty"`
	Sender           SenderInfo             `json:"sender,omitempty"`
	Message          []MessageDataInfo      `json:"message,omitempty"`
}

type CurrentTalkativeInfo struct {
	UserId   int64  `json:"user_id,omitempty"`
	Nickname string `json:"nickname,omitempty"`
	Avatar   string `json:"avatar,omitempty"`
	DayCount int32  `json:"day_count,omitempty"`
}

type HonorInfo struct {
	UserId      int64  `json:"user_id,omitempty"`
	Nickname    string `json:"nickname,omitempty"`
	Avatar      string `json:"avatar,omitempty"`
	Description string `json:"Description,omitempty"`
}

// SegmentInfo 消息段
type SegmentInfo struct {
	Type string          `json:"type,omitempty"`
	Data SegmentDataInfo `json:"data,omitempty"`
}

type SegmentDataInfo struct {
	Type     string `json:"type,omitempty"`
	QQ       string `json:"qq,omitempty"`
	Id       int64  `json:"id,omitempty"`
	UserId   int64  `json:"user_id,omitempty"`
	Nickname string `json:"nickname,omitempty"`
	Content  string `json:"content,omitempty"`
	Url      string `json:"url,omitempty"`
	Lat      string `json:"lat,omitempty"`
	Lon      string `json:"lon,omitempty"`
	Title    string `json:"title,omitempty"`
	Audio    string `json:"audio,omitempty"`
	Image    string `json:"image,omitempty"`
	Video    string `json:"video,omitempty"`
	Data     string `json:"data,omitempty"`
}

type ScheduledTaskInfo struct {
	Name string `json:"Name,omitempty"`
	Desc string `json:"desc,omitempty"`
	Task func() `json:"task,omitempty"`
	Cron string `json:"cron,omitempty"`
}

var WSP WindStandardProtocolAPI
