package main

import (
	"fmt"
	"log"

	"github.com/scagogogo/go-mod-parser/pkg"
	"github.com/scagogogo/go-mod-parser/pkg/module"
)

func main() {
	// 示例: 从字符串解析go.mod内容
	fmt.Println("从字符串解析go.mod示例")
	goModContent := `module example.com/mymodule

go 1.19

require (
	github.com/gin-gonic/gin v1.8.1
	github.com/sirupsen/logrus v1.9.0 // indirect
)

replace github.com/gin-gonic/gin => github.com/custom/gin v1.8.2

exclude github.com/sirupsen/logrus v1.8.0

retract v1.0.0 // 严重bug
`

	// 解析go.mod内容
	mod, err := pkg.ParseGoModContent(goModContent)
	if err != nil {
		log.Fatalf("解析go.mod内容失败: %v", err)
	}

	// 打印基本信息
	fmt.Printf("模块名: %s\n", mod.Name)
	fmt.Printf("Go版本: %s\n", mod.GoVersion)

	// 使用辅助函数检查依赖
	fmt.Println("\n检查特定依赖:")
	checkDependency(mod, "github.com/gin-gonic/gin")
	checkDependency(mod, "github.com/sirupsen/logrus")
	checkDependency(mod, "github.com/unknown/pkg") // 不存在的依赖

	// 检查版本是否被撤回
	fmt.Println("\n检查版本撤回:")
	checkRetraction(mod, "v1.0.0")
	checkRetraction(mod, "v2.0.0") // 未撤回的版本
}

// 检查特定依赖
func checkDependency(mod *module.Module, path string) {
	if pkg.HasRequire(mod, path) {
		req := pkg.GetRequire(mod, path)
		fmt.Printf("✓ 依赖 %s 存在:\n", path)
		fmt.Printf("  - 版本: %s\n", req.Version)
		fmt.Printf("  - 间接依赖: %v\n", req.Indirect)

		// 检查该依赖是否有替换规则
		if pkg.HasReplace(mod, path) {
			rep := pkg.GetReplace(mod, path)
			fmt.Printf("  - 替换为: %s %s\n", rep.New.Path, rep.New.Version)
		}
	} else {
		fmt.Printf("✗ 依赖 %s 不存在\n", path)
	}
}

// 检查版本是否被撤回
func checkRetraction(mod *module.Module, version string) {
	if pkg.HasRetract(mod, version) {
		fmt.Printf("✓ 版本 %s 已被撤回\n", version)
	} else {
		fmt.Printf("✗ 版本 %s 未被撤回\n", version)
	}
}
