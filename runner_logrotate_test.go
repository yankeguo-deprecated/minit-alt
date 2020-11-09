package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func TestRotationMarkRemove(t *testing.T) {
	o, m := rotationMarkExtract(filepath.Join("test", "hello.ROT2020-02-32.log"))
	assert.Equal(t, filepath.Join("test", "hello.log"), o)
	assert.Equal(t, "2020-02-32", m)
	o, m = rotationMarkExtract("hello.ROT2020022.log")
	assert.Equal(t, "hello.log", o)
	assert.Equal(t, "2020022", m)
	o, m = rotationMarkExtract("hello.ROT2020022.log")
	assert.Equal(t, "hello.log", o)
	assert.Equal(t, "2020022", m)
	o, m = rotationMarkExtract(filepath.Join("hello", "helloOT2020022.log"))
	assert.Equal(t, filepath.Join("hello", "helloOT2020022.log"), o)
	assert.Equal(t, "", m)
}

func TestRotationMarkAdd(t *testing.T) {
	assert.Equal(t, filepath.Join("test", "hello.ROT*.log"), rotationMarkAdd(filepath.Join("test", "hello.log"), "*"))
	assert.Equal(t, filepath.Join("test", "hello.ROT*"), rotationMarkAdd(filepath.Join("test", "hello"), "*"))
	assert.Equal(t, ".ROT*.hello", rotationMarkAdd(".hello", "*"))
	assert.Equal(t, "hello.ROT000000000011.log", rotationMarkAdd("hello.log", fmt.Sprintf("%012d", 11)))
}
