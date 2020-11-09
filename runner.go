package main

import (
	"context"
	"github.com/acicn/minit/pkg/mlog"
)

type RunnerLevel int

const (
	RunnerL1 RunnerLevel = iota + 1
	RunnerL2
	RunnerL3
)

type RunnerFactory struct {
	Level  RunnerLevel
	Create func(unit Unit, logger *mlog.Logger) (Runner, error)
}

var (
	RunnerFactories = map[string]*RunnerFactory{
		"render": {
			Level: RunnerL1,
			Create: func(unit Unit, logger *mlog.Logger) (Runner, error) {
				return NewRenderRunner(unit, logger)
			},
		},
		"once": {
			Level: RunnerL2,
			Create: func(unit Unit, logger *mlog.Logger) (Runner, error) {
				return NewOnceRunner(unit, logger)
			},
		},
		"daemon": {
			Level: RunnerL3,
			Create: func(unit Unit, logger *mlog.Logger) (Runner, error) {
				return NewDaemonRunner(unit, logger)
			},
		},
		"cron": {
			Level: RunnerL3,
			Create: func(unit Unit, logger *mlog.Logger) (Runner, error) {
				return NewCronRunner(unit, logger)
			},
		},
		"logrotate": {
			Level: RunnerL3,
			Create: func(unit Unit, logger *mlog.Logger) (Runner, error) {
				return NewLogrotateRunner(unit, logger)
			},
		},
	}
)

type Runner interface {
	Run(ctx context.Context)
}
