package zFile

import (
	"bufio"
	"errors"
	"io"
	"os"
	"path/filepath"
)

// ReadFile 读取文件内容
func ReadFile(filePath string) ([]byte, error) {
	return os.ReadFile(filePath)
}

// WriteFile 写入文件内容
func WriteFile(filePath string, content []byte, perm os.FileMode) error {
	return os.WriteFile(filePath, content, perm)
}

// AppendFile 追加文件内容
func AppendFile(filePath string, content []byte, perm os.FileMode) error {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, perm)
	if err != nil {
		return err
	}
	defer file.Close()
	
	_, err = file.Write(content)
	return err
}

// ReadLines 按行读取文件
func ReadLines(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	
	return lines, scanner.Err()
}

// WriteLines 按行写入文件
func WriteLines(filePath string, lines []string, perm os.FileMode) error {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, perm)
	if err != nil {
		return err
	}
	defer file.Close()
	
	writer := bufio.NewWriter(file)
	for _, line := range lines {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}
	
	return writer.Flush()
}

// CopyFile 复制文件
func CopyFile(src, dst string, perm os.FileMode) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	
	dstFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return err
	}
	defer dstFile.Close()
	
	_, err = io.Copy(dstFile, srcFile)
	return err
}

// MoveFile 移动文件
func MoveFile(src, dst string) error {
	return os.Rename(src, dst)
}

// DeleteFile 删除文件
func DeleteFile(filePath string) error {
	return os.Remove(filePath)
}

// CreateDir 创建目录
func CreateDir(dirPath string, perm os.FileMode) error {
	return os.Mkdir(dirPath, perm)
}

// CreateDirAll 递归创建目录
func CreateDirAll(dirPath string, perm os.FileMode) error {
	return os.MkdirAll(dirPath, perm)
}

// RemoveDir 删除目录
func RemoveDir(dirPath string) error {
	return os.Remove(dirPath)
}

// RemoveDirAll 递归删除目录
func RemoveDirAll(dirPath string) error {
	return os.RemoveAll(dirPath)
}

// ListFiles 列出目录中的文件
func ListFiles(dirPath string) ([]string, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}
	
	var files []string
	for _, entry := range entries {
		if !entry.IsDir() {
			files = append(files, entry.Name())
		}
	}
	
	return files, nil
}

// ListFilesRecursive 递归列出目录中的文件
func ListFilesRecursive(dirPath string) ([]string, error) {
	var files []string
	
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if !info.IsDir() {
			files = append(files, path)
		}
		
		return nil
	})
	
	return files, err
}

// GetFileInfo 获取文件信息
func GetFileInfo(filePath string) (os.FileInfo, error) {
	return os.Stat(filePath)
}

// GetFileSize 获取文件大小
func GetFileSize(filePath string) (int64, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return 0, err
	}
	
	if info.IsDir() {
		return 0, errors.New("path is a directory, not a file")
	}
	
	return info.Size(), nil
}

// GetFileModTime 获取文件修改时间
func GetFileModTime(filePath string) (int64, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return 0, err
	}
	
	return info.ModTime().Unix(), nil
}

// IsFile 检查路径是否为文件
func IsFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	
	return !info.IsDir()
}

// IsDir 检查路径是否为目录
func IsDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	
	return info.IsDir()
}

// Exists 检查路径是否存在
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}
