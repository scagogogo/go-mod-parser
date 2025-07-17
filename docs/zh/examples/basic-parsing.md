# 基础解析

本节演示使用 Go Mod Parser 进行基本解析操作。

## 简单文件解析

最基本的操作是从磁盘解析 go.mod 文件：

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/go-mod-parser/pkg"
)

func main() {
    // 解析 go.mod 文件
    mod, err := pkg.ParseGoModFile("go.mod")
    if err != nil {
        log.Fatalf("解析 go.mod 失败: %v", err)
    }
    
    // 显示基本信息
    fmt.Printf("模块: %s\n", mod.Name)
    fmt.Printf("Go 版本: %s\n", mod.GoVersion)
    fmt.Printf("依赖数量: %d\n", len(mod.Requires))
}
```

## 内容解析

直接从字符串解析 go.mod 内容：

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/go-mod-parser/pkg"
)

func main() {
    content := `module github.com/example/project

go 1.21

require (
    github.com/gin-gonic/gin v1.9.1
    github.com/stretchr/testify v1.8.4 // indirect
)
`
    
    mod, err := pkg.ParseGoModContent(content)
    if err != nil {
        log.Fatalf("解析内容失败: %v", err)
    }
    
    fmt.Printf("模块: %s\n", mod.Name)
    fmt.Printf("Go 版本: %s\n", mod.GoVersion)
    
    for _, req := range mod.Requires {
        fmt.Printf("依赖: %s %s\n", req.Path, req.Version)
        if req.Indirect {
            fmt.Println("  (间接)")
        }
    }
}
```

## 错误处理

解析时始终处理错误：

```go
package main

import (
    "fmt"
    "log"
    "os"
    
    "github.com/scagogogo/go-mod-parser/pkg"
)

func main() {
    mod, err := pkg.ParseGoModFile("go.mod")
    if err != nil {
        if os.IsNotExist(err) {
            log.Fatal("go.mod 文件未找到")
        } else {
            log.Fatalf("解析错误: %v", err)
        }
    }
    
    fmt.Printf("成功解析: %s\n", mod.Name)
}
```

## 完整示例

这是一个演示基本解析和综合输出的完整示例：

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/go-mod-parser/pkg"
)

func main() {
    // 解析 go.mod 文件
    mod, err := pkg.ParseGoModFile("go.mod")
    if err != nil {
        log.Fatalf("解析 go.mod 错误: %v", err)
    }
    
    // 打印模块信息
    fmt.Printf("📦 模块: %s\n", mod.Name)
    fmt.Printf("🐹 Go 版本: %s\n", mod.GoVersion)
    
    // 打印依赖
    if len(mod.Requires) > 0 {
        fmt.Printf("\n📋 依赖 (%d):\n", len(mod.Requires))
        for i, req := range mod.Requires {
            fmt.Printf("%d. %s %s", i+1, req.Path, req.Version)
            if req.Indirect {
                fmt.Printf(" (间接)")
            }
            fmt.Println()
        }
    } else {
        fmt.Println("\n📋 未找到依赖")
    }
    
    // 打印替换指令
    if len(mod.Replaces) > 0 {
        fmt.Printf("\n🔄 替换指令 (%d):\n", len(mod.Replaces))
        for i, rep := range mod.Replaces {
            fmt.Printf("%d. %s => %s", i+1, rep.Old.Path, rep.New.Path)
            if rep.New.Version != "" {
                fmt.Printf(" %s", rep.New.Version)
            }
            fmt.Println()
        }
    }
    
    // 打印排除指令
    if len(mod.Excludes) > 0 {
        fmt.Printf("\n🚫 排除指令 (%d):\n", len(mod.Excludes))
        for i, exc := range mod.Excludes {
            fmt.Printf("%d. %s %s\n", i+1, exc.Path, exc.Version)
        }
    }
    
    // 打印撤回指令
    if len(mod.Retracts) > 0 {
        fmt.Printf("\n⚠️  撤回指令 (%d):\n", len(mod.Retracts))
        for i, ret := range mod.Retracts {
            fmt.Printf("%d. ", i+1)
            if ret.Version != "" {
                fmt.Printf("%s", ret.Version)
            } else {
                fmt.Printf("[%s, %s]", ret.VersionLow, ret.VersionHigh)
            }
            if ret.Rationale != "" {
                fmt.Printf(" (%s)", ret.Rationale)
            }
            fmt.Println()
        }
    }
}
```

## 下一步

- [文件发现](/zh/examples/file-discovery) - 了解自动发现功能
- [依赖分析](/zh/examples/dependency-analysis) - 详细分析依赖
- [高级用法](/zh/examples/advanced-usage) - 复杂场景和模式
