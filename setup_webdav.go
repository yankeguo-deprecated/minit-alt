package main

import (
	"fmt"
	"golang.org/x/net/webdav"
	"net/http"
	"os"
	"strings"
	"time"
)

func SetupWebDAV() (err error) {
	envRoot := strings.TrimSpace(os.Getenv("MINIT_WEBDAV_ROOT"))
	if envRoot == "" {
		return
	}
	if err = os.MkdirAll(envRoot, 0755); err != nil {
		err = fmt.Errorf("无法初始化 WebDAV 根目录: %s: %s", envRoot, err.Error())
		return
	}
	envPort := strings.TrimSpace(os.Getenv("MINIT_WEBDAV_PORT"))
	if envPort == "" {
		envPort = "7486"
	}
	log.Printf("启动 WebDAV 服务: 路径 %s 端口 %s", envRoot, envPort)
	h := &webdav.Handler{
		FileSystem: webdav.Dir(envRoot),
		LockSystem: webdav.NewMemLS(),
		Logger: func(req *http.Request, err error) {
			if err != nil {
				log.Printf("WebDAV: %s %s: %s", req.Method, req.URL.Path, err.Error())
			} else {
				log.Printf("WebDAV: %s %s", req.Method, req.URL.Path)
			}
		},
	}
	envUsername := strings.TrimSpace(os.Getenv("MINIT_WEBDAV_USERNAME"))
	envPassword := strings.TrimSpace(os.Getenv("MINIT_WEBDAV_PASSWORD"))
	s := http.Server{
		Addr: ":" + envPort,
		Handler: http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			if envUsername != "" && envPassword != "" {
				if username, password, ok := req.BasicAuth(); !ok || username != envUsername || password != envPassword {
					rw.Header().Add("WWW-Authenticate", `Basic realm=Minit WebDAV`)
					rw.WriteHeader(http.StatusUnauthorized)
					return
				}
			}
			h.ServeHTTP(rw, req)
		}),
	}
	go func() {
		for {
			if err := s.ListenAndServe(); err != nil {
				log.Printf("无法启动 WebDAV 服务器: %s", err.Error())
			}
			time.Sleep(time.Second * 10)
		}
	}()
	return
}
