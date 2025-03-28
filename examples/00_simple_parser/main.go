package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/scagogogo/go-mod-parser/pkg"
	"github.com/scagogogo/go-mod-parser/pkg/module"
)

func main() {
	var filePath string
	var findInDir bool
	var checkDep string
	var checkRetract string
	var prettyPrint bool

	flag.StringVar(&filePath, "file", "", "go.mod文件路径，如果不指定则自动查找")
	flag.BoolVar(&findInDir, "find", false, "是否在当前目录及父目录中查找go.mod文件")
	flag.StringVar(&checkDep, "dep", "", "检查特定依赖是否存在及其信息")
	flag.StringVar(&checkRetract, "retract", "", "检查特定版本是否已被撤回")
	flag.BoolVar(&prettyPrint, "pretty", false, "是否以易读方式打印而非JSON格式")
	flag.Parse()

	var mod *module.Module
	var err error

	switch {
	case filePath != "":
		// 解析指定的go.mod文件
		mod, err = pkg.ParseGoModFile(filePath)
		if err != nil {
			log.Fatalf("解析go.mod文件失败: %v", err)
		}
	case findInDir:
		// 在当前目录及父目录中查找并解析go.mod文件
		mod, err = pkg.FindAndParseGoModInCurrentDir()
		if err != nil {
			log.Fatalf("查找并解析go.mod文件失败: %v", err)
		}
	default:
		// 没有指定选项，打印使用说明
		fmt.Println("请指定go.mod文件路径或使用-find选项")
		flag.Usage()
		os.Exit(1)
	}

	// 检查特定依赖
	if checkDep != "" {
		if pkg.HasRequire(mod, checkDep) {
			req := pkg.GetRequire(mod, checkDep)
			fmt.Printf("找到依赖 %s:\n", checkDep)
			fmt.Printf("  版本: %s\n", req.Version)
			fmt.Printf("  间接依赖: %v\n", req.Indirect)
		} else {
			fmt.Printf("依赖 %s 不存在\n", checkDep)
		}
		return
	}

	// 检查版本是否被撤回
	if checkRetract != "" {
		if pkg.HasRetract(mod, checkRetract) {
			fmt.Printf("版本 %s 已被撤回\n", checkRetract)
		} else {
			fmt.Printf("版本 %s 未被撤回\n", checkRetract)
		}
		return
	}

	// 普通打印或JSON格式打印
	if prettyPrint {
		printModulePretty(mod)
	} else {
		// 将解析结果转换为JSON并打印
		jsonData, err := json.MarshalIndent(mod, "", "  ")
		if err != nil {
			log.Fatalf("转换JSON失败: %v", err)
		}
		fmt.Println(string(jsonData))
	}
}

// printModulePretty 以易读方式打印Module信息
func printModulePretty(mod *module.Module) {
	fmt.Printf("模块名: %s\n", mod.Name)
	fmt.Printf("Go版本: %s\n", mod.GoVersion)

	if len(mod.Requires) > 0 {
		fmt.Println("\n依赖项:")
		for _, req := range mod.Requires {
			indirect := ""
			if req.Indirect {
				indirect = " // indirect"
			}
			fmt.Printf("- %s %s%s\n", req.Path, req.Version, indirect)
		}
	}

	if len(mod.Replaces) > 0 {
		fmt.Println("\n替换规则:")
		for _, rep := range mod.Replaces {
			oldVersion := ""
			if rep.Old.Version != "" {
				oldVersion = " " + rep.Old.Version
			}
			fmt.Printf("- %s%s => %s %s\n", rep.Old.Path, oldVersion, rep.New.Path, rep.New.Version)
		}
	}

	if len(mod.Excludes) > 0 {
		fmt.Println("\n排除项:")
		for _, exc := range mod.Excludes {
			fmt.Printf("- %s %s\n", exc.Path, exc.Version)
		}
	}

	if len(mod.Retracts) > 0 {
		fmt.Println("\n撤回版本:")
		for _, ret := range mod.Retracts {
			if ret.Version != "" {
				fmt.Printf("- %s", ret.Version)
			} else {
				fmt.Printf("- [%s, %s]", ret.VersionLow, ret.VersionHigh)
			}
			if ret.Rationale != "" {
				fmt.Printf(" // %s", ret.Rationale)
			}
			fmt.Println()
		}
	}
}
