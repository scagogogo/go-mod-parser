# 快速开始

本指南将帮助你在几分钟内开始使用 Go Mod Parser。

## 基本用法

### 解析 go.mod 文件

最常见的用例是解析现有的 go.mod 文件：

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/go-mod-parser/pkg"
)

func main() {
    // 通过路径解析 go.mod 文件
    mod, err := pkg.ParseGoModFile("path/to/go.mod")
    if err != nil {
        log.Fatalf("解析 go.mod 失败: %v", err)
    }
    
    // 访问基本信息
    fmt.Printf("模块名称: %s\n", mod.Name)
    fmt.Printf("Go 版本: %s\n", mod.GoVersion)
    
    // 列出所有依赖
    fmt.Println("\n依赖项:")
    for _, req := range mod.Requires {
        indirect := ""
        if req.Indirect {
            indirect = " // indirect"
        }
        fmt.Printf("- %s %s%s\n", req.Path, req.Version, indirect)
    }
}
```

### 解析 go.mod 内容

你也可以直接从字符串解析 go.mod 内容：

```go
content := `module github.com/example/myproject

go 1.21

require (
    github.com/gin-gonic/gin v1.9.1
    github.com/stretchr/testify v1.8.4 // indirect
)

replace github.com/old/package => github.com/new/package v1.0.0

exclude github.com/problematic/package v1.0.0

retract v1.0.1 // 安全漏洞
`

mod, err := pkg.ParseGoModContent(content)
if err != nil {
    log.Fatalf("解析内容失败: %v", err)
}

fmt.Printf("解析的模块: %s\n", mod.Name)
```

### 自动发现 go.mod 文件

库可以自动在当前目录或父目录中查找 go.mod 文件：

```go
// 在当前目录或父目录中查找并解析 go.mod
mod, err := pkg.FindAndParseGoModInCurrentDir()
if err != nil {
    log.Fatalf("查找并解析 go.mod 失败: %v", err)
}

fmt.Printf("找到模块: %s\n", mod.Name)

// 或指定起始目录
mod, err = pkg.FindAndParseGoModFile("/path/to/project")
if err != nil {
    log.Fatalf("查找 go.mod 失败: %v", err)
}
```

## 处理依赖

### 检查依赖是否存在

```go
// 检查特定依赖是否存在
if pkg.HasRequire(mod, "github.com/gin-gonic/gin") {
    fmt.Println("项目使用 Gin 框架")
    
    // 获取依赖的详细信息
    req := pkg.GetRequire(mod, "github.com/gin-gonic/gin")
    fmt.Printf("版本: %s\n", req.Version)
    fmt.Printf("间接依赖: %v\n", req.Indirect)
}
```

### 分析替换指令

```go
// 检查替换指令
if pkg.HasReplace(mod, "github.com/old/package") {
    replace := pkg.GetReplace(mod, "github.com/old/package")
    fmt.Printf("包 %s 被替换为 %s %s\n", 
        replace.Old.Path, replace.New.Path, replace.New.Version)
}
```

### 检查排除的包

```go
// 检查特定版本是否被排除
if pkg.HasExclude(mod, "github.com/problematic/package", "v1.0.0") {
    fmt.Println("problematic package 的 v1.0.0 版本被排除")
}
```

### 检查撤回的版本

```go
// 检查版本是否被撤回
if pkg.HasRetract(mod, "v1.0.1") {
    fmt.Println("版本 v1.0.1 已被撤回")
}
```

## 完整示例

这是一个展示大部分功能的完整示例：

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/go-mod-parser/pkg"
)

func main() {
    // 解析 go.mod 文件
    mod, err := pkg.FindAndParseGoModInCurrentDir()
    if err != nil {
        log.Fatalf("错误: %v", err)
    }
    
    // 打印基本信息
    fmt.Printf("📦 模块: %s\n", mod.Name)
    fmt.Printf("🐹 Go 版本: %s\n", mod.GoVersion)
    
    // 分析依赖
    fmt.Printf("\n📋 依赖项 (%d):\n", len(mod.Requires))
    for _, req := range mod.Requires {
        status := "直接"
        if req.Indirect {
            status = "间接"
        }
        fmt.Printf("  • %s %s (%s)\n", req.Path, req.Version, status)
    }
    
    // 显示替换指令
    if len(mod.Replaces) > 0 {
        fmt.Printf("\n🔄 替换指令 (%d):\n", len(mod.Replaces))
        for _, rep := range mod.Replaces {
            fmt.Printf("  • %s => %s %s\n", 
                rep.Old.Path, rep.New.Path, rep.New.Version)
        }
    }
    
    // 显示排除的包
    if len(mod.Excludes) > 0 {
        fmt.Printf("\n🚫 排除的包 (%d):\n", len(mod.Excludes))
        for _, exc := range mod.Excludes {
            fmt.Printf("  • %s %s\n", exc.Path, exc.Version)
        }
    }
    
    // 显示撤回的版本
    if len(mod.Retracts) > 0 {
        fmt.Printf("\n⚠️  撤回的版本 (%d):\n", len(mod.Retracts))
        for _, ret := range mod.Retracts {
            if ret.Version != "" {
                fmt.Printf("  • %s", ret.Version)
            } else {
                fmt.Printf("  • [%s, %s]", ret.VersionLow, ret.VersionHigh)
            }
            if ret.Rationale != "" {
                fmt.Printf(" (%s)", ret.Rationale)
            }
            fmt.Println()
        }
    }
}
```

## 实用技巧

### 框架检测

```go
func detectFramework(mod *module.Module) {
    frameworks := map[string]string{
        "github.com/gin-gonic/gin":     "Gin",
        "github.com/gorilla/mux":       "Gorilla Mux",
        "github.com/labstack/echo/v4":  "Echo",
        "github.com/gofiber/fiber/v2":  "Fiber",
        "github.com/beego/beego/v2":    "Beego",
    }
    
    fmt.Println("🔍 框架检测:")
    found := false
    for path, name := range frameworks {
        if pkg.HasRequire(mod, path) {
            req := pkg.GetRequire(mod, path)
            fmt.Printf("  ✓ %s %s\n", name, req.Version)
            found = true
        }
    }
    
    if !found {
        fmt.Println("  未检测到常见的 Web 框架")
    }
}
```

### 依赖统计

```go
func analyzeStats(mod *module.Module) {
    direct := 0
    indirect := 0
    
    for _, req := range mod.Requires {
        if req.Indirect {
            indirect++
        } else {
            direct++
        }
    }
    
    fmt.Printf("📊 依赖统计:\n")
    fmt.Printf("  直接依赖: %d\n", direct)
    fmt.Printf("  间接依赖: %d\n", indirect)
    fmt.Printf("  总计: %d\n", len(mod.Requires))
    fmt.Printf("  替换规则: %d\n", len(mod.Replaces))
    fmt.Printf("  排除规则: %d\n", len(mod.Excludes))
    fmt.Printf("  撤回版本: %d\n", len(mod.Retracts))
}
```

### 安全检查

```go
func securityCheck(mod *module.Module) {
    fmt.Println("🔒 安全检查:")
    
    issues := 0
    
    // 检查使用的撤回版本
    for _, req := range mod.Requires {
        if pkg.HasRetract(mod, req.Version) {
            fmt.Printf("  ⚠️  使用了撤回版本: %s %s\n", req.Path, req.Version)
            issues++
        }
    }
    
    // 检查本地路径替换
    for _, rep := range mod.Replaces {
        if strings.HasPrefix(rep.New.Path, "./") || strings.HasPrefix(rep.New.Path, "../") {
            fmt.Printf("  🔍 检测到本地替换: %s => %s\n", rep.Old.Path, rep.New.Path)
            issues++
        }
    }
    
    if issues == 0 {
        fmt.Println("  ✅ 未发现安全问题")
    } else {
        fmt.Printf("  发现 %d 个潜在安全问题\n", issues)
    }
}
```

## 下一步

- 探索 [API 参考](/zh/api/) 获取详细文档
- 查看更多 [示例](/zh/examples/) 了解高级用法模式
- 了解库使用的 [数据结构](/zh/api/data-structures)
