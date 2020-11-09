package main

import (
	"fmt"
	"github.com/acicn/minit/pkg/mlog"
	"github.com/acicn/minit/pkg/shellquote"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
)

var (
	childPids                 = map[int]bool{}
	childPidsLock sync.Locker = &sync.Mutex{}
)

type ExecuteOptions struct {
	Dir     string   `yaml:"dir"`     // 所有涉及命令执行的单元，指定命令执行时的当前目录
	Shell   string   `yaml:"shell"`   // 使用 shell 来执行命令，比如 'bash'
	Command []string `yaml:"command"` // 所有涉及命令执行的单元，指定命令执行的内容
}

func addPid(pid int) {
	childPidsLock.Lock()
	defer childPidsLock.Unlock()
	childPids[pid] = true
}

func removePid(pid int) {
	childPidsLock.Lock()
	defer childPidsLock.Unlock()
	delete(childPids, pid)
}

func notifyPIDs(sig os.Signal) {
	childPidsLock.Lock()
	defer childPidsLock.Unlock()
	for pid, found := range childPids {
		if found {
			if process, _ := os.FindProcess(pid); process != nil {
				_ = process.Signal(sig)
			}
		}
	}
}

func execute(opts ExecuteOptions, logger *mlog.Logger) (err error) {
	argv := make([]string, 0)

	// 构建 argv
	if opts.Shell != "" {
		if argv, err = shellquote.Split(opts.Shell); err != nil {
			err = fmt.Errorf("无法处理 shell 参数，请检查: %s", err.Error())
			return
		}
	} else {
		for _, arg := range opts.Command {
			argv = append(argv, os.ExpandEnv(arg))
		}
	}

	// 构建 cmd
	var outPipe, errPipe io.ReadCloser
	cmd := exec.Command(argv[0], argv[1:]...)
	if opts.Shell != "" {
		cmd.Stdin = strings.NewReader(strings.Join(opts.Command, "\n"))
	}
	cmd.Dir = opts.Dir
	// 阻止信号传递
	setupCmdSysProcAttr(cmd)

	if outPipe, err = cmd.StdoutPipe(); err != nil {
		return
	}
	if errPipe, err = cmd.StderrPipe(); err != nil {
		return
	}

	// 执行
	if err = cmd.Start(); err != nil {
		return
	}

	// 记录 Pid
	addPid(cmd.Process.Pid)

	// 串流
	go logger.StreamOut(outPipe)
	go logger.StreamErr(errPipe)

	// 等待退出
	if err = cmd.Wait(); err != nil {
		logger.Errorf("进程退出: %s", err.Error())
		err = nil
	} else {
		logger.Printf("进程退出")
	}

	// 移除 Pid
	removePid(cmd.Process.Pid)

	return
}
