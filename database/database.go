package database

import (
	"ProjectWIND/LOG"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

const address = "./data/database/datamaps.wdb"
const core = "./data/core.json"

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

func (this *Database) addDatamap(id string) {
	// 创建新数据表
	db := newDatamap(id)
	this.Datamaps[id] = db
}

// func hash(word string) string {
// 	hash, err := bcrypt.GenerateFromPassword([]byte(word), bcrypt.DefaultCost)
// 	if err != nil {
// 		LOG.Error("[Error]Error while hash password: %v", err)
// 		return ""
// 	}
// 	return string(hash)
// }

// func hashCheck(word string, hash string) bool {
// 	err := bcrypt.CompareHashAndPassword(hash, []byte(word))
// 	return err == nil
// }

func folderCheck(filename string) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		err := os.MkdirAll(filename, 0755)
		if err != nil {
			LOG.Fatal("[Error]Error occured while create folder: %v", err)
		}
	}
}

func fileCheck(filename string) {
	// 检查并创建文件
	dir := filepath.Dir(filename)
	folderCheck(dir)
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		file, err := os.Create(filename)
		if err != nil {
			LOG.Fatal("[Error]Error occured while create file: %v", err)
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				LOG.Fatal("[Error]Error occured while close file: %v", err)
			}
		}(file)
	}
}

func writeContent(f *os.File, str string) error {
	// 写入内容到文件
	if f == nil {
		// log.Printf("[Error]file is nil")
		LOG.Error("[Error]file is nil")
		return errors.New("file is nil")
	}
	_, err := f.Write([]byte(str))
	if err != nil {
		LOG.Error("[Error]Error while write content to file: %v", err)
		return err
	}
	return nil
}

func printContent(file string) (string, error) {
	// 读取文件内容
	bytes, err := os.ReadFile(file)
	if err == nil {
		return string(bytes), nil
	} else {
		return "", err
	}
}

func getCorePassword() string {
	// 获取核心密码
	filename := core
	fileCheck(filename)
	dataJson, err := printContent(filename)
	if err != nil {
		LOG.Error("[Error]Error while read file %s: %v", filename, err)
		return ""
	}
	config := make(map[string]string)
	err = json.Unmarshal([]byte(dataJson), &config)
	if err != nil {
		LOG.Error("[Error]Error while unmarshal data: %v", err)
		return ""
	}
	password, ok := config["password"]
	if !ok {
		LOG.Warn("[Warning]Password not found in core.json")
		return ""
	}
	return password
}

func saveData(db *Database) error {
	// 保存数据到文件
	dataJson, err := json.Marshal(db)
	if err != nil {
		LOG.Error("[Error]:Error while marshal data: %v", err)
		return err
	}
	filename := address
	file, err := os.Create(filename)
	if err != nil {
		LOG.Error("[Error]:Error while create file %s: %v", filename, err)
		return err
	}
	writeContent(file, string(dataJson))
	return nil
}

func loadData(db *Database) error {
	// 读取配置文件
	filename := address
	fileCheck(filename)
	dataJson, err := printContent(filename)
	if err != nil {
		// log.Printf("[Error]:Error while read file %s: %v", filename, err)
		LOG.Error("[Error]:Error while read file %s: %v", filename, err)
		return err
	}
	err = json.Unmarshal([]byte(dataJson), db)
	if err != nil {
		// log.Printf("[Error]:Error while unmarshal data: %v", err)
		LOG.Warn("[Warning]:Error while unmarshal data: %v", err)
		return err
	}
	return nil
}

var DB *Database

func dataSet(datamap string, unit string, id string, key string, value interface{}, isAllowed bool, isMaster bool) {
	// 修改数据
	dm, ok := DB.Datamaps[datamap]
	if !ok {
		// 创建新数据表
		DB.addDatamap(datamap)
		dm = DB.Datamaps[datamap]
	}
	if !isAllowed && !isMaster && dm.Permission != "private" {
		LOG.Warn("[Warning]:Permission denied")
		return
	}
	if !isMaster && dm.Permission == "master" {
		LOG.Warn("[Warning]:Permission denied")
		return
	}
	switch unit {
	case "config":
		switch id {
		case "number":
			valueInt64, ok := value.(int64) // 断言value为int64类型
			if !ok {
				LOG.Error("[Error]:Value cannot be asserted to int64")
				return
			}
			dm.Configs.Number[key] = valueInt64 // 使用断言后的int64值
		case "string":
			valueStr, ok := value.(string) // 断言value为string类型
			if !ok {
				LOG.Error("[Error]:Value cannot be asserted to string")
				return
			}
			dm.Configs.String[key] = valueStr // 使用断言后的string值
		case "float":
			valueFloat64, ok := value.(float64) // 断言value为float64类型
			if !ok {
				LOG.Error("[Error]:Value cannot be asserted to float64")
				return
			}
			dm.Configs.Float[key] = valueFloat64 // 使用断言后的float64值
		case "number_slice":
			valueInt64Slice, ok := value.([]int64) // 断言value为[]int64类型
			if !ok {
				LOG.Error("[Error]:Value cannot be asserted to []int64")
				return
			}
			dm.Configs.Number_Slice[key] = valueInt64Slice // 使用断言后的[]int64值
		case "string_slice":
			valueStrSlice, ok := value.([]string) // 断言value为[]string类型
			if !ok {
				LOG.Error("[Error]:Value cannot be asserted to []string")
				return
			}
			dm.Configs.String_Slice[key] = valueStrSlice // 使用断言后的[]string值
		case "hash":
			valueStr, ok := value.(string) // 断言value为string类型
			if !ok {
				LOG.Error("[Error]:Value cannot be asserted to string")
				return
			}
			dm.Configs.Hash = valueStr // 使用断言后的string值
		default:
			LOG.Error("[Error]:Invalid id %s", id)
		}
	case "user":
		valueStr, ok := value.(string) // 断言value为string类型
		if !ok {
			LOG.Error("[Error]:Value cannot be asserted to string")
			return
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
		user.Data[key] = valueStr // 使用断言后的string值
	case "group":
		valueStr, ok := value.(string) // 断言value为string类型
		if !ok {
			LOG.Error("[Error]:Value cannot be asserted to string")
			return
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
		group.Data[key] = valueStr // 使用断言后的string值
	case "global":
		valueStr, ok := value.(string) // 断言value为string类型
		if !ok {
			LOG.Error("[Error]:Value cannot be asserted to string")
			return
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
		LOG.Error("[Error]:Invalid unit %s", unit)
	}
}

func dataGet(datamap string, unit string, id string, key string, isAllowed bool, isMaster bool) (interface{}, bool) {
	dm, ok := DB.Datamaps[datamap]
	if !ok {
		LOG.Warn("[Warning]:Datamap %s not found", datamap)
		return "", false
	}
	if !isAllowed && !isMaster && dm.Permission != "private" {
		LOG.Warn("[Warning]:Permission denied")
		return "", false
	}
	if !isMaster && dm.Permission == "master" {
		LOG.Warn("[Warning]:Permission denied")
		return "", false
	}
	switch unit {
	case "config":
		switch id {
		case "number":
			value, ok := dm.Configs.Number[key]
			if !ok {
				LOG.Warn("[Warning]:Config %s not found", key)
				return 0, false
			}
			return value, true
		case "string":
			value, ok := dm.Configs.String[key]
			if !ok {
				LOG.Warn("[Warning]:Config %s not found", key)
				return "", false
			}
			return value, true
		case "float":
			value, ok := dm.Configs.Float[key]
			if !ok {
				LOG.Warn("[Warning]:Config %s not found", key)
				return 0.0, false
			}
			return value, true
		case "number_slice":
			value, ok := dm.Configs.Number_Slice[key]
			if !ok {
				LOG.Warn("[Warning]:Config %s not found", key)
				return []int64{}, false
			}
			return value, true
		case "string_slice":
			value, ok := dm.Configs.String_Slice[key]
			if !ok {
				LOG.Warn("[Warning]:Config %s not found", key)
				return []string{}, false
			}
			return value, true
		case "hash":
			return dm.Configs.Hash, true
		default:
			LOG.Error("[Error]:Invalid id %s", id)
			return "", false
		}
	case "user":
		user, ok := dm.Users[id]
		if !ok {
			LOG.Warn("[Warning]:User %s not found", id)
			return "", false
		}
		if user.Data == nil {
			LOG.Warn("[Warning]:User %s's data is nil", id)
			return "", false
		}
		value, ok := user.Data[key]
		if !ok {
			LOG.Warn("[Warning]:User %s's data %s not found", id, key)
			return "", false
		}
		return value, true
	case "group":
		group, ok := dm.Groups[id]
		if !ok {
			LOG.Warn("[Warning]:Group %s not found", id)
			return "", false
		}
		if group.Data == nil {
			LOG.Warn("[Warning]:Group %s's data is nil", id)
			return "", false
		}
		value, ok := group.Data[key]
		if !ok {
			LOG.Warn("[Warning]:Group %s's data %s not found", id, key)
			return "", false
		}
		return value, true
	case "global":
		global, ok := dm.Global[id]
		if !ok {
			LOG.Warn("[Warning]:Global %s not found", id)
			return "", false
		}
		if global.Data == nil {
			LOG.Warn("[Warning]:Global data of %s is nil", id)
			return "", false
		}
		value, ok := global.Data[key]
		if !ok {
			LOG.Warn("[Warning]:Global data of %s's %s not found", id, key)
			return "", false
		}
		return value, true
	default:
		LOG.Error("[Error]:Invalid unit %s", unit)
		return "", false
	}
}

func initializeDatabase() *Database {
	// 启动并检查程序
	LOG.Info("Starting database ...")
	db := newDatabase()
	loadData(&db)
	LOG.Info("Database started successfully.")
	return &db
}

func Start() {
	DB = initializeDatabase()
	// 创建一个通道用于接收信号
	dataChan := make(chan os.Signal, 1)
	// 监听指定的信号，如SIGINT (Ctrl+C) 和 SIGTERM
	signal.Notify(dataChan, syscall.SIGINT, syscall.SIGTERM)

	// 定义一个Ticker用于每1小时触发一次保存操作
	saveTicker := time.NewTicker(600 * time.Second)
	defer saveTicker.Stop()

	// 启动一个goroutine等待信号和定时保存
	go func() {
		for {
			select {
			case <-dataChan:
				// 接收到信号，保存数据并退出程序
				fmt.Println("")
				LOG.Info("Received signal, saving data and exiting...")
				saveData(DB)
				os.Exit(0)
			case <-saveTicker.C:
				// 定时保存数据
				LOG.Info("Saving data automatically...")
				saveData(DB)
			}
		}
	}()

	select {} // 阻塞
}

func CreatePublicDatamap(appName string, id string) {
	// 查询权限
	hash := getCorePassword()
	if hash == "" {
		// 删除数据表哈希
		dataSet(appName, "config", "hash", "", "", true, true)
	}
	datahash, ok := dataGet(appName, "config", "hash", "", true, true)
	if !ok {
		LOG.Error("[Error]:Error while get hash of %s", appName)
		return
	}
	if hash != datahash {
		LOG.Warn("[Warning]:App %s is not allowed to create public datamap", appName)
		return
	}

	// 创建公开数据表
	db, ok := DB.Datamaps[id]
	if !ok {
		db = newDatamap(id)
		db.Permission = "public"
		DB.Datamaps[id] = db
	} else {
		LOG.Info("[Info]:Datamap %s already exists", id)
	}
}

func CreateMasterDatamap(id string) {
	// 创建核心数据表
	db, ok := DB.Datamaps[id]
	if !ok {
		db = newDatamap(id)
		db.Permission = "master"
		DB.Datamaps[id] = db
	} else {
		LOG.Info("[Info]:Datamap %s already exists", id)
	}
}

// 修改数据（核心）
func MasterSet(datamap string, unit string, id string, key string, value interface{}) {
	dataSet(datamap, unit, id, key, value, true, true)
}

// 查询数据（核心）
func MasterGet(datamap string, unit string, id string, key string) (interface{}, bool) {
	return dataGet(datamap, unit, id, key, true, true)
}

func Get(appName string, datamap string, unit string, id string, key string, isGettingConfig bool) (interface{}, bool) {
	// 查询数据
	if unit == "config" && id == "hash" {
		// app不允许访问hash数据
		LOG.Error("[Error]:App %s is not allowed to access hash data", appName)
		return "", false
	}
	if !isGettingConfig && unit == "config" {
		// 不允许在非config数据表中访问config数据
		LOG.Error("[Error]:App %s is not allowed to access config data", appName)
		return "", false
	}
	if appName != datamap {
		// 需要master密码来访问其他app的数据
		hash := getCorePassword()
		if hash == "" {
			// 删除数据表哈希
			dataSet(appName, "config", "hash", "", "", true, true)
		}
		datahash, ok := dataGet(appName, "config", "hash", "", true, true)
		if !ok {
			LOG.Error("[Error]:Error while get hash of %s", appName)
		}
		if hash != datahash {
			LOG.Warn("[Warning]:App %s is not allowed to access data of %s", appName, datamap)
			return dataGet(appName, unit, id, key, false, false)
		}

	}
	return dataGet(appName, unit, id, key, true, true)
}

func Set(appName string, datamap string, unit string, id string, key string, value interface{}) {
	// 修改数据
	if unit == "config" {
		// app不允许修改config数据
		LOG.Error("[Error]:App %s is not allowed to modify config data", appName)
		return
	}
	if appName != datamap {
		// 需要master密码来访问其他app的数据
		hash := getCorePassword()
		if hash == "" {
			// 删除数据表哈希
			dataSet(appName, "config", "hash", "", "", true, true)
		}
		datahash, ok := dataGet(appName, "config", "hash", "", true, true)
		if !ok {
			LOG.Error("[Error]:Error while get hash of %s", appName)
		}
		if hash != datahash {
			LOG.Warn("[Warning]:App %s is not allowed to access data of %s", appName, datamap)
			dataSet(appName, unit, id, key, value, false, false)
		}
	}
	dataSet(appName, unit, id, key, value, true, false)
}
