package mlog

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestLog(t *testing.T) {
	log, err := NewLogger(os.TempDir(), "test", "test")
	require.NoError(t, err)
	log.Print("hello", "world")
	log.Printf("hello, %s", "world")
	log.Error("error", "world")
	log.Errorf("error, %s", "world")
}
