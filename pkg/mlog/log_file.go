package mlog

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

type LogFile struct {
	name string
	dir  string

	currentFile *os.File
	currentSize int64

	maxSize  int64
	maxCount int64

	l sync.Locker
}

func (l *LogFile) currentFileName() string {
	return filepath.Join(l.dir, l.name+".log")
}

func (l *LogFile) archiveFileName(id int64) string {
	return filepath.Join(l.dir, fmt.Sprintf("%s.%d.log", l.name, id))
}

func (l *LogFile) nextArchiveId() (id int64, err error) {
	var fis []os.FileInfo
	if fis, err = ioutil.ReadDir(l.dir); err != nil {
		return
	}

	for _, fi := range fis {
		if strings.HasPrefix(fi.Name(), l.name+".") &&
			strings.HasSuffix(fi.Name(), ".log") {
			idStr := strings.TrimSuffix(strings.TrimPrefix(fi.Name(), l.name+"."), ".log")
			nid, _ := strconv.ParseInt(idStr, 10, 64)
			if nid > id {
				id = nid
			}
		}
	}

	id += 1

	if l.maxCount > 0 && id > l.maxCount {
		id = 1
	}
	return
}

func (l *LogFile) open() (err error) {
	var file *os.File
	if file, err = os.OpenFile(
		l.currentFileName(),
		os.O_WRONLY|os.O_CREATE|os.O_APPEND,
		0644,
	); err != nil {
		return
	}

	var info os.FileInfo
	if info, err = file.Stat(); err != nil {
		file.Close()
		return
	}

	l.currentFile = file
	l.currentSize = info.Size()

	return
}

func (l *LogFile) reallocate() (err error) {
	l.l.Lock()
	defer l.l.Unlock()

	var info os.FileInfo

	if info, err = l.currentFile.Stat(); err != nil {
		return
	}

	if info.Size() <= l.maxSize {
		return
	}

	var id int64
	if id, err = l.nextArchiveId(); err != nil {
		return
	}

	// try remove existed, in case id looped due to maxCount
	_ = os.Remove(l.archiveFileName(id))

	if err = os.Rename(l.currentFileName(), l.archiveFileName(id)); err != nil {
		return
	}

	if err = l.open(); err != nil {
		return
	}

	return nil
}

func (l *LogFile) Write(p []byte) (n int, err error) {
	if n, err = l.currentFile.Write(p); err != nil {
		return
	}

	atomic.AddInt64(&l.currentSize, int64(n))

	if l.currentSize > l.maxSize {
		if err = l.reallocate(); err != nil {
			return
		}
	}

	return
}

func (l *LogFile) Close() error {
	return l.currentFile.Close()
}

func NewLogFile(dir, name string, maxSize int64, maxCount int64) (lf *LogFile, err error) {
	lf = &LogFile{
		dir:      dir,
		name:     name,
		maxSize:  maxSize,
		maxCount: maxCount,
		l:        &sync.Mutex{},
	}
	if err = lf.open(); err != nil {
		return
	}
	return
}
