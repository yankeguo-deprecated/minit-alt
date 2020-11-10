package main

import (
	"bytes"
	"io/ioutil"
)

const (
	BannerFile = "/etc/banner.minit.txt"
)

func setupBanner() {
	var err error
	var buf []byte
	if buf, err = ioutil.ReadFile(BannerFile); err != nil {
		return
	}
	lines := bytes.Split(buf, []byte{'\n'})
	for _, line := range lines {
		log.Print(string(line))
	}
	return
}
