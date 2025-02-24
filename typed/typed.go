package typed

type CoreConfigInfo struct {
	CoreName     string     `json:"core_name"`
	Protocols    []Protocol `json:"protocols"`
	WebUIPort    uint16     `json:"webui_port"`
	PasswordHash string     `json:"password_hash"`
	ServiceName  string     `json:"service_name"`
}

type Protocol struct {
	ProtocolName     string `json:"protocol_name"`
	ProtocolPlatform string `json:"protocol_platform"`
	ProtocolAddr     string `json:"protocol_addr"`
	Token            string `json:"token"`
	Enable           bool   `json:"enable"`
}
