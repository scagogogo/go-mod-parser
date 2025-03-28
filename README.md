# Go Mod Parser

[![Go Reference](https://pkg.go.dev/badge/github.com/scagogogo/go-mod-parser.svg)](https://pkg.go.dev/github.com/scagogogo/go-mod-parser)
[![Go Report Card](https://goreportcard.com/badge/github.com/scagogogo/go-mod-parser)](https://goreportcard.com/report/github.com/scagogogo/go-mod-parser)
[![License](https://img.shields.io/github/license/scagogogo/go-mod-parser)](LICENSE)
[![Tests](https://github.com/scagogogo/go-mod-parser/actions/workflows/go-test.yml/badge.svg)](https://github.com/scagogogo/go-mod-parser/actions/workflows/go-test.yml)

Go Mod Parser 是一个功能完整、使用简便的 `go.mod` 文件解析库，它将 go.mod 文件转换为结构化的 Go 对象，使得依赖管理和模块分析变得更加容易。无论是构建依赖分析工具、模块管理系统，还是需要检查项目依赖的 CI/CD 流程，本库都能提供可靠的支持。

## 特性

- ✅ **完整支持所有指令** - 解析 `module`、`go`、`require`、`replace`、`exclude` 和 `retract` 指令
- 🧩 **结构化数据** - 将 go.mod 文件转换为易于使用的 Go 结构体
- 🔍 **自动查找** - 能在项目及父目录中自动定位 go.mod 文件
- 🔄 **依赖分析** - 提供丰富的辅助函数用于分析模块依赖关系
- 📝 **注释支持** - 正确处理 `// indirect` 标记和其他注释
- 🧪 **测试完善** - 完整的单元测试覆盖确保解析的准确性
- 📚 **示例丰富** - 多个实用示例帮助快速上手

## 安装

```bash
go get github.com/scagogogo/go-mod-parser
```

## 快速开始

### 解析指定的 go.mod 文件

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/go-mod-parser/pkg"
)

func main() {
    // 解析指定路径的 go.mod 文件
    mod, err := pkg.ParseGoModFile("path/to/go.mod")
    if err != nil {
        log.Fatalf("解析go.mod文件失败: %v", err)
    }
    
    // 访问解析结果
    fmt.Printf("模块名: %s\n", mod.Name)
    fmt.Printf("Go版本: %s\n", mod.GoVersion)
    
    // 打印所有依赖
    fmt.Println("依赖项:")
    for _, req := range mod.Requires {
        indirect := ""
        if req.Indirect {
            indirect = " // indirect"
        }
        fmt.Printf("- %s %s%s\n", req.Path, req.Version, indirect)
    }
}
```

### 自动查找并解析 go.mod 文件

```go
// 在当前目录及其父目录中查找并解析 go.mod 文件
mod, err := pkg.FindAndParseGoModInCurrentDir()
if err != nil {
    log.Fatalf("查找并解析go.mod文件失败: %v", err)
}

fmt.Printf("找到并解析模块: %s\n", mod.Name)
```

### 解析 go.mod 内容字符串

```go
content := `module github.com/example/module

go 1.21

require github.com/stretchr/testify v1.8.4
`

// 解析 go.mod 内容
mod, err := pkg.ParseGoModContent(content)
if err != nil {
    log.Fatalf("解析go.mod内容失败: %v", err)
}

fmt.Printf("模块名: %s\n", mod.Name)
```

## 主要功能

### 1. 完整解析 go.mod 文件结构

解析 go.mod 文件中的所有标准指令：

- **module** - 模块声明
- **go** - Go 版本要求
- **require** - 依赖声明（包括 indirect 标记）
- **replace** - 替换规则
- **exclude** - 排除规则
- **retract** - 版本撤回（支持单版本和版本范围）

### 2. 丰富的辅助函数

```go
// 检查特定依赖
if pkg.HasRequire(mod, "github.com/stretchr/testify") {
    req := pkg.GetRequire(mod, "github.com/stretchr/testify")
    fmt.Printf("依赖版本: %s (间接依赖: %v)\n", req.Version, req.Indirect)
}

// 检查替换规则
if pkg.HasReplace(mod, "github.com/old/pkg") {
    rep := pkg.GetReplace(mod, "github.com/old/pkg")
    fmt.Printf("替换: %s => %s %s\n", rep.Old.Path, rep.New.Path, rep.New.Version)
}

// 检查排除规则
if pkg.HasExclude(mod, "github.com/problematic/pkg", "v1.0.0") {
    fmt.Println("该版本已被排除")
}

// 检查版本撤回
if pkg.HasRetract(mod, "v1.0.0") {
    fmt.Println("该版本已被撤回")
}
```

### 3. 完整的 API

详见 [pkg.go.dev 文档](https://pkg.go.dev/github.com/scagogogo/go-mod-parser)

| 函数 | 描述 |
|------|------|
| `ParseGoModFile(path)` | 解析指定路径的 go.mod 文件 |
| `ParseGoModContent(content)` | 解析 go.mod 文件内容字符串 |
| `FindAndParseGoModFile(dir)` | 在指定目录及其父目录中查找并解析 go.mod 文件 |
| `FindAndParseGoModInCurrentDir()` | 在当前目录及其父目录中查找并解析 go.mod 文件 |
| `HasRequire(mod, path)` | 检查模块是否有特定的依赖 |
| `GetRequire(mod, path)` | 获取模块的特定依赖 |
| `HasReplace(mod, path)` | 检查模块是否有特定的替换规则 |
| `GetReplace(mod, path)` | 获取模块的特定替换规则 |
| `HasExclude(mod, path, version)` | 检查模块是否有特定的排除规则 |
| `HasRetract(mod, version)` | 检查模块是否有特定的撤回版本 |

## 示例

项目包含多个完整的示例，展示不同使用场景：

- [00_simple_parser](examples/00_simple_parser) - 简单命令行工具示例
- [01_basic_parsing](examples/01_basic_parsing) - 基础解析示例
- [02_find_and_parse](examples/02_find_and_parse) - 查找和解析示例
- [03_check_dependencies](examples/03_check_dependencies) - 依赖检查示例
- [04_replaces_and_excludes](examples/04_replaces_and_excludes) - 替换和排除规则示例
- [05_retract_versions](examples/05_retract_versions) - 版本撤回示例
- [06_programmatic_api](examples/06_programmatic_api) - 编程 API 示例

详细说明请查看 [examples/README.md](examples/README.md)。

## 项目结构

```
pkg/
├── api.go             # 主要公共 API
├── module/            # 模块数据结构定义
├── parser/            # go.mod 文件解析逻辑
└── utils/             # 工具函数
```

## 应用场景

- 构建依赖分析工具
- 模块版本管理系统
- CI/CD 流程中的依赖检查
- Go 项目构建工具
- 模块关系可视化
- 依赖更新推荐系统

## 参考文档

以下是关于 Go 模块和 go.mod 文件格式的官方参考文档：

1. [Go Modules Reference](https://go.dev/ref/mod) - Go 模块系统的权威参考
2. [Go Modules Wiki](https://github.com/golang/go/wiki/Modules) - 更多技术细节和示例
3. [Go 命令文档](https://go.dev/doc/modules/gomod-ref) - go.mod 文件格式详细参考
4. [Go Modules: retract directive](https://go.dev/doc/modules/version-numbers#retract) - retract 指令说明
5. [Go 语言规范](https://go.dev/ref/spec) - Go 语言官方规范

## 贡献

欢迎提交 Issue 和 Pull Request 来改进这个项目。请确保提交前运行测试并保持代码风格一致。

```bash
# 运行测试
go test -v ./...
```

## 许可证

本项目基于 [MIT 许可证](LICENSE) 开源。 