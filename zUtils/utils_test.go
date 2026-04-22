package zUtils

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCurrentDirectory(t *testing.T) {
	dir, err := GetCurrentDirectory()
	assert.NoError(t, err)
	assert.NotEmpty(t, dir)
}

func TestRecover_NoPanic(t *testing.T) {
	var err error
	func() {
		defer RecoverToErr(&err)
	}()
	assert.Nil(t, err)
}

func TestRecoverToErr_WithPanic(t *testing.T) {
	var err error
	func() {
		defer RecoverToErr(&err)
		panic(errors.New("test panic"))
	}()
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "test panic")
}

func TestRecoverToErr_WithStringPanic(t *testing.T) {
	var err error
	func() {
		defer RecoverToErr(&err)
		panic("string panic")
	}()
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "unknown error")
}
