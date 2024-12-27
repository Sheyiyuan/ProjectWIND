package database

import (
	"ProjectWIND/LOG"
	"encoding/json"
	"errors"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

type unit struct {
	Id   string
	Data map[string]string
}

type User unit
type Group unit
type Global unit

type Database struct {
	Id     string
	Users  map[string]User
	Groups map[string]Group
	Global map[string]Global
	//...
	// Others map[string]map[string]unit
}

func newDatabase(id string) Database {
	// 创建数据库
	db := &Database{
		Id:     id,
		Users:  make(map[string]User),
		Groups: make(map[string]Group),
	}
	return *db
}

func folderCheck(filename string) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		err := os.MkdirAll(filename, 0755)
		if err != nil {
			LOG.FATAL("[ERROR]Error occurred while create folder: %v", err)
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
			LOG.FATAL("[ERROR]Error occurred while create file: %v", err)
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				LOG.FATAL("[ERROR]Error occurred while close file: %v", err)
			}
		}(file)
	}
}

func writeContent(f *os.File, str string) error {
	// 写入内容到文件
	if f == nil {
		// log.Printf("[ERROR]file is nil")
		LOG.ERROR("[ERROR]file is nil")
		return errors.New("file is nil")
	}
	_, err := f.Write([]byte(str))
	if err != nil {
		LOG.ERROR("[ERROR]Error while write content to file: %v", err)
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

func SaveData(db *Database) error {
	// 保存数据到文件
	dataJson, err := json.Marshal(db)
	if err != nil {
		LOG.ERROR("[ERROR]:Error while marshal data: %v", err)
		return err
	}
	filename := "./data/database/" + db.Id + ".txt"
	file, err := os.Create(filename)
	if err != nil {
		LOG.ERROR("[ERROR]:Error while create file %s: %v", filename, err)
		return err
	}
	err = writeContent(file, string(dataJson))
	if err != nil {
		return err
	}
	return nil
}

func loadData(db *Database) error {
	// 读取配置文件
	filename := "./data/database/" + db.Id + ".txt"
	fileCheck(filename)
	dataJson, err := printContent(filename)
	if err != nil {
		// log.Printf("[ERROR]:Error while read file %s: %v", filename, err)
		LOG.ERROR("[ERROR]:Error while read file %s: %v", filename, err)
		return err
	}
	err = json.Unmarshal([]byte(dataJson), db)
	if err != nil {
		// log.Printf("[ERROR]:Error while unmarshal data: %v", err)
		LOG.WARN("[WARNING]:Error while unmarshal data: %v", err)
		return err
	}
	return nil
}

func DataGet(db *Database, category string, id string, key string) (string, bool) {
	// 查询数据
	switch category {
	case "user":
		user, ok := db.Users[id]
		if !ok {
			LOG.WARN("[WARNING]:User %s not found", id)
			return "", false
		}
		if user.Data == nil {
			LOG.WARN("[WARNING]:User %s's data is nil", id)
			return "", false
		}
		value, ok := user.Data[key]
		if !ok {
			LOG.WARN("[WARNING]:User %s's data %s not found", id, key)
			return "", false
		}
		return value, true
	case "group":
		group, ok := db.Groups[id]
		if !ok {
			LOG.WARN("[WARNING]:Group %s not found", id)
			return "", false
		}
		if group.Data == nil {
			LOG.WARN("[WARNING]:Group %s's data is nil", id)
			return "", false
		}
		value, ok := group.Data[key]
		if !ok {
			LOG.WARN("[WARNING]:Group %s's data %s not found", id, key)
			return "", false
		}
		return value, true
	case "global":
		global, ok := db.Global[id]
		if !ok {
			LOG.WARN("[WARNING]:Global %s not found", id)
			return "", false
		}
		if global.Data == nil {
			LOG.WARN("[WARNING]:Global %s's data is nil", id)
			return "", false
		}
		value, ok := global.Data[key]
		if !ok {
			LOG.WARN("[WARNING]:Global %s's data %s not found", id, key)
			return "", false
		}
		return value, true
	default:
		LOG.ERROR("[ERROR]:Invalid category %s", category)
		return "", false
	}
}

func DataSet(db *Database, category string, id string, key string, value string) {
	// 修改数据
	switch category {
	case "user":
		user, ok := db.Users[id]
		if !ok {
			db.Users[id] = User{
				Id:   id,
				Data: make(map[string]string),
			}
			user = db.Users[id]
		}
		if user.Data == nil {
			user.Data = make(map[string]string)
		}
		user.Data[key] = value
	case "group":
		group, ok := db.Groups[id]
		if !ok {
			db.Groups[id] = Group{
				Id:   id,
				Data: make(map[string]string),
			}
			group = db.Groups[id]
		}
		if group.Data == nil {
			group.Data = make(map[string]string)
		}
		group.Data[key] = value
	case "global":
		global, ok := db.Global[id]
		if !ok {
			db.Global[id] = Global{
				Id:   id,
				Data: make(map[string]string),
			}
			global = db.Global[id]
		}
		if global.Data == nil {
			global.Data = make(map[string]string)
		}
		global.Data[key] = value
	default:
		LOG.ERROR("[ERROR]:Invalid category %s", category)
	}
}

func keepDatabase(db *Database) {
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
				LOG.INFO("Received signal, saving data and exiting...")
				err := SaveData(db)
				if err != nil {
					return
				}
				os.Exit(0)
			case <-saveTicker.C:
				// 定时保存数据
				LOG.INFO("Saving data automatically...")
				err := SaveData(db)
				if err != nil {
					return
				}
			}
		}
	}()
	select {} // 阻塞主goroutine
}

func Start() *Database {
	// 启动并检查程序
	LOG.INFO("Starting database ...")
	db := newDatabase("datamap")
	err := loadData(&db)
	if err != nil {
		return nil
	}
	LOG.INFO("Database started successfully.")
	keepDatabase(&db)
	return &db
}
