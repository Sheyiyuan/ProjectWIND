package LOG

import (
	"fmt"
	"log"
)

func DEBUG(text string, msg ...interface{}) {
	msgText := fmt.Sprintf(text, msg...)
	log.Println("[DEBUG] ", msgText)
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
