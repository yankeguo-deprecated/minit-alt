package main

import (
	"flag"
	"fmt"
	"github.com/acicn/minit/pkg/shellquote"
	"gopkg.in/yaml.v2"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type FilterMode int

const (
	FilterNone FilterMode = iota
	FilterWhitelist
	FilterBlacklist
)

const (
	DefaultGroup = "default"
)

var (
	filters    = map[string]bool{}
	filterMode FilterMode
)

func init() {
	filterMode = FilterNone

	var raw string
	if raw = os.Getenv("MINIT_ENABLE"); raw != "" {
		filterMode = FilterWhitelist
	} else if raw = os.Getenv("MINIT_DISABLE"); raw != "" {
		filterMode = FilterBlacklist
	}

	names := strings.Split(raw, ",")
	for _, name := range names {
		name = strings.TrimSpace(name)
		if name == "@" {
			continue
		}
		filters[name] = true
	}
}

type Unit struct {
	ExecuteOptions `yaml:",inline"`

	Name  string `yaml:"name"`  // 单元名
	Group string `yaml:"group"` // 单元分组
	Kind  string `yaml:"kind"`  // 单元类型
	Count int    `yaml:"count"` // 单元副本数量

	Raw bool `yaml:"raw"` // 不对渲染文件进行空白行处理

	Files []string `yaml:"files"` // render, logrotate, logcollect 单元，通配符指定要处理的文件

	Cron string `yaml:"cron"` // cron 单元, 定时表达式
	Mode string `yaml:"mode"` // logrotate 单元，模式 daily 或者 size
	Keep int    `yaml:"keep"` // logrotate 单元，保留天数/份数
}

func (u Unit) CanonicalName() string {
	return u.Kind + "/" + u.Name
}

func LoadArgsMain() (unit Unit, ok bool, err error) {
	args := flag.Args()
	for i, arg := range args {
		if arg == "--" {
			args = args[i+1:]
			break
		}
	}
	if len(args) == 0 {
		return
	}
	unit = Unit{
		Name:  "arg-main",
		Group: DefaultGroup,
		Kind:  KindDaemon,
		ExecuteOptions: ExecuteOptions{
			Command: args,
		},
	}
	ok = true
	return
}

func LoadEnvMain() (unit Unit, ok bool, err error) {
	cmd := strings.TrimSpace(os.Getenv("MINIT_MAIN"))
	if cmd == "" {
		return
	}
	name := strings.TrimSpace(os.Getenv("MINIT_MAIN_NAME"))
	if name == "" {
		name = "env-main"
	}
	group := strings.TrimSpace(os.Getenv("MINIT_MAIN_GROUP"))
	if group == "" {
		group = DefaultGroup
	}
	kind := KindDaemon
	if once, _ := strconv.ParseBool(strings.TrimSpace(os.Getenv("MINIT_MAIN_ONCE"))); once {
		kind = KindOnce
		return
	}
	var command []string
	if command, err = shellquote.Split(cmd); err != nil {
		return
	}
	unit = Unit{
		Name:  name,
		Group: group,
		Kind:  kind,
		ExecuteOptions: ExecuteOptions{
			Command: command,
			Dir:     strings.TrimSpace(os.Getenv("MINIT_MAIN_DIR")),
		},
	}
	ok = true
	return
}

func LoadDir(dir string) (units []Unit, err error) {
	var files []string
	for _, ext := range []string{"*.yml", "*.yaml"} {
		if files, err = filepath.Glob(filepath.Join(dir, ext)); err != nil {
			return
		}
		for _, file := range files {
			var units0 []Unit
			if units0, err = LoadFile(file); err != nil {
				return
			}
			units = append(units, units0...)
		}
	}
	return
}

func LoadFile(fn string) (units []Unit, err error) {
	var f *os.File
	if f, err = os.Open(fn); err != nil {
		return
	}
	defer f.Close()

	dec := yaml.NewDecoder(f)
	for {
		var unit Unit
		if err = dec.Decode(&unit); err != nil {
			if err == io.EOF {
				err = nil
			} else {
				err = fmt.Errorf("无法解析文件 %s: %s", fn, err.Error())
			}
			return
		}

		// 清理下空格
		unit.Name = strings.TrimSpace(unit.Name)
		unit.Kind = strings.TrimSpace(unit.Kind)
		unit.Cron = strings.TrimSpace(unit.Cron)
		unit.Dir = strings.TrimSpace(unit.Dir)
		unit.Group = strings.TrimSpace(unit.Group)

		// 默认组名
		if unit.Group == "" {
			unit.Group = DefaultGroup
		}

		// 打开关闭
		switch filterMode {
		case FilterNone:
		case FilterWhitelist:
			if !filters[unit.Name] && !filters["@"+unit.Group] {
				log.Printf("取消单元载入: %s", unit.Name)
				continue
			}
		case FilterBlacklist:
			if filters[unit.Name] || filters["@"+unit.Group] {
				log.Printf("取消单元载入: %s", unit.Name)
				continue
			}
		}

		// 重复型
		if unit.Count > 0 {
			for i := 0; i < unit.Count; i++ {
				subUnit := unit
				subUnit.Name = fmt.Sprintf("%s-%d", unit.Name, i+1)
				units = append(units, subUnit)
			}
		} else {
			units = append(units, unit)
		}
	}
}
