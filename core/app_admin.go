package core

import (
	"ProjectWIND/LOG"
	"ProjectWIND/wba"
	"os"
	"path/filepath"
	"strings"

	"github.com/dop251/goja"
)

var CmdMap = make([]map[string]wba.Cmd, 4)
var AppMap = make(map[wba.AppKey]wba.AppInfo)

// ReloadApps 重新加载应用
func ReloadApps() (total int, success int) {
	// 清空AppMap和CmdMap
	CmdMap = make([]map[string]wba.Cmd, 4)
	AppMap = make(map[wba.AppKey]wba.AppInfo)
	appsDir := "./data/app/"
	appFiles, err := os.ReadDir(appsDir)
	total = 0
	success = 0
	if err != nil {
		LOG.Error("加载应用所在目录失败:%v", err)
		return
	}

	for _, file := range appFiles {
		totalDelta, successDelta := reloadAPP(file, appsDir)
		total += totalDelta
		success += successDelta
	}
	CmdMap[0] = AppCore.CmdMap
	return total, success
}

// reloadAPP 重新加载单个应用
func reloadAPP(file os.DirEntry, appsDir string) (totalDelta int, successDelta int) {
	if file.IsDir() {
		return 0, 0
	}

	ext := filepath.Ext(file.Name())
	if ext == ".js" {
		pluginPath := filepath.Join(appsDir, file.Name())
		jsCode, err := os.ReadFile(pluginPath)
		if err != nil {
			LOG.Error("读取应用 %s 时发生错误: %v", pluginPath, err)
			return 1, 0
		}

		runtime := goja.New()
		_, err = runtime.RunString(string(jsCode))
		if err != nil {
			LOG.Error("执行应用 %s 失败: %v", pluginPath, err)
			return 1, 0
		}

		// 创建JS可用的wbaObj对象
		wbaObj := runtime.NewObject()
		wsp := runtime.NewObject()
		wsd := runtime.NewObject()
		wst := runtime.NewObject()
		_ = runtime.Set("wba", wbaObj)
		_ = wbaObj.Set("NewApp", wba.NewApp)
		_ = wbaObj.Set("WithName", wba.WithName)
		_ = wbaObj.Set("WithAuthor", wba.WithAuthor)
		_ = wbaObj.Set("WithVersion", wba.WithVersion)
		_ = wbaObj.Set("WithDescription", wba.WithDescription)
		_ = wbaObj.Set("WithWebUrl", wba.WithWebUrl)
		_ = wbaObj.Set("WithLicense", wba.WithLicense)
		_ = wbaObj.Set("WithAppType", wba.WithAppType)
		_ = wbaObj.Set("WithRule", wba.WithRule)
		_ = wbaObj.Set("wsp", wsp)
		_ = wbaObj.Set("wsd", wsd)
		_ = wbaObj.Set("wst", wst)
		//WSP注册
		_ = wsp.Set("UnsafelySendMsg", AppApi.UnsafelySendMsg)
		_ = wsp.Set("UnsafelySendPrivateMsg", AppApi.UnsafelySendPrivateMsg)
		_ = wsp.Set("UnsafelySendGroupMsg", AppApi.UnsafelySendGroupMsg)
		_ = wsp.Set("SendMsg", AppApi.SendMsg)
		_ = wsp.Set("SendPrivateMsg", AppApi.SendPrivateMsg)
		_ = wsp.Set("SendGroupMsg", AppApi.SendGroupMsg)
		_ = wsp.Set("UnsafelyDeleteMsg", AppApi.UnsafelyDeleteMsg)
		_ = wsp.Set("DeleteMsg", AppApi.DeleteMsg)
		_ = wsp.Set("SendLike", AppApi.SendLike)
		_ = wsp.Set("SetGroupKick", AppApi.SetGroupKick)
		_ = wsp.Set("SetGroupBan", AppApi.SetGroupBan)
		_ = wsp.Set("SetGroupWholeBan", AppApi.SetGroupWholeBan)
		_ = wsp.Set("SetGroupAdmin", AppApi.SetGroupAdmin)
		_ = wsp.Set("SetGroupLeave", AppApi.SetGroupLeave)
		_ = wsp.Set("SetGroupCard", AppApi.SetGroupCard)
		_ = wsp.Set("SetGroupName", AppApi.SetGroupName)
		_ = wsp.Set("SetGroupSpecialTitle", AppApi.SetGroupSpecialTitle)
		_ = wsp.Set("SetFriendAddRequest", AppApi.SetFriendAddRequest)
		_ = wsp.Set("SetGroupAddRequest", AppApi.SetGroupAddRequest)
		_ = wsp.Set("GetLoginInfo", AppApi.GetLoginInfo)
		_ = wsp.Set("GetVersionInfo", AppApi.GetVersionInfo)
		_ = wsp.Set("GetMsg", AppApi.GetMsg)
		_ = wsp.Set("GetGroupInfo", AppApi.GetGroupInfo)
		_ = wsp.Set("GetForwardMsg", AppApi.GetForwardMsg)
		_ = wsp.Set("GetStrangerInfo", AppApi.GetStrangerInfo)
		_ = wsp.Set("GetGroupList", AppApi.GetGroupList)
		_ = wsp.Set("GetGroupMemberList", AppApi.GetGroupMemberList)
		_ = wsp.Set("GetFriendList", AppApi.GetFriendList)
		_ = wsp.Set("GetGroupMemberInfo", AppApi.GetGroupMemberInfo)
		_ = wsp.Set("GetGroupHonorInfo", AppApi.GetGroupHonorInfo)
		_ = wsp.Set("GetStatus", AppApi.GetStatus)
		_ = wsp.Set("GetCookies", AppApi.GetCookies)
		_ = wsp.Set("GetCSRFToken", AppApi.GetCSRFToken)
		_ = wsp.Set("GetCredentials", AppApi.GetCredentials)
		_ = wsp.Set("GetImage", AppApi.GetImage)
		_ = wsp.Set("GetRecord", AppApi.GetRecord)
		_ = wsp.Set("CanSendImage", AppApi.CanSendImage)
		_ = wsp.Set("CanSendRecord", AppApi.CanSendRecord)
		_ = wsp.Set("SetRestart", AppApi.SetRestart)
		_ = wsp.Set("CleanCache", AppApi.CleanCache)
		_ = wsp.Set("GetVersionInfo", AppApi.GetVersionInfo)
		//WST注册
		_ = wst.Set("LogWith", AppApi.LogWith)
		_ = wst.Set("Log", AppApi.Log)
		_ = wst.Set("MsgMarshal", AppApi.MsgUnmarshal)
		//WSD注册
		_ = wsd.Set("SetUserVariable", DatabaseApi.SetUserVariable)
		_ = wsd.Set("SetGroupVariable", DatabaseApi.SetGroupVariable)
		_ = wsd.Set("SetOutUserVariable", DatabaseApi.SetOutUserVariable)
		_ = wsd.Set("SetOutGroupVariable", DatabaseApi.SetOutGroupVariable)
		_ = wsd.Set("UnsafelySetUserVariable", DatabaseApi.UnsafelySetUserVariable)
		_ = wsd.Set("UnsafelySetGroupVariable", DatabaseApi.UnsafelySetGroupVariable)
		_ = wsd.Set("UnsafelySetGlobalVariable", DatabaseApi.UnsafelySetGlobalVariable)
		_ = wsd.Set("UnsafelySetOutUserVariable", DatabaseApi.UnsafelySetOutUserVariable)
		_ = wsd.Set("UnsafelySetOutGroupVariable", DatabaseApi.UnsafelySetOutGroupVariable)
		_ = wsd.Set("UnsafelySetOutGlobalVariable", DatabaseApi.UnsafelySetOutGlobalVariable)
		_ = wsd.Set("GetUserVariable", DatabaseApi.GetUserVariable)
		_ = wsd.Set("GetGroupVariable", DatabaseApi.GetGroupVariable)
		_ = wsd.Set("GetOutUserVariable", DatabaseApi.GetOutUserVariable)
		_ = wsd.Set("GetOutGroupVariable", DatabaseApi.GetOutGroupVariable)
		_ = wsd.Set("UnsafelyGetUserVariable", DatabaseApi.UnsafelyGetUserVariable)
		_ = wsd.Set("UnsafelyGetGroupVariable", DatabaseApi.UnsafelyGetGroupVariable)
		_ = wsd.Set("UnsafelyGetGlobalVariable", DatabaseApi.UnsafelyGetGlobalVariable)
		_ = wsd.Set("UnsafelyGetOutUserVariable", DatabaseApi.UnsafelyGetOutUserVariable)
		_ = wsd.Set("UnsafelyGetOutGroupVariable", DatabaseApi.UnsafelyGetOutGroupVariable)
		_ = wsd.Set("UnsafelyGetOutGlobalVariable", DatabaseApi.UnsafelyGetOutGlobalVariable)
		_ = wsd.Set("GetIntConfig", DatabaseApi.GetIntConfig)
		_ = wsd.Set("GetFloatConfig", DatabaseApi.GetFloatConfig)
		_ = wsd.Set("GetStringConfig", DatabaseApi.GetStringConfig)
		_ = wsd.Set("GetIntSliceConfig", DatabaseApi.GetIntSliceConfig)
		_ = wsd.Set("GetStringSliceConfig", DatabaseApi.GetStringSliceConfig)
		_ = wsd.Set("UnsafelyCreatePublicDatamap", DatabaseApi.UnsafelyCreatePublicDatamap)

		// 获取AppInit函数
		appInitVal := runtime.Get("AppInit")
		if appInitVal == nil || goja.IsUndefined(appInitVal) {
			LOG.Error("应用 %s 缺少AppInit函数", pluginPath)
			return 1, 0
		}
		appInitFunc, ok := goja.AssertFunction(appInitVal)
		if !ok {
			LOG.Error("应用 %s 的AppInit不是有效函数", pluginPath)
			return 1, 0
		}

		// 调用AppInit获取应用实例
		jsApp, err := appInitFunc(goja.Undefined())
		if err != nil {
			LOG.Error("初始化应用实例失败: %v", err)
			return 1, 0
		}

		// 调用Init方法
		initVal := jsApp.ToObject(runtime).Get("Init")
		initFunc, ok := goja.AssertFunction(initVal)
		if !ok {
			LOG.Error("应用 %s 缺少有效的Init方法", pluginPath)
			return 1, 0
		}

		_, err = initFunc(wbaObj)
		if err != nil {
			LOG.Trace("应用初始化失败: %v", err)
			return 1, 0
		}

		// 调用Get方法
		getVal := jsApp.ToObject(runtime).Get("Get")
		getFunc, ok := goja.AssertFunction(getVal)
		if !ok {
			LOG.Error("应用 %s 缺少有效的Get方法", pluginPath)
			return 1, 0
		}

		appInfoVal, err := getFunc(jsApp)
		if err != nil {
			LOG.Error("获取应用信息失败: %v", err)
			return 1, 0
		}

		// 转换应用信息
		var appInfo wba.AppInfo
		if err := runtime.ExportTo(appInfoVal, &appInfo); err != nil {
			LOG.Error("应用信息转换失败: %v", err)
			return 1, 0
		}

		AppMap[wba.AppKey{Name: appInfo.AppKey.Name, Type: appInfo.AppKey.Type, Version: appInfo.AppKey.Version, Level: checkAppLevel(appInfo)}] = appInfo
		cmdIndex := AppTypeToInt(appInfo.AppKey.Type)
		// 合并命令
		CmdMap[cmdIndex] = mergeMaps(CmdMap[cmdIndex], appInfo.CmdMap)

		// 注册定时任务
		for _, task := range appInfo.ScheduledTasks {
			RegisterCron(appInfo.AppKey.Name, task)
		}

		LOG.Info("JS应用 %s 加载成功", pluginPath)
		return 1, 1
	}
	return 0, 0
}

func mergeMaps(map1, map2 map[string]wba.Cmd) map[string]wba.Cmd {
	// 合并map1和map2到map3中
	map3 := make(map[string]wba.Cmd)
	for key, value := range map1 {
		map3[key] = value
	}
	for key, value := range map2 {
		map3[key] = value
	}
	return map3
}

func AppTypeToInt(appType string) int32 {
	appType = strings.ToLower(appType)
	switch appType {
	case "system":
		return 1
	case "rule":
		return 2
	default:
		return 3
	}
}

func checkAppLevel(appInfo wba.AppInfo) int32 {
	return 0
}
