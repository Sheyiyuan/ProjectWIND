package wba

import (
	"fmt"
)

type WindStandardTools interface {
	// MsgUnmarshal 解析消息JSON字符串为 MessageEventInfo 结构体。
	//
	// 参数:
	//
	//  - message: 要解析的消息JSON字符串。
	//
	// 返回值:
	//
	//  - msg: 解析后的消息结构体。
	MsgUnmarshal(message string) (msg MessageEventInfo)

	// LogWith 使用指定日志级别记录日志，支持可变参数占位符。
	//
	// 参数:
	//  - level: 日志级别: "trace", "debug", "info", "notice", "warn", "error"。
	//  - log: 日志内容。
	//  - args: 可变参数，用于格式化日志内容。
	LogWith(level string, log string, args ...interface{})

	// Log 记录日志，级别为 "info"，支持可变参数占位符。
	//
	// 参数:
	//  - log: 日志内容。
	//  - args: 可变参数，用于格式化日志内容。
	Log(log string, args ...interface{})

	// VersionLabelAnalysis 解析版本标签为 VersionInfo 结构体。
	//
	// 参数:
	//  - versionLabel: 版本标签字符串，格式为 "大版本.小版本.修复版本"。
	// 返回值:
	//  - versionInfo: 解析后的版本信息结构体。
	VersionLabelAnalysis(versionLabel VersionLabel) (versionInfo VersionInfo)

	// VersionCompare 比较两个版本标签的大小。
	//
	// 参数:
	//  - version1: 第一个版本标签字符串。
	//  - version2: 第二个版本标签字符串。
	// 返回值:
	//  - result: 比较结果，-1表示version1小于version2，0表示version1等于version2，1表示version1大于version2。
	VersionCompare(version1, version2 VersionLabel) (result int)

	// SessionLabelAnalysis 解析会话标签为 SessionInfo 结构体。
	//
	// 参数:
	//  - sessionLabel: 会话标签字符串，格式为 "平台:会话类型-会话ID"。
	// 返回值:
	//  - sessionInfo: 解析后的会话信息结构体。
	SessionLabelAnalysis(sessionLabel SessionLabel) (sessionInfo SessionInfo)
}

func (v VersionInfo) String() string {
	return fmt.Sprintf("%d.%d.%d", v.BigVersion, v.SmallVersion, v.FixVersion)
}
