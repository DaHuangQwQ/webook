package domain

import (
	"github.com/robfig/cron/v3"
	"time"
)

type Job struct {
	Id       int64
	Name     string
	Executor string
	// 通用任务的抽象
	Cfg        string
	Cron       string
	CancelFunc func() error
}

var parser = cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)

func (job *Job) NextTime() time.Time {
	parse, _ := parser.Parse(job.Cron)
	return parse.Next(time.Now())
}
