package core

import (
	"ProjectWIND/LOG"
	"ProjectWIND/wba"
	"os"
	"path/filepath"
	"plugin"
)

var CmdMap = make(map[string]wba.Cmd)

func ReloadApps() (total int, success int) {
	appsDir := "./data/app/"
	appFiles, err := os.ReadDir(appsDir)
	total = 0
	success = 0
	if err != nil {
		LOG.ERROR("加载应用所在目录失败:%v", err)
		return
	}

	for _, file := range appFiles {
		totalDelta, successDelta := reloadAPP(file, appsDir)
		total += totalDelta
		success += successDelta
	}
	CmdMap = mergeMaps(CmdMap, AppCore.CmdMap)
	return total, success
}

func reloadAPP(file os.DirEntry, appsDir string) (totalDelta int, successDelta int) {
	if file.IsDir() {
		return 0, 0
	}

	ext := filepath.Ext(file.Name())
	if ext == ".so" || (ext == ".dll" && os.PathSeparator == '\\') {
		pluginPath := filepath.Join(appsDir, file.Name())
		p, err := plugin.Open(pluginPath)
		if err != nil {
			LOG.ERROR("打开应用 %s 时发生错误: %v", pluginPath, err)
			return 1, 0
		}

		initSymbol, err := p.Lookup("Application")
		if err != nil {
			LOG.ERROR("找不到应用 %s 提供的 Application 接口: %v", pluginPath, err)
			return 1, 0
		}

		app, ok := initSymbol.(wba.APP)
		if !ok {
			LOG.ERROR("应用 %s 提供的 Application 接口不是 wba.APP 类型", pluginPath)
			return 1, 0
		}

		err = app.Init(&AppApi)
		if err != nil {
			LOG.ERROR("初始化应用 %s 失败: %v", pluginPath, err)
		}

		CmdMap = mergeMaps(CmdMap, app.Get().CmdMap)
		LOG.INFO("应用 %s 加载成功", pluginPath)
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
