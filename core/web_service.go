package core

import (
	"ProjectWIND/LOG"
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func WebServer(port string) *server.Hertz {
	// 设置自定义日志记录器
	hlog.SetLevel(hlog.LevelFatal)
	h := server.Default(server.WithHostPorts("0.0.0.0:" + port))
	LOG.Info("WebUI已启动，监听端口：%v", port)
	h.Use(LoggingMiddleware())
	h.GET("/index", func(ctx context.Context, c *app.RequestContext) {
		//返回webui/index.html
		c.File("./webui/index.html")
	})
	h.GET("/app/api/:appName/*action", func(ctx context.Context, c *app.RequestContext) {
		appName := c.Param("appName")
		action := c.Param("action")
		message := appName + "-" + action + "\n"
		c.String(consts.StatusOK, message)
	})
	return h
}

// LoggingMiddleware 是一个中间件函数，用于记录请求信息
func LoggingMiddleware() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		// 获取客户端 IP
		clientIP := ctx.ClientIP()
		// 获取请求完整路径
		fullPath := ctx.Request.URI().PathOriginal()
		// 继续处理请求
		ctx.Next(c)
		// 获取请求处理状态码
		statusCode := ctx.Response.StatusCode()
		// 在请求处理后记录结束信息
		LOG.Debug("收到网络请求 | IP: %s | Path: %s | Status: %d ", clientIP, fullPath, statusCode)
	}
}
