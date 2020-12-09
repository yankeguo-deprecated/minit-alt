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

type LoggerOptions struct {
	Dir      string
	Name     string
	Filename string
}

type Logger struct {
	namePrefix []byte

	outFile *LogFile
	errFile *LogFile

	out io.Writer
	err io.Writer
}

func NewLogger(opts LoggerOptions) (logger *Logger, err error) {
	logger = &Logger{
		namePrefix: []byte(" [" + opts.Name + "] "),
	}
	if logger.outFile, err = NewLogFile(opts.Dir, opts.Filename+".out", 64*1024*1024, 5); err != nil {
		return
	}
	if logger.errFile, err = NewLogFile(opts.Dir, opts.Filename+".err", 64*1024*1024, 5); err != nil {
		return
	}
	logger.out = io.MultiWriter(os.Stdout, logger.outFile)
	logger.err = io.MultiWriter(os.Stderr, logger.errFile)
	return
}

func (l *Logger) Close() error {
	if l.outFile != nil {
		_ = l.outFile.Close()
	}
	if l.errFile != nil {
		_ = l.errFile.Close()
	}
	return nil
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
