package mlog

import (
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

func TestLog(t *testing.T) {
	os.MkdirAll(filepath.Join("testdata", "logger"), 0755)
	log, err := NewLogger(LoggerOptions{
		Dir:      filepath.Join("testdata", "logger"),
		Name:     "test",
		Filename: "test",
	})
	require.NoError(t, err)
	log.Print("hello", "world")
	log.Printf("hello, %s", "world")
	log.Error("error", "world")
	log.Errorf("error, %s", "world")
}
