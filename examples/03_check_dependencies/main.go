package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/scagogogo/go-mod-parser/pkg"
	"github.com/scagogogo/go-mod-parser/pkg/module"
)

func main() {
	var goModPath string
	var checkDep string
	var listIndirect bool
	var listDirect bool
	var findPrefix string

	// 解析命令行参数
	flag.StringVar(&goModPath, "f", "", "go.mod文件路径")
	flag.StringVar(&checkDep, "dep", "", "检查特定依赖是否存在")
	flag.BoolVar(&listIndirect, "indirect", false, "列出所有间接依赖")
	flag.BoolVar(&listDirect, "direct", false, "列出所有直接依赖")
	flag.StringVar(&findPrefix, "prefix", "", "查找特定前缀的依赖")
	flag.Parse()

	if goModPath == "" {
		// 如果未指定go.mod文件路径，尝试在当前目录及父目录中查找
		var err error
		mod, err := pkg.FindAndParseGoModInCurrentDir()
		if err != nil {
			log.Fatalf("查找并解析go.mod文件失败: %v", err)
		}
		analyzeModule(mod, checkDep, listIndirect, listDirect, findPrefix)
	} else {
		// 解析指定的go.mod文件
		mod, err := pkg.ParseGoModFile(goModPath)
		if err != nil {
			log.Fatalf("解析go.mod文件失败: %v", err)
		}
		analyzeModule(mod, checkDep, listIndirect, listDirect, findPrefix)
	}
}

// analyzeModule 分析模块依赖项
func analyzeModule(mod *module.Module, checkDep string, listIndirect, listDirect bool, findPrefix string) {
	// 检查特定依赖
	if checkDep != "" {
		if pkg.HasRequire(mod, checkDep) {
			req := pkg.GetRequire(mod, checkDep)
			fmt.Printf("依赖 %s:\n", checkDep)
			fmt.Printf("  版本: %s\n", req.Version)
			fmt.Printf("  间接依赖: %v\n", req.Indirect)

			// 检查是否存在替换规则
			replacement, found := findReplacement(mod, checkDep)
			if found {
				fmt.Printf("  替换为: %s %s\n", replacement.New.Path, replacement.New.Version)
			}

			// 检查是否存在排除规则
			if isExcluded(mod, checkDep, req.Version) {
				fmt.Printf("  警告: 此版本已被排除\n")
			}
		} else {
			fmt.Printf("依赖 %s 不存在\n", checkDep)
		}
		return
	}

	// 列出所有直接依赖
	if listDirect {
		fmt.Println("直接依赖:")
		count := 0
		for _, req := range mod.Requires {
			if !req.Indirect {
				fmt.Printf("  %s %s\n", req.Path, req.Version)
				count++
			}
		}
		fmt.Printf("共 %d 个直接依赖\n", count)
	}

	// 列出所有间接依赖
	if listIndirect {
		fmt.Println("\n间接依赖:")
		count := 0
		for _, req := range mod.Requires {
			if req.Indirect {
				fmt.Printf("  %s %s\n", req.Path, req.Version)
				count++
			}
		}
		fmt.Printf("共 %d 个间接依赖\n", count)
	}

	// 查找带有特定前缀的依赖
	if findPrefix != "" {
		fmt.Printf("\n前缀为 '%s' 的依赖:\n", findPrefix)
		count := 0
		for _, req := range mod.Requires {
			if strings.HasPrefix(req.Path, findPrefix) {
				indirect := ""
				if req.Indirect {
					indirect = " (间接)"
				}
				fmt.Printf("  %s %s%s\n", req.Path, req.Version, indirect)
				count++
			}
		}
		if count == 0 {
			fmt.Printf("未找到前缀为 '%s' 的依赖\n", findPrefix)
		} else {
			fmt.Printf("共找到 %d 个依赖\n", count)
		}
	}

	// 如果没有其他选项，显示摘要
	if !listDirect && !listIndirect && findPrefix == "" && checkDep == "" {
		directCount := 0
		indirectCount := 0

		for _, req := range mod.Requires {
			if req.Indirect {
				indirectCount++
			} else {
				directCount++
			}
		}

		fmt.Printf("模块 %s 依赖分析:\n", mod.Name)
		fmt.Printf("直接依赖: %d\n", directCount)
		fmt.Printf("间接依赖: %d\n", indirectCount)
		fmt.Printf("总依赖数: %d\n", directCount+indirectCount)
		fmt.Println("\n使用 -direct 查看所有直接依赖")
		fmt.Println("使用 -indirect 查看所有间接依赖")
		fmt.Println("使用 -dep <路径> 检查特定依赖")
		fmt.Println("使用 -prefix <前缀> 查找特定前缀的依赖")
	}
}

// findReplacement 查找依赖的替换规则
func findReplacement(mod *module.Module, path string) (*module.Replace, bool) {
	for _, r := range mod.Replaces {
		if r.Old.Path == path {
			return r, true
		}
	}
	return nil, false
}

// isExcluded 检查依赖是否被排除
func isExcluded(mod *module.Module, path, version string) bool {
	for _, exc := range mod.Excludes {
		if exc.Path == path && exc.Version == version {
			return true
		}
	}
	return false
}
