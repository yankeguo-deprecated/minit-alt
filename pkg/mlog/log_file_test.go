package mlog

import (
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

func TestNewLogFile(t *testing.T) {
	_ = os.RemoveAll(filepath.Join("testdata", "logs"))
	_ = os.MkdirAll(filepath.Join("testdata", "logs"), 0755)
	f, err := NewLogFile(filepath.Join("testdata", "logs"), "test", 10, 0)
	require.NoError(t, err)
	_, err = f.Write([]byte("hello, world, hello, world, hello, world"))
	require.NoError(t, err)
	_, err = f.Write([]byte("hello, world, hello, world, hello, world"))
	require.NoError(t, err)
	_, err = f.Write([]byte("hello, world, hello, world, hello, world"))
	require.NoError(t, err)
	err = f.Close()
	require.NoError(t, err)
	f, err = NewLogFile(filepath.Join("testdata", "logs"), "test-maxcount", 10, 2)
	require.NoError(t, err)
	_, err = f.Write([]byte("hello, world, hello, world, hello, world"))
	require.NoError(t, err)
	_, err = f.Write([]byte("hello, world, hello, world, hello, world"))
	require.NoError(t, err)
	_, err = f.Write([]byte("hello, world, hello, world, hello, world"))
	require.NoError(t, err)
	err = f.Close()
	require.NoError(t, err)
}
