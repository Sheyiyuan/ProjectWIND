package database

import (
	"ProjectWIND/LOG"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

const address = "./data/database/datamaps.wdb"
const core = "./data/core.json"

type Errno struct {
	Code      int
	Condition string
}

type unit struct {
	Id   string
	Data map[string]string
}

type User unit
type Group unit
type Global unit

type Configs struct {
	Number       map[string]int64
	String       map[string]string
	Float        map[string]float64
	Number_Slice map[string][]int64
	String_Slice map[string][]string
	Hash         string
}

type Datamap struct {
	Id         string
	Permission string
	Users      map[string]User
	Groups     map[string]Group
	Global     map[string]Global
	Configs    Configs
}

type Database struct {
	Datamaps map[string]Datamap
}

func (eno *Errno) Log() {
	// 记录错误日志
	switch eno.Code / 100 {
	case 1:
		LOG.Fatal(eno.Condition)
	case 2:
		LOG.Error(eno.Condition)
	case 3:
		LOG.Warn(eno.Condition)
	case 4:
		LOG.Info(eno.Condition)
	case 5:
		LOG.Notice(eno.Condition)
	case 6:
		LOG.Debug(eno.Condition)
	}
}

func newDatamap(id string) Datamap {
	// 创建数据表
	db := &Datamap{
		Id:         id,
		Permission: "private",
		Users:      make(map[string]User),
		Groups:     make(map[string]Group),
		Global:     make(map[string]Global),
		Configs: Configs{
			Number:       make(map[string]int64),
			String:       make(map[string]string),
			Float:        make(map[string]float64),
			Number_Slice: make(map[string][]int64),
			String_Slice: make(map[string][]string),
			Hash:         "",
		},
	}
	return *db
}

func newDatabase() Database {
	// 创建数据库
	db := &Database{
		Datamaps: make(map[string]Datamap),
	}
	return *db
}

func (dbh *Database) addDatamap(id string) {
	// 创建新数据表
	db := newDatamap(id)
	dbh.Datamaps[id] = db
}

func folderCheck(filename string) (Error Errno) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		err := os.MkdirAll(filename, 0755)
		if err != nil {
			return Errno{101, fmt.Sprintf("创建文件夹时出错：%v", err)}
		}
	}
	return Errno{0, ""}
}

func fileCheck(filename string) (Error Errno) {
	// 检查并创建文件
	dir := filepath.Dir(filename)
	folderCheck(dir)
	eno := Errno{0, ""}
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		file, err := os.Create(filename)
		if err != nil {
			return Errno{101, fmt.Sprintf("创建文件时出错: %v", err)}
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				eno = Errno{302, fmt.Sprintf("关闭文件时出错: %v", err)}
			}
		}(file)
	}
	return eno
}

func writeContent(f *os.File, str string) (errno Errno) {
	// 写入内容到文件
	if f == nil {
		return Errno{101, "文件不存在"}
	}
	_, err := f.Write([]byte(str))
	if err != nil {
		return Errno{101, fmt.Sprintf("无法写入到文件: %v", err)}
	}
	return Errno{0, ""}
}

func printContent(file string) (text string, errno Errno) {
	// 读取文件内容
	bytes, err := os.ReadFile(file)
	if err == nil {
		return string(bytes), Errno{0, ""}
	} else {
		return "", Errno{101, fmt.Sprintf("读取文件时出错: %v", err)}
	}
}

func getCorePassword() (pass string, errno Errno) {
	// 获取核心密码
	filename := core
	fileCheck(filename)
	dataJson, err := printContent(filename)
	if err.Code != 0 {
		return "", Errno{101, fmt.Sprintf("读取文件时出错 %s: %v", filename, err)}
	}
	config := make(map[string]string)
	err2 := json.Unmarshal([]byte(dataJson), config)
	if err2 != nil {
		return "", Errno{201, fmt.Sprintf("反序列化时出错: %v", err2)}
	}
	password, ok := config["password"]
	if !ok {
		return "", Errno{601, "core.json中未找到配置密码项"}
	}
	return password, Errno{0, ""}
}

func saveData(db *Database) (errno Errno) {
	// 保存数据到文件
	dataJson, err := json.Marshal(db)
	if err != nil {
		return Errno{201, fmt.Sprintf("序列化数据时出错: %v", err)}
	}
	filename := address
	file, err := os.Create(filename)
	if err != nil {
		return Errno{101, fmt.Sprintf("创建文件时出错 %s: %v", filename, err)}
	}
	err2 := writeContent(file, string(dataJson))
	if err2.Code != 0 {
		return err2
	}
	return Errno{0, ""}
}

func loadData(db *Database) (errno Errno) {
	// 读取配置文件
	filename := address
	fileCheck(filename)
	dataJson, err := printContent(filename)
	if err.Code != 0 {
		return err
	}
	err2 := json.Unmarshal([]byte(dataJson), db)
	if err2 != nil {
		return Errno{201, fmt.Sprintf("反序列化数据时出错: %v", err2)}
	}
	return Errno{0, ""}
}

var DB *Database

func dataSet(datamap string, unit string, id string, key string, value interface{}, isAllowed bool, isMaster bool) (errno Errno) {
	// 修改数据
	dm, ok := DB.Datamaps[datamap]
	if !ok {
		// 创建新数据表
		DB.addDatamap(datamap)
		dm = DB.Datamaps[datamap]
	}
	if !isAllowed && !isMaster && dm.Permission != "private" {
		return Errno{301, "访问权限不足"}
	}
	if !isMaster && dm.Permission == "master" {
		return Errno{301, "访问权限不足"}
	}
	switch unit {
	case "config":
		switch id {
		case "number":
			valueInt64, ok := value.(int64)
			if !ok {
				return Errno{303, "设置的配置值无法被断言为int64类型"}
			}
			dm.Configs.Number[key] = valueInt64
		case "string":
			valueStr, ok := value.(string)
			if !ok {
				return Errno{303, "设置的配置值无法被断言为string类型"}
			}
			dm.Configs.String[key] = valueStr
		case "float":
			valueFloat64, ok := value.(float64)
			if !ok {
				return Errno{303, "设置的配置值无法被断言为float64类型"}
			}
			dm.Configs.Float[key] = valueFloat64
		case "number_slice":
			valueInt64Slice, ok := value.([]int64)
			if !ok {
				return Errno{303, "设置的配置值无法被断言为[]int64类型"}
			}
			dm.Configs.Number_Slice[key] = valueInt64Slice
		case "string_slice":
			valueStrSlice, ok := value.([]string)
			if !ok {
				return Errno{303, "设置的配置值无法被断言为[]string类型"}
			}
			dm.Configs.String_Slice[key] = valueStrSlice
		case "hash":
			valueStr, ok := value.(string)
			if !ok {
				return Errno{303, "设置的配置值无法被断言为string类型"}
			}
			dm.Configs.Hash = valueStr
		default:
			return Errno{304, "不合法的配置项类型"}
		}
	case "user":
		valueStr, ok := value.(string)
		if !ok {
			return Errno{303, "变量值无法被断言为string类型"}
		}
		user, ok := dm.Users[id]
		if !ok {
			dm.Users[id] = User{
				Id:   id,
				Data: make(map[string]string),
			}
			user = dm.Users[id]
		}
		if user.Data == nil {
			user.Data = make(map[string]string)
		}
		user.Data[key] = valueStr
	case "group":
		valueStr, ok := value.(string)
		if !ok {
			return Errno{303, "变量值无法被断言为string类型"}
		}
		group, ok := dm.Groups[id]
		if !ok {
			dm.Groups[id] = Group{
				Id:   id,
				Data: make(map[string]string),
			}
			group = dm.Groups[id]
		}
		if group.Data == nil {
			group.Data = make(map[string]string)
		}
		group.Data[key] = valueStr
	case "global":
		valueStr, ok := value.(string)
		if !ok {
			return Errno{303, "变量值无法被断言为string类型"}
		}
		global, ok := dm.Global[id]
		if !ok {
			dm.Global[id] = Global{
				Id:   id,
				Data: make(map[string]string),
			}
			global = dm.Global[id]
		}
		if global.Data == nil {
			global.Data = make(map[string]string)
		}
		global.Data[key] = valueStr // 使用断言后的string值
	default:
		return Errno{304, "不合法的数据单元"}
	}
	return Errno{0, ""}
}

func dataGet(datamap string, unit string, id string, key string, isAllowed bool, isMaster bool) (res interface{}, errno Errno) {
	dm, ok := DB.Datamaps[datamap]
	if !ok {
		return "", Errno{601, fmt.Sprintf("数据表 %s 不存在", datamap)}
	}
	if !isAllowed && !isMaster && dm.Permission != "private" {
		return "", Errno{301, "访问权限不足"}
	}
	if !isMaster && dm.Permission == "master" {
		return "", Errno{301, "访问权限不足"}
	}
	switch unit {
	case "config":
		switch id {
		case "number":
			value, ok := dm.Configs.Number[key]
			if !ok {
				return 0, Errno{601, fmt.Sprintf("配置项不存在%s", key)}
			}
			return value, Errno{0, ""}
		case "string":
			value, ok := dm.Configs.String[key]
			if !ok {
				return "", Errno{601, fmt.Sprintf("配置项不存在%s", key)}
			}
			return value, Errno{0, ""}
		case "float":
			value, ok := dm.Configs.Float[key]
			if !ok {
				return 0.0, Errno{601, fmt.Sprintf("配置项不存在%s", key)}
			}
			return value, Errno{0, ""}
		case "number_slice":
			value, ok := dm.Configs.Number_Slice[key]
			if !ok {
				return []int64{}, Errno{601, fmt.Sprintf("配置项不存在%s", key)}
			}
			return value, Errno{0, ""}
		case "string_slice":
			value, ok := dm.Configs.String_Slice[key]
			if !ok {
				return []string{}, Errno{601, fmt.Sprintf("配置项不存在%s", key)}
			}
			return value, Errno{0, ""}
		case "hash":
			return dm.Configs.Hash, Errno{0, ""}
		default:
			return "", Errno{304, "不合法的配置项类型"}
		}
	case "user":
		user, ok := dm.Users[id]
		if !ok {
			return "", Errno{601, fmt.Sprintf("用户 %s 不存在", id)}
		}
		if user.Data == nil {
			return "", Errno{601, fmt.Sprintf("用户 %s 的数据显示为nil", id)}
		}
		value, ok := user.Data[key]
		if !ok {
			return "", Errno{601, fmt.Sprintf("用户 %s 的数据中键 %s 不存在", id, key)}
		}
		return value, Errno{0, ""}
	case "group":
		group, ok := dm.Groups[id]
		if !ok {
			return "", Errno{601, fmt.Sprintf("群组 %s 不存在", id)}
		}
		if group.Data == nil {
			return "", Errno{601, fmt.Sprintf("群组 %s 的数据显示为nil", id)}
		}
		value, ok := group.Data[key]
		if !ok {
			return "", Errno{601, fmt.Sprintf("群组 %s 的数据中键 %s 不存在", id, key)}
		}
		return value, Errno{0, ""}
	case "global":
		global, ok := dm.Global[id]
		if !ok {
			return "", Errno{601, fmt.Sprintf("全局变量 %s 不存在", id)}
		}
		if global.Data == nil {
			return "", Errno{601, fmt.Sprintf("全局变量 %s 的数据显示为nil", id)}
		}
		value, ok := global.Data[key]
		if !ok {
			return "", Errno{601, fmt.Sprintf("全局变量 %s 的数据中键 %s 不存在", id, key)}
		}
		return value, Errno{0, ""}
	default:
		return "", Errno{304, "不合法的数据单元"}
	}
}

func initializeDatabase() *Database {
	// 启动并检查程序
	LOG.Info("正在启动数据库")
	db := newDatabase()
	loadData(&db)
	LOG.Info("数据库启动完成")
	return &db
}

func Start() {
	DB = initializeDatabase()
	// 创建一个通道用于接收信号
	dataChan := make(chan os.Signal, 1)
	// 监听指定的信号，如SIGINT (Ctrl+C) 和 SIGTERM
	signal.Notify(dataChan, syscall.SIGINT, syscall.SIGTERM)

	// 定义一个Ticker用于每1小时触发一次保存操作
	saveTicker := time.NewTicker(3600 * time.Second)
	defer saveTicker.Stop()

	// 启动一个goroutine等待信号和定时保存
	go func() {
		for {
			select {
			case <-dataChan:
				// 接收到信号，保存数据并退出程序
				LOG.Info("关闭中，正在保存数据")
				saveData(DB)
				os.Exit(0)
			case <-saveTicker.C:
				// 定时保存数据
				LOG.Info("自动保存数据")
				saveData(DB)
			}
		}
	}()

	select {}
}

// 修改数据表权限（核心）
func MasterSetDatamapPermission(datamap string, premission string) {
	db, ok := DB.Datamaps[datamap]
	if !ok {
		eno := Errno{601, fmt.Sprintf("数据表 %s 不存在", datamap)}
		eno.Log()
		return
	}
	db.Permission = premission
	DB.Datamaps[datamap] = db
}

func CreatePublicDatamap(appName string, id string) (errno Errno) {
	// 查询权限
	hash, eno := getCorePassword()
	if eno.Code != 0 {
		return eno
	}
	if hash == "" {
		// 删除数据表哈希
		dataSet(appName, "config", "hash", "", "", true, true)
	}
	datahash, eno := dataGet(appName, "config", "hash", "", true, true)
	if eno.Code != 0 {
		return eno
	}
	if hash != datahash {
		eno := Errno{301, "应用没有创建公开数据表的权限"}
		return eno
	}

	// 创建公开数据表
	db, ok := DB.Datamaps[id]
	if !ok {
		db = newDatamap(id)
		db.Permission = "public"
		DB.Datamaps[id] = db
	} else {
		eno := Errno{601, fmt.Sprintf("数据表 %s 已经存在", id)}
		return eno
	}
	return Errno{0, ""}
}

func MasterCreatePublicDatamap(id string) {
	// 创建核心数据表
	db, ok := DB.Datamaps[id]
	if !ok {
		db = newDatamap(id)
		db.Permission = "master"
		DB.Datamaps[id] = db
	} else {
		eno := Errno{601, fmt.Sprintf("数据表 %s 已经存在", id)}
		eno.Log()
	}
}

func MasterCreateMasterDatamap(id string) {
	// 创建公开数据表
	db, ok := DB.Datamaps[id]
	if !ok {
		db = newDatamap(id)
		db.Permission = "public"
		DB.Datamaps[id] = db
	} else {
		eno := Errno{601, fmt.Sprintf("数据表 %s 已经存在", id)}
		eno.Log()
	}
}

// 修改数据（核心）
func MasterSet(datamap string, unit string, id string, key string, value interface{}) {
	eno := dataSet(datamap, unit, id, key, value, true, true)
	if eno.Code != 0 {
		eno.Log()
	}
}

// 查询数据（核心）
func MasterGet(datamap string, unit string, id string, key string) (interface{}, bool) {
	val, eno := dataGet(datamap, unit, id, key, true, true)
	if eno.Code != 0 {
		eno.Log()
		return "", false
	}
	return val, true
}

func Get(appName string, datamap string, unit string, id string, key string, isGettingConfig bool) (value interface{}, errno Errno) {
	// 查询数据
	if unit == "config" && id == "hash" {
		// app不允许访问hash数据
		eno := Errno{301, fmt.Sprintf("应用 %s 不允许访问配置项信息", appName)}
		eno.Log()
	}
	if !isGettingConfig && unit == "config" {
		// 不允许在非config数据表中访问config数据
		eno := Errno{301, fmt.Sprintf("应用 %s 不能在常规读写中访问配置项信息，请使用配置项读取功能", appName)}
		eno.Log()
	}
	if appName != datamap {
		// 需要master密码来访问其他app的数据
		hash, eno := getCorePassword()
		if eno.Code != 0 {
			return "", eno
		}
		if hash == "" {
			// 删除数据表哈希
			dataSet(appName, "config", "hash", "", "", true, true)
		}
		datahash, eno := dataGet(appName, "config", "hash", "", true, true)
		if eno.Code != 0 {
			eno.Log()
		}
		if hash != datahash {
			value, eno := dataGet(appName, unit, id, key, false, false)
			if eno.Code != 0 {
				return value, eno
			}
			return value, Errno{0, ""}
		}
	}
	value, eno := dataGet(appName, unit, id, key, true, false)
	if eno.Code != 0 {
		return value, eno
	}
	return value, Errno{0, ""}
}

func Set(appName string, datamap string, unit string, id string, key string, value interface{}) (errno Errno) {
	// 修改数据
	if unit == "config" {
		// app不允许修改config数据
		return Errno{301, fmt.Sprintf("应用 %s 不允许修改配置项信息", appName)}
	}
	if appName != datamap {
		// 需要master密码来访问其他app的数据
		hash, eno := getCorePassword()
		if eno.Code != 0 {
			return eno
		}
		if hash == "" {
			// 删除数据表哈希
			eno := dataSet(appName, "config", "hash", "", "", true, true)
			if eno.Code != 0 {
				return eno
			}
		}
		datahash, eno := dataGet(appName, "config", "hash", "", true, true)
		if eno.Code != 0 {
			return eno
		}
		if hash != datahash {
			eno := dataSet(appName, unit, id, key, value, false, false)
			if eno.Code != 0 {
				return eno
			}
			return Errno{0, ""}
		}
	}
	eno := dataSet(appName, unit, id, key, value, true, false)
	if eno.Code != 0 {
		return eno
	}
	return Errno{0, ""}
}
