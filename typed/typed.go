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
