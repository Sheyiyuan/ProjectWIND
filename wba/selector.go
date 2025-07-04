package wba

import (
	"ProjectWIND/LOG"
	"encoding/json"
	"fmt"
)

type CmdAgentSelector struct {
	Session     SessionInfo
	Selectors   [4]map[SelectorLabel]Selector
	cmdTable    map[CmdLabel]map[CmdSetLabel]struct{}
	cmdSetTable map[CmdSetLabel]map[OptionLabel]struct{}
	optionTable map[OptionLabel][]map[SelectorLabel]struct{}
}

type Selector struct {
	Level         Priority
	Name          SelectorLabel
	Options       map[OptionLabel]Option
	currentOption OptionLabel
}

type Option struct {
	Name    string
	CmdSets map[CmdSetLabel]CmdSet
}

type CmdSet struct {
	Name    CmdSetLabel
	Cmds    map[CmdLabel]CmdInfo
	Enabled bool
}

type CmdInfo struct {
	Name    CmdLabel
	Cmd     Cmd
	Enabled bool
}

// NewCmdAgentSelector 创建一个新的 CmdAgentSelector 实例并初始化。
//
// 参数:
//
//   - session: 所属会话信息。
func NewCmdAgentSelector() CmdAgentSelector {
	c := CmdAgentSelector{
		Session:     SessionInfo{},
		Selectors:   [4]map[SelectorLabel]Selector{},
		cmdTable:    make(map[CmdLabel]map[CmdSetLabel]struct{}),
		cmdSetTable: make(map[CmdSetLabel]map[OptionLabel]struct{}),
		optionTable: make(map[OptionLabel][]map[SelectorLabel]struct{}),
	}
	for i := 0; i < 4; i++ {
		c.Selectors[i] = make(map[SelectorLabel]Selector)
		c.optionTable[fmt.Sprintf("%d", i)] = make([]map[SelectorLabel]struct{}, 4)
	}
	return c
}

func (a *CmdAgentSelector) LoadConfig(session SessionInfo, config string) error {
	a.Session = session
	conf := CmdAgentSelector{}
	err := json.Unmarshal([]byte(config), &conf)
	if err != nil {
		_ = json.Unmarshal([]byte("{}"), &conf)
	}
	for i := 0; i < 4; i++ {
		for selectorName, selector := range a.Selectors[i] {
			selector.currentOption = conf.Selectors[i][selectorName].currentOption
			for optionName, option := range selector.Options {
				for cmdSetName, cmdSet := range option.CmdSets {
					for cmdName, cmd := range cmdSet.Cmds {
						cmd.Enabled = conf.Selectors[i][selectorName].Options[optionName].CmdSets[cmdSetName].Cmds[cmdName].Enabled
					}
				}
			}
		}
	}
	return err
}

// AddCmdMap 添加一个命令映射(map[AppKey]CmdList)到 CmdAgentSelector 中。
//
// 参数:
//
//   - cmdMap: 命令映射，键为 AppKey，值为 CmdList。
func (a *CmdAgentSelector) AddCmdMap(cmdMap map[AppKey]CmdList) {
	for appKey, cmds := range cmdMap {
		a.AddCmdSet(appKey, cmds)
	}
}

// AddCmdSet 添加一个命令集(即APP)到 CmdAgentSelector 中。
//
// 参数:
//
//   - appKey: 应用键，包含应用名称、级别、选择器和选项。
//   - cmds: 命令列表，键为命令名称，值为 CmdInfo。
func (a *CmdAgentSelector) AddCmdSet(appKey AppKey, cmds CmdList) {

	for cmdName, cmd := range cmds {
		// 确保所有必要的map已初始化
		if a.cmdTable == nil {
			a.cmdTable = make(map[CmdLabel]map[CmdSetLabel]struct{})
		}
		if a.cmdSetTable == nil {
			a.cmdSetTable = make(map[CmdSetLabel]map[OptionLabel]struct{})
		}
		if a.optionTable == nil {
			a.optionTable = make(map[OptionLabel][]map[SelectorLabel]struct{})
		}

		// 初始化嵌套map
		if a.cmdTable[cmdName] == nil {
			a.cmdTable[cmdName] = make(map[CmdSetLabel]struct{})
		}
		if a.cmdSetTable[appKey.Name] == nil {
			a.cmdSetTable[appKey.Name] = make(map[OptionLabel]struct{})
		}
		if a.optionTable[appKey.Option] == nil {
			a.optionTable[appKey.Option] = make([]map[SelectorLabel]struct{}, 4)
		}
		if a.optionTable[appKey.Option][appKey.Level] == nil {
			a.optionTable[appKey.Option][appKey.Level] = make(map[SelectorLabel]struct{})
		}

		// 确保Selector层级的map已初始化
		if a.Selectors[appKey.Level] == nil {
			a.Selectors[appKey.Level] = make(map[SelectorLabel]Selector)
		}
		if a.Selectors[appKey.Level][appKey.Selector].Options == nil {
			selector := a.Selectors[appKey.Level][appKey.Selector]
			selector.Options = make(map[OptionLabel]Option)
			a.Selectors[appKey.Level][appKey.Selector] = selector
		}
		if a.Selectors[appKey.Level][appKey.Selector].Options[appKey.Option].CmdSets == nil {
			selector := a.Selectors[appKey.Level][appKey.Selector]
			option := selector.Options[appKey.Option]
			option.CmdSets = make(map[CmdSetLabel]CmdSet)
			selector.Options[appKey.Option] = option
			a.Selectors[appKey.Level][appKey.Selector] = selector
		}
		if a.Selectors[appKey.Level][appKey.Selector].Options[appKey.Option].CmdSets[appKey.Name].Cmds == nil {
			selector := a.Selectors[appKey.Level][appKey.Selector]
			option := selector.Options[appKey.Option]
			cmdSet := option.CmdSets[appKey.Name]
			cmdSet.Cmds = make(map[CmdLabel]CmdInfo)
			option.CmdSets[appKey.Name] = cmdSet
			selector.Options[appKey.Option] = option
			a.Selectors[appKey.Level][appKey.Selector] = selector
		}
		a.cmdTable[cmdName][appKey.Name] = struct{}{}
		a.cmdSetTable[appKey.Name][appKey.Option] = struct{}{}
		a.optionTable[appKey.Option][appKey.Level][appKey.Selector] = struct{}{}
		a.Selectors[appKey.Level][appKey.Selector].Options[appKey.Option].CmdSets[appKey.Name].Cmds[cmdName] = CmdInfo{
			Cmd:     cmd,
			Enabled: true,
		}
		LOG.Debug("add cmd: %s, app: %s, option: %s, cmdSet: %s", cmdName, appKey.Name, appKey.Option, appKey.Selector)
	}
	a.Selectors[appKey.Level][appKey.Selector].Options[appKey.Option].CmdSets[appKey.Name] = CmdSet{
		Name:    appKey.Name,
		Cmds:    a.Selectors[appKey.Level][appKey.Selector].Options[appKey.Option].CmdSets[appKey.Name].Cmds,
		Enabled: true,
	}
	if a.Selectors[appKey.Level][appKey.Selector].currentOption == "" {
		a.Selectors[appKey.Level][appKey.Selector] = Selector{
			Level:         appKey.Level,
			Name:          appKey.Selector,
			Options:       a.Selectors[appKey.Level][appKey.Selector].Options,
			currentOption: appKey.Option,
		}
	}
}

// FindCmd 根据命令名称和是否只查找启用的命令，查找符合条件的命令。
//
// 参数:
//
//   - cmdName: 要查找的命令名称。
//   - onlyEnabled: 是否只查找启用的命令。
//
// 返回值:
//
//   - result: 找到的命令列表（越靠前的优先级越高）。
//   - count: 找到的命令数量。
func (a *CmdAgentSelector) FindCmd(cmdName string, onlyEnabled bool) (result []Cmd, count int) {
	result = make([]Cmd, 0)
	count = 0
	cmdSets := a.cmdTable[cmdName]
	for cmdSet := range cmdSets {
		options := a.cmdSetTable[cmdSet]
		for option := range options {
			selectors := a.optionTable[option]
			for i := 0; i < len(a.Selectors); i++ {
				for selector := range selectors[i] {
					if option != a.Selectors[i][selector].currentOption && onlyEnabled {
						continue
					}
					cmdInfo := a.Selectors[i][selector].Options[option].CmdSets[cmdSet].Cmds[cmdName]
					if onlyEnabled && !cmdInfo.Enabled {
						continue
					}
					result = append(result, cmdInfo.Cmd)
					count++
				}

			}

		}
	}
	return result, count
}

// FindCmdSet 根据命令集名称和是否只查找启用的命令集，查找符合条件的命令集。
//
// 参数:
//
//   - cmdSetName: 要查找的命令集名称。
//   - onlyEnabled: 是否只查找启用的命令集。
//
// 返回值:
//
//   - result: 找到的命令集列表（越靠前的优先级越高）。
//   - count: 找到的命令集数量。
func (a *CmdAgentSelector) FindCmdSet(cmdSetName string, onlyEnabled bool) (result []CmdSet, count int) {
	result = make([]CmdSet, 0)
	count = 0
	options := a.cmdSetTable[cmdSetName]
	for option := range options {
		selectors := a.optionTable[option]
		for i := 0; i < len(a.Selectors); i++ {
			for selector := range selectors[i] {
				if option != a.Selectors[i][selector].currentOption && onlyEnabled {
					continue
				}
				cmdSet := a.Selectors[i][selector].Options[option].CmdSets[cmdSetName]
				result = append(result, cmdSet)
				count++
			}
		}
	}
	return result, count
}

// FindOption 根据选项名称，查找符合条件的选项。
//
// 参数:
//
//   - optionName: 要查找的选项名称。
//
// 返回值:
//
//   - result: 找到的选项列表（越靠前的优先级越高）。
//   - count: 找到的选项数量。
func (a *CmdAgentSelector) FindOption(optionName string) (result []Option, count int) {
	result = make([]Option, 0)
	count = 0
	selectors := a.optionTable[optionName]
	for i := 0; i < len(a.Selectors); i++ {
		for selector := range selectors[i] {
			if optionName != a.Selectors[i][selector].currentOption {
				continue
			}
			option := a.Selectors[i][selector].Options[optionName]
			result = append(result, option)
			count++
		}
	}
	return result, count
}

// FindSelector 根据选择器名称，查找符合条件的选择器。
//
// 参数:
//
//   - selectorName: 要查找的选择器名称。
//
// 返回值:
//
//   - result: 找到的选择器列表（越靠前的优先级越高）。
//   - count: 找到的选择器数量。
func (a *CmdAgentSelector) FindSelector(selectorName string) (result []Selector, count int) {
	result = make([]Selector, 0)
	count = 0
	selectors := a.optionTable[selectorName]
	for i := 0; i < len(a.Selectors); i++ {
		for selector := range selectors[i] {
			if selectorName != a.Selectors[i][selector].currentOption {
				continue
			}
			selector := a.Selectors[i][selector]
			result = append(result, selector)
			count++
		}
	}
	return result, count
}

// GetCurrentOption 获取当前选择器的当前选项。
//
// 参数:
//
//   - level: 选择器的级别。
//   - selectorName: 选择器的名称。
//
// 返回值:
//
//   - result: 当前选项，当isExist为false时，此结果为 Option 对应的空值。
//   - isExist: 是否存在当前选项。
func (a *CmdAgentSelector) GetCurrentOption(level int, selectorName string) (result Option, isExist bool) {
	result, ok := a.Selectors[level][selectorName].Options[a.Selectors[level][selectorName].currentOption]
	if !ok {
		return Option{}, false
	}
	return result, true
}

// ToggleCmd 切换命令的启用状态。
//
// 参数:
//
//   - enabled: 要设置的启用状态。
//   - level: 选择器的级别。
//   - selectorName: 选择器的名称。
//   - optionName: 选项的名称。
//   - cmdSetName: 命令集的名称。
//   - cmdName: 命令的名称。
//
// 返回值:
//
//   - state: 命令的启用状态。
//   - preState: 命令的之前的启用状态。
//   - err: 错误信息，当发生错误时返回。
func (a *CmdAgentSelector) ToggleCmd(enabled bool, level Priority, selectorName SelectorLabel, optionName OptionLabel, cmdSetName CmdSetLabel, cmdName CmdLabel) (state bool, preState bool, err error) {
	if level < 0 || level > 3 {
		return false, false, fmt.Errorf("level must be between 0 and 3")
	}
	selector, ok := a.Selectors[level][selectorName]
	if !ok {
		return false, false, fmt.Errorf("selector %s not found", selectorName)
	}
	option, ok := selector.Options[optionName]
	if !ok {
		return false, false, fmt.Errorf("option %s not found", optionName)
	}
	cmdSet, ok := option.CmdSets[cmdName]
	if !ok {
		return false, false, fmt.Errorf("cmdSet %s not found", cmdName)
	}
	cmdInfo, ok := cmdSet.Cmds[cmdName]
	if !ok {
		return false, false, fmt.Errorf("cmd %s not found", cmdName)
	}
	preState = cmdInfo.Enabled
	cmdInfo.Enabled = enabled
	a.Selectors[level][selectorName].Options[optionName].CmdSets[cmdSetName].Cmds[cmdName] = cmdInfo
	return enabled, preState, nil
}
