package core

import (
	"ProjectWIND/LOG"
	"ProjectWIND/wba"
	"fmt"
	"github.com/dop251/goja"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
)

var runtimePool = sync.Pool{
	New: func() interface{} {
		rt := goja.New()
		rt.SetFieldNameMapper(CamelCaseFieldNameMapper{})
		_ = rt.Set("console", map[string]interface{}{
			"log": func(v ...interface{}) {
				LOG.Info("JS log: %v", v...)
			},
			"error": func(v ...interface{}) {
				LOG.Error("JS error: %v", v...)
			},
		})
		// 创建JS可用的wbaObj对象
		windObj := rt.NewObject()
		_ = rt.Set("wind", windObj)
		_ = windObj.Set("newApp", wba.NewApp)
		_ = windObj.Set("withName", wba.WithSelector)
		_ = windObj.Set("withDescription", wba.WithDescription)
		_ = windObj.Set("withWebUrl", wba.WithWebUrl)
		_ = windObj.Set("withLicense", wba.WithLicense)
		_ = rt.Set("wsp", ProtocolApi)
		_ = rt.Set("wsd", DatabaseApi)
		_ = rt.Set("wst", ToolsApi)
		return rt
	},
}

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

		runtime := runtimePool.Get().(*goja.Runtime)
		defer runtimePool.Put(runtime)

		// 重置runtime状态
		runtime.ClearInterrupt()
		runtime.SetRandSource(nil)

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

		// 获取app对象
		var jsApp goja.Value
		func() {
			defer func() {
				if r := recover(); r != nil {
					LOG.Error("获取app对象时发生panic: %v", r)
				}
			}()
			jsApp = runtime.Get("app")
		}()

		if jsApp == nil || goja.IsUndefined(jsApp) {
			LOG.Error("应用 %s 缺少app对象", pluginPath)
			return 1, 0
		}

		// 应用信息转换
		var appInfo wba.AppInfo
		err = func() (err error) {
			defer func() {
				if r := recover(); r != nil {
					err = fmt.Errorf("应用信息转换时发生panic: %v", r)
				}
			}()
			return runtime.ExportTo(jsApp, &appInfo)
		}()
		if err != nil {
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
