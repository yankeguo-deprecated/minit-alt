package main

import (
	"context"
	"fmt"
	"github.com/acicn/minit/pkg/mlog"
	"github.com/robfig/cron/v3"
)

const KindCron = "cron"

type CronRunner struct {
	Unit
	logger *mlog.Logger
}

func (r *CronRunner) Run(ctx context.Context) {
	r.logger.Printf("控制器启动")
	defer r.logger.Printf("控制器退出")

	cr := cron.New(cron.WithLogger(cron.PrintfLogger(r.logger)))
	_, err := cr.AddFunc(r.Cron, func() {
		r.logger.Printf("定时任务触发")
		_ = execute(r.ExecuteOptions, r.logger)
		r.logger.Printf("定时任务结束")
	})
	if err != nil {
		// 已经检查过表达式了，不应该报错
		panic(err)
	}

	cr.Start()

	<-ctx.Done()
	<-cr.Stop().Done()
}

func NewCronRunner(unit Unit, logger *mlog.Logger) (Runner, error) {
	if len(unit.Command) == 0 {
		return nil, fmt.Errorf("没有指定命令，检查 command 字段")
	}
	if len(unit.Cron) == 0 {
		return nil, fmt.Errorf("没有指定 cron 表达式，检查 cron 字段")
	}
	if _, err := cron.ParseStandard(unit.Cron); err != nil {
		return nil, fmt.Errorf("cron 表达式语法错误，检查 cron 字段: %s", err.Error())
	}
	return &CronRunner{
		Unit:   unit,
		logger: logger,
	}, nil
}
