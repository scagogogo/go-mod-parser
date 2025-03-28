package module

import (
	"io"
	"os"
)

// Module 表示一个go.mod文件的内容
type Module struct {
	// Name 模块名称
	Name string

	// GoVersion go版本
	GoVersion string

	// Requires 依赖项
	Requires []*Require

	// Replaces 替换规则
	Replaces []*Replace

	// Excludes 排除规则
	Excludes []*Exclude

	// Retracts 撤回版本
	Retracts []*Retract
}

// Require 表示一个require指令
type Require struct {
	// Path 模块路径
	Path string

	// Version 模块版本
	Version string

	// Indirect 是否为间接依赖
	Indirect bool
}

// Replace 表示一个replace指令
type Replace struct {
	// Old 替换前的模块信息
	Old *ReplaceItem

	// New 替换后的模块信息
	New *ReplaceItem
}

// ReplaceItem 表示replace指令中的模块信息
type ReplaceItem struct {
	// Path 模块路径
	Path string

	// Version 模块版本，可能为空
	Version string
}

// Exclude 表示一个exclude指令
type Exclude struct {
	// Path 模块路径
	Path string

	// Version 模块版本
	Version string
}

// Retract 表示一个retract指令
type Retract struct {
	// Version 单个撤回版本
	Version string

	// VersionRange 撤回版本范围，格式为 [low,high]
	VersionLow  string
	VersionHigh string

	// Rationale 撤回理由
	Rationale string
}

// OpenAndProcess 打开文件并使用处理函数处理
func OpenAndProcess(path string, process func(io.Reader) (*Module, error)) (*Module, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return process(file)
}
