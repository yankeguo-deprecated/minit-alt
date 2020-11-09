package main

import (
	"context"
	"fmt"
	"github.com/acicn/minit/pkg/mlog"
)

type OnceRunner struct {
	Unit
	logger *mlog.Logger
}

func (r *OnceRunner) Run(ctx context.Context) {
	r.logger.Printf("控制器启动")
	defer r.logger.Printf("控制器退出")
	if err := execute(r.ExecuteOptions, r.logger); err != nil {
		r.logger.Errorf("启动失败: %s", err.Error())
		return
	}
}

func NewOnceRunner(unit Unit, logger *mlog.Logger) (Runner, error) {
	if len(unit.Command) == 0 {
		return nil, fmt.Errorf("没有指定命令，检查 command 字段")
	}
	return &OnceRunner{
		Unit:   unit,
		logger: logger,
	}, nil
}
