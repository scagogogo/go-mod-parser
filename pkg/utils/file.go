package utils

import (
	"errors"
	"os"
	"path/filepath"
)

var (
	// ErrGoModNotFound 表示在当前目录及父目录中未找到go.mod文件
	ErrGoModNotFound = errors.New("go.mod file not found")
)

// FindGoModFile 在指定目录及其父目录中查找go.mod文件
func FindGoModFile(dir string) (string, error) {
	if dir == "" {
		var err error
		dir, err = os.Getwd()
		if err != nil {
			return "", err
		}
	}

	dir, err := filepath.Abs(dir)
	if err != nil {
		return "", err
	}

	for {
		path := filepath.Join(dir, "go.mod")
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}

		// 已经到达根目录
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return "", ErrGoModNotFound
}

// IsFile 检查指定路径是否是文件
func IsFile(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

// IsDir 检查指定路径是否是目录
func IsDir(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

// Exists 检查指定路径是否存在
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
