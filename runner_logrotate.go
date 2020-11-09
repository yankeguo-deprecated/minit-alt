package main

import (
	"context"
	"fmt"
	"github.com/acicn/minit/pkg/mlog"
	"github.com/robfig/cron/v3"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

// filename mark
// daily: FILENAME.ROT2020-06-02.EXT
// filesize: FILENAME.ROT000000000001.EXT (%012d)

const (
	RotationModeDaily    = "daily"
	RotationModeFilesize = "filesize"

	RotationCron = "@every 1m"

	RotationDailyDateLayout = "2006-01-02"
	RotationFilesize        = 256 * 1024 * 1024

	Rot = "ROT"
)

var (
	RotationMarkDailyPattern    = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	RotationMarkFilesizePattern = regexp.MustCompile(`^\d+$`)
)

func rotationMarkExtract(filename string) (original string, mark string) {
	dir, base := filepath.Dir(filename), filepath.Base(filename)
	bs := strings.Split(base, ".")
	if len(bs) < 2 {
		original = filename
		return
	}
	if !strings.HasPrefix(bs[len(bs)-2], Rot) {
		original = filename
		return
	}
	mark = bs[len(bs)-2][len(Rot):]
	original = filepath.Join(dir, strings.Join(append(bs[:len(bs)-2], bs[len(bs)-1]), "."))
	return
}

func rotationMarkAdd(filename string, mark string) string {
	dir, base := filepath.Dir(filename), filepath.Base(filename)
	bs := strings.Split(base, ".")
	if len(bs) < 2 {
		return filepath.Join(dir, base+"."+Rot+mark)
	}
	return filepath.Join(dir, strings.Join(append(bs[:len(bs)-1], Rot+mark, bs[len(bs)-1]), "."))
}

type rotationFile struct {
	original string
	marks    map[string]bool
}

type LogrotateRunner struct {
	Unit
	logger *mlog.Logger
}

func (l *LogrotateRunner) Run(ctx context.Context) {
	l.logger.Printf("控制器启动")
	defer l.logger.Printf("控制器退出")

	cr := cron.New(cron.WithLogger(cron.PrintfLogger(l.logger)))
	_, err := cr.AddFunc(RotationCron, func() {
		l.logger.Printf("开始日志轮转")
		defer l.logger.Printf("结束日志轮转")
		l.rotate()
	})
	if err != nil {
		panic(err)
	}

	cr.Start()
	<-ctx.Done()
	<-cr.Stop().Done()
}

func (l *LogrotateRunner) collectRotationFiles() []*rotationFile {
	rfs := map[string]*rotationFile{}

	for _, fPat := range l.Files {
		matches, _ := filepath.Glob(fPat)
		for _, match := range matches {
			filename, _ := filepath.Abs(match)
			if filename != "" {
				orig, mark := rotationMarkExtract(filename)
				rf := rfs[orig]
				if rf == nil {
					rf = &rotationFile{original: orig, marks: map[string]bool{}}
					rfs[filename] = rf
				}
				if mark != "" {
					rf.marks[mark] = true
				}
			}
		}
	}

	ret := make([]*rotationFile, 0, len(rfs))
	for _, rf := range rfs {
		ret = append(ret, rf)
	}
	return ret
}

func (l *LogrotateRunner) rotate() {
	now := time.Now()
	bod := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	boy := bod.Add(-time.Hour * 24)
	moy := boy.Format(RotationDailyDateLayout)

	// 遍历所有通配符，建立文件组
	rfs := l.collectRotationFiles()

	// 遍历所有 rotationFile
	for _, rf := range rfs {
		// 删除不符合规则的 ROT 文件
		for mark, ok := range rf.marks {
			if !ok {
				continue
			}
			switch l.Mode {
			case RotationModeDaily:
				if RotationMarkDailyPattern.MatchString(mark) {
					continue
				}
			case RotationModeFilesize:
				if RotationMarkFilesizePattern.MatchString(mark) {
					continue
				}
			default:
				continue
			}
			rf.marks[mark] = false
			_ = os.Remove(rotationMarkAdd(rf.original, mark))
		}

		// 排序
		marks := make([]string, 0, len(rf.marks))
		for mark, ok := range rf.marks {
			if !ok {
				continue
			}
			marks = append(marks, mark)
		}
		sort.Strings(marks)

		// 进行数量限制
		if l.Keep > 0 && len(marks) > l.Keep {
			for _, mark := range marks[0 : len(marks)-l.Keep] {
				_ = os.Remove(rotationMarkAdd(rf.original, mark))
			}
			marks = marks[len(marks)-l.Keep:]
		}

		// 进行轮转
		switch l.Mode {
		case RotationModeDaily:
			foy := rotationMarkAdd(rf.original, moy)
			if _, err := os.Stat(foy); err == nil {
				l.logger.Printf("昨日文件已经存在: %s", rf.original)
				continue
			} else if !os.IsNotExist(err) {
				l.logger.Printf("未知错误: %s: %s", rf.original, err.Error())
				continue
			}
			_ = os.Rename(rf.original, rotationMarkAdd(rf.original, moy))
		case RotationModeFilesize:
			if fi, err := os.Stat(rf.original); err != nil {
				l.logger.Printf("无法检测文件: %s: %s", rf.original, err.Error())
				continue
			} else {
				if fi.Size() < RotationFilesize {
					continue
				}
				var id int64
				if len(marks) > 0 {
					var err error
					if id, err = strconv.ParseInt(marks[len(marks)-1], 10, 64); err != nil {
						l.logger.Printf("无法解析最大编号: %s: %s", rf.original, err.Error())
						continue
					}
				}
				id = id + 1
				_ = os.Rename(rf.original, rotationMarkAdd(rf.original, fmt.Sprintf("%012d", id)))
			}
		}
	}

	if len(l.Command) > 0 {
		_ = execute(l.ExecuteOptions, l.logger)
	}
}

func NewLogrotateRunner(unit Unit, logger *mlog.Logger) (Runner, error) {
	switch unit.Mode {
	case RotationModeDaily:
	case RotationModeFilesize:
	default:
		return nil, fmt.Errorf("未知的 logrotate 模式: %s", unit.Mode)
	}
	return &LogrotateRunner{
		Unit:   unit,
		logger: logger,
	}, nil
}
