package LOG

import (
	"fmt"
	"log"
	"runtime"
)

func DEBUG(text string, msg ...interface{}) {
	pc, file, line, ok := runtime.Caller(2)
	if ok {
		funcName := runtime.FuncForPC(pc).Name()
		log.Printf("[DEBUG]  [%s:%d %s()] %s\n", file, line, funcName, fmt.Sprintf(text, msg...))
	} else {
		log.Printf("[DEBUG]  %s\n", fmt.Sprintf(text, msg...))
	}
}

func INFO(text string, msg ...interface{}) {
	msgText := fmt.Sprintf(text, msg...)
	log.Println("[INFO]  ", msgText)
}

func WARN(text string, msg ...interface{}) {
	msgText := fmt.Sprintf(text, msg...)
	log.Println("[WARN]  ", msgText)
}

func ERROR(text string, msg ...interface{}) {
	msgText := fmt.Sprintf(text, msg...)
	log.Println("[ERROR] ", msgText)
}

func FATAL(text string, msg ...interface{}) {
	msgText := fmt.Sprintf(text, msg...)
	log.Fatalln("[FATAL] ", msgText)
}
