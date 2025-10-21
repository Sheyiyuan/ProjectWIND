# WIND
像风一样自由

<div style="text-align: center;"><img src="./logo.png" alt="ProjectWIND" width="50%" ></div>

license: [MIT](./LICENSE)

🚧🚧🚧🚧施工中ing…… 🚧🚧🚧🚧
> warning: 该项目正在开发中，请勿用于生产环境

## WIND 是什么？

WIND ( *WIND is not dice* ) 是一个基于 Go 语言开发的bot框架，旨在以高自由度完成各种功能。

### 特点
- 基于 Golang 的核心，轻量高效，跨平台支持
- 基于onebot11协议，可对接各种平台
- 利用goja实现的js插件系统，高效灵活的插件实现
- 定时任务、API事件响应等多种触发方式，满足各种需求
- 同时支持webUI和Terminal交互，方便开发调试和维护管理
- 支持多用户系统，灵活管理用户权限，高效资源共享，一个核心运行多个bot

---

TODO:
- [ ] 底层框架
  - [x] 单协议通信
  - [ ] 文件资源管理
  - [ ] 多用户系统
  - [ ] 数据库管理
  - [x] 日志输出
  - [x] 配置文件
- [ ] 插件管理
  - [x] 基础事件处理
  - [x] 指令管理
  - [ ] 定时任务管理
  - [ ] API管理
- [ ] web ui
- [ ] 文档编写
- [ ] 手册编写

---
# 项目依赖与致谢

本项目的开发收到了以下项目的启发和支持，在此对其开发者表示诚挚的感谢：

- [sealdice](https://github.com/sealdice/sealdice-core/) 海豹 TRPG 骰点核心，开源跑团辅助工具
- [OlivOS](https://github.com/OlivOS-Team/OlivOS) 一个强大的跨平台交互栈与机器人框架
- [NapCatQQ](https://github.com/NapNeko/NapCatQQ) 基于 NTQQ 的现代协议端框架

本项目在开发过程中使用了以下优秀的外部库，在此对其开发者表示诚挚的感谢：

## 1. goja
- **库名称**：goja
- **仓库地址**：[https://github.com/dop251/goja](https://github.com/dop251/goja)
- **用途说明**：goja 作为一款强大的 JavaScript 解释器，在本项目中承担了处理动态的 JavaScript 脚本逻辑，为项目提供了灵活的脚本扩展能力，使得我们能够在项目中实现js插件的功能。它极大地丰富了项目的功能和灵活性，让跨平台的插件开发变得更加容易。

## 2. gocron
- **库名称**：gocron
- **仓库地址**：[https://github.com/go-co-op/gocron](https://github.com/go-co-op/gocron)
- **用途说明**：gocron 是一个出色的任务调度库。在本项目里，它被用作定时任务的调度器，确保了项目中的各种定时任务能够精准、可靠地执行。其简洁易用的 API 设计大大降低了我们实现复杂任务调度逻辑的难度，为项目的稳定运行提供了有力保障。

## 3. hertz
- **库名称**：hertz
- **仓库地址**：[https://github.com/cloudwego/hertz](https://github.com/cloudwego/hertz)
- **用途说明**：hertz 是一个基于 Go 语言开发的高性能 HTTP 路由器。在本项目中，它被用作项目的 HTTP 服务器，为项目提供了快速、高效的 HTTP 请求处理能力。让我们能够灵活地对 HTTP 请求进行处理。

非常感谢以上项目团队的开源贡献，使得我们的项目开发能够借助这些优秀的工具快速推进，为用户带来更好的体验。

---
后面没有了，开发者很懒，什么都没写。