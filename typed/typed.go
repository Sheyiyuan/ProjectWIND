package typed

type CoreConfigInfo struct {
	CoreName     string `json:"core_name"`
	ProtocolAddr string `json:"protocol_addr"`
	Token        string `json:"token"`
	WebUIPort    uint16 `json:"webui_port"`
	PasswordHash string `json:"password_hash"`
	ServiceName  string `json:"service_name"`
}
