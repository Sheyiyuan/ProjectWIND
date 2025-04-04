package wba

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
