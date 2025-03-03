package typed

type CoreConfigInfo struct {
	CoreName     string   `json:"core_name"`
	Protocol     Protocol `json:"protocol"`
	WebUIPort    uint16   `json:"webui_port"`
	PasswordHash string   `json:"password_hash"`
	ServiceName  string   `json:"service_name"`
}

type Protocol struct {
	Name     string `json:"protocol_name"`
	Platform string `json:"protocol_platform"`
	Addr     string `json:"protocol_addr"`
	Token    string `json:"token"`
	Enable   bool   `json:"enable"`
}

type AppKey struct {
	AppName    string `json:"app_name"`
	AppType    string `json:"app_type"`
	AppLevel   int32  `json:"app_level"`
	AppVersion string `json:"app_version"`
}

type SessionWorkSpace struct {
	SessionId   string          `json:"session_id"`
	SessionType string          `json:"session_type"`
	Rule        string          `json:"rule"`
	Enable      bool            `json:"enable"`
	AppEnable   map[AppKey]bool `json:"app_enable"`
	CmdEnable   map[string]bool `json:"cmd_enable"`
	WorkLevel   int32           `json:"work_level"`
}
