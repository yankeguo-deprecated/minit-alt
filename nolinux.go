//+build !linux

package main

import "os/exec"

func setupCmdSysProcAttr(*exec.Cmd) {
}

func setupTHP() error {
	return nil
}

func setupSysctl() error {
	return nil
}

func setupRLimits() error {
	return nil
}
