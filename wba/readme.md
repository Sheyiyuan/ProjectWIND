# WBA 文档

WBA是用于对接WIND核心和APP通信的接口，规定了应用结构规范，并提供了与wind核心交互的接口定义。

## 目录
- [1. 应用结构规范](#1-应用结构规范)
  - [1.1 应用结构](#11-应用结构)
  - [1.2 Application结构体的成员](#12-application结构体的成员)
  - [1.3 Application的注册](#13-application的注册)
- [2. 接口定义](#2-接口定义)
  - [2.1 protocol模块](#21-protocol模块)
  - [2.2 Log模块](#22-log模块)
  - [2.3 数据库模块](#23-数据库模块)
  - [2.4 文件管理模块](#24-文件管理模块)
  - [2.5 远程控制模块](#25-远程控制模块)


## 1. 应用结构规范

### 1.1 应用结构

wba会向wind核心提供两个用于调用应用逻辑的方法：

- init(Wind):

    初始化wba，传入wind实例，应用可以通过传入的wind实例调用wind提供的接口，具体参见[wind接口定义](#2-接口定义)。

- Get()

    核心会通过该方法获取应用中定义的对象，该方法返回[Application 结构体](#12-application结构体的成员)

### 1.2 Application结构体的成员
| 字段名                 | 类型                           | 默认值        | 说明                                                                                                                                                                                                                |
|---------------------|------------------------------|:-----------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| Name                | string                       | -          | 应用名称，会被用于帮助文档的检索                                                                                                                                                                                                  |
| Version             | string                       | -          | 应用版本号                                                                                                                                                                                                             |
| Author              | string                       | -          | 作者                                                                                                                                                                                                                |
| Description         | string                       | -          | 应用详细描述，会被写入到帮助文档中                                                                                                                                                                                                 |
| Namespace           | string                       | `"PUBLIC"` | 命名空间，用于区分不同应用                                                                                                                                                                                                     |
| Homepage            | string                       | -          | 应用主页地址                                                                                                                                                                                                            |
| License             | string                       | `"MIT"`    | 应用许可证                                                                                                                                                                                                             |
| AppType             | string                       | `"fun"`    | 应用类型,可选值有以下不同级别：<ul><li>`fun`：娱乐级应用</li><li>`assist`：规则辅助级应用</li><li>`rule`：规则级应用</li><li>`system`：系统级应用</ul> <br/>应用调用时的优先级从低到高依次为`fun`、`assist`、`rule`、`system`。同等级的应用，如有关联的高级别应用，则优先级由高级别的应用决定。否则会按照注册顺序决定优先级。 |
| Rule                | string                       | -          | 仅在AppType为`assist`和`rule`时有效，用于确认规则类型                                                                                                                                                                             |
| CmdMap              | map[string] Cmd              | -          | 存储命令名称和命令实例的映射map                                                                                                                                                                                                 |
| MessageEventHandler | func(msg MessageEventInfo)   | -          | 在收到消息事件时触发的事件处理函数                                                                                                                                                                                                 |
| NoticeEventHandler  | func(msg NoticeEventInfo)    | -          | 在收到通知事件时触发的事件处理函数                                                                                                                                                                                                 |
| RequestEventHandler | func(msg RequestEventInfo)   | -          | 在收到请求事件时触发的事件处理函数                                                                                                                                                                                                 |
| MetaEventHandler    | func(msg MetaEventInfo)      | -          | 在收到元事件时触发的事件处理函数                                                                                                                                                                                                  |
| ScheduledTasks      | map[string]ScheduledTaskInfo | -          | 定时任务列表，key为任务名称，value为任务信息结构体                                                                                                                                                                                     |
| API                 | map[string]interface{}       | -          | 应用对外提供的API列表，可被其他应用调用，其中，key为API名称，value为API实例                                                                                                                                                                    |


#### Cmd结构体的成员
| 字段名   | 类型     | 说明                                             |
|-------|--------|------------------------------------------------|
| NAME  | string | 命令的名称                                          |
| DESC  | string | 命令的描述，会被用于帮助信息的展示                              |
| SOLVE | func   | 解决命令的函数，接收参数为命令参数的字符串切片和 `MessageEventInfo` 类型 |

#### ScheduledTaskInfo结构体的成员
| 字段名  | 类型     | 说明             |
|------|--------|----------------|
| Name | string | 任务名称           |
| Desc | string | 任务描述           |
| Task | func() | 任务执行的函数        |
| Cron | string | 定时任务的 Cron 表达式 |

### 1.3 Application的注册

下面给出一个注册应用的示例代码：

#### 注册应用
```go
package main

import (
	"github.com/wind/wba"
)

func appInit() wba.AppInfo {
	// 写入应用信息
	app := wba.NewApp(
		wba.WithName("app_demo"),                           // 应用名称
		wba.WithAuthor("WIND"),                             // 作者
		wba.WithVersion("1.0.0"),                           // 版本
		wba.WithDescription("This is a demo application"),  // 应用描述
		wba.WithNamespace("app_demo"),                      // 命名空间, 私有数据库请使用应用的名称, 公共数据库请使用"PUBLIC"
		wba.WithWebUrl("https://github.com/wind/app_demo"), // 应用主页
		wba.WithLicense("MIT"),                             // 应用许可证
	)

	return app
}

// Application 向核心暴露的应用接口,标识符为Application, 不可修改
var Application = appInit()
```

定义初始化函数 `appInit()` ，返回 `wba.AppInfo` 类型，该类型包含应用的基本信息，包括名称、版本、作者、描述、命名空间、主页、许可证等。并使用固定标识符进行接口暴露，该标识符为 `Application`，不可修改。

在初始化时，使用 `wba.NewApp()` 方法创建 `wba.AppInfo` 实例，并使用 `wba.WithName()`、`wba.WithAuthor()`、`wba.WithVersion()`、`wba.WithDescription()`、`wba.WithNamespace()`、`wba.WithWebUrl()`、`wba.WithLicense()`、`wba.WithAppType()`、`wba.WithRule()` 方法设置应用的基本信息。未设置的字段将使用默认值。

对于其他的接口，如命令、定时任务等，wba提供了一些函数和方法，用于设置这些接口。

#### 注册命令
```go
	// 定义命令
	cmdTest := wba.NewCmd(
		//命令名称
		"app",
		//命令介绍
		"插件测试",
		func(args []string, msg wba.MessageEventInfo) {
			val := args[0]
			log.Println("app_demo cmdTest", val)
			switch val {
			case "help":
				{
					wba.Wind.SendMsg(msg, "app_demo help", false)
				}
			default:
				{
					wba.Wind.SendMsg(msg, "Hello, wind app!", false)
					return
				}
			}
		},
	)

	// 将命令添加到应用命令列表中
	app.AddCmd("app", cmdTest)
	//可以为同一个命令注册多个名称
	app.AddCmd("test", cmdTest)
```

使用 `wba.NewCmd()` 方法创建 `wba.Cmd` 实例，并设置命令的名称、描述、解决函数。之后使用 `app.AddCmd()` 方法将命令添加到应用命令列表中。

#### 注册定时任务
```go
	// 定义定时任务
	taskTest := wba.ScheduledTaskinfo{
		Name: "task_test",
		Desc: "定时任务测试",
		Task: func() {
			log.Println("task_test")
		},
		Cron: "*/1 * * * *",
	}

	// 将定时任务添加到应用定时任务列表中
	app.AddScheduledTask(taskTest)
	}
	
```

使用 `wba.ScheduledTaskinfo` 结构体定义定时任务，并设置任务名称、描述、执行函数、定时任务的 Cron 表达式。之后使用 `app.AddScheduledTask()` 方法将定时任务添加到应用定时任务列表中。



## 2. 接口定义

wind在初始化时，会调用 `init(Wind)` 方法，传入 `Wind` 实例，应用可以通过该实例调用wind提供的接口。

wind实例提供的接口可以分为下面几个部分：

### 2.1 Protocol模块

所有API请求按照OneBot11标准，使用JSON格式进行数据交换。api命名为由原文档中蛇形命名法改为双驼峰命名法。

关于API的详细信息，请参考[OneBot文档](https://github.com/botuniverse/onebot-11/blob/master/README.md)。

```go
	SendMsg(msg MessageEventInfo, message string, autoEscape bool)
	SendPrivateMsg(msg MessageEventInfo, message string, autoEscape bool)
	SendGroupMsg(msg MessageEventInfo, message string, autoEscape bool)
	DeleteMsg(msg MessageEventInfo)
	SendLike(userId int64, times int)
	SetGroupKick(groupId int64, userId int64, rejectAddRequest bool)
	SetGroupBan(groupId int64, userId int64, duration int32)
	SetGroupWholeBan(groupId int64, enable bool)
	SetGroupAdmin(groupId int64, userId int64, enable bool)
	SetGroupLeave(groupId int64, isDismiss bool)
	SetGroupCard(groupId int64, userId int64, card string)
	SetGroupName(groupId int64, groupName string)
	SetGroupSpecialTitle(groupId int64, userId int64, specialTitle string, duration int32)
	SetFriendAddRequest(flag string, approve bool, remark string)
	SetGroupAddRequest(flag string, subType string, approve bool, reason string)
	GetLoginInfo() APIResponseInfo
	GetVersionInfo() APIResponseInfo
	GetMsg(msgId int32) APIResponseInfo
	GetForwardMsg(msgId string) APIResponseInfo
	GetGroupList() APIResponseInfo
	GetGroupMemberList(groupId int64) APIResponseInfo
	GetGroupMemberInfo(groupId int64, userId int64, noCache bool) APIResponseInfo
	GetFriendList() APIResponseInfo
	GetStrangerInfo(userId int64, noCache bool) APIResponseInfo
	GetGroupInfo(groupId int64, noCache bool) APIResponseInfo
	GetGroupHonorInfo(groupId int64, Type string) APIResponseInfo
	GetStatus() APIResponseInfo
	GetCookies(domain string) APIResponseInfo
	GetCSRFToken() APIResponseInfo
	GetCredentials(domain string) APIResponseInfo
	GetImage(file string) APIResponseInfo
	GetRecord(file string, outFormat string) APIResponseInfo
	CanSendImage() APIResponseInfo
	CanSendRecord() APIResponseInfo
	SetRestart(delay int32)
	CleanCache()
```

### 2.2 Log模块

1. 日志模块使用go-logging库，日志级别分为DEBUG、INFO、WARN、ERROR。

2. 日志模块提供LogWith方法，可以自定义日志级别，调用级别为DEBUG时，会打印输出调用者的文件名、函数名、行号。

3. 日志模块提供Log方法，默认日志级别为INFO。

### 2.3 数据库模块



### 2.4 文件管理模块



### 2.5 远程控制模块
