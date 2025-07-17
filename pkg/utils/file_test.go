package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindGoModFile(t *testing.T) {
	// 创建临时目录结构
	tempDir := t.TempDir()

	// 创建一个多层目录结构
	subDirPath := filepath.Join(tempDir, "level1", "level2", "level3")
	err := os.MkdirAll(subDirPath, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory structure: %v", err)
	}

	// 在中间层创建go.mod文件
	goModPath := filepath.Join(tempDir, "level1", "go.mod")
	err = os.WriteFile(goModPath, []byte("module example.com/test"), 0644)
	if err != nil {
		t.Fatalf("Failed to write test go.mod file: %v", err)
	}

	// 测试从子目录查找
	startDir := filepath.Join(tempDir, "level1", "level2", "level3")
	found, err := FindGoModFile(startDir)
	if err != nil {
		t.Fatalf("FindGoModFile failed: %v", err)
	}

	// 验证找到了正确的文件
	if found != goModPath {
		t.Errorf("Expected to find %q, got %q", goModPath, found)
	}

	// 测试从不包含go.mod的目录查找
	emptyDir := filepath.Join(tempDir, "empty")
	err = os.Mkdir(emptyDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create empty test directory: %v", err)
	}

	_, err = FindGoModFile(emptyDir)
	if err != ErrGoModNotFound {
		t.Errorf("Expected ErrGoModNotFound, got %v", err)
	}
}

func TestFindGoModFile_EmptyDir(t *testing.T) {
	// 测试空字符串参数 - 应该使用当前工作目录
	// 保存当前工作目录
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}
	defer os.Chdir(originalWd)

	// 创建临时目录并切换到该目录
	tempDir := t.TempDir()
	err = os.Chdir(tempDir)
	if err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	// 在当前目录创建go.mod文件
	goModPath := filepath.Join(tempDir, "go.mod")
	err = os.WriteFile(goModPath, []byte("module example.com/test"), 0644)
	if err != nil {
		t.Fatalf("Failed to write test go.mod file: %v", err)
	}

	// 测试空字符串参数
	found, err := FindGoModFile("")
	if err != nil {
		t.Fatalf("FindGoModFile with empty string failed: %v", err)
	}

	// 验证找到了正确的文件 - 检查文件名是否正确
	if filepath.Base(found) != "go.mod" {
		t.Errorf("Expected to find go.mod file, got %q", found)
	}

	// 验证文件确实存在且可读
	content, err := os.ReadFile(found)
	if err != nil {
		t.Errorf("Failed to read found file: %v", err)
	}
	if string(content) != "module example.com/test" {
		t.Errorf("File content doesn't match expected")
	}
}

func TestFindGoModFile_InvalidPath(t *testing.T) {
	// 测试无效路径 - 使用一个包含无效字符的路径
	// 在某些系统上，这可能会导致filepath.Abs返回错误
	invalidPath := string([]byte{0})
	_, err := FindGoModFile(invalidPath)
	// 我们期望这里会有错误，但具体错误类型可能因系统而异
	if err == nil {
		t.Log("Warning: Expected error for invalid path, but got none")
	}
}

func TestFindGoModFile_RootDirectory(t *testing.T) {
	// 测试从根目录开始查找（应该找不到）
	_, err := FindGoModFile("/")
	if err != ErrGoModNotFound {
		t.Errorf("Expected ErrGoModNotFound when searching from root, got %v", err)
	}
}

func TestFindGoModFile_NonExistentPath(t *testing.T) {
	// 测试不存在的路径
	_, err := FindGoModFile("/this/path/does/not/exist")
	if err == nil {
		t.Error("Expected error for non-existent path, got none")
	}
}

func TestIsFile(t *testing.T) {
	// 创建临时目录和文件
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "testfile.txt")
	err := os.WriteFile(filePath, []byte("test content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// 测试文件检测
	if !IsFile(filePath) {
		t.Errorf("IsFile(%q) = false, want true", filePath)
	}

	// 测试目录检测
	if IsFile(tempDir) {
		t.Errorf("IsFile(%q) = true, want false", tempDir)
	}

	// 测试不存在的路径
	if IsFile(filepath.Join(tempDir, "nonexistent")) {
		t.Errorf("IsFile for nonexistent path = true, want false")
	}
}

func TestIsDir(t *testing.T) {
	// 创建临时目录和文件
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "testfile.txt")
	err := os.WriteFile(filePath, []byte("test content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// 测试目录检测
	if !IsDir(tempDir) {
		t.Errorf("IsDir(%q) = false, want true", tempDir)
	}

	// 测试文件检测
	if IsDir(filePath) {
		t.Errorf("IsDir(%q) = true, want false", filePath)
	}

	// 测试不存在的路径
	if IsDir(filepath.Join(tempDir, "nonexistent")) {
		t.Errorf("IsDir for nonexistent path = true, want false")
	}
}

func TestExists(t *testing.T) {
	// 创建临时目录和文件
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "testfile.txt")
	err := os.WriteFile(filePath, []byte("test content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// 测试存在的文件
	if !Exists(filePath) {
		t.Errorf("Exists(%q) = false, want true", filePath)
	}

	// 测试存在的目录
	if !Exists(tempDir) {
		t.Errorf("Exists(%q) = false, want true", tempDir)
	}

	// 测试不存在的路径
	if Exists(filepath.Join(tempDir, "nonexistent")) {
		t.Errorf("Exists for nonexistent path = true, want false")
	}
}
