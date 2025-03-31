package wba

type WindStandardTools interface {
	MsgUnmarshal(message string) (msg MessageEventInfo)
}
