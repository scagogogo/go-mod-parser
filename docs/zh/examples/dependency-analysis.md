# 依赖分析

本节演示如何分析 go.mod 文件中的依赖、替换、排除和撤回。

## 基本依赖检查

检查特定依赖是否存在：

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
        log.Fatalf("解析 go.mod 失败: %v", err)
    }
    
    // 检查特定依赖
    dependencies := []string{
        "github.com/gin-gonic/gin",
        "github.com/gorilla/mux",
        "github.com/stretchr/testify",
    }
    
    fmt.Println("依赖检查:")
    for _, dep := range dependencies {
        if pkg.HasRequire(mod, dep) {
            req := pkg.GetRequire(mod, dep)
            fmt.Printf("✅ %s %s", dep, req.Version)
            if req.Indirect {
                fmt.Printf(" (间接)")
            }
            fmt.Println()
        } else {
            fmt.Printf("❌ %s (未找到)\n", dep)
        }
    }
}
```

## 框架检测

检测流行的 Go 框架：

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/go-mod-parser/pkg"
)

func detectFrameworks(mod *module.Module) {
    frameworks := map[string]string{
        "github.com/gin-gonic/gin":     "Gin Web 框架",
        "github.com/gorilla/mux":       "Gorilla Mux 路由器",
        "github.com/labstack/echo/v4":  "Echo Web 框架",
        "github.com/gofiber/fiber/v2":  "Fiber Web 框架",
        "github.com/beego/beego/v2":    "Beego 框架",
        "github.com/revel/revel":       "Revel 框架",
    }
    
    fmt.Println("🔍 框架检测:")
    found := false
    
    for path, name := range frameworks {
        if pkg.HasRequire(mod, path) {
            req := pkg.GetRequire(mod, path)
            fmt.Printf("  ✅ %s (%s)\n", name, req.Version)
            found = true
        }
    }
    
    if !found {
        fmt.Println("  ❌ 未检测到流行的 Web 框架")
    }
}

func main() {
    mod, err := pkg.FindAndParseGoModInCurrentDir()
    if err != nil {
        log.Fatalf("解析 go.mod 失败: %v", err)
    }
    
    detectFrameworks(mod)
}
```

## 依赖统计

分析依赖模式：

```go
package main

import (
    "fmt"
    "log"
    "sort"
    "strings"
    
    "github.com/scagogogo/go-mod-parser/pkg"
)

func analyzeDependencyStats(mod *module.Module) {
    direct := 0
    indirect := 0
    domains := make(map[string]int)
    
    for _, req := range mod.Requires {
        if req.Indirect {
            indirect++
        } else {
            direct++
        }
        
        // 从模块路径提取域名
        parts := strings.Split(req.Path, "/")
        if len(parts) > 0 {
            domain := parts[0]
            domains[domain]++
        }
    }
    
    fmt.Printf("📊 依赖统计:\n")
    fmt.Printf("  总依赖数: %d\n", len(mod.Requires))
    fmt.Printf("  直接依赖: %d\n", direct)
    fmt.Printf("  间接依赖: %d\n", indirect)
    fmt.Printf("  替换指令: %d\n", len(mod.Replaces))
    fmt.Printf("  排除指令: %d\n", len(mod.Excludes))
    fmt.Printf("  撤回指令: %d\n", len(mod.Retracts))
    
    // 顶级域名
    fmt.Println("\n🌐 顶级依赖域名:")
    type domainCount struct {
        domain string
        count  int
    }
    
    var sortedDomains []domainCount
    for domain, count := range domains {
        sortedDomains = append(sortedDomains, domainCount{domain, count})
    }
    
    sort.Slice(sortedDomains, func(i, j int) bool {
        return sortedDomains[i].count > sortedDomains[j].count
    })
    
    for i, dc := range sortedDomains {
        if i >= 5 { // 显示前 5 个
            break
        }
        fmt.Printf("  %d. %s (%d 个依赖)\n", i+1, dc.domain, dc.count)
    }
}

func main() {
    mod, err := pkg.FindAndParseGoModInCurrentDir()
    if err != nil {
        log.Fatalf("解析 go.mod 失败: %v", err)
    }
    
    analyzeDependencyStats(mod)
}
```

## 替换指令分析

分析替换模式：

```go
package main

import (
    "fmt"
    "log"
    "strings"
    
    "github.com/scagogogo/go-mod-parser/pkg"
)

func analyzeReplaces(mod *module.Module) {
    if len(mod.Replaces) == 0 {
        fmt.Println("🔄 未找到替换指令")
        return
    }
    
    fmt.Printf("🔄 替换指令分析 (共 %d 个):\n", len(mod.Replaces))
    
    localReplaces := 0
    moduleReplaces := 0
    
    for i, rep := range mod.Replaces {
        fmt.Printf("\n%d. %s", i+1, rep.Old.Path)
        if rep.Old.Version != "" {
            fmt.Printf(" %s", rep.Old.Version)
        }
        fmt.Printf(" => %s", rep.New.Path)
        if rep.New.Version != "" {
            fmt.Printf(" %s", rep.New.Version)
        }
        
        if strings.HasPrefix(rep.New.Path, "./") || strings.HasPrefix(rep.New.Path, "../") {
            fmt.Printf(" (本地路径)")
            localReplaces++
        } else {
            fmt.Printf(" (模块)")
            moduleReplaces++
        }
    }
    
    fmt.Printf("\n\n总结:\n")
    fmt.Printf("  本地路径替换: %d\n", localReplaces)
    fmt.Printf("  模块替换: %d\n", moduleReplaces)
}

func main() {
    mod, err := pkg.FindAndParseGoModInCurrentDir()
    if err != nil {
        log.Fatalf("解析 go.mod 失败: %v", err)
    }
    
    analyzeReplaces(mod)
}
```

## 安全分析

检查潜在的安全问题：

```go
package main

import (
    "fmt"
    "log"
    "strings"
    
    "github.com/scagogogo/go-mod-parser/pkg"
)

func securityAnalysis(mod *module.Module) {
    fmt.Println("🔒 安全分析:")
    
    issues := 0
    
    // 检查使用中的撤回版本
    fmt.Println("\n⚠️  撤回版本检查:")
    retractedFound := false
    for _, req := range mod.Requires {
        if pkg.HasRetract(mod, req.Version) {
            fmt.Printf("  ❌ 使用撤回版本: %s %s\n", req.Path, req.Version)
            issues++
            retractedFound = true
        }
    }
    if !retractedFound {
        fmt.Println("  ✅ 未使用撤回版本")
    }
    
    // 检查本地路径替换
    fmt.Println("\n🔍 本地路径替换检查:")
    localFound := false
    for _, rep := range mod.Replaces {
        if strings.HasPrefix(rep.New.Path, "./") || strings.HasPrefix(rep.New.Path, "../") {
            fmt.Printf("  ⚠️  本地替换: %s => %s\n", rep.Old.Path, rep.New.Path)
            issues++
            localFound = true
        }
    }
    if !localFound {
        fmt.Println("  ✅ 无本地路径替换")
    }
    
    // 检查生产环境中的开发依赖
    fmt.Println("\n🧪 开发依赖检查:")
    devDeps := []string{"testify", "mock", "test", "debug", "dev"}
    devFound := false
    for _, req := range mod.Requires {
        if !req.Indirect {
            for _, devKeyword := range devDeps {
                if strings.Contains(strings.ToLower(req.Path), devKeyword) {
                    fmt.Printf("  ⚠️  潜在开发依赖作为直接依赖: %s\n", req.Path)
                    devFound = true
                    break
                }
            }
        }
    }
    if !devFound {
        fmt.Println("  ✅ 无明显的开发依赖作为直接依赖")
    }
    
    fmt.Printf("\n📋 总结: 发现 %d 个潜在安全问题\n", issues)
}

func main() {
    mod, err := pkg.FindAndParseGoModInCurrentDir()
    if err != nil {
        log.Fatalf("解析 go.mod 失败: %v", err)
    }
    
    securityAnalysis(mod)
}
```

## 版本分析

分析版本模式：

```go
package main

import (
    "fmt"
    "log"
    "regexp"
    "strings"
    
    "github.com/scagogogo/go-mod-parser/pkg"
)

func analyzeVersions(mod *module.Module) {
    fmt.Println("📋 版本分析:")
    
    semverPattern := regexp.MustCompile(`^v\d+\.\d+\.\d+`)
    preReleasePattern := regexp.MustCompile(`-\w+`)
    
    semverCount := 0
    preReleaseCount := 0
    pseudoVersionCount := 0
    
    fmt.Println("\n依赖版本:")
    for _, req := range mod.Requires {
        version := req.Version
        fmt.Printf("  %s: %s", req.Path, version)
        
        if semverPattern.MatchString(version) {
            semverCount++
            if preReleasePattern.MatchString(version) {
                fmt.Printf(" (预发布)")
                preReleaseCount++
            } else {
                fmt.Printf(" (稳定)")
            }
        } else if strings.Contains(version, "-") && len(version) > 20 {
            fmt.Printf(" (伪版本)")
            pseudoVersionCount++
        } else {
            fmt.Printf(" (其他)")
        }
        
        if req.Indirect {
            fmt.Printf(" [间接]")
        }
        fmt.Println()
    }
    
    fmt.Printf("\n版本总结:\n")
    fmt.Printf("  语义版本: %d\n", semverCount)
    fmt.Printf("  预发布版本: %d\n", preReleaseCount)
    fmt.Printf("  伪版本: %d\n", pseudoVersionCount)
    fmt.Printf("  其他版本: %d\n", len(mod.Requires)-semverCount-pseudoVersionCount)
}

func main() {
    mod, err := pkg.FindAndParseGoModInCurrentDir()
    if err != nil {
        log.Fatalf("解析 go.mod 失败: %v", err)
    }
    
    analyzeVersions(mod)
}
```

## 综合分析工具

结合所有功能的完整分析工具：

```go
package main

import (
    "fmt"
    "log"
    "strings"
    
    "github.com/scagogogo/go-mod-parser/pkg"
)

func comprehensiveAnalysis(mod *module.Module) {
    fmt.Printf("📦 模块: %s\n", mod.Name)
    fmt.Printf("🐹 Go 版本: %s\n", mod.GoVersion)
    fmt.Println(strings.Repeat("=", 50))
    
    // 基本统计
    analyzeDependencyStats(mod)
    fmt.Println()
    
    // 框架检测
    detectFrameworks(mod)
    fmt.Println()
    
    // 替换分析
    analyzeReplaces(mod)
    fmt.Println()
    
    // 安全分析
    securityAnalysis(mod)
    fmt.Println()
    
    // 版本分析
    analyzeVersions(mod)
}

func main() {
    mod, err := pkg.FindAndParseGoModInCurrentDir()
    if err != nil {
        log.Fatalf("解析 go.mod 失败: %v", err)
    }
    
    comprehensiveAnalysis(mod)
}
```

## 下一步

- [高级用法](/zh/examples/advanced-usage) - 复杂分析模式
- [文件发现](/zh/examples/file-discovery) - 自动发现功能
- [基础解析](/zh/examples/basic-parsing) - 基本解析操作
