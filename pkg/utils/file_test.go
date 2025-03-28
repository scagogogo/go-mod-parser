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
