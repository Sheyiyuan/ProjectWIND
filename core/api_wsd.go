package core

import (
	"ProjectWIND/database"
	"ProjectWIND/wba"
	"fmt"
)

//database模块
// //数据库部分允许字符串变量的读写操作，允许读取配置项操作

type databaseInfo struct{}

func (dbi *databaseInfo) varSet(app wba.AppInfo, datamap string, unit string, id string, key string, value string) {
	database.Set(app.AppKey.Name, datamap, unit, id, key, value)
}

func (dbi *databaseInfo) SetUserVariable(app wba.AppInfo, msg wba.MessageEventInfo, key string, value string) {
	id := fmt.Sprintf("%d", msg.UserId)
	dbi.varSet(app, app.AppKey.Name, "user", id, key, value)
}

func (dbi *databaseInfo) SetGroupVariable(app wba.AppInfo, msg wba.MessageEventInfo, key string, value string) {
	var id string
	if msg.MessageType == "group" {
		id = "group_" + fmt.Sprintf("%d", msg.GroupId)
	}
	if msg.MessageType == "private" {
		id = "user_" + fmt.Sprintf("%d", msg.UserId)
	}
	dbi.varSet(app, app.AppKey.Name, "group", id, key, value)
}

func (dbi *databaseInfo) SetOutUserVariable(app wba.AppInfo, datamap string, msg wba.MessageEventInfo, key string, value string) {
	id := fmt.Sprintf("%d", msg.UserId)
	dbi.varSet(app, datamap, "user", id, key, value)
}

func (dbi *databaseInfo) SetOutGroupVariable(app wba.AppInfo, datamap string, msg wba.MessageEventInfo, key string, value string) {
	var id string
	if msg.MessageType == "group" {
		id = "group_" + fmt.Sprintf("%d", msg.GroupId)
	}
	if msg.MessageType == "private" {
		id = "user_" + fmt.Sprintf("%d", msg.UserId)
	}
	dbi.varSet(app, datamap, "group", id, key, value)
}

func (dbi *databaseInfo) UnsafelySetUserVariable(app wba.AppInfo, id string, key string, value string) {
	dbi.varSet(app, app.AppKey.Name, "user", id, key, value)
}

func (dbi *databaseInfo) UnsafelySetGroupVariable(app wba.AppInfo, id string, key string, value string) {
	dbi.varSet(app, app.AppKey.Name, "group", id, key, value)
}

func (dbi *databaseInfo) UnsafelySetGlobalVariable(app wba.AppInfo, id string, key string, value string) {
	dbi.varSet(app, app.AppKey.Name, "global", id, key, value)
}

func (dbi *databaseInfo) UnsafelySetOutUserVariable(app wba.AppInfo, datamap string, id string, key string, value string) {
	dbi.varSet(app, datamap, "user", id, key, value)
}

func (dbi *databaseInfo) UnsafelySetOutGroupVariable(app wba.AppInfo, datamap string, id string, key string, value string) {
	dbi.varSet(app, datamap, "group", id, key, value)
}

func (dbi *databaseInfo) UnsafelySetOutGlobalVariable(app wba.AppInfo, datamap string, id string, key string, value string) {
	dbi.varSet(app, datamap, "global", id, key, value)
}

func (dbi *databaseInfo) varGet(app wba.AppInfo, datamap string, unit string, id string, key string) (string, bool) {
	res, ok := database.Get(app.AppKey.Name, datamap, unit, id, key, false)
	if !ok {
		return "", false
	}
	resStr, ok := res.(string)
	if !ok {
		return "", false
	}
	return resStr, true
}

func (dbi *databaseInfo) GetUserVariable(app wba.AppInfo, msg wba.MessageEventInfo, key string) (string, bool) {
	id := fmt.Sprintf("%d", msg.UserId)
	return dbi.varGet(app, app.AppKey.Name, "user", id, key)
}

func (dbi *databaseInfo) GetGroupVariable(app wba.AppInfo, msg wba.MessageEventInfo, key string) (string, bool) {
	var id string
	if msg.MessageType == "group" {
		id = "group_" + fmt.Sprintf("%d", msg.GroupId)
	}
	if msg.MessageType == "private" {
		id = "user_" + fmt.Sprintf("%d", msg.UserId)
	}
	return dbi.varGet(app, app.AppKey.Name, "group", id, key)
}

func (dbi *databaseInfo) GetOutUserVariable(app wba.AppInfo, datamap string, msg wba.MessageEventInfo, key string) (string, bool) {
	id := fmt.Sprintf("%d", msg.UserId)
	return dbi.varGet(app, datamap, "user", id, key)
}

func (dbi *databaseInfo) GetOutGroupVariable(app wba.AppInfo, datamap string, msg wba.MessageEventInfo, key string) (string, bool) {
	var id string
	if msg.MessageType == "group" {
		id = "group_" + fmt.Sprintf("%d", msg.GroupId)
	}
	if msg.MessageType == "private" {
		id = "user_" + fmt.Sprintf("%d", msg.UserId)
	}
	return dbi.varGet(app, datamap, "group", id, key)
}

func (dbi *databaseInfo) UnsafelyGetUserVariable(app wba.AppInfo, id string, key string) (string, bool) {
	return dbi.varGet(app, app.AppKey.Name, "user", id, key)
}

func (dbi *databaseInfo) UnsafelyGetGroupVariable(app wba.AppInfo, id string, key string) (string, bool) {
	return dbi.varGet(app, app.AppKey.Name, "group", id, key)
}

func (dbi *databaseInfo) UnsafelyGetGlobalVariable(app wba.AppInfo, id string, key string) (string, bool) {
	return dbi.varGet(app, app.AppKey.Name, "global", id, key)
}

func (dbi *databaseInfo) UnsafelyGetOutUserVariable(app wba.AppInfo, datamap string, id string, key string) (string, bool) {
	return dbi.varGet(app, datamap, "user", id, key)
}

func (dbi *databaseInfo) UnsafelyGetOutGroupVariable(app wba.AppInfo, datamap string, id string, key string) (string, bool) {
	return dbi.varGet(app, datamap, "group", id, key)
}

func (dbi *databaseInfo) UnsafelyGetOutGlobalVariable(app wba.AppInfo, datamap string, id string, key string) (string, bool) {
	return dbi.varGet(app, datamap, "global", id, key)
}

func (dbi *databaseInfo) GetIntConfig(app wba.AppInfo, datamap string, key string) (int64, bool) {
	res, ok := database.Get(app.AppKey.Name, datamap, "config", "number", key, true)
	if !ok {
		return 0, false
	}
	resInt, ok := res.(int64)
	if !ok {
		return 0, false
	}
	return resInt, true
}

func (dbi *databaseInfo) GetStringConfig(app wba.AppInfo, datamap string, key string) (string, bool) {
	res, ok := database.Get(app.AppKey.Name, datamap, "config", "string", key, true)
	if !ok {
		return "", false
	}
	resStr, ok := res.(string)
	if !ok {
		return "", false
	}
	return resStr, true
}

func (dbi *databaseInfo) GetFloatConfig(app wba.AppInfo, datamap string, key string) (float64, bool) {
	res, ok := database.Get(app.AppKey.Name, datamap, "config", "float", key, true)
	if !ok {
		return 0, false
	}
	resFloat, ok := res.(float64)
	if !ok {
		return 0, false
	}
	return resFloat, true
}

func (dbi *databaseInfo) GetIntSliceConfig(app wba.AppInfo, datamap string, key string) ([]int64, bool) {
	res, ok := database.Get(app.AppKey.Name, datamap, "config", "number_slice", key, true)
	if !ok {
		return nil, false
	}
	resSlice, ok := res.([]int64)
	if !ok {
		return nil, false
	}
	return resSlice, true
}

func (dbi *databaseInfo) GetStringSliceConfig(app wba.AppInfo, datamap string, key string) ([]string, bool) {
	res, ok := database.Get(app.AppKey.Name, datamap, "config", "string_slice", key, true)
	if !ok {
		return nil, false
	}
	resSlice, ok := res.([]string)
	if !ok {
		return nil, false
	}
	return resSlice, true
}

func (dbi *databaseInfo) UnsafelyCreatePublicDatamap(app wba.AppInfo, datamapId string) {
	appName := app.AppKey.Name
	database.CreatePublicDatamap(appName, datamapId)
}

var DatabaseApi databaseInfo
