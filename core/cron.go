package core

import (
	"ProjectWIND/LOG"
	"ProjectWIND/wba"
	"github.com/robfig/cron/v3"
)

func RegisterCron(task wba.ScheduledTaskInfo) {
	// 注册定时任务
	c := cron.New(cron.WithSeconds())
	_, err := c.AddFunc(task.Cron, task.Task)
	if err != nil {
		LOG.ERROR("添加定时任务 %s 时出错%v:", task.Name, err)
	}
	c.Start()
	LOG.INFO("定时任务 %s 注册成功", task.Name)
}
