# 示例

本节提供了全面的示例，演示如何在各种场景中使用 Go Mod Parser。每个示例都包含完整的、可运行的代码和解释。

## 概览

示例按复杂度和用例组织：

- **[基础解析](/zh/examples/basic-parsing)** - 简单的解析操作
- **[文件发现](/zh/examples/file-discovery)** - 自动发现 go.mod 文件
- **[依赖分析](/zh/examples/dependency-analysis)** - 分析依赖和关系
- **[高级用法](/zh/examples/advanced-usage)** - 复杂场景和最佳实践

## 快速示例

### 解析 go.mod 文件

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/go-mod-parser/pkg"
)

func main() {
    mod, err := pkg.ParseGoModFile("go.mod")
    if err != nil {
        log.Fatalf("解析失败: %v", err)
    }
    
    fmt.Printf("模块: %s\n", mod.Name)
    fmt.Printf("Go 版本: %s\n", mod.GoVersion)
    fmt.Printf("依赖数量: %d\n", len(mod.Requires))
}
```

### 检查特定依赖

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/go-mod-parser/pkg"
)

func main() {
    mod, err := pkg.FindAndParseGoModInCurrentDir()
    if err != nil {
        log.Fatalf("查找 go.mod 失败: %v", err)
    }
    
    // 检查流行框架
    frameworks := []string{
        "github.com/gin-gonic/gin",
        "github.com/gorilla/mux", 
        "github.com/labstack/echo/v4",
        "github.com/gofiber/fiber/v2",
    }
    
    fmt.Println("框架检测:")
    for _, framework := range frameworks {
        if pkg.HasRequire(mod, framework) {
            req := pkg.GetRequire(mod, framework)
            fmt.Printf("✓ %s %s\n", framework, req.Version)
        } else {
            fmt.Printf("✗ %s\n", framework)
        }
    }
}
```

### 分析替换指令

```go
package main

import (
    "fmt"
    "log"
    "strings"
    
    "github.com/scagogogo/go-mod-parser/pkg"
)

func main() {
    mod, err := pkg.FindAndParseGoModInCurrentDir()
    if err != nil {
        log.Fatalf("解析失败: %v", err)
    }
    
    if len(mod.Replaces) == 0 {
        fmt.Println("未找到替换指令")
        return
    }
    
    fmt.Printf("找到 %d 个替换指令:\n\n", len(mod.Replaces))
    
    for i, rep := range mod.Replaces {
        fmt.Printf("%d. %s", i+1, rep.Old.Path)
        if rep.Old.Version != "" {
            fmt.Printf(" %s", rep.Old.Version)
        }
        fmt.Printf(" => %s", rep.New.Path)
        if rep.New.Version != "" {
            fmt.Printf(" %s", rep.New.Version)
        }
        
        // 确定替换类型
        if strings.HasPrefix(rep.New.Path, "./") || strings.HasPrefix(rep.New.Path, "../") {
            fmt.Printf(" (本地路径)")
        } else {
            fmt.Printf(" (模块)")
        }
        fmt.Println()
    }
}
```

## 示例项目

### go.mod 分析的 CLI 工具

```go
package main

import (
    "flag"
    "fmt"
    "log"
    "os"
    "path/filepath"
    "strings"
    
    "github.com/scagogogo/go-mod-parser/pkg"
)

func main() {
    var (
        path = flag.String("path", "", "go.mod 文件或目录的路径")
        verbose = flag.Bool("verbose", false, "详细输出")
        checkSecurity = flag.Bool("security", false, "检查安全问题")
    )
    flag.Parse()
    
    var mod *module.Module
    var err error
    
    if *path != "" {
        if strings.HasSuffix(*path, "go.mod") {
            mod, err = pkg.ParseGoModFile(*path)
        } else {
            mod, err = pkg.FindAndParseGoModFile(*path)
        }
    } else {
        mod, err = pkg.FindAndParseGoModInCurrentDir()
    }
    
    if err != nil {
        log.Fatalf("错误: %v", err)
    }
    
    // 基本信息
    fmt.Printf("📦 模块: %s\n", mod.Name)
    fmt.Printf("🐹 Go 版本: %s\n", mod.GoVersion)
    fmt.Printf("📋 依赖: %d\n", len(mod.Requires))
    
    if *verbose {
        printVerboseInfo(mod)
    }
    
    if *checkSecurity {
        checkSecurityIssues(mod)
    }
}

func printVerboseInfo(mod *module.Module) {
    // 打印依赖
    if len(mod.Requires) > 0 {
        fmt.Println("\n📋 依赖:")
        for _, req := range mod.Requires {
            fmt.Printf("  %s %s", req.Path, req.Version)
            if req.Indirect {
                fmt.Printf(" (间接)")
            }
            fmt.Println()
        }
    }
    
    // 打印替换
    if len(mod.Replaces) > 0 {
        fmt.Println("\n🔄 替换:")
        for _, rep := range mod.Replaces {
            fmt.Printf("  %s => %s", rep.Old.Path, rep.New.Path)
            if rep.New.Version != "" {
                fmt.Printf(" %s", rep.New.Version)
            }
            fmt.Println()
        }
    }
    
    // 打印排除
    if len(mod.Excludes) > 0 {
        fmt.Println("\n🚫 排除:")
        for _, exc := range mod.Excludes {
            fmt.Printf("  %s %s\n", exc.Path, exc.Version)
        }
    }
    
    // 打印撤回
    if len(mod.Retracts) > 0 {
        fmt.Println("\n⚠️  撤回:")
        for _, ret := range mod.Retracts {
            if ret.Version != "" {
                fmt.Printf("  %s", ret.Version)
            } else {
                fmt.Printf("  [%s, %s]", ret.VersionLow, ret.VersionHigh)
            }
            if ret.Rationale != "" {
                fmt.Printf(" (%s)", ret.Rationale)
            }
            fmt.Println()
        }
    }
}

func checkSecurityIssues(mod *module.Module) {
    fmt.Println("\n🔒 安全分析:")
    
    issues := 0
    
    // 检查使用中的撤回版本
    for _, req := range mod.Requires {
        if pkg.HasRetract(mod, req.Version) {
            fmt.Printf("⚠️  使用撤回版本: %s %s\n", req.Path, req.Version)
            issues++
        }
    }
    
    // 检查本地路径替换
    for _, rep := range mod.Replaces {
        if strings.HasPrefix(rep.New.Path, "./") || strings.HasPrefix(rep.New.Path, "../") {
            fmt.Printf("🔍 检测到本地替换: %s => %s\n", rep.Old.Path, rep.New.Path)
            issues++
        }
    }
    
    if issues == 0 {
        fmt.Println("✅ 未发现安全问题")
    } else {
        fmt.Printf("发现 %d 个潜在安全问题\n", issues)
    }
}
```

### 依赖比较工具

```go
package main

import (
    "fmt"
    "log"
    "os"
    
    "github.com/scagogogo/go-mod-parser/pkg"
)

func main() {
    if len(os.Args) != 3 {
        fmt.Println("用法: compare <go.mod1> <go.mod2>")
        os.Exit(1)
    }
    
    mod1, err := pkg.ParseGoModFile(os.Args[1])
    if err != nil {
        log.Fatalf("解析 %s 失败: %v", os.Args[1], err)
    }
    
    mod2, err := pkg.ParseGoModFile(os.Args[2])
    if err != nil {
        log.Fatalf("解析 %s 失败: %v", os.Args[2], err)
    }
    
    compareDependencies(mod1, mod2)
}

func compareDependencies(mod1, mod2 *module.Module) {
    fmt.Printf("比较 %s vs %s\n\n", mod1.Name, mod2.Name)
    
    // 创建映射以便查找
    deps1 := make(map[string]string)
    deps2 := make(map[string]string)
    
    for _, req := range mod1.Requires {
        deps1[req.Path] = req.Version
    }
    
    for _, req := range mod2.Requires {
        deps2[req.Path] = req.Version
    }
    
    // 查找共同依赖
    fmt.Println("📋 共同依赖:")
    for path, version1 := range deps1 {
        if version2, exists := deps2[path]; exists {
            if version1 == version2 {
                fmt.Printf("  ✓ %s %s\n", path, version1)
            } else {
                fmt.Printf("  ⚠️  %s: %s vs %s\n", path, version1, version2)
            }
        }
    }
    
    // 查找仅在 mod1 中的依赖
    fmt.Println("\n📦 仅在 " + mod1.Name + " 中:")
    for path, version := range deps1 {
        if _, exists := deps2[path]; !exists {
            fmt.Printf("  + %s %s\n", path, version)
        }
    }
    
    // 查找仅在 mod2 中的依赖
    fmt.Println("\n📦 仅在 " + mod2.Name + " 中:")
    for path, version := range deps2 {
        if _, exists := deps1[path]; !exists {
            fmt.Printf("  + %s %s\n", path, version)
        }
    }
}
```

## 运行示例

1. **保存代码** 到 `.go` 文件
2. **初始化 Go 模块**（如果尚未完成）：
   ```bash
   go mod init example
   ```
3. **添加依赖**：
   ```bash
   go get github.com/scagogogo/go-mod-parser
   ```
4. **运行示例**：
   ```bash
   go run main.go
   ```

## 下一步

探索详细的示例分类：

- **[基础解析](/zh/examples/basic-parsing)** - 从这里开始学习简单用例
- **[文件发现](/zh/examples/file-discovery)** - 了解自动发现功能  
- **[依赖分析](/zh/examples/dependency-analysis)** - 高级依赖分析
- **[高级用法](/zh/examples/advanced-usage)** - 复杂场景和模式

每个部分都包含多个示例，提供完整的源代码和解释。
