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
		LOG.ERROR("Error reading apps directory:%v", err)
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
			LOG.ERROR("Error opening app %s: %v", pluginPath, err)
			return 1, 0
		}

		initSymbol, err := p.Lookup("Application")
		if err != nil {
			LOG.ERROR("Error finding interface Application in app %s: %v", pluginPath, err)
			return 1, 0
		}

		app, ok := initSymbol.(wba.APP)
		if !ok {
			LOG.ERROR("init symbol in app %s is not a right type", pluginPath)
			return 1, 0
		}

		err = app.Init(&AppApi)
		if err != nil {
			LOG.ERROR("Error initializing app %s: %v", pluginPath, err)
		}

		CmdMap = mergeMaps(CmdMap, app.Get().CmdMap)
		LOG.INFO("App %s initialized successfully", pluginPath)
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
