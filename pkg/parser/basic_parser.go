package parser

import (
	"strings"

	"github.com/scagogogo/go-mod-parser/pkg/module"
)

// parseModuleName 解析模块名称
func parseModuleName(mod *module.Module, line string) (bool, error) {
	if matches := moduleRegexp.FindStringSubmatch(line); len(matches) == 2 {
		mod.Name = matches[1]
		return true, nil
	}
	return false, nil
}

// parseGoVersion 解析Go版本
func parseGoVersion(mod *module.Module, line string) (bool, error) {
	if matches := goRegexp.FindStringSubmatch(line); len(matches) == 2 {
		mod.GoVersion = matches[1]
		return true, nil
	}
	return false, nil
}

// isIndirect 检查一行是否包含indirect注释
func isIndirect(line string) bool {
	return indirectCommentRegexp.MatchString(line)
}

// blockStartsWith 检查一行是否以指定关键词开始一个块
func blockStartsWith(line, keyword string) bool {
	return strings.HasPrefix(strings.TrimSpace(line), keyword) && strings.HasSuffix(strings.TrimSpace(line), "(")
}

// isBlockEnd 检查一行是否为块结束
func isBlockEnd(line string) bool {
	return strings.TrimSpace(line) == ")"
}
