//go:build linux
// +build linux

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func setupSysctl() (err error) {
	cfgs := strings.Split(os.Getenv("MINIT_SYSCTL"), ",")
	for _, cfg := range cfgs {
		splits := strings.Split(cfg, "=")
		if len(splits) != 2 {
			continue
		}
		k, v := strings.TrimSpace(splits[0]), strings.TrimSpace(splits[1])
		if k == "" {
			continue
		}
		ks := []string{"/proc", "sys"}
		ks = append(ks, strings.Split(k, ".")...)
		filename := filepath.Join(ks...)
		log.Printf("写入内核参数 %s=%s", k, v)
		if err = ioutil.WriteFile(filename, []byte(v), 0644); err != nil {
			err = fmt.Errorf("无法写入内核参数 %s=%s: %s", k, v, err.Error())
			return
		}
	}
	return
}
