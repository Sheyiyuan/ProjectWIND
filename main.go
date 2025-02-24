package main

import (
	"ProjectWIND/core"
	"fmt"
	"os"
)

var CoreOs = core.GetOS()

func main() {
	//如果没有参数，则启动WebUI
	if len(os.Args) <= 1 {
		initCore()
		// 启动WebUI
		h := core.WebServer("8080")
		go h.Spin()
		// 连接到数据库
		go startDatabase()
		// 连接到协议端
		go startProtocol()
		select {}
	}
	cmdArgs := os.Args[1:]
	if cmdArgs[0] == "-h" || cmdArgs[0] == "--help" {
		fmt.Printf(logo)
		fmt.Printf("%v\n", helpDoc)
		return
	}
	if cmdArgs[0] == "-r" || cmdArgs[0] == "--run" {
		// 启动服务
		startWebUI()
		return
	}
	if cmdArgs[0] == "-i" || cmdArgs[0] == "--init" {
		// 初始化项目
		initCore()
		return
	}
	if cmdArgs[0] == "-v" || cmdArgs[0] == "--version" {
		// 显示版本信息
		fmt.Printf(logo)
		fmt.Printf("%v \n架构：%v\n", version, CoreOs.String())
		return
	}
	if cmdArgs[0] == "-s" || cmdArgs[0] == "--service" {
		// 注册Linux服务并启动
		registerService()
		return
	}
	fmt.Println("未知命令，请使用-h查看帮助。")
	return
}
