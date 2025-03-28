package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/scagogogo/go-mod-parser/pkg"
	"github.com/scagogogo/go-mod-parser/pkg/utils"
)

func main() {
	var startDir string

	// 获取开始查找的目录
	if len(os.Args) > 1 {
		startDir = os.Args[1]
	} else {
		// 如果未提供目录，使用当前目录
		currentDir, err := os.Getwd()
		if err != nil {
			log.Fatalf("无法获取当前目录: %v", err)
		}
		startDir = currentDir
	}

	fmt.Printf("从目录 '%s' 开始查找go.mod文件...\n", startDir)

	// 查找go.mod文件
	goModPath, err := utils.FindGoModFile(startDir)
	if err != nil {
		log.Fatalf("查找go.mod文件失败: %v", err)
	}

	fmt.Printf("找到go.mod文件: %s\n", goModPath)

	// 解析找到的go.mod文件
	mod, err := pkg.ParseGoModFile(goModPath)
	if err != nil {
		log.Fatalf("解析go.mod文件失败: %v", err)
	}

	// 打印模块信息
	fmt.Printf("\n模块名称: %s\n", mod.Name)
	fmt.Printf("Go版本: %s\n", mod.GoVersion)
	fmt.Printf("依赖项数量: %d\n", len(mod.Requires))
	fmt.Printf("替换规则数量: %d\n", len(mod.Replaces))
	fmt.Printf("排除项数量: %d\n", len(mod.Excludes))
	fmt.Printf("撤回版本数量: %d\n", len(mod.Retracts))

	// 打印go.mod文件所在目录的其他Go文件
	modDir := filepath.Dir(goModPath)
	fmt.Printf("\n目录 '%s' 中的Go文件:\n", modDir)

	err = filepath.Walk(modDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".go" {
			relPath, _ := filepath.Rel(modDir, path)
			fmt.Printf("  %s\n", relPath)
		}
		return nil
	})

	if err != nil {
		log.Printf("查找Go文件时出错: %v", err)
	}
}
