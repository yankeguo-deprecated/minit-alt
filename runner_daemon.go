package main

import (
	"context"
	"fmt"
	"github.com/acicn/minit/pkg/mlog"
	"time"
)

type DaemonRunner struct {
	Unit
	logger *mlog.Logger
}

func (r *DaemonRunner) Run(ctx context.Context) {
	r.logger.Printf("控制器启动")
	defer r.logger.Printf("控制器退出")
forLoop:
	for {
		// 检查 ctx 是否已经结束
		if ctx.Err() != nil {
			break forLoop
		}

		var err error
		if err = execute(r.ExecuteOptions, r.logger); err != nil {
			r.logger.Errorf("启动失败: %s", err.Error())
		}

		// 检查 ctx 是否已经结束
		if ctx.Err() != nil {
			break forLoop
		}

		// 重试
		r.logger.Printf("5s 后重启")

		timer := time.NewTimer(time.Second * 5)
		select {
		case <-timer.C:
		case <-ctx.Done():
			break forLoop
		}
	}
}

func NewDaemonRunner(unit Unit, logger *mlog.Logger) (Runner, error) {
	if len(unit.Command) == 0 {
		return nil, fmt.Errorf("没有指定命令，检查 command 字段")
	}
	return &DaemonRunner{
		Unit:   unit,
		logger: logger,
	}, nil
}
