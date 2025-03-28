package pkg

import (
	"github.com/scagogogo/go-mod-parser/pkg/module"
	"github.com/scagogogo/go-mod-parser/pkg/parser"
)

// ParseGoModFile 解析指定路径的go.mod文件
func ParseGoModFile(path string) (*module.Module, error) {
	return parser.ParseGoModFile(path)
}

// ParseGoModContent 解析go.mod文件内容
func ParseGoModContent(content string) (*module.Module, error) {
	return parser.ParseGoModContent(content)
}

// FindAndParseGoModFile 在指定目录及其父目录中查找并解析go.mod文件
func FindAndParseGoModFile(dir string) (*module.Module, error) {
	return parser.FindAndParseGoModFile(dir)
}

// FindAndParseGoModInCurrentDir 在当前目录及其父目录中查找并解析go.mod文件
func FindAndParseGoModInCurrentDir() (*module.Module, error) {
	return parser.FindAndParseGoModInCurrentDir()
}

// 以下是便捷函数，帮助用户检查和访问go.mod文件的不同部分

// HasRequire 检查模块是否有特定的依赖
func HasRequire(mod *module.Module, path string) bool {
	return parser.HasRequire(mod, path)
}

// GetRequire 获取模块的特定依赖
func GetRequire(mod *module.Module, path string) *module.Require {
	return parser.GetRequire(mod, path)
}

// HasReplace 检查模块是否有特定的替换规则
func HasReplace(mod *module.Module, path string) bool {
	return parser.HasReplace(mod, path)
}

// GetReplace 获取模块的特定替换规则
func GetReplace(mod *module.Module, path string) *module.Replace {
	return parser.GetReplace(mod, path)
}

// HasExclude 检查模块是否有特定的排除规则
func HasExclude(mod *module.Module, path, version string) bool {
	return parser.HasExclude(mod, path, version)
}

// HasRetract 检查模块是否有特定的撤回版本
func HasRetract(mod *module.Module, version string) bool {
	return parser.HasRetract(mod, version)
}
