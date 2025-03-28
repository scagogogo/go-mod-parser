package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/scagogogo/go-mod-parser/pkg"
	"github.com/scagogogo/go-mod-parser/pkg/module"
)

func main() {
	var goModPath string
	var listReplaces bool
	var listExcludes bool
	var checkLocalReplaces bool

	// 解析命令行参数
	flag.StringVar(&goModPath, "f", "", "go.mod文件路径")
	flag.BoolVar(&listReplaces, "replaces", false, "列出所有替换规则")
	flag.BoolVar(&listExcludes, "excludes", false, "列出所有排除规则")
	flag.BoolVar(&checkLocalReplaces, "local", false, "检查本地替换")
	flag.Parse()

	// 解析go.mod文件
	var mod *module.Module
	var err error

	if goModPath == "" {
		// 如果未指定go.mod文件路径，尝试在当前目录及父目录中查找
		mod, err = pkg.FindAndParseGoModInCurrentDir()
		if err != nil {
			log.Fatalf("查找并解析go.mod文件失败: %v", err)
		}
	} else {
		// 解析指定的go.mod文件
		mod, err = pkg.ParseGoModFile(goModPath)
		if err != nil {
			log.Fatalf("解析go.mod文件失败: %v", err)
		}
	}

	// 如果没有指定选项，默认显示替换和排除规则
	if !listReplaces && !listExcludes && !checkLocalReplaces {
		listReplaces = true
		listExcludes = true
	}

	// 列出所有替换规则
	if listReplaces {
		fmt.Println("替换规则:")
		if len(mod.Replaces) == 0 {
			fmt.Println("  无替换规则")
		} else {
			for _, rep := range mod.Replaces {
				oldVersion := ""
				if rep.Old.Version != "" {
					oldVersion = " " + rep.Old.Version
				}
				fmt.Printf("  %s%s => %s %s\n", rep.Old.Path, oldVersion, rep.New.Path, rep.New.Version)
			}
		}
	}

	// 列出所有排除规则
	if listExcludes {
		fmt.Println("\n排除规则:")
		if len(mod.Excludes) == 0 {
			fmt.Println("  无排除规则")
		} else {
			for _, exc := range mod.Excludes {
				fmt.Printf("  %s %s\n", exc.Path, exc.Version)
			}
		}
	}

	// 检查本地替换
	if checkLocalReplaces {
		fmt.Println("\n本地替换:")
		localCount := 0
		for _, rep := range mod.Replaces {
			if isLocalReplacement(rep.New.Path) {
				oldVersion := ""
				if rep.Old.Version != "" {
					oldVersion = " " + rep.Old.Version
				}
				fmt.Printf("  %s%s => %s\n", rep.Old.Path, oldVersion, rep.New.Path)

				// 检查路径是否存在
				if _, err := os.Stat(rep.New.Path); os.IsNotExist(err) {
					fmt.Printf("    警告: 本地路径不存在\n")
				}

				localCount++
			}
		}

		if localCount == 0 {
			fmt.Println("  无本地替换")
		} else {
			fmt.Printf("\n总共 %d 个本地替换\n", localCount)
		}
	}
}

// isLocalReplacement 检查是否为本地路径替换
func isLocalReplacement(path string) bool {
	// 本地路径通常以./, ../, /开头
	return len(path) > 0 && (path[0] == '.' || path[0] == '/')
}
