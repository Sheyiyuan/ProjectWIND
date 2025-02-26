package core

import (
	"ProjectWIND/LOG"
	"ProjectWIND/wba"
	"github.com/robfig/cron/v3"
)

func RegisterCron(appNames string, task wba.ScheduledTaskInfo) {
	// 注册定时任务
	c := cron.New(cron.WithSeconds())
	_, err := c.AddFunc(task.Cron, task.Task)
	if err != nil {
		LOG.Error("添加定时任务 [%s]%s 时出错%v:", appNames, task.Name, err)
	}
	c.Start()
	LOG.Info("定时任务 [%s]%s 注册成功", appNames, task.Name)
}
