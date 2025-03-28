package main

import (
	"fmt"
	"log"
	"os"

	"github.com/scagogogo/go-mod-parser/pkg"
)

func main() {
	// 检查命令行参数
	if len(os.Args) != 2 {
		fmt.Println("使用: go run main.go <go.mod文件路径>")
		os.Exit(1)
	}

	// 获取go.mod文件路径
	goModPath := os.Args[1]

	// 解析go.mod文件
	mod, err := pkg.ParseGoModFile(goModPath)
	if err != nil {
		log.Fatalf("解析go.mod文件失败: %v", err)
	}

	// 打印模块基本信息
	fmt.Printf("模块名称: %s\n", mod.Name)
	fmt.Printf("Go版本: %s\n", mod.GoVersion)

	// 打印依赖项
	fmt.Println("\n依赖项:")
	if len(mod.Requires) == 0 {
		fmt.Println("  无依赖项")
	} else {
		for _, req := range mod.Requires {
			indirect := ""
			if req.Indirect {
				indirect = " // indirect"
			}
			fmt.Printf("  %s %s%s\n", req.Path, req.Version, indirect)
		}
	}
}
