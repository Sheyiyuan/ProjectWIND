package core

import (
	"ProjectWIND/LOG"
	"ProjectWIND/wba"
	"github.com/dop251/goja"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

type CamelCaseFieldNameMapper struct{}

func (CamelCaseFieldNameMapper) FieldName(_ reflect.Type, f reflect.StructField) string {
	name := f.Name
	if len(name) == 0 {
		return name
	}
	// 首字母小写
	return strings.ToLower(name[:1]) + name[1:]
}

func (CamelCaseFieldNameMapper) MethodName(_ reflect.Type, m reflect.Method) string {
	name := m.Name
	if len(name) == 0 {
		return name
	}
	// 首字母小写
	return strings.ToLower(name[:1]) + name[1:]
}

var GlobalCmdAgentSelector = wba.NewCmdAgentSelector()
var CmdMap = make(map[wba.AppKey]wba.CmdList)
var AppMap = make(map[wba.AppKey]wba.AppInfo)
var ScheduledTaskMap = make(map[wba.AppKey]map[string]wba.ScheduledTaskInfo)

// ReloadApps 重新加载应用
func ReloadApps() (total int, success int) {
	// 清空AppMap和CmdMap
	CmdMap = make(map[wba.AppKey]wba.CmdList)
	AppMap = make(map[wba.AppKey]wba.AppInfo)
	ScheduledTaskMap = make(map[wba.AppKey]map[string]wba.ScheduledTaskInfo)
	GlobalCmdAgentSelector = wba.NewCmdAgentSelector()
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
	CmdMap[AppCore.AppKey] = AppCore.CmdMap
	GlobalCmdAgentSelector.AddCmdMap(CmdMap)
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
		runtime.SetFieldNameMapper(CamelCaseFieldNameMapper{})
		runtime.Set("console", map[string]interface{}{
			"log": func(v ...interface{}) {
				LOG.Info("JS log: %v", v...)
			},
			"error": func(v ...interface{}) {
				LOG.Error("JS error: %v", v...)
			},
		})

		// 添加错误捕获
		safeRun := func(fn func() error) {
			defer func() {
				if r := recover(); r != nil {
					LOG.Error("JS执行错误: %v", r)
				}
			}()
			if err := fn(); err != nil {
				LOG.Error("JS执行错误: %v", err)
			}
		}

		// 修改JS代码执行部分
		safeRun(func() error {
			_, err := runtime.RunString(string(jsCode))
			return err
		})

		// 创建JS可用的wbaObj对象
		wbaObj := runtime.NewObject()
		//wsp := runtime.NewObject()
		//wsd := runtime.NewObject()
		//wst := runtime.NewObject()
		_ = runtime.Set("wba", wbaObj)
		_ = wbaObj.Set("newApp", wba.NewApp)
		_ = wbaObj.Set("withName", wba.WithSelector)
		_ = wbaObj.Set("withDescription", wba.WithDescription)
		_ = wbaObj.Set("withWebUrl", wba.WithWebUrl)
		_ = wbaObj.Set("withLicense", wba.WithLicense)
		_ = wbaObj.Set("wsp", ProtocolApi)
		_ = wbaObj.Set("wsd", DatabaseApi)
		_ = wbaObj.Set("wst", ToolsApi)
		//_ = wbaObj.Set("wsp", wsp)
		//_ = wbaObj.Set("wsd", wsd)
		//_ = wbaObj.Set("wst", wst)
		////WSP注册
		//_ = wsp.Set("unsafelySendMsg", ProtocolApi.UnsafelySendMsg)
		//_ = wsp.Set("unsafelySendPrivateMsg", ProtocolApi.UnsafelySendPrivateMsg)
		//_ = wsp.Set("unsafelySendGroupMsg", ProtocolApi.UnsafelySendGroupMsg)
		//_ = wsp.Set("sendMsg", ProtocolApi.SendMsg)
		//_ = wsp.Set("sendPrivateMsg", ProtocolApi.SendPrivateMsg)
		//_ = wsp.Set("sendGroupMsg", ProtocolApi.SendGroupMsg)
		//_ = wsp.Set("unsafelyDeleteMsg", ProtocolApi.UnsafelyDeleteMsg)
		//_ = wsp.Set("deleteMsg", ProtocolApi.DeleteMsg)
		//_ = wsp.Set("sendLike", ProtocolApi.SendLike)
		//_ = wsp.Set("setGroupKick", ProtocolApi.SetGroupKick)
		//_ = wsp.Set("setGroupBan", ProtocolApi.SetGroupBan)
		//_ = wsp.Set("setGroupWholeBan", ProtocolApi.SetGroupWholeBan)
		//_ = wsp.Set("setGroupAdmin", ProtocolApi.SetGroupAdmin)
		//_ = wsp.Set("setGroupLeave", ProtocolApi.SetGroupLeave)
		//_ = wsp.Set("setGroupCard", ProtocolApi.SetGroupCard)
		//_ = wsp.Set("setGroupName", ProtocolApi.SetGroupName)
		//_ = wsp.Set("setGroupSpecialTitle", ProtocolApi.SetGroupSpecialTitle)
		//_ = wsp.Set("setFriendAddRequest", ProtocolApi.SetFriendAddRequest)
		//_ = wsp.Set("setGroupAddRequest", ProtocolApi.SetGroupAddRequest)
		//_ = wsp.Set("getLoginInfo", ProtocolApi.GetLoginInfo)
		//_ = wsp.Set("getVersionInfo", ProtocolApi.GetVersionInfo)
		//_ = wsp.Set("getMsg", ProtocolApi.GetMsg)
		//_ = wsp.Set("getGroupInfo", ProtocolApi.GetGroupInfo)
		//_ = wsp.Set("getForwardMsg", ProtocolApi.GetForwardMsg)
		//_ = wsp.Set("getStrangerInfo", ProtocolApi.GetStrangerInfo)
		//_ = wsp.Set("getGroupList", ProtocolApi.GetGroupList)
		//_ = wsp.Set("getGroupMemberList", ProtocolApi.GetGroupMemberList)
		//_ = wsp.Set("getFriendList", ProtocolApi.GetFriendList)
		//_ = wsp.Set("getGroupMemberInfo", ProtocolApi.GetGroupMemberInfo)
		//_ = wsp.Set("getGroupHonorInfo", ProtocolApi.GetGroupHonorInfo)
		//_ = wsp.Set("getStatus", ProtocolApi.GetStatus)
		//_ = wsp.Set("getCookies", ProtocolApi.GetCookies)
		//_ = wsp.Set("getCSRFToken", ProtocolApi.GetCSRFToken)
		//_ = wsp.Set("getCredentials", ProtocolApi.GetCredentials)
		//_ = wsp.Set("getImage", ProtocolApi.GetImage)
		//_ = wsp.Set("getRecord", ProtocolApi.GetRecord)
		//_ = wsp.Set("canSendImage", ProtocolApi.CanSendImage)
		//_ = wsp.Set("canSendRecord", ProtocolApi.CanSendRecord)
		//_ = wsp.Set("cetRestart", ProtocolApi.SetRestart)
		//_ = wsp.Set("cleanCache", ProtocolApi.CleanCache)
		//_ = wsp.Set("getVersionInfo", ProtocolApi.GetVersionInfo)
		////WST注册
		//_ = wst.Set("logWith", ToolsApi.LogWith)
		//_ = wst.Set("log", ToolsApi.Log)
		//_ = wst.Set("msgMarshal", ToolsApi.MsgUnmarshal)
		////WSD注册
		//_ = wsd.Set("setUserVariable", DatabaseApi.SetUserVariable)
		//_ = wsd.Set("setGroupVariable", DatabaseApi.SetGroupVariable)
		//_ = wsd.Set("setOutUserVariable", DatabaseApi.SetOutUserVariable)
		//_ = wsd.Set("setOutGroupVariable", DatabaseApi.SetOutGroupVariable)
		//_ = wsd.Set("unsafelySetUserVariable", DatabaseApi.UnsafelySetUserVariable)
		//_ = wsd.Set("unsafelySetGroupVariable", DatabaseApi.UnsafelySetGroupVariable)
		//_ = wsd.Set("unsafelySetGlobalVariable", DatabaseApi.UnsafelySetGlobalVariable)
		//_ = wsd.Set("unsafelySetOutUserVariable", DatabaseApi.UnsafelySetOutUserVariable)
		//_ = wsd.Set("unsafelySetOutGroupVariable", DatabaseApi.UnsafelySetOutGroupVariable)
		//_ = wsd.Set("unsafelySetOutGlobalVariable", DatabaseApi.UnsafelySetOutGlobalVariable)
		//_ = wsd.Set("getUserVariable", DatabaseApi.GetUserVariable)
		//_ = wsd.Set("getGroupVariable", DatabaseApi.GetGroupVariable)
		//_ = wsd.Set("getOutUserVariable", DatabaseApi.GetOutUserVariable)
		//_ = wsd.Set("getOutGroupVariable", DatabaseApi.GetOutGroupVariable)
		//_ = wsd.Set("unsafelyGetUserVariable", DatabaseApi.UnsafelyGetUserVariable)
		//_ = wsd.Set("unsafelyGetGroupVariable", DatabaseApi.UnsafelyGetGroupVariable)
		//_ = wsd.Set("unsafelyGetGlobalVariable", DatabaseApi.UnsafelyGetGlobalVariable)
		//_ = wsd.Set("unsafelyGetOutUserVariable", DatabaseApi.UnsafelyGetOutUserVariable)
		//_ = wsd.Set("unsafelyGetOutGroupVariable", DatabaseApi.UnsafelyGetOutGroupVariable)
		//_ = wsd.Set("unsafelyGetOutGlobalVariable", DatabaseApi.UnsafelyGetOutGlobalVariable)
		//_ = wsd.Set("getIntConfig", DatabaseApi.GetIntConfig)
		//_ = wsd.Set("getFloatConfig", DatabaseApi.GetFloatConfig)
		//_ = wsd.Set("getStringConfig", DatabaseApi.GetStringConfig)
		//_ = wsd.Set("getIntSliceConfig", DatabaseApi.GetIntSliceConfig)
		//_ = wsd.Set("getStringSliceConfig", DatabaseApi.GetStringSliceConfig)
		//_ = wsd.Set("unsafelyCreatePublicDatamap", DatabaseApi.UnsafelyCreatePublicDatamap)

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

		//// 调用Init方法
		//initVal := jsApp.ToObject(runtime).Get("Init")
		//initFunc, ok := goja.AssertFunction(initVal)
		//if !ok {
		//	LOG.Error("应用 %s 缺少有效的Init方法 %#v", pluginPath, initFunc)
		//	return 1, 0
		//}
		//
		//_, err = initFunc(wbaObj)
		//if err != nil {
		//	LOG.Trace("应用初始化失败: %v", err)
		//	return 1, 0
		//}
		//
		//// 调用Get方法
		//getVal := jsApp.ToObject(runtime).Get("Get")
		//getFunc, ok := goja.AssertFunction(getVal)
		//if !ok {
		//	LOG.Error("应用 %s 缺少有效的Get方法", pluginPath)
		//	return 1, 0
		//}
		//
		//appInfoVal, err := getFunc(jsApp)
		//if err != nil {
		//	LOG.Error("获取应用信息失败: %v", err)
		//	return 1, 0
		//}

		// 转换应用信息
		var appInfo wba.AppInfo
		if err := runtime.ExportTo(jsApp, &appInfo); err != nil {
			LOG.Error("应用信息转换失败: %v", err)
			return 1, 0
		}
		// 初始化map字段
		if appInfo.CmdMap == nil {
			appInfo.CmdMap = make(map[string]wba.Cmd)
		}
		if appInfo.ScheduledTasks == nil {
			appInfo.ScheduledTasks = make(map[string]wba.ScheduledTaskInfo)
		}

		AppMap[appInfo.AppKey] = appInfo

		CmdMap[appInfo.AppKey] = appInfo.CmdMap

		ScheduledTaskMap[appInfo.AppKey] = appInfo.ScheduledTasks

		// 注册定时任务
		for _, task := range appInfo.ScheduledTasks {
			taskCopy := task
			RegisterCron(appInfo.AppKey.Name, wba.ScheduledTaskInfo{
				Name: taskCopy.Name,
				Desc: taskCopy.Desc,
				Cron: taskCopy.Cron,
				Task: func() {
					safeRun(func() error {
						taskCopy.Task()
						return nil
					})
				},
			})
		}

		LOG.Info("JS应用 %s 加载成功", pluginPath)
		return 1, 1
	}
	return 0, 0
}
