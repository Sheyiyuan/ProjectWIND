package wba

import (
	"errors"
	"fmt"
)

type APP interface {
	Get() AppInfo
	Run(cmd string, args []string, msg MessageEventInfo) error
	Init(api WindAPI) error
}

type WindAPI interface {
	SendMsg(msg MessageEventInfo, message string, autoEscape bool)
	SendPrivateMsg(msg MessageEventInfo, message string, autoEscape bool)
	SendGroupMsg(msg MessageEventInfo, message string, autoEscape bool)
	DeleteMsg(msg MessageEventInfo)
	SendLike(userId int64, times int)
	SetGroupKick(groupId int64, userId int64, rejectAddRequest bool)
	SetGroupBan(groupId int64, userId int64, duration int32)
	SetGroupWholeBan(groupId int64, enable bool)
	SetGroupAdmin(groupId int64, userId int64, enable bool)
	SetGroupLeave(groupId int64, isDismiss bool)
	SetGroupCard(groupId int64, userId int64, card string)
	SetGroupName(groupId int64, groupName string)
	SetGroupSpecialTitle(groupId int64, userId int64, specialTitle string, duration int32)
	SetFriendAddRequest(flag string, approve bool, remark string)
	SetGroupAddRequest(flag string, subType string, approve bool, reason string)
	GetLoginInfo() APIResponseInfo
}

type AppInfo struct {
	name        string
	version     string
	author      string
	description string
	namespace   string
	webUrl      string
	license     string
	appType     string
	rule        string
	CmdMap      map[string]Cmd
}

func (ei AppInfo) Get() AppInfo {
	return ei
}

func (ei *AppInfo) Run(cmd string, args []string, msg MessageEventInfo) error {
	_, ok := ei.CmdMap[cmd]
	if !ok {
		return errors.New("cmd not found")
	}
	ei.CmdMap[cmd].SOLVE(args, msg)
	return nil
}

func (ei *AppInfo) Init(api WindAPI) error {
	Wind = api
	return nil
}

func (ei *AppInfo) AddCmd(name string, cmd Cmd) error {
	ei.CmdMap[name] = cmd
	return nil
}

type AppInfoOption func(ei *AppInfo)

func WithName(name string) AppInfoOption {
	return func(ei *AppInfo) {
		ei.name = name
	}
}

func WithVersion(version string) AppInfoOption {
	return func(ei *AppInfo) {
		ei.version = version
	}
}

func WithAuthor(author string) AppInfoOption {
	return func(ei *AppInfo) {
		ei.author = author
	}
}

func WithDescription(description string) AppInfoOption {
	return func(ei *AppInfo) {
		ei.description = description
	}
}

func WithNamespace(namespace string) AppInfoOption {
	return func(ei *AppInfo) {
		ei.namespace = namespace
	}
}

func WithWebUrl(webUrl string) AppInfoOption {
	return func(ei *AppInfo) {
		ei.webUrl = webUrl
	}
}

func WithLicense(license string) AppInfoOption {
	return func(ei *AppInfo) {
		ei.license = license
	}
}

func WithAppType(appType string) AppInfoOption {
	return func(ei *AppInfo) {
		ei.appType = appType
	}
}

func WithRule(rule string) AppInfoOption {
	return func(ei *AppInfo) {
		ei.rule = fmt.Sprintf("rule_%s", rule)
	}
}

func NewApp(opts ...AppInfoOption) AppInfo {
	Ext := AppInfo{
		name:        "Wind",
		version:     "v1.0.0",
		author:      "Wind",
		description: "A simple and easy-to-use bot framework",
		namespace:   "wind",
		webUrl:      "https://github.com/Sheyiyuan/wind_app_model",
		license:     "MIT",
		appType:     "fun",
		rule:        "none",
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
	Name  string `json:"name,omitempty"`
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
	Name string `json:"name,omitempty"`
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
	Description string `json:"description,omitempty"`
}

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

var Wind WindAPI
