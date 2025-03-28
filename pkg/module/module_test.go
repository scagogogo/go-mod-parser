package module

import (
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestOpenAndProcess(t *testing.T) {
	// 创建临时文件
	tempDir := t.TempDir()
	testFilePath := filepath.Join(tempDir, "test.mod")
	testContent := "module github.com/example/module\n\ngo 1.21\n"

	// 写入测试内容
	err := os.WriteFile(testFilePath, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// 定义处理函数
	processor := func(r io.Reader) (*Module, error) {
		// 读取内容并验证
		data, err := io.ReadAll(r)
		if err != nil {
			return nil, err
		}

		if string(data) != testContent {
			t.Errorf("Content mismatch: got %q, want %q", string(data), testContent)
		}

		// 返回一个模拟的Module对象
		return &Module{
			Name:      "github.com/example/module",
			GoVersion: "1.21",
		}, nil
	}

	// 调用被测函数
	mod, err := OpenAndProcess(testFilePath, processor)
	if err != nil {
		t.Fatalf("OpenAndProcess failed: %v", err)
	}

	// 验证结果
	if mod.Name != "github.com/example/module" {
		t.Errorf("Expected module name 'github.com/example/module', got %q", mod.Name)
	}
	if mod.GoVersion != "1.21" {
		t.Errorf("Expected Go version '1.21', got %q", mod.GoVersion)
	}
}

func TestOpenAndProcess_ErrorCase(t *testing.T) {
	// 使用不存在的文件路径
	_, err := OpenAndProcess("/non/existent/file.mod", func(r io.Reader) (*Module, error) {
		return nil, nil
	})

	// 应该返回错误
	if err == nil {
		t.Error("Expected error when opening non-existent file, got nil")
	}
}
