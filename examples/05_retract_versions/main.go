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
	var listRetracts bool
	var checkVersion string
	var printRationale bool

	// 解析命令行参数
	flag.StringVar(&goModPath, "f", "", "go.mod文件路径")
	flag.BoolVar(&listRetracts, "list", false, "列出所有撤回版本")
	flag.StringVar(&checkVersion, "check", "", "检查特定版本是否被撤回")
	flag.BoolVar(&printRationale, "why", false, "显示撤回原因")
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

	// 检查特定版本是否被撤回
	if checkVersion != "" {
		if pkg.HasRetract(mod, checkVersion) {
			fmt.Printf("版本 %s 已被撤回\n", checkVersion)

			// 如果需要显示撤回原因
			if printRationale {
				rationales := getRetractRationales(mod, checkVersion)
				if len(rationales) > 0 {
					fmt.Println("撤回原因:")
					for _, rationale := range rationales {
						fmt.Printf("  %s\n", rationale)
					}
				} else {
					fmt.Println("无撤回原因")
				}
			}
		} else {
			fmt.Printf("版本 %s 未被撤回\n", checkVersion)
		}
		return
	}

	// 列出所有撤回版本
	if listRetracts || (!listRetracts && checkVersion == "") {
		fmt.Println("撤回版本列表:")
		if len(mod.Retracts) == 0 {
			fmt.Println("  无撤回版本")
		} else {
			for _, ret := range mod.Retracts {
				var versionInfo string
				if ret.Version != "" {
					versionInfo = ret.Version
				} else {
					versionInfo = fmt.Sprintf("[%s, %s]", ret.VersionLow, ret.VersionHigh)
				}

				if printRationale && ret.Rationale != "" {
					fmt.Printf("  %s // %s\n", versionInfo, ret.Rationale)
				} else {
					fmt.Printf("  %s\n", versionInfo)
				}
			}
		}

		fmt.Printf("\n总共 %d 个撤回声明\n", len(mod.Retracts))
	}
}

// getRetractRationales 获取版本的所有撤回原因
func getRetractRationales(mod *module.Module, version string) []string {
	var rationales []string

	for _, ret := range mod.Retracts {
		// 检查单个版本的撤回
		if ret.Version != "" && ret.Version == version && ret.Rationale != "" {
			rationales = append(rationales, ret.Rationale)
			continue
		}

		// 检查版本范围的撤回
		if ret.Version == "" && isVersionInRange(version, ret.VersionLow, ret.VersionHigh) && ret.Rationale != "" {
			rationales = append(rationales, ret.Rationale)
		}
	}

	// 移除重复的原因
	return uniqueStrings(rationales)
}

// isVersionInRange 检查版本是否在范围内
// 注意：这是一个简化的实现，实际应该使用semver库进行完整的语义化版本比较
func isVersionInRange(version, low, high string) bool {
	// 移除版本号前缀 'v' 以便于比较
	version = strings.TrimPrefix(version, "v")
	low = strings.TrimPrefix(low, "v")
	high = strings.TrimPrefix(high, "v")

	// 简单比较字符串（非严格的语义化版本比较）
	return version >= low && version <= high
}

// uniqueStrings 移除字符串切片中的重复项
func uniqueStrings(strings []string) []string {
	keys := make(map[string]bool)
	var list []string
	for _, entry := range strings {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
