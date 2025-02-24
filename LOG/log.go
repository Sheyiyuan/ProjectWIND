package LOG

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"io"
	"log"
	"runtime"
)

func Trace(text string, msg ...interface{}) {
	pc, file, line, ok := runtime.Caller(3)
	if !ok {
		pc, file, line, ok = runtime.Caller(2)
	}
	if ok {
		funcName := runtime.FuncForPC(pc).Name()
		log.Printf("[Trace]  [%s:%d %s()] %s\n", file, line, funcName, fmt.Sprintf(text, msg...))
	} else {
		log.Printf("[Trace]  %s\n", fmt.Sprintf(text, msg...))
	}
}

func Debug(text string, msg ...interface{}) {
	log.Printf("[Debug]  %s\n", fmt.Sprintf(text, msg...))
}

func Info(text string, msg ...interface{}) {
	msgText := fmt.Sprintf(text, msg...)
	log.Println("[Info]  ", msgText)
}

func Notice(text string, msg ...interface{}) {
	msgText := fmt.Sprintf(text, msg...)
	log.Println("[Notice]", msgText)
}

func Warn(text string, msg ...interface{}) {
	msgText := fmt.Sprintf(text, msg...)
	log.Println("[Warn]  ", msgText)
}

func Error(text string, msg ...interface{}) {
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
