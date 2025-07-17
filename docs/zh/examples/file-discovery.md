# 文件发现

Go Mod Parser 提供强大的自动发现功能，可以自动在项目结构中定位 go.mod 文件。

## 从当前目录自动发现

查找并解析 go.mod 文件的最简单方法：

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/go-mod-parser/pkg"
)

func main() {
    // 从当前目录开始查找并解析 go.mod
    mod, err := pkg.FindAndParseGoModInCurrentDir()
    if err != nil {
        log.Fatalf("查找 go.mod 失败: %v", err)
    }
    
    fmt.Printf("找到模块: %s\n", mod.Name)
}
```

## 从特定目录自动发现

从特定目录开始搜索：

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/go-mod-parser/pkg"
)

func main() {
    // 从特定目录开始查找并解析 go.mod
    projectDir := "/path/to/project/subdirectory"
    mod, err := pkg.FindAndParseGoModFile(projectDir)
    if err != nil {
        log.Fatalf("在 %s 中查找 go.mod 失败: %v", projectDir, err)
    }
    
    fmt.Printf("找到模块: %s\n", mod.Name)
}
```

## 自动发现的工作原理

自动发现过程：

1. 在指定目录（或当前目录）中开始
2. 在当前目录中查找 `go.mod` 文件
3. 如果未找到，移动到父目录
4. 重复直到找到 go.mod 文件或到达根目录
5. 如果未找到 go.mod 文件则返回错误

## 带回退的健壮发现

```go
package main

import (
    "fmt"
    "log"
    "os"
    "path/filepath"
    
    "github.com/scagogogo/go-mod-parser/pkg"
)

func findGoModWithFallbacks() (*module.Module, error) {
    // 首先尝试当前目录
    mod, err := pkg.FindAndParseGoModInCurrentDir()
    if err == nil {
        fmt.Println("在当前目录树中找到 go.mod")
        return mod, nil
    }
    
    // 尝试特定路径作为回退
    fallbackPaths := []string{
        ".",
        "..",
        "../..",
        filepath.Join(os.Getenv("HOME"), "go", "src", "myproject"),
    }
    
    for _, path := range fallbackPaths {
        fmt.Printf("尝试回退路径: %s\n", path)
        mod, err := pkg.FindAndParseGoModFile(path)
        if err == nil {
            fmt.Printf("在以下位置找到 go.mod: %s\n", path)
            return mod, nil
        }
    }
    
    return nil, fmt.Errorf("在任何位置都未找到 go.mod 文件")
}

func main() {
    mod, err := findGoModWithFallbacks()
    if err != nil {
        log.Fatalf("发现失败: %v", err)
    }
    
    fmt.Printf("成功找到模块: %s\n", mod.Name)
}
```

## 处理单体仓库

在单体仓库场景中，你可能有多个 go.mod 文件：

```go
package main

import (
    "fmt"
    "log"
    "os"
    "path/filepath"
    
    "github.com/scagogogo/go-mod-parser/pkg"
)

func findAllGoModFiles(rootDir string) (map[string]*module.Module, error) {
    modules := make(map[string]*module.Module)
    
    err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        
        if info.Name() == "go.mod" {
            mod, parseErr := pkg.ParseGoModFile(path)
            if parseErr != nil {
                fmt.Printf("警告: 解析 %s 失败: %v\n", path, parseErr)
                return nil // 继续遍历
            }
            
            modules[path] = mod
            fmt.Printf("找到模块: %s 位于 %s\n", mod.Name, path)
        }
        
        return nil
    })
    
    return modules, err
}

func main() {
    rootDir := "." // 或指定你的单体仓库根目录
    
    modules, err := findAllGoModFiles(rootDir)
    if err != nil {
        log.Fatalf("遍历目录错误: %v", err)
    }
    
    fmt.Printf("\n找到 %d 个模块:\n", len(modules))
    for path, mod := range modules {
        fmt.Printf("- %s: %s\n", mod.Name, path)
    }
}
```

## 带验证的发现

验证发现的 go.mod 文件：

```go
package main

import (
    "fmt"
    "log"
    "strings"
    
    "github.com/scagogogo/go-mod-parser/pkg"
)

func validateModule(mod *module.Module) []string {
    var issues []string
    
    if mod.Name == "" {
        issues = append(issues, "模块名称为空")
    }
    
    if mod.GoVersion == "" {
        issues = append(issues, "未指定 go 版本")
    }
    
    if !strings.Contains(mod.Name, ".") {
        issues = append(issues, "模块名称应包含域名")
    }
    
    return issues
}

func discoverAndValidate() {
    mod, err := pkg.FindAndParseGoModInCurrentDir()
    if err != nil {
        log.Fatalf("发现失败: %v", err)
    }
    
    fmt.Printf("发现模块: %s\n", mod.Name)
    
    issues := validateModule(mod)
    if len(issues) > 0 {
        fmt.Println("\n验证问题:")
        for _, issue := range issues {
            fmt.Printf("- %s\n", issue)
        }
    } else {
        fmt.Println("✅ 模块验证通过")
    }
}

func main() {
    discoverAndValidate()
}
```

## 性能考虑

在大型目录树中获得更好性能：

```go
package main

import (
    "fmt"
    "log"
    "os"
    "path/filepath"
    "time"
    
    "github.com/scagogogo/go-mod-parser/pkg"
)

func timedDiscovery(startDir string) {
    start := time.Now()
    
    mod, err := pkg.FindAndParseGoModFile(startDir)
    
    duration := time.Since(start)
    
    if err != nil {
        log.Printf("发现在 %v 内失败: %v", duration, err)
        return
    }
    
    fmt.Printf("在 %v 内找到 %s\n", duration, mod.Name)
}

func efficientDiscovery(startDir string, maxDepth int) (*module.Module, error) {
    currentDir := startDir
    depth := 0
    
    for depth < maxDepth {
        goModPath := filepath.Join(currentDir, "go.mod")
        
        if _, err := os.Stat(goModPath); err == nil {
            return pkg.ParseGoModFile(goModPath)
        }
        
        parentDir := filepath.Dir(currentDir)
        if parentDir == currentDir {
            // 到达根目录
            break
        }
        
        currentDir = parentDir
        depth++
    }
    
    return nil, fmt.Errorf("在 %d 层内未找到 go.mod", maxDepth)
}

func main() {
    // 标准发现
    fmt.Println("标准发现:")
    timedDiscovery(".")
    
    // 限制深度发现
    fmt.Println("\n限制深度发现:")
    mod, err := efficientDiscovery(".", 5)
    if err != nil {
        log.Printf("限制发现失败: %v", err)
    } else {
        fmt.Printf("限制搜索找到: %s\n", mod.Name)
    }
}
```

## 下一步

- [依赖分析](/zh/examples/dependency-analysis) - 分析发现的模块
- [高级用法](/zh/examples/advanced-usage) - 复杂发现模式
- [基础解析](/zh/examples/basic-parsing) - 学习基本解析操作
