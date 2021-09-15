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
	namePrefix string

	fileOut *LogFile
	fileErr *LogFile

	consoleOut io.Writer
	consoleErr io.Writer
}

func NewLogger(opts LoggerOptions) (logger *Logger, err error) {
	logger = &Logger{
		namePrefix: " [" + opts.Name + "] ",
	}
	if logger.fileOut, err = NewLogFile(opts.Dir, opts.Filename+".out", 64*1024*1024, 5); err != nil {
		return
	}
	if logger.fileErr, err = NewLogFile(opts.Dir, opts.Filename+".err", 64*1024*1024, 5); err != nil {
		return
	}
	logger.consoleOut = os.Stdout
	logger.consoleErr = os.Stderr
	return
}

func (l *Logger) Close() error {
	if l.fileOut != nil {
		_ = l.fileOut.Close()
	}
	if l.fileErr != nil {
		_ = l.fileErr.Close()
	}
	return nil
}

func (l *Logger) Print(items ...interface{}) {
	l.AppendOut(append([]byte(fmt.Sprint(items...)), '\n'))
}

func (l *Logger) Error(items ...interface{}) {
	l.AppendErr(append([]byte(fmt.Sprint(items...)), '\n'))
}

func (l *Logger) Printf(pattern string, items ...interface{}) {
	l.AppendOut(append([]byte(fmt.Sprintf(pattern, items...)), '\n'))
}

func (l *Logger) Errorf(pattern string, items ...interface{}) {
	l.AppendErr(append([]byte(fmt.Sprintf(pattern, items...)), '\n'))
}

func (l *Logger) StreamOut(r io.Reader) {
	br := bufio.NewReader(r)
	for {
		b, err := br.ReadBytes('\n')
		if err == nil {
			l.AppendOut(b)
		} else {
			if len(b) != 0 {
				l.AppendOut(append(b, '\n'))
			}
			break
		}
	}
}

func (l *Logger) StreamErr(r io.Reader) {
	br := bufio.NewReader(r)
	for {
		b, err := br.ReadBytes('\n')
		if err == nil {
			l.AppendErr(b)
		} else {
			if len(b) != 0 {
				l.AppendErr(append(b, '\n'))
			}
			break
		}
	}
}

func (l *Logger) AppendOut(b []byte) {
	bc, bf := l.FormatLine(b)
	_, _ = l.consoleOut.Write(bc)
	_, _ = l.fileOut.Write(bf)
}

func (l *Logger) AppendErr(b []byte) {
	bc, bf := l.FormatLine(b)
	_, _ = l.consoleErr.Write(bc)
	_, _ = l.fileErr.Write(bf)
}

func (l *Logger) FormatLine(b []byte) (bufConsole []byte, bufFile []byte) {
	bufFile = b
	bufConsole = append([]byte(loggerNow().Format(LoggerDateLayout)+l.namePrefix), b...)
	return
}
