package core

import (
	"ProjectWIND/database"
	"ProjectWIND/wba"
	"fmt"
)

// database模块
// 数据库部分允许字符串变量的读写操作，允许读取配置项操作

type databaseInfo struct{}

func (dbi *databaseInfo) varSet(app wba.AppInfo, datamap string, unit string, id string, key string, value string) (errno database.Errno) {
	return database.Set(app.AppKey.Name, datamap, unit, id, key, value)
}

func (dbi *databaseInfo) SetUserVariable(app wba.AppInfo, msg wba.MessageEventInfo, key string, value string) (isSuccess bool) {
	id := fmt.Sprintf("%d", msg.UserId)
	eno := dbi.varSet(app, app.AppKey.Name, "user", id, key, value)
	if eno.Code != 0 {
		eno.Log()
		return false
	}
	return true
}

func (dbi *databaseInfo) SetGroupVariable(app wba.AppInfo, msg wba.MessageEventInfo, key string, value string) (isSuccess bool) {
	var id string
	if msg.MessageType == "group" {
		id = "group_" + fmt.Sprintf("%d", msg.GroupId)
	}
	if msg.MessageType == "private" {
		id = "user_" + fmt.Sprintf("%d", msg.UserId)
	}
	eno := dbi.varSet(app, app.AppKey.Name, "group", id, key, value)
	if eno.Code != 0 {
		eno.Log()
		return false
	}
	return true
}

func (dbi *databaseInfo) SetOutUserVariable(app wba.AppInfo, datamap string, msg wba.MessageEventInfo, key string, value string) (isSuccess bool) {
	id := fmt.Sprintf("%d", msg.UserId)
	eno := dbi.varSet(app, datamap, "user", id, key, value)
	if eno.Code != 0 {
		eno.Log()
		return false
	}
	return true
}

func (dbi *databaseInfo) SetOutGroupVariable(app wba.AppInfo, datamap string, msg wba.MessageEventInfo, key string, value string) (isSuccess bool) {
	var id string
	if msg.MessageType == "group" {
		id = "group_" + fmt.Sprintf("%d", msg.GroupId)
	}
	if msg.MessageType == "private" {
		id = "user_" + fmt.Sprintf("%d", msg.UserId)
	}
	eno := dbi.varSet(app, datamap, "group", id, key, value)
	if eno.Code != 0 {
		eno.Log()
		return false
	}
	return true
}

func (dbi *databaseInfo) UnsafelySetUserVariable(app wba.AppInfo, id string, key string, value string) (isSuccess bool) {
	eno := dbi.varSet(app, app.AppKey.Name, "user", id, key, value)
	if eno.Code != 0 {
		eno.Log()
		return false
	}
	return true
}

func (dbi *databaseInfo) UnsafelySetGroupVariable(app wba.AppInfo, id string, key string, value string) (isSuccess bool) {
	eno := dbi.varSet(app, app.AppKey.Name, "group", id, key, value)
	if eno.Code != 0 {
		eno.Log()
		return false
	}
	return true
}

func (dbi *databaseInfo) UnsafelySetGlobalVariable(app wba.AppInfo, id string, key string, value string) (isSuccess bool) {
	eno := dbi.varSet(app, app.AppKey.Name, "global", id, key, value)
	if eno.Code != 0 {
		eno.Log()
		return false
	}
	return true
}

func (dbi *databaseInfo) UnsafelySetOutUserVariable(app wba.AppInfo, datamap string, id string, key string, value string) (isSuccess bool) {
	eno := dbi.varSet(app, datamap, "user", id, key, value)
	if eno.Code != 0 {
		eno.Log()
		return false
	}
	return true
}

func (dbi *databaseInfo) UnsafelySetOutGroupVariable(app wba.AppInfo, datamap string, id string, key string, value string) (isSuccess bool) {
	eno := dbi.varSet(app, datamap, "group", id, key, value)
	if eno.Code != 0 {
		eno.Log()
		return false
	}
	return true
}

func (dbi *databaseInfo) UnsafelySetOutGlobalVariable(app wba.AppInfo, datamap string, id string, key string, value string) (isSuccess bool) {
	eno := dbi.varSet(app, datamap, "global", id, key, value)
	if eno.Code != 0 {
		eno.Log()
		return false
	}
	return true
}

func (dbi *databaseInfo) varGet(app wba.AppInfo, datamap string, unit string, id string, key string) (value string, errno database.Errno) {
	res, eno := database.Get(app.AppKey.Name, datamap, unit, id, key, false)
	if eno.Code != 0 {
		return "", eno
	}
	resStr, ok := res.(string)
	if !ok {
		return "", database.Errno{Code: 303, Condition: "获取数据时类型断言错误"}
	}
	return resStr, eno
}

func (dbi *databaseInfo) GetUserVariable(app wba.AppInfo, msg wba.MessageEventInfo, key string) (value string, isSuccess bool) {
	id := fmt.Sprintf("%d", msg.UserId)
	res, eno := dbi.varGet(app, app.AppKey.Name, "user", id, key)
	if eno.Code != 0 {
		eno.Log()
		return "", false
	}
	if res == "" {
		return res, false
	}
	return res, true
}

func (dbi *databaseInfo) GetGroupVariable(app wba.AppInfo, msg wba.MessageEventInfo, key string) (vaule string, isSuccess bool) {
	var id string
	if msg.MessageType == "group" {
		id = "group_" + fmt.Sprintf("%d", msg.GroupId)
	}
	if msg.MessageType == "private" {
		id = "user_" + fmt.Sprintf("%d", msg.UserId)
	}
	res, eno := dbi.varGet(app, app.AppKey.Name, "group", id, key)
	if eno.Code != 0 {
		eno.Log()
		return "", false
	}
	if res == "" {
		return res, false
	}
	return res, true
}

func (dbi *databaseInfo) GetOutUserVariable(app wba.AppInfo, datamap string, msg wba.MessageEventInfo, key string) (vaule string, isSuccess bool) {
	id := fmt.Sprintf("%d", msg.UserId)
	res, eno := dbi.varGet(app, datamap, "user", id, key)
	if eno.Code != 0 {
		eno.Log()
		return "", false
	}
	if res == "" {
		return res, false
	}
	return res, true
}

func (dbi *databaseInfo) GetOutGroupVariable(app wba.AppInfo, datamap string, msg wba.MessageEventInfo, key string) (vaule string, isSuccess bool) {
	var id string
	if msg.MessageType == "group" {
		id = "group_" + fmt.Sprintf("%d", msg.GroupId)
	}
	if msg.MessageType == "private" {
		id = "user_" + fmt.Sprintf("%d", msg.UserId)
	}
	res, eno := dbi.varGet(app, datamap, "group", id, key)
	if eno.Code != 0 {
		eno.Log()
		return "", false
	}
	if res == "" {
		return res, false
	}
	return res, true
}

func (dbi *databaseInfo) UnsafelyGetUserVariable(app wba.AppInfo, id string, key string) (vaule string, isSuccess bool) {
	res, eno := dbi.varGet(app, app.AppKey.Name, "user", id, key)
	if eno.Code != 0 {
		eno.Log()
		return "", false
	}
	if res == "" {
		return res, false
	}
	return res, true
}

func (dbi *databaseInfo) UnsafelyGetGroupVariable(app wba.AppInfo, id string, key string) (vaule string, isSuccess bool) {
	res, eno := dbi.varGet(app, app.AppKey.Name, "group", id, key)
	if eno.Code != 0 {
		eno.Log()
		return "", false
	}
	if res == "" {
		return res, false
	}
	return res, true
}

func (dbi *databaseInfo) UnsafelyGetGlobalVariable(app wba.AppInfo, id string, key string) (vaule string, isSuccess bool) {
	res, eno := dbi.varGet(app, app.AppKey.Name, "global", id, key)
	if eno.Code != 0 {
		eno.Log()
		return "", false
	}
	if res == "" {
		return res, false
	}
	return res, true
}

func (dbi *databaseInfo) UnsafelyGetOutUserVariable(app wba.AppInfo, datamap string, id string, key string) (vaule string, isSuccess bool) {
	res, eno := dbi.varGet(app, datamap, "user", id, key)
	if eno.Code != 0 {
		eno.Log()
		return "", false
	}
	if res == "" {
		return res, false
	}
	return res, true
}

func (dbi *databaseInfo) UnsafelyGetOutGroupVariable(app wba.AppInfo, datamap string, id string, key string) (vaule string, isSuccess bool) {
	res, eno := dbi.varGet(app, datamap, "group", id, key)
	if eno.Code != 0 {
		eno.Log()
		return "", false
	}
	if res == "" {
		return res, false
	}
	return res, true
}

func (dbi *databaseInfo) UnsafelyGetOutGlobalVariable(app wba.AppInfo, datamap string, id string, key string) (vaule string, isSuccess bool) {
	res, eno := dbi.varGet(app, datamap, "global", id, key)
	if eno.Code != 0 {
		eno.Log()
		return "", false
	}
	if res == "" {
		return res, false
	}
	return res, true
}

func (dbi *databaseInfo) GetIntConfig(app wba.AppInfo, datamap string, key string) (value int64, isSuccess bool) {
	res, eno := database.Get(app.AppKey.Name, datamap, "config", "number", key, true)
	if eno.Code != 0 {
		return 0, false
	}
	resInt, ok := res.(int64)
	if !ok {
		eno = database.Errno{Code: 303, Condition: "获取配置时类型断言错误"}
		eno.Log()
		return 0, false
	}
	return resInt, true
}

func (dbi *databaseInfo) GetStringConfig(app wba.AppInfo, datamap string, key string) (value string, isSuccess bool) {
	res, eno := database.Get(app.AppKey.Name, datamap, "config", "string", key, true)
	if eno.Code != 0 {
		eno.Log()
		return "", false
	}
	resStr, ok := res.(string)
	if !ok {
		eno := database.Errno{Code: 303, Condition: "获取配置时类型断言错误"}
		eno.Log()
		return "", false
	}
	return resStr, true
}

func (dbi *databaseInfo) GetFloatConfig(app wba.AppInfo, datamap string, key string) (value float64, isSuccess bool) {
	res, eno := database.Get(app.AppKey.Name, datamap, "config", "float", key, true)
	if eno.Code != 0 {
		eno.Log()
		return 0, false
	}
	resFloat, ok := res.(float64)
	if !ok {
		eno := database.Errno{Code: 303, Condition: "获取配置时类型断言错误"}
		eno.Log()
		return 0, false
	}
	return resFloat, true
}

func (dbi *databaseInfo) GetIntSliceConfig(app wba.AppInfo, datamap string, key string) (value []int64, isSuccess bool) {
	res, eno := database.Get(app.AppKey.Name, datamap, "config", "number_slice", key, true)
	if eno.Code != 0 {
		eno.Log()
		return nil, false
	}
	resSlice, ok := res.([]int64)
	if !ok {
		eno := database.Errno{Code: 303, Condition: "获取配置时类型断言错误"}
		eno.Log()
		return nil, false
	}
	return resSlice, true
}

func (dbi *databaseInfo) GetStringSliceConfig(app wba.AppInfo, datamap string, key string) (value []string, isSuccess bool) {
	res, eno := database.Get(app.AppKey.Name, datamap, "config", "string_slice", key, true)
	if eno.Code != 0 {
		eno.Log()
		return nil, false
	}
	resSlice, ok := res.([]string)
	if !ok {
		eno := database.Errno{Code: 303, Condition: "获取配置时类型断言错误"}
		eno.Log()
		return nil, false
	}
	return resSlice, true
}

func (dbi *databaseInfo) UnsafelyCreatePublicDatamap(app wba.AppInfo, datamapId string) (isSuccess bool) {
	appName := app.AppKey.Name
	eno := database.CreatePublicDatamap(appName, datamapId)
	if eno.Code != 0 {
		eno.Log()
		return false
	}
	return true
}

var DatabaseApi databaseInfo
