//go:build linux
// +build linux

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const (
	controlFileTHP = "/sys/kernel/mm/transparent_hugepage/enabled"
)

func setupTHP() (err error) {
	val := strings.TrimSpace(os.Getenv("MINIT_THP"))
	if val == "" {
		return
	}
	var buf []byte
	if buf, err = ioutil.ReadFile(controlFileTHP); err != nil {
		err = fmt.Errorf("无法读取透明大页配置文件 %s: %s", controlFileTHP, err.Error())
		return
	}
	log.Printf("当前透明大页配置: %s", bytes.TrimSpace(buf))
	log.Printf("写入透明大页配置: %s", val)
	if err = ioutil.WriteFile(controlFileTHP, []byte(val), 644); err != nil {
		err = fmt.Errorf("无法写入透明大页配置文件 %s: %s", controlFileTHP, err.Error())
		return
	}
	if buf, err = ioutil.ReadFile(controlFileTHP); err != nil {
		err = fmt.Errorf("无法读取透明大页配置文件 %s: %s", controlFileTHP, err.Error())
		return
	}
	log.Printf("当前透明大页配置: %s", bytes.TrimSpace(buf))
	return
}
