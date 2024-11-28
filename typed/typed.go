package typed

type ConfigInfo struct {
	CoreName     string            `json:"core_name"`
	ProtocolAddr map[string]string `json:"protocol_addr"`
	WebUIPort    uint16            `json:"webui_port"`
	PasswordHash string            `json:"password_hash"`
	ServiceName  string            `json:"service_name"`
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

type APIRequestInfo struct {
	Action string     `json:"action,omitempty"`
	Params ParamsInfo `json:"params"`
	Echo   string     `json:"echo,omitempty"`
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
}
