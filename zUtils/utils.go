package zUtils

import (
	"errors"
	"os"
	"runtime/debug"
)

func GetCurrentDirectory() (string, error) {
	return os.Getwd()
}

func Recover() error {
	if err := recover(); err != nil {
		var errMsg string
		if e, ok := err.(error); ok {
			errMsg = e.Error()
		} else {
			errMsg = "unknown error"
		}
		return errors.New(errMsg + "\n" + string(debug.Stack()))
	}
	return nil
}
