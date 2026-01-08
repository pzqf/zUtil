package zFile

import (
	"os"
	"testing"
)

func TestReadWriteFile(t *testing.T) {
	filePath := "test.txt"
	content := []byte("Hello, zFile!")
	
	// 测试写入文件
	err := WriteFile(filePath, content, 0644)
	if err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}
	
	// 测试读取文件
	readContent, err := ReadFile(filePath)
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}
	
	if string(readContent) != string(content) {
		t.Errorf("Content mismatch: expected %q, got %q", content, readContent)
	}
	
	// 清理
	os.Remove(filePath)
}

func TestAppendFile(t *testing.T) {
	filePath := "test_append.txt"
	content := []byte("Hello, ")
	appendContent := []byte("zFile!")
	
	// 测试写入文件
	err := WriteFile(filePath, content, 0644)
	if err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}
	
	// 测试追加文件
	err = AppendFile(filePath, appendContent, 0644)
	if err != nil {
		t.Fatalf("AppendFile failed: %v", err)
	}
	
	// 测试读取文件
	readContent, err := ReadFile(filePath)
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}
	
	expectedContent := []byte("Hello, zFile!")
	if string(readContent) != string(expectedContent) {
		t.Errorf("Content mismatch: expected %q, got %q", expectedContent, readContent)
	}
	
	// 清理
	os.Remove(filePath)
}

func TestReadWriteLines(t *testing.T) {
	filePath := "test_lines.txt"
	lines := []string{
		"Line 1",
		"Line 2",
		"Line 3",
	}
	
	// 测试写入行
	err := WriteLines(filePath, lines, 0644)
	if err != nil {
		t.Fatalf("WriteLines failed: %v", err)
	}
	
	// 测试读取行
	readLines, err := ReadLines(filePath)
	if err != nil {
		t.Fatalf("ReadLines failed: %v", err)
	}
	
	if len(readLines) != len(lines) {
		t.Errorf("Line count mismatch: expected %d, got %d", len(lines), len(readLines))
	}
	
	for i, line := range lines {
		if readLines[i] != line {
			t.Errorf("Line %d mismatch: expected %q, got %q", i+1, line, readLines[i])
		}
	}
	
	// 清理
	os.Remove(filePath)
}

func TestCopyMoveFile(t *testing.T) {
	srcPath := "test_src.txt"
	dstPath := "test_dst.txt"
	content := []byte("Hello, CopyMove!")
	
	// 创建源文件
	err := WriteFile(srcPath, content, 0644)
	if err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}
	
	// 测试复制文件
	err = CopyFile(srcPath, dstPath, 0644)
	if err != nil {
		t.Fatalf("CopyFile failed: %v", err)
	}
	
	// 测试读取复制后的文件
	readContent, err := ReadFile(dstPath)
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}
	
	if string(readContent) != string(content) {
		t.Errorf("Content mismatch: expected %q, got %q", content, readContent)
	}
	
	// 测试移动文件
	movePath := "test_move.txt"
	err = MoveFile(dstPath, movePath)
	if err != nil {
		t.Fatalf("MoveFile failed: %v", err)
	}
	
	// 检查移动后的文件是否存在
	if !Exists(movePath) {
		t.Errorf("MoveFile failed: file not found at %s", movePath)
	}
	
	// 清理
	os.Remove(srcPath)
	os.Remove(movePath)
}

func TestDirOperations(t *testing.T) {
	dirPath := "test_dir"
	
	// 测试创建目录
	err := CreateDir(dirPath, 0755)
	if err != nil {
		t.Fatalf("CreateDir failed: %v", err)
	}
	
	// 测试目录是否存在
	if !Exists(dirPath) {
		t.Errorf("CreateDir failed: directory not found")
	}
	
	// 测试是否为目录
	if !IsDir(dirPath) {
		t.Errorf("IsDir failed: expected true")
	}
	
	// 测试创建子目录
	subDirPath := dirPath + "/sub_dir"
	err = CreateDirAll(subDirPath, 0755)
	if err != nil {
		t.Fatalf("CreateDirAll failed: %v", err)
	}
	
	// 测试递归删除目录
	err = RemoveDirAll(dirPath)
	if err != nil {
		t.Fatalf("RemoveDirAll failed: %v", err)
	}
	
	// 检查目录是否已删除
	if Exists(dirPath) {
		t.Errorf("RemoveDirAll failed: directory still exists")
	}
}

func TestListFiles(t *testing.T) {
	dirPath := "test_list_dir"
	file1 := "file1.txt"
	file2 := "file2.txt"
	
	// 创建目录
	err := CreateDir(dirPath, 0755)
	if err != nil {
		t.Fatalf("CreateDir failed: %v", err)
	}
	
	// 创建文件
	WriteFile(dirPath+"/"+file1, []byte("Test 1"), 0644)
	WriteFile(dirPath+"/"+file2, []byte("Test 2"), 0644)
	
	// 测试列出文件
	files, err := ListFiles(dirPath)
	if err != nil {
		t.Fatalf("ListFiles failed: %v", err)
	}
	
	if len(files) != 2 {
		t.Errorf("ListFiles failed: expected 2 files, got %d", len(files))
	}
	
	// 清理
	RemoveDirAll(dirPath)
}
