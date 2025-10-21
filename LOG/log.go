package LOG

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"io"
	"log"
	"runtime"
)

var logLevel = 2

func SetLogLevel(level int) {
	if level < 0 {
		level = 0
	}
	if level > 5 {
		level = 5
	}
	logLevel = level
}

func Trace(text string, msg ...interface{}) {
	if logLevel > 0 {
		return
	}
	var pc uintptr
	var file string
	var line int
	var ok bool
	// 从第 2 层开始循环查找调用栈，直到找不到调用信息为止
	for i := 2; ; i++ {
		pc, file, line, ok = runtime.Caller(i)
		if !ok || line <= 0 {
			pc, file, line, ok = runtime.Caller(i - 1)
			break
		}
	}
	funcName := runtime.FuncForPC(pc).Name()
	log.Printf("[Trace]  [%s:%d %s()] %s\n", file, line, funcName, fmt.Sprintf(text, msg...))
}

func Debug(text string, msg ...interface{}) {
	if logLevel > 1 {
		return
	}
	log.Printf("[Debug]  %s\n", fmt.Sprintf(text, msg...))
}

func Info(text string, msg ...interface{}) {
	if logLevel > 2 {
		return
	}
	msgText := fmt.Sprintf(text, msg...)
	log.Println("[Info]  ", msgText)
}

func Notice(text string, msg ...interface{}) {
	if logLevel > 3 {
		return
	}
	msgText := fmt.Sprintf(text, msg...)
	log.Println("[Notice]", msgText)
}

func Warn(text string, msg ...interface{}) {
	if logLevel > 4 {
		return
	}
	msgText := fmt.Sprintf(text, msg...)
	log.Println("[Warn]  ", msgText)
}

func Error(text string, msg ...interface{}) {
	if logLevel > 5 {
		return
	}
	msgText := fmt.Sprintf(text, msg...)
	log.Println("[Error] ", msgText)
}

func Fatal(text string, msg ...interface{}) {
	msgText := fmt.Sprintf(text, msg...)
	log.Fatalln("[Fatal] ", msgText)
}

// CustomLogger 是一个实现了 hlog.Logger 接口的自定义日志记录器
type CustomLogger struct{}

func (c *CustomLogger) Trace(v ...interface{}) {
	Trace(fmt.Sprint(v...))
}

func (c *CustomLogger) Debug(v ...interface{}) {
	Debug(fmt.Sprint(v...))
}

func (c *CustomLogger) Info(v ...interface{}) {
	Info(fmt.Sprint(v...))
}

func (c *CustomLogger) Notice(v ...interface{}) {
	Info(fmt.Sprint(v...))
}

func (c *CustomLogger) Warn(v ...interface{}) {
	Warn(fmt.Sprint(v...))
}

func (c *CustomLogger) Error(v ...interface{}) {
	Error(fmt.Sprint(v...))
}

func (c *CustomLogger) Fatal(v ...interface{}) {
	Fatal(fmt.Sprint(v...))
}

func (c *CustomLogger) CtxTracef(ctx context.Context, format string, v ...interface{}) {

}

func (c *CustomLogger) CtxDebugf(ctx context.Context, format string, v ...interface{}) {

}

func (c *CustomLogger) CtxInfof(ctx context.Context, format string, v ...interface{}) {

}

func (c *CustomLogger) CtxNoticef(ctx context.Context, format string, v ...interface{}) {

}

func (c *CustomLogger) CtxWarnf(ctx context.Context, format string, v ...interface{}) {

}

func (c *CustomLogger) CtxErrorf(ctx context.Context, format string, v ...interface{}) {

}

func (c *CustomLogger) CtxFatalf(ctx context.Context, format string, v ...interface{}) {

}

func (c *CustomLogger) SetLevel(level hlog.Level) {

}

func (c *CustomLogger) SetOutput(writer io.Writer) {

}

// Tracef 实现 hlog.Logger 接口的 Tracef 方法
func (c *CustomLogger) Tracef(format string, args ...interface{}) {
	Trace(format, args...)
}

// Debugf 实现 hlog.Logger 接口的 Debugf 方法
func (c *CustomLogger) Debugf(format string, args ...interface{}) {
	Debug(format, args...)
}

// Infof 实现 hlog.Logger 接口的 Infof 方法
func (c *CustomLogger) Infof(format string, args ...interface{}) {
	Info(format, args...)
}

// Warnf 实现 hlog.Logger 接口的 Warnf 方法
func (c *CustomLogger) Warnf(format string, args ...interface{}) {
	Warn(format, args...)
}

// Errorf 实现 hlog.Logger 接口的 Errorf 方法
func (c *CustomLogger) Errorf(format string, args ...interface{}) {
	Error(format, args...)
}

// Fatalf 实现 hlog.Logger 接口的 Fatalf 方法
func (c *CustomLogger) Fatalf(format string, args ...interface{}) {
	Fatal(format, args...)
}

func (c *CustomLogger) Noticef(format string, args ...interface{}) {
	Info(format, args...)
}
