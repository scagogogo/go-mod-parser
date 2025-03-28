package parser

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/scagogogo/go-mod-parser/pkg/module"
	"github.com/scagogogo/go-mod-parser/pkg/utils"
)

// ParseFromReader 从io.Reader解析go.mod文件
func ParseFromReader(r io.Reader) (*module.Module, error) {
	mod := &module.Module{
		Requires: make([]*module.Require, 0),
		Replaces: make([]*module.Replace, 0),
		Excludes: make([]*module.Exclude, 0),
		Retracts: make([]*module.Retract, 0),
	}

	scanner := bufio.NewScanner(r)
	lineNum := 0
	inBlock := false
	blockType := ""

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// 跳过空行或注释
		if line == "" || (strings.HasPrefix(line, "//") && !strings.Contains(line, "indirect")) {
			continue
		}

		// 检查块结束
		if inBlock && line == ")" {
			inBlock = false
			continue
		}

		// 检查块开始
		if strings.HasSuffix(line, "(") {
			inBlock = true
			blockType = strings.TrimSpace(strings.TrimSuffix(line, "("))
			continue
		}

		if inBlock {
			if err := handleBlockLine(mod, blockType, line); err != nil {
				return nil, fmt.Errorf("line %d: %w", lineNum, err)
			}
			continue
		}

		// 非块内容解析
		if err := handleSingleLine(mod, line); err != nil {
			return nil, fmt.Errorf("line %d: %w", lineNum, err)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return mod, nil
}

// ParseFromString 从字符串解析go.mod文件
func ParseFromString(s string) (*module.Module, error) {
	return ParseFromReader(strings.NewReader(s))
}

// ParseFromFile 从文件解析go.mod文件
func ParseFromFile(path string) (*module.Module, error) {
	return module.OpenAndProcess(path, ParseFromReader)
}

// handleSingleLine 处理单行语句
func handleSingleLine(mod *module.Module, line string) error {
	// 尝试解析模块名
	if handled, err := parseModuleName(mod, line); err != nil {
		return err
	} else if handled {
		return nil
	}

	// 尝试解析Go版本
	if handled, err := parseGoVersion(mod, line); err != nil {
		return err
	} else if handled {
		return nil
	}

	// 尝试解析单行require
	if handled, err := parseRequireSingleLine(mod, line); err != nil {
		return err
	} else if handled {
		return nil
	}

	// 尝试解析单行replace
	if handled, err := parseReplaceSingleLine(mod, line); err != nil {
		return err
	} else if handled {
		return nil
	}

	// 尝试解析单行exclude
	if handled, err := parseExcludeSingleLine(mod, line); err != nil {
		return err
	} else if handled {
		return nil
	}

	// 尝试解析单行retract
	if handled, err := parseRetractSingleLine(mod, line); err != nil {
		return err
	} else if handled {
		return nil
	}

	return fmt.Errorf("unrecognized line format: %s", line)
}

// handleBlockLine 处理块内的语句
func handleBlockLine(mod *module.Module, blockType, line string) error {
	switch blockType {
	case "require":
		return parseRequireBlockLine(mod, line)
	case "replace":
		return parseReplaceBlockLine(mod, line)
	case "exclude":
		return parseExcludeBlockLine(mod, line)
	case "retract":
		return parseRetractBlockLine(mod, line)
	default:
		return fmt.Errorf("unknown block type: %s", blockType)
	}
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

// ParseGoModFile 解析指定路径的go.mod文件
func ParseGoModFile(path string) (*module.Module, error) {
	return ParseFromFile(path)
}

// ParseGoModContent 解析go.mod文件内容
func ParseGoModContent(content string) (*module.Module, error) {
	return ParseFromString(content)
}
