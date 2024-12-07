package main

import (
	"ProjectWIND/LOG"
	"ProjectWIND/core"
	"ProjectWIND/typed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

func initCore() string {
	// 初始化日志记录器
	log.SetFlags(log.Ldate | log.Ltime)
	log.SetPrefix("[WIND] ")

	LOG.INFO("正在初始化WIND配置文件...")

	err := checkAndUpdateConfig("./data/core.json")
	if err != nil {
		LOG.FATAL("Failed to initialize core.json file: %v", err)
	}
	// 创建日志文件
	logFile := fmt.Sprintf("./data/log/WIND_CORE_%s.log", time.Now().Format("20060102150405"))
	_, err = os.Stat(logFile)
	if os.IsNotExist(err) {
		file, err := os.Create(logFile)
		if err != nil {
			LOG.FATAL("Failed to create log file: %v", err)
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				LOG.FATAL("Failed to close log file: %v", err)
			}
		}(file)
	}

	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		LOG.FATAL("Failed to open log file: %v", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			LOG.FATAL("Failed to close log file: %v", err)
		}
	}(file)

	// 设置日志输出到文件
	log.SetOutput(io.MultiWriter(os.Stdout, file))
	LOG.INFO("WIND配置文件初始化完成！")
	return logFile
}

func checkAndUpdateConfig(configPath string) error {
	// 检查并创建必要的目录和文件
	if _, err := os.Stat("./data/"); os.IsNotExist(err) {
		// 如果不存在，则创建该文件夹
		err := os.Mkdir("./data/", 0755)
		if err != nil {
			LOG.FATAL("Failed to create data folder: %v", err)
		}
	}

	// 检查./data/文件夹中是否存在core.json文件
	if _, err := os.Stat("./data/core.json"); os.IsNotExist(err) {
		// 如果不存在，则创建该文件
		file, err := os.Create("./data/core.json")
		if err != nil {
			LOG.FATAL("Failed to create core.json file: %v", err)
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				LOG.FATAL("Failed to close core.json file: %v", err)
			}
		}(file)
	}

	// 检查并更新配置文件
	var coreConfig typed.CoreConfigInfo

	// 定义默认配置
	var defaultConfig typed.CoreConfigInfo
	defaultConfig.CoreName = "windCore"
	defaultConfig.WebUIPort = 3211
	defaultConfig.ProtocolAddr = ""
	defaultConfig.ServiceName = "wind"
	// 读取配置文件
	file, err := os.Open(configPath)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			LOG.FATAL("Failed to close core.json file: %v", err)
		}
	}(file)

	// 解码JSON配置
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&coreConfig)
	if err != nil {
		if !errors.Is(err, io.EOF) {
			return err
		}
	}

	// 检查并更新配置
	if coreConfig.ProtocolAddr == "" {
		coreConfig.ProtocolAddr = defaultConfig.ProtocolAddr
	}
	if coreConfig.WebUIPort == 0 {
		coreConfig.WebUIPort = defaultConfig.WebUIPort
	}
	if coreConfig.CoreName == "" {
		coreConfig.CoreName = defaultConfig.CoreName
	}
	if coreConfig.ServiceName == "" {
		coreConfig.ServiceName = defaultConfig.ServiceName
	}
	if coreConfig.PasswordHash == "" {
		coreConfig.PasswordHash = ""
	}
	if coreConfig.Token == "" {
		coreConfig.Token = ""
	}

	formattedJSON, err := json.MarshalIndent(coreConfig, "", "  ")
	if err != nil {
		return err
	}

	// 将格式化后的JSON字符串写入文件
	file, err = os.Create("./data/core.json")
	if err != nil {
		LOG.FATAL("Error creating core.json file:%v", err)
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			LOG.FATAL("Failed to close core.json file: %v", err)
		}
	}(file)

	_, err = file.Write(formattedJSON)
	if err != nil {
		LOG.FATAL("Error writing to core.json file: %v", err)
		return err
	}

	checkDataFolderExistence := func(dataAddress string) error {
		// 检查./data/文件夹中是否存在dataAddress文件夹
		if _, err := os.Stat(dataAddress); os.IsNotExist(err) {
			err := os.Mkdir(dataAddress, 0755)
			if err != nil {
				return err
			}
		}
		return nil
	}

	err = checkDataFolderExistence("./data/app/")
	if err != nil {
		LOG.FATAL("Failed to create app folder: %v", err)
		return err
	}
	err = checkDataFolderExistence("./data/images/")
	if err != nil {
		LOG.FATAL("Failed to create images folder: %v", err)
		return err
	}
	err = checkDataFolderExistence("./data/database/")
	if err != nil {
		LOG.FATAL("Failed to create database folder: %v", err)
		return err
	}
	err = checkDataFolderExistence("./data/log/")
	if err != nil {
		LOG.FATAL("Failed to create log folder: %v", err)
		return err
	}

	return nil
}

func startWebUI() {
	{
		//初始化
		logFile := initCore()
		// 设置日志输出到文件
		log.SetFlags(log.Ldate | log.Ltime)
		log.SetPrefix("[WIND] ")
		// 打开日志文件
		file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			LOG.FATAL("Failed to open log file: %v", err)
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				LOG.FATAL("Failed to close log file: %v", err)
			}
		}(file)
		// 设置日志输出到文件
		log.SetOutput(io.MultiWriter(os.Stdout, file))

		LOG.INFO("正在启动WIND核心服务...")
		// 启动 WebSocket 处理程序

		//TODO: 这里要添加webUI的启动代码
	}
}

func registerService() {
	//初始化
	logFile := initCore()
	// 设置日志输出到文件
	log.SetFlags(log.Ldate | log.Ltime)
	log.SetPrefix("[WIND] ")
	// 打开日志文件
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		LOG.FATAL("Failed to create log file: %v", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			LOG.FATAL("Failed to close log file: %v", err)
		}
	}(file)
	// 设置日志输出到文件
	log.SetOutput(io.MultiWriter(os.Stdout, file))
	//TODO: 这里要添加注册服务的代码
}

func startProtocol() {
	//初始化
	logFile := initCore()
	// 设置日志输出到文件
	log.SetFlags(log.Ldate | log.Ltime)
	log.SetPrefix("[WIND] ")
	// 打开日志文件
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		LOG.FATAL("Failed to create log file: %v", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			LOG.FATAL("Failed to close log file: %v", err)
		}
	}(file)
	// 设置日志输出到文件
	log.SetOutput(io.MultiWriter(os.Stdout, file))
	ReloadApps()
	//从配置文件中读取配置信息
	LOG.INFO("正在启动WIND协议服务...")
	var config typed.CoreConfigInfo
	file, err = os.Open("./data/core.json")
	if err != nil {
		LOG.FATAL("Failed to open core.json file: %v", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			LOG.FATAL("Failed to close core.json file: %v", err)
		}
	}(file)

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		LOG.FATAL("Failed to decode config file when linking to protocol: %v", err)
	}
	//获取协议地址
	protocolAddr := config.ProtocolAddr
	//获取token
	token := config.Token
	//链接协议
	// 启动 WebSocket 处理程序
	LOG.INFO("正在启动WebSocket链接程序...")
	err = core.WebSocketHandler(protocolAddr, token)
	if err != nil {
		// 如果发生错误，记录错误并退出程序
		LOG.FATAL("Failed to start WebSocket link program: %v", err)
	}
	return
}

func AutoSave() {
	for {
		LOG.INFO("自动保存")
		//TODO: 这里要添加自动保存的代码
		time.Sleep(time.Second * 60)
	}
}

func ReloadApps() {
	LOG.INFO("正在重新加载应用...")
	total, success := core.ReloadApps()
	LOG.INFO("应用重新加载完成，共加载%d个应用，成功加载%d个应用。", total, success)
}
