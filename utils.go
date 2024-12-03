package main

import (
	"ProjectWIND/LOG"
	"ProjectWIND/protocol"
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
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.SetPrefix("[WIND] ")

	LOG.INFO("正在初始化WIND配置文件...")

	err := checkAndUpdateConfig("./data/config.json")
	if err != nil {
		LOG.FATAL("Failed to initialize config file: %v", err)
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

	// 检查./data/文件夹中是否存在config.json文件
	if _, err := os.Stat("./data/config.json"); os.IsNotExist(err) {
		// 如果不存在，则创建该文件
		file, err := os.Create("./data/config.json")
		if err != nil {
			LOG.FATAL("Failed to create config file: %v", err)
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				LOG.FATAL("Failed to close config file: %v", err)
			}
		}(file)
	}

	// 检查并更新配置文件
	var config typed.ConfigInfo

	// 定义默认配置
	var defaultConfig typed.ConfigInfo
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
			LOG.FATAL("Failed to close config file: %v", err)
		}
	}(file)

	// 解码JSON配置
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		if !errors.Is(err, io.EOF) {
			return err
		}
	}

	// 检查并更新配置
	if config.ProtocolAddr == "" {
		config.ProtocolAddr = defaultConfig.ProtocolAddr
	}
	if config.WebUIPort == 0 {
		config.WebUIPort = defaultConfig.WebUIPort
	}
	if config.CoreName == "" {
		config.CoreName = defaultConfig.CoreName
	}
	if config.ServiceName == "" {
		config.ServiceName = defaultConfig.ServiceName
	}
	if config.PasswordHash == "" {
		config.PasswordHash = ""
	}

	formattedJSON, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	// 将格式化后的JSON字符串写入文件
	file, err = os.Create("./data/config.json")
	if err != nil {
		LOG.FATAL("Error creating config file:%v", err)
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			LOG.FATAL("Failed to close config file: %v", err)
		}
	}(file)

	_, err = file.Write(formattedJSON)
	if err != nil {
		LOG.FATAL("Error writing to config file: %v", err)
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
		log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
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
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
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
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
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
	//从配置文件中读取配置信息
	LOG.INFO("正在启动WIND协议服务...")
	var config typed.ConfigInfo
	file, err = os.Open("./data/config.json")
	if err != nil {
		LOG.FATAL("Failed to open config file: %v", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			LOG.FATAL("Failed to close config file: %v", err)
		}
	}(file)

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		LOG.FATAL("Failed to decode config file when linking to protocol: %v", err)
	}
	//获取协议地址
	protocolAddr := config.ProtocolAddr
	//链接协议
	// 启动 WebSocket 处理程序
	LOG.INFO("正在启动WebSocket链接程序...")
	err = protocol.WebSocketHandler(protocolAddr)
	if err != nil {
		// 如果发生错误，记录错误并退出程序
		LOG.FATAL("Failed to start WebSocket link program: %v", err)
	}
	return
}
