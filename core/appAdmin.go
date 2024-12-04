package core

import (
	"ProjectWIND/LOG"
	"os"
	"path/filepath"
	"plugin"
)

func reloadApps() {
	appsDir := "./data/app/"
	appFiles, err := os.ReadDir(appsDir)
	total := 0
	success := 0
	if err != nil {
		LOG.ERROR("Error reading apps directory:%v", err)
		return
	}

	for _, file := range appFiles {
		totalDelta, successDelta := reloadAPP(file, appsDir)
		total += totalDelta
		success += successDelta
	}
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
			LOG.ERROR("Error opening app %s: %v\n", pluginPath, err)
			return 1, 0
		}

		initSymbol, err := p.Lookup("init")
		if err != nil {
			LOG.ERROR("Error finding init function in app %s: %v\n", pluginPath, err)
			return 1, 0
		}

		initFunc, ok := initSymbol.(func())
		if !ok {
			LOG.ERROR("init symbol in app %s is not a function\n", pluginPath)
			return 1, 0
		}

		initFunc()
		LOG.INFO("App %s initialized successfully\n", pluginPath)
		return 1, 1
	}
	return 0, 0
}
