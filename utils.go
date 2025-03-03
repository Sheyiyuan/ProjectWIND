package main

import (
	"ProjectWIND/LOG"
	"ProjectWIND/core"
	"ProjectWIND/database"
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

	LOG.Info("正在初始化WIND配置文件...")

	err := checkAndUpdateConfig("./data/core.json")
	if err != nil {
		LOG.Fatal("初始化时，加载配置文件 ./data/core.json 失败: %v", err)
	}
	// 创建日志文件
	logFile := fmt.Sprintf("./data/log/WIND_CORE_%s.log", time.Now().Format("20060102150405"))
	_, err = os.Stat(logFile)
	if os.IsNotExist(err) {
		file, err := os.Create(logFile)
		if err != nil {
			LOG.Fatal("初始化时，创建日志文件失败: %v", err)
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				LOG.Fatal("无法关闭日志文件: %v", err)
			}
		}(file)
	}

	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		LOG.Fatal("初始化时，无法打开日志文件: %v", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			LOG.Fatal("无法关闭日志文件: %v", err)
		}
	}(file)

	// 设置日志输出到文件
	log.SetOutput(io.MultiWriter(os.Stdout, file))
	LOG.Info("WIND配置文件初始化完成！")
	return logFile
}

func checkAndUpdateConfig(configPath string) error {
	// 检查并创建必要的目录和文件
	if _, err := os.Stat("./data/"); os.IsNotExist(err) {
		// 如果不存在，则创建该文件夹
		err := os.Mkdir("./data/", 0755)
		if err != nil {
			LOG.Fatal("初始化时，创建data文件夹失败: %v", err)
		}
	}

	// 检查./data/文件夹中是否存在core.json文件
	if _, err := os.Stat("./data/core.json"); os.IsNotExist(err) {
		// 如果不存在，则创建该文件
		file, err := os.Create("./data/core.json")
		if err != nil {
			LOG.Fatal("初始化时，创建 ./data/core.json 配置文件失败: %v", err)
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				LOG.Fatal("关闭 ./data/core.json 配置文件失败: %v", err)
			}
		}(file)
	}

	// 检查并更新配置文件
	var coreConfig typed.CoreConfigInfo

	var defaultProtocol typed.Protocol
	defaultProtocol.ProtocolName = "EXAMPLE"
	defaultProtocol.ProtocolPlatform = "在这里输入协议平台"
	defaultProtocol.ProtocolAddr = "在这里输入协议地址，如'ws://127.0.0.1:8080'"
	defaultProtocol.Token = "在这里输入协议的Token"
	defaultProtocol.Enable = true

	var defaultConfig typed.CoreConfigInfo
	defaultConfig.CoreName = "windCore"
	defaultConfig.WebUIPort = 3211
	defaultConfig.Protocols = []typed.Protocol{defaultProtocol}
	defaultConfig.ServiceName = "wind"
	// 读取配置文件
	file, err := os.Open(configPath)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			LOG.Fatal("无法关闭配置文件 ./data/core.json: %v", err)
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
	//if coreConfig.ProtocolAddr == "" {
	//	coreConfig.ProtocolAddr = defaultConfig.ProtocolAddr
	//}
	if coreConfig.Protocols == nil || len(coreConfig.Protocols) == 0 {
		coreConfig.Protocols = defaultConfig.Protocols
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
	//if coreConfig.Token == "" {
	//	coreConfig.Token = ""
	//}

	formattedJSON, err := json.MarshalIndent(coreConfig, "", "  ")
	if err != nil {
		return err
	}

	// 将格式化后的JSON字符串写入文件
	file, err = os.Create("./data/core.json")
	if err != nil {
		LOG.Fatal("初始化时，创建 ./data/core.json 配置文件失败: %v", err)
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			LOG.Fatal("无法关闭配置文件 ./data/core.json: %v", err)
		}
	}(file)

	_, err = file.Write(formattedJSON)
	if err != nil {
		LOG.Fatal("初始化时，写入 ./data/core.json 配置文件失败: %v", err)
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
		LOG.Fatal("创建应用文件夹 ./data/app/ 失败: %v", err)
		return err
	}
	err = checkDataFolderExistence("./data/images/")
	if err != nil {
		LOG.Fatal("创建图片文件夹 ./data/images/ 失败: %v", err)
		return err
	}
	err = checkDataFolderExistence("./data/files/")
	if err != nil {
		LOG.Fatal("创建文件文件夹 ./data/files/ 失败: %v", err)
		return err
	}
	err = checkDataFolderExistence("./data/videos/")
	if err != nil {
		LOG.Fatal("创建视频文件夹 ./data/videos/ 失败: %v", err)
		return err
	}
	err = checkDataFolderExistence("./data/audios/")
	if err != nil {
		LOG.Fatal("创建音频文件夹 ./data/audios/ 失败: %v", err)
		return err
	}
	err = checkDataFolderExistence("./data/database/")
	if err != nil {
		LOG.Fatal("创建数据库文件夹 ./data/database/ 失败: %v", err)
		return err
	}
	err = checkDataFolderExistence("./data/log/")
	if err != nil {
		LOG.Fatal("创建日志文件夹 ./data/log/ 失败: %v", err)
		return err
	}
	err = checkDataFolderExistence("./data/app/configs/")
	if err != nil {
		LOG.Fatal("创建应用配置文件夹 ./data/app/configs/ 失败: %v", err)
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
			LOG.Fatal("打开日志文件失败: %v", err)
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				LOG.Fatal("无法关闭日志文件: %v", err)
			}
		}(file)
		// 设置日志输出到文件
		log.SetOutput(io.MultiWriter(os.Stdout, file))

		LOG.Info("正在启动WIND核心服务...")
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
		LOG.Fatal("无法打开日志文件: %v", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			LOG.Fatal("无法关闭日志文件: %v", err)
		}
	}(file)
	// 设置日志输出到文件
	log.SetOutput(io.MultiWriter(os.Stdout, file))
	//在/etc/systemd/system/下创建服务文件

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
		LOG.Fatal("无法打开日志文件: %v", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			LOG.Fatal("无法关闭日志文件: %v", err)
		}
	}(file)
	// 设置日志输出到文件
	log.SetOutput(io.MultiWriter(os.Stdout, file))
	ReloadApps()
	//从配置文件中读取配置信息
	LOG.Info("正在启动WIND协议服务...")
	var config typed.CoreConfigInfo
	file, err = os.Open("./data/core.json")
	if err != nil {
		LOG.Fatal("无法打开配置文件 ./data/core.json: %v", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			LOG.Fatal("无法关闭配置文件 ./data/core.json: %v", err)
		}
	}(file)

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		LOG.Fatal("连接协议时，解析配置文件 ./data/core.json 失败: %v", err)
	}
	LOG.Info("正在启动WebSocket链接程序...")
	protocolNum := 0
	breakNum := 0
	UnenableProtocolNum := 0
	for _, protocol := range config.Protocols {
		protocolName := protocol.ProtocolName
		if protocolName == "EXAMPLE" {
			continue
		}
		if protocolName == "" {
			LOG.Warn("连接协议 %s 时，协议名称为空，跳过该协议", protocolName)
			breakNum++
			continue
		}
		//获取协议地址
		protocolAddr := protocol.ProtocolAddr
		if protocolAddr == "" {
			LOG.Warn("连接协议 %s 时，协议地址为空，跳过该协议", protocolName)
			breakNum++
			continue
		}
		if protocol.Enable == false {
			LOG.Warn("连接协议 %s 时，协议已禁用，跳过该协议", protocolName)
			UnenableProtocolNum++
			continue
		}
		//获取token
		token := protocol.Token
		// 启动 WebSocket 处理程序
		go func() {
			err := core.WebSocketHandler(protocolAddr, token)
			if err != nil {
				LOG.Error("连接协议时，启动 WebSocket 处理程序失败: %v", err)
			}
		}()
		protocolNum++
	}
	LOG.Info(" %d 个协议服务启动完成, %d 个协议服务已禁用, %d 个协议服务因为配置错误被跳过。", protocolNum, UnenableProtocolNum, breakNum)
	select {}
}

func ReloadApps() {
	LOG.Info("正在重新加载应用...")
	total, success := core.ReloadApps()
	LOG.Info("应用重新加载完成，共加载%d个应用，成功加载%d个应用。", total, success)
}

func startDatabase() {
	go database.Start()
	time.Sleep(time.Second * 1)
	// 读写测试
	// for i := 0; i < 10; i++ {
	// 	data, ok := database.Get("go", "go", "user", "test", "test"+fmt.Sprintf("%d", i), false)
	// 	if !ok {
	// 		LOG.Error("Failed to get data from database")
	// 		continue
	// 	}
	// 	LOG.Info("Get data from database: %v", data)
	// 	time.Sleep(time.Second * 1)
	// }
	// time.Sleep(time.Second * 1)
	// for i := 0; i < 10; i++ {
	// 	database.Set("go", "go", "user", "test", "test"+fmt.Sprintf("%d", i), "test"+fmt.Sprintf("%d", 1000+i))
	// 	time.Sleep(time.Second * 1)
	// }
	// time.Sleep(time.Second * 1)
	// for i := 0; i < 10; i++ {
	// 	data, ok := database.Get("go", "go", "user", "test", "test"+fmt.Sprintf("%d", i), false)
	// 	if !ok {
	// 		LOG.Error("Failed to get data from database")
	// 		continue
	// 	}
	// 	LOG.Info("Get data from database: %v", data)
	// 	time.Sleep(time.Second * 1)
	// }
	select {}
}
