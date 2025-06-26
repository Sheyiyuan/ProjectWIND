package core

import (
	"ProjectWIND/LOG"
	"ProjectWIND/wba"
	"github.com/dop251/goja"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
)

var GlobalAppList = make(map[string]wba.AppInfo)

func RegisterApp(app wba.AppInfo, force ...bool) (code int) {
	var isForce bool = false
	if len(force) > 0 {
		isForce = force[0]
	}
	if _, ok := GlobalAppList[app.AppKey.Name]; ok {
		if isForce {
			LOG.Notice("应用 %v 已存在，正在覆盖", app.AppKey.Name)
			GlobalAppList[app.AppKey.Name] = app
			return 2
		} else {
			LOG.Notice("应用 %v 已存在,跳过注册", app.AppKey.Name)
			return 1
		}
	}
	GlobalAppList[app.AppKey.Name] = app
	LOG.Info("应用 %v 注册成功", app.AppKey.Name)
	return 0
}

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
		_ = windObj.Set("withDescription", wba.WithDescription)
		_ = windObj.Set("withWebUrl", wba.WithWebUrl)
		_ = windObj.Set("withLicense", wba.WithLicense)
		_ = windObj.Set("registerApp", RegisterApp)
		_ = rt.Set("wsp", ProtocolApi)
		_ = rt.Set("wsd", DatabaseApi)
		_ = rt.Set("wst", ToolsApi)
		return rt
	},
}

func GetRuntimeCopy() *goja.Runtime {
	// 从池中获取基础runtime
	base := runtimePool.Get().(*goja.Runtime)

	// 创建新runtime并复制必要配置
	rt := goja.New()
	rt.SetFieldNameMapper(CamelCaseFieldNameMapper{})

	// 复制核心对象
	_ = rt.Set("console", base.Get("console"))
	_ = rt.Set("wind", base.Get("wind"))
	_ = rt.Set("wsp", base.Get("wsp"))
	_ = rt.Set("wsd", base.Get("wsd"))
	_ = rt.Set("wst", base.Get("wst"))

	// 放回基础runtime
	runtimePool.Put(base)

	return rt
}

func createNewRuntime() *goja.Runtime {
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
	windObj := rt.NewObject()
	_ = rt.Set("wind", windObj)
	_ = windObj.Set("newApp", wba.NewApp)
	_ = windObj.Set("withDescription", wba.WithDescription)
	_ = windObj.Set("withWebUrl", wba.WithWebUrl)
	_ = windObj.Set("withLicense", wba.WithLicense)
	_ = windObj.Set("registerApp", RegisterApp)
	_ = rt.Set("wsp", ProtocolApi)
	_ = rt.Set("wsd", DatabaseApi)
	_ = rt.Set("wst", ToolsApi)
	_ = rt.Set("wsc", CalculatorApi)
	return rt
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
func ReloadApps() (total int, registered int) {
	// 清空AppMap和CmdMap
	CmdMap = make(map[wba.AppKey]wba.CmdList)
	AppMap = make(map[wba.AppKey]wba.AppInfo)
	ScheduledTaskMap = make(map[wba.AppKey]map[string]wba.ScheduledTaskInfo)
	GlobalCmdAgentSelector = wba.NewCmdAgentSelector()
	appsDir := "./data/app/"
	appFiles, err := os.ReadDir(appsDir)
	total = 0
	registered = 0
	if err != nil {
		LOG.Error("加载应用所在目录失败:%v", err)
		return
	}

	for _, file := range appFiles {
		totalDelta, registeredDelta := registerAppToList(file, appsDir)
		total += totalDelta
		registered += registeredDelta
	}

	for _, app := range GlobalAppList {
		loadApp(app)
	}

	CmdMap[AppCore.AppKey] = AppCore.CmdMap
	GlobalCmdAgentSelector.AddCmdMap(CmdMap)
	return total, registered
}

func loadApp(appInfo wba.AppInfo) {
	// 配置错误捕获
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

	LOG.Info("JS应用 %s 加载成功", appInfo.AppKey.Name)
}

// registerAppToList 注册应用到列表中
func registerAppToList(file os.DirEntry, appsDir string) (totalDelta int, successDelta int) {
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

		runtime := createNewRuntime()

		// 配置错误捕获
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

		// 执行JS代码
		safeRun(func() error {
			_, err := runtime.RunString(string(jsCode))
			return err
		})

		return 1, 1
	}
	return 0, 0
}
