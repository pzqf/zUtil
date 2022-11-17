package zUtils

import (
	"bytes"
	"encoding/gob"
	"errors"
	"os"
	"runtime/debug"
)

func DeepCopy(dst, src interface{}) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)
}

func GetCurrentDirectory() string {
	if dir, err := os.Getwd(); err == nil {
		return dir
	}
	return ""
}

func IsExistPath(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

func IsDir(dir string) bool {
	f, e := os.Stat(dir)
	if e != nil {
		return false
	}
	return f.IsDir()
}

func Recover() error {
	if err := recover(); err != nil {
		return errors.New(err.(error).Error() + "/n" + string(debug.Stack()))
	}
	return nil
}
