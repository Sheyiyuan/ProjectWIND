package wba

import (
	"fmt"
	"strconv"
	"strings"
)

type WindStandardTools interface {
	// MsgUnmarshal 解析消息JSON字符串为 MessageEventInfo 结构体。
	//
	// 参数:
	//
	//  - message: 要解析的消息字符串。
	//
	// 返回值:
	//
	//  - msg: 解析后的消息结构体。
	MsgUnmarshal(message string) (msg MessageEventInfo)
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

func (v VersionInfo) String() string {
	return fmt.Sprintf("%d.%d.%d", v.BigVersion, v.SmallVersion, v.FixVersion)
}

func VersionLabelAnalysis(versionLabel VersionLabel) VersionInfo {
	version := strings.Split(versionLabel, ".")
	bigVersion, err := strconv.ParseUint(version[0], 10, 8)
	if err != nil {
		return VersionInfo{}
	}
	smallVersion, err := strconv.ParseUint(version[1], 10, 8)
	if err != nil {
		return VersionInfo{}
	}
	fixVersion, err := strconv.ParseUint(version[2], 10, 8)
	if err != nil {
		return VersionInfo{}
	}
	return VersionInfo{
		BigVersion:   uint8(bigVersion),
		SmallVersion: uint8(smallVersion),
		FixVersion:   uint8(fixVersion),
	}
}

func SessionLabelAnalysis(sessionLabel SessionLabel) SessionInfo {
	platform := strings.Split(sessionLabel, ":")[0]
	sessionTypeAndId := strings.Split(sessionLabel, ":")[1]
	sessionType := strings.Split(sessionTypeAndId, "-")[0]
	sessionId := strings.Split(sessionTypeAndId, "-")[1]
	Id, err := strconv.ParseInt(sessionId, 10, 64)
	if err != nil {
		return SessionInfo{}
	}
	return SessionInfo{
		Platform:    platform,
		SessionType: sessionType,
		SessionId:   Id,
	}
}
