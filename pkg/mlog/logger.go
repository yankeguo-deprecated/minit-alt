package mlog

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

const (
	LoggerDateLayout = "15:04:05.000"
)

var (
	loggerBuffers = &sync.Pool{
		New: func() interface{} {
			return &bytes.Buffer{}
		},
	}

	loggerNow = time.Now
)

type Logger struct {
	namePrefix []byte

	out io.Writer
	err io.Writer
}

func NewLogger(dir, name, filename string) (logger *Logger, err error) {
	logger = &Logger{
		namePrefix: []byte(" [" + name + "] "),
	}
	var outFile, errFile *LogFile
	if outFile, err = NewLogFile(dir, filename+".out", 64*1024*1024, 5); err != nil {
		return
	}
	if errFile, err = NewLogFile(dir, filename+".err", 64*1024*1024, 5); err != nil {
		return
	}
	logger.out = io.MultiWriter(os.Stdout, outFile)
	logger.err = io.MultiWriter(os.Stderr, errFile)
	return
}

func (l *Logger) Print(items ...interface{}) {
	appendLogLine(l.namePrefix, append([]byte(fmt.Sprint(items...)), '\n'), l.out)
}

func (l *Logger) Error(items ...interface{}) {
	appendLogLine(l.namePrefix, append([]byte(fmt.Sprint(items...)), '\n'), l.err)
}

func (l *Logger) Printf(pattern string, items ...interface{}) {
	appendLogLine(l.namePrefix, append([]byte(fmt.Sprintf(pattern, items...)), '\n'), l.out)
}

func (l *Logger) Errorf(pattern string, items ...interface{}) {
	appendLogLine(l.namePrefix, append([]byte(fmt.Sprintf(pattern, items...)), '\n'), l.err)
}

func (l *Logger) StreamOut(r io.Reader) {
	streamLogLine(l.namePrefix, r, l.out)
}

func (l *Logger) StreamErr(r io.Reader) {
	streamLogLine(l.namePrefix, r, l.err)
}

func appendLogLine(name, b []byte, w io.Writer) {
	buf := loggerBuffers.Get().(*bytes.Buffer)
	buf.WriteString(loggerNow().Format(LoggerDateLayout))
	buf.Write(name)
	buf.Write(b)
	_, _ = w.Write(buf.Bytes())
	buf.Reset()
	loggerBuffers.Put(buf)
}
func streamLogLine(name []byte, r io.Reader, w io.Writer) {
	br := bufio.NewReader(r)
	for {
		b, err := br.ReadBytes('\n')
		if err == nil {
			appendLogLine(name, b, w)
		} else {
			if len(b) != 0 {
				appendLogLine(name, append(b, '\n'), w)
			}
			break
		}
	}
}
