package zError

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	err := New("something went wrong")
	assert.NotNil(t, err)
	assert.Equal(t, 0, err.GetCode())
	assert.Equal(t, "something went wrong", err.GetMessage())
	assert.Contains(t, err.Error(), "something went wrong")
}

func TestNewWithCode(t *testing.T) {
	err := NewWithCode(1001, "player not found")
	assert.NotNil(t, err)
	assert.Equal(t, 1001, err.GetCode())
	assert.Equal(t, "player not found", err.GetMessage())
	assert.Contains(t, err.Error(), "1001")
	assert.Contains(t, err.Error(), "player not found")
}

func TestErrorf(t *testing.T) {
	err := Errorf("player %d not found in map %d", 1001, 2001)
	assert.NotNil(t, err)
	assert.Equal(t, 0, err.GetCode())
	assert.Equal(t, "player 1001 not found in map 2001", err.GetMessage())
}

func TestBaseError_Error(t *testing.T) {
	err := NewWithCode(500, "internal error")
	errStr := err.Error()
	assert.Contains(t, errStr, "500")
	assert.Contains(t, errStr, "internal error")
}

func TestErrorInterface(t *testing.T) {
	var err Error = NewWithCode(1001, "test error")
	assert.Equal(t, 1001, err.GetCode())
	assert.Equal(t, "test error", err.GetMessage())
}
