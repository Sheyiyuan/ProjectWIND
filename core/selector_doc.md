# Selector设计思路与使用说明

## 问题背景

在骰子执行指令时，我们经常遇到这样一些问题：
- 不同的插件可能使用相同的指令（很多跑团对应规则都有同名的指令），但是参数和功能完全不同，导致指令无法正常执行。并且当其中有系统级指令时，可能会导致其他严重问题，此时需要一个机制来区分指令。
- 当我们在进行trpg时日志开启后，可能需要屏蔽某些无关的指令和自定义回复，此时需要一个机制来过滤指令和回复（即权限控制）。
- 在不同的群聊中，各种插件的开启情况各不相同，此时需要一个机制来管理不同组群插件的开启和关闭情况（包括前两条中的权限控制问题）。

为解决上述问题，我们设计了一个指令选择器，用于管理和选择指令。

## 设计思路

选择器的设计思路如下：
- 每个组群有一个独立的总选择器(Agent Selector)，用于管理该组群的指令和回复。
- 每个选择器有若干个层级(Level)，在查找可用指令(Cmd)时，从最高层级(0)开始查找，直到找到可用指令为止。
- 每个层级有若干个选择器(Selector)，每个选择器中包含若干个选项(Option)，每个选项对应一个指令集(CmdSet)。每个选择器只能同时选择一个选项，选择后其余选项中的指令集将被屏蔽。
- 指令集是由若干个应用(App)中的指令组成的，一个应用只能归属于一个指令集，一个指令集可以包含多个应用的指令。

数据结构如下：
```go
type AgentSelector struct {
    Session   wba.Session
    Selectors []Selector
}

type Selector struct {
    Level    int
    Name     string
    Options  []*Option
}

type Option struct {
    Name    string
    CmdSets map[string]CmdSet
}

type CmdSet struct {
    map[string]wba.Command
}

type Cmmand struct {
    Name string
    Desc string
    App  AppKey
    Solve func(*wba.Context)
}
```