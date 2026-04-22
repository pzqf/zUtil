package zUtils

import (
	"errors"
	"os"
	"runtime/debug"
)

func GetCurrentDirectory() (string, error) {
	return os.Getwd()
}

// Recover 捕获 panic 并返回错误信息
// 注意：此函数设计为在 defer 中直接调用，如: defer zUtils.RecoverToErr(&err)
// 直接调用 Recover() 不会捕获 panic，因为 recover() 只在 defer 直接调用链中生效
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

// RecoverToErr 捕获 panic 并将错误写入 errPtr
// 用法: var err error; defer zUtils.RecoverToErr(&err)
func RecoverToErr(errPtr *error) {
	if err := recover(); err != nil {
		var errMsg string
		if e, ok := err.(error); ok {
			errMsg = e.Error()
		} else {
			errMsg = "unknown error"
		}
		*errPtr = errors.New(errMsg + "\n" + string(debug.Stack()))
	}
}
