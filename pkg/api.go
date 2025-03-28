package pkg

import (
	"github.com/scagogogo/go-mod-parser/pkg/module"
	"github.com/scagogogo/go-mod-parser/pkg/parser"
	"github.com/scagogogo/go-mod-parser/pkg/utils"
)

// ParseGoModFile 解析指定路径的go.mod文件
func ParseGoModFile(path string) (*module.Module, error) {
	return parser.ParseFromFile(path)
}

// ParseGoModContent 解析go.mod文件内容
func ParseGoModContent(content string) (*module.Module, error) {
	return parser.ParseFromString(content)
}

// FindAndParseGoModFile 在指定目录及其父目录中查找并解析go.mod文件
func FindAndParseGoModFile(dir string) (*module.Module, error) {
	path, err := utils.FindGoModFile(dir)
	if err != nil {
		return nil, err
	}
	return ParseGoModFile(path)
}

// FindAndParseGoModInCurrentDir 在当前目录及其父目录中查找并解析go.mod文件
func FindAndParseGoModInCurrentDir() (*module.Module, error) {
	return FindAndParseGoModFile("")
}

// 以下是便捷函数，帮助用户检查和访问go.mod文件的不同部分

// HasRequire 检查模块是否有特定的依赖
func HasRequire(mod *module.Module, path string) bool {
	for _, req := range mod.Requires {
		if req.Path == path {
			return true
		}
	}
	return false
}

// GetRequire 获取模块的特定依赖
func GetRequire(mod *module.Module, path string) *module.Require {
	for _, req := range mod.Requires {
		if req.Path == path {
			return req
		}
	}
	return nil
}

// HasReplace 检查模块是否有特定的替换规则
func HasReplace(mod *module.Module, path string) bool {
	for _, rep := range mod.Replaces {
		if rep.Old.Path == path {
			return true
		}
	}
	return false
}

// GetReplace 获取模块的特定替换规则
func GetReplace(mod *module.Module, path string) *module.Replace {
	for _, rep := range mod.Replaces {
		if rep.Old.Path == path {
			return rep
		}
	}
	return nil
}

// HasExclude 检查模块是否有特定的排除规则
func HasExclude(mod *module.Module, path, version string) bool {
	for _, exc := range mod.Excludes {
		if exc.Path == path && exc.Version == version {
			return true
		}
	}
	return false
}

// HasRetract 检查模块是否有特定的撤回版本
func HasRetract(mod *module.Module, version string) bool {
	for _, ret := range mod.Retracts {
		// 检查单个版本
		if ret.Version == version {
			return true
		}
		// 检查版本范围
		if ret.VersionLow != "" && ret.VersionHigh != "" {
			// 这里我们只做简单的字符串比较，实际上需要更复杂的版本比较逻辑
			if version >= ret.VersionLow && version <= ret.VersionHigh {
				return true
			}
		}
	}
	return false
}
