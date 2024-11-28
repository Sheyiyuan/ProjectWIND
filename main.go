package main

import (
	"ProjectWIND/protocol"
	"ProjectWIND/typed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

func init() {
	// 初始化日志记录器
	log.SetFlags(log.Ldate | log.Ltime)
	log.SetPrefix("[WIND] ")

	// 检查并创建必要的目录和文件
	if _, err := os.Stat("./data/"); os.IsNotExist(err) {
		// 如果不存在，则创建该文件夹
		err := os.Mkdir("./data/", 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	// 检查./data/文件夹中是否存在config.json文件
	if _, err := os.Stat("./data/config.json"); os.IsNotExist(err) {
		// 如果不存在，则创建该文件
		file, err := os.Create("./data/config.json")
		if err != nil {
			log.Fatal(err)
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				log.Fatal(err)
			}
		}(file)
	}
	// 检查并更新配置文件
	err := checkAndUpdateConfig("./data/config.json")
	if err != nil {
		log.Fatal(err)
	}

	// 检查./data/文件夹中是否存在app/文件夹
	if _, err := os.Stat("./data/app/"); os.IsNotExist(err) {
		// 如果不存在，则创建该文件夹
		err := os.Mkdir("./data/app/", 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	// 检查./data/文件夹中是否存在database/文件夹
	if _, err := os.Stat("./data/database/"); os.IsNotExist(err) {
		// 如果不存在，则创建该文件夹
		err := os.Mkdir("./data/database/", 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	// 检查./data/文件夹中是否存在log/文件夹
	if _, err := os.Stat("./data/log/"); os.IsNotExist(err) {
		// 如果不存在，则创建该文件夹
		err := os.Mkdir("./data/log/", 0755)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	//读取参数
	if len(os.Args) == 1 || os.Args[0] == "help" || os.Args[0] == "-h" || os.Args[0] == "--help" {
		fmt.Println(`Usage: input command "start"or"-s" to start the server.`)
		return
	}
	command := os.Args[1]
	switch command {
	case "start":
	case "-s":
		{
			log.Println("Starting ProjectWIND...")
			// 启动 WebSocket 处理程序
			log.Println("Starting WebSocket handler...")
			_, err := protocol.WebSocketHandler()
			if err != nil {
				// 如果发生错误，记录错误并退出程序
				log.Fatal(err)
			}
			return
		}
	default:
		{
			fmt.Println("Invalid command.")
			return
		}
	}
}
func checkAndUpdateConfig(configPath string) error {
	var config typed.ConfigInfo
	// 定义默认配置
	var defaultConfig typed.ConfigInfo
	defaultConfig.CoreName = "windCore"
	defaultConfig.WebUIPort = 3211
	defaultConfig.ProtocolAddr = make(map[string]string)
	defaultConfig.ServiceName = "wind"
	// 读取配置文件
	file, err := os.Open(configPath)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Printf("Failed to close config file: %v", err)
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
	if config.ProtocolAddr == nil {
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
		fmt.Println("Error formatting JSON:", err)
		return err
	}

	// 将格式化后的JSON字符串写入文件
	file, err = os.Create("./data/config.json")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Printf("Failed to close config file: %v", err)
		}
	}(file)

	_, err = file.Write(formattedJSON)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return err
	}
	return nil
}
