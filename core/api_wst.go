package core

import (
	"ProjectWIND/LOG"
	"ProjectWIND/wba"
	"encoding/json"
	"strconv"
	"strings"
)

type toolsAPI struct{}

//二、LOG模块

/*
关于LOG模块的说明

1.日志模块级别分为TRACE、DEBUG、INFO、NOTICE、WARN、ERROR五个级别，默认级别为INFO。

2.日志模块提供LogWith方法，可以自定义日志级别。

3.日志模块提供Log方法，默认日志级别为INFO。
*/

// LogWith 打印日志(带级别)
//
// 参数：
//   - level: 日志级别，支持"TRACE"、"DEBUG"、"INFO"、"NOTICE"、 "WARN"、"ERROR"
//   - content: 日志内容
//   - args: 可选参数，用于格式化日志内容
//
// 返回值：
//   - 无
func (t *toolsAPI) LogWith(level string, content string, args ...interface{}) {
	level = strings.ToLower(level)
	switch level {
	case "trace":
		LOG.Trace(content, args...)
		return
	case "debug":
		LOG.Debug(content, args...)
		return
	case "notice":
		LOG.Notice(content, args...)
		return
	case "warn":
		LOG.Warn(content, args...)
		return
	case "error":
		LOG.Error(content, args...)
		return
	default:
		LOG.Info(content, args...)
		return
	}
}

// Log 打印日志
//
// 参数：
//   - content: 日志内容
//   - args: 可选参数，用于格式化日志内容
func (t *toolsAPI) Log(content string, args ...interface{}) {
	LOG.Info(content, args...)
}

// MsgUnmarshal 解析消息
//
// 参数：
//   - messageJSON: 从数据库中获取的JSON序列化后的消息字符串
//
// 返回值：
//   - wba.MessageEventInfo: 解析后的消息事件信息，解析失败时返回空结构体
func (t *toolsAPI) MsgUnmarshal(messageJSON string) (msg wba.MessageEventInfo) {
	err := json.Unmarshal([]byte(messageJSON), &msg)
	if err != nil {
		return wba.MessageEventInfo{}
	}
	return msg
}

func (t *toolsAPI) SessionLabelAnalysis(sessionLabel wba.SessionLabel) wba.SessionInfo {
	platform := strings.Split(sessionLabel, ":")[0]
	sessionTypeAndId := strings.Split(sessionLabel, ":")[1]
	sessionType := strings.Split(sessionTypeAndId, "-")[0]
	sessionId := strings.Split(sessionTypeAndId, "-")[1]
	Id, err := strconv.ParseInt(sessionId, 10, 64)
	if err != nil {
		return wba.SessionInfo{}
	}
	return wba.SessionInfo{
		Platform:    platform,
		SessionType: sessionType,
		SessionId:   Id,
	}
}

func (t *toolsAPI) VersionCompare(version1, version2 wba.VersionLabel) int {
	version1Info := VersionLabelAnalysis(version1)
	version2Info := VersionLabelAnalysis(version2)
	if version1Info.BigVersion < version2Info.BigVersion {
		return -1
	} else if version1Info.BigVersion > version2Info.BigVersion {
		return 1
	} else {
		if version1Info.SmallVersion < version2Info.SmallVersion {
			return -1
		} else if version1Info.SmallVersion > version2Info.SmallVersion {
			return 1
		} else {
			if version1Info.FixVersion < version2Info.FixVersion {
				return -1
			} else if version1Info.FixVersion > version2Info.FixVersion {
				return 1
			} else {
				return 0
			}
		}
	}
}

func VersionLabelAnalysis(versionLabel wba.VersionLabel) wba.VersionInfo {
	version := strings.Split(versionLabel, ".")
	bigVersion, err := strconv.ParseUint(version[0], 10, 8)
	if err != nil {
		return wba.VersionInfo{}
	}
	smallVersion, err := strconv.ParseUint(version[1], 10, 8)
	if err != nil {
		return wba.VersionInfo{}
	}
	fixVersion, err := strconv.ParseUint(version[2], 10, 8)
	if err != nil {
		return wba.VersionInfo{}
	}
	return wba.VersionInfo{
		BigVersion:   uint8(bigVersion),
		SmallVersion: uint8(smallVersion),
		FixVersion:   uint8(fixVersion),
	}
}

var ToolsApi toolsAPI
