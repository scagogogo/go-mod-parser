# 错误处理

Go Mod Parser 提供全面的错误处理，帮助你诊断和处理解析 go.mod 文件时的各种失败场景。

## 错误类型

### 文件系统错误

这些错误在访问文件或目录时出现问题时发生。

#### 文件未找到

```go
mod, err := pkg.ParseGoModFile("/nonexistent/go.mod")
if err != nil {
    // 错误: open /nonexistent/go.mod: no such file or directory
    fmt.Printf("错误: %v\n", err)
}
```

#### 权限被拒绝

```go
mod, err := pkg.ParseGoModFile("/root/go.mod")
if err != nil {
    // 错误: open /root/go.mod: permission denied
    fmt.Printf("错误: %v\n", err)
}
```

#### 目录未找到（自动发现）

```go
mod, err := pkg.FindAndParseGoModFile("/nonexistent/directory")
if err != nil {
    // 错误: go.mod file not found
    fmt.Printf("错误: %v\n", err)
}
```

### 解析错误

这些错误在 go.mod 文件内容格式错误或包含无效语法时发生。

#### 无效模块声明

```go
content := `invalid module declaration`
mod, err := pkg.ParseGoModContent(content)
if err != nil {
    // 错误: line 1: unrecognized line format: invalid module declaration
    fmt.Printf("解析错误: %v\n", err)
}
```

#### 无效 require 格式

```go
content := `module example.com/test

require invalid-require-format
`
mod, err := pkg.ParseGoModContent(content)
if err != nil {
    // 错误: line 3: invalid require declaration
    fmt.Printf("解析错误: %v\n", err)
}
```

#### 无效 replace 格式

```go
content := `module example.com/test

replace invalid => format
`
mod, err := pkg.ParseGoModContent(content)
if err != nil {
    // 错误: line 3: invalid replace declaration
    fmt.Printf("解析错误: %v\n", err)
}
```

### 块解析错误

在块语句（require、replace、exclude、retract 块）内发生的错误。

```go
content := `module example.com/test

require (
    github.com/example/pkg
)
`
mod, err := pkg.ParseGoModContent(content)
if err != nil {
    // 错误: line 4: invalid require declaration
    fmt.Printf("块解析错误: %v\n", err)
}
```

## 错误处理模式

### 基本错误检查

调用解析函数时始终检查错误：

```go
mod, err := pkg.ParseGoModFile("go.mod")
if err != nil {
    log.Fatalf("解析 go.mod 失败: %v", err)
}
// 继续使用 mod...
```

### 优雅错误处理

```go
func parseGoModSafely(path string) (*module.Module, error) {
    mod, err := pkg.ParseGoModFile(path)
    if err != nil {
        // 记录错误但不崩溃
        log.Printf("警告: 解析 %s 失败: %v", path, err)
        return nil, err
    }
    return mod, nil
}
```

### 错误类型检测

```go
func handleParseError(err error) {
    errStr := err.Error()
    
    switch {
    case strings.Contains(errStr, "no such file"):
        fmt.Println("文件未找到 - 检查路径")
    case strings.Contains(errStr, "permission denied"):
        fmt.Println("权限被拒绝 - 检查文件权限")
    case strings.Contains(errStr, "go.mod file not found"):
        fmt.Println("在目录树中未找到 go.mod")
    case strings.Contains(errStr, "line"):
        fmt.Println("解析错误 - 检查 go.mod 语法")
    default:
        fmt.Printf("未知错误: %v\n", err)
    }
}
```

### 回退策略

```go
func parseWithFallbacks(primaryPath, fallbackPath string) (*module.Module, error) {
    // 首先尝试主路径
    mod, err := pkg.ParseGoModFile(primaryPath)
    if err == nil {
        return mod, nil
    }
    
    log.Printf("主解析失败: %v", err)
    
    // 尝试回退路径
    mod, err = pkg.ParseGoModFile(fallbackPath)
    if err == nil {
        log.Printf("成功解析回退: %s", fallbackPath)
        return mod, nil
    }
    
    // 最后尝试自动发现
    mod, err = pkg.FindAndParseGoModInCurrentDir()
    if err == nil {
        log.Println("成功自动发现 go.mod")
        return mod, nil
    }
    
    return nil, fmt.Errorf("所有解析尝试都失败了")
}
```

### 解析后验证

```go
func validateParsedModule(mod *module.Module) error {
    if mod.Name == "" {
        return fmt.Errorf("模块名称为空")
    }
    
    if mod.GoVersion == "" {
        return fmt.Errorf("未指定 go 版本")
    }
    
    // 检查重复依赖
    seen := make(map[string]bool)
    for _, req := range mod.Requires {
        if seen[req.Path] {
            return fmt.Errorf("重复依赖: %s", req.Path)
        }
        seen[req.Path] = true
    }
    
    return nil
}

func parseAndValidate(path string) (*module.Module, error) {
    mod, err := pkg.ParseGoModFile(path)
    if err != nil {
        return nil, fmt.Errorf("解析错误: %w", err)
    }
    
    if err := validateParsedModule(mod); err != nil {
        return nil, fmt.Errorf("验证错误: %w", err)
    }
    
    return mod, nil
}
```

## 健壮的解析函数

这是一个处理多种错误场景的综合示例：

```go
package main

import (
    "fmt"
    "log"
    "os"
    "path/filepath"
    "strings"
    
    "github.com/scagogogo/go-mod-parser/pkg"
    "github.com/scagogogo/go-mod-parser/pkg/module"
)

func robustParseGoMod(path string) (*module.Module, error) {
    // 规范化路径
    absPath, err := filepath.Abs(path)
    if err != nil {
        return nil, fmt.Errorf("解析路径 %s 失败: %w", path, err)
    }
    
    // 检查文件是否存在
    if _, err := os.Stat(absPath); os.IsNotExist(err) {
        return nil, fmt.Errorf("go.mod 文件未找到: %s", absPath)
    } else if err != nil {
        return nil, fmt.Errorf("访问文件 %s 失败: %w", absPath, err)
    }
    
    // 尝试解析
    mod, err := pkg.ParseGoModFile(absPath)
    if err != nil {
        return nil, fmt.Errorf("解析 %s 失败: %w", absPath, err)
    }
    
    // 基本验证
    if mod.Name == "" {
        return nil, fmt.Errorf("解析的模块名称为空")
    }
    
    return mod, nil
}

func parseWithRetry(paths []string) (*module.Module, error) {
    var lastErr error
    
    for _, path := range paths {
        mod, err := robustParseGoMod(path)
        if err == nil {
            return mod, nil
        }
        
        log.Printf("解析 %s 失败: %v", path, err)
        lastErr = err
    }
    
    return nil, fmt.Errorf("所有解析尝试都失败了，最后错误: %w", lastErr)
}

func main() {
    // 尝试多个可能的位置
    candidates := []string{
        "go.mod",
        "./go.mod", 
        "../go.mod",
        "../../go.mod",
    }
    
    mod, err := parseWithRetry(candidates)
    if err != nil {
        // 最终回退：自动发现
        mod, err = pkg.FindAndParseGoModInCurrentDir()
        if err != nil {
            log.Fatalf("无法找到或解析任何 go.mod 文件: %v", err)
        }
        log.Println("成功使用自动发现")
    }
    
    fmt.Printf("成功解析模块: %s\n", mod.Name)
}
```

## 错误预防

### 输入验证

```go
func validateInputs(path string) error {
    if path == "" {
        return fmt.Errorf("路径不能为空")
    }
    
    if !strings.HasSuffix(path, "go.mod") && !isDirectory(path) {
        return fmt.Errorf("路径必须是 go.mod 文件或目录")
    }
    
    return nil
}

func isDirectory(path string) bool {
    info, err := os.Stat(path)
    return err == nil && info.IsDir()
}
```

### 内容验证

```go
func validateGoModContent(content string) error {
    if content == "" {
        return fmt.Errorf("内容不能为空")
    }
    
    if !strings.Contains(content, "module ") {
        return fmt.Errorf("内容似乎不是有效的 go.mod 文件")
    }
    
    return nil
}

func parseContentSafely(content string) (*module.Module, error) {
    if err := validateGoModContent(content); err != nil {
        return nil, err
    }
    
    return pkg.ParseGoModContent(content)
}
```

## 最佳实践

1. **始终检查错误**: 永远不要忽略解析函数的错误
2. **提供上下文**: 使用 `fmt.Errorf` 包装错误并添加额外上下文
3. **适当记录**: 为不同错误类型使用适当的日志级别
4. **实现回退**: 为关键解析操作准备备用策略
5. **验证结果**: 检查解析数据的一致性和完整性
6. **处理边缘情况**: 考虑空文件、权限问题和网络问题

## 相关文档

- [核心函数](/zh/api/core-functions) - 可能返回错误的函数
- [数据结构](/zh/api/data-structures) - 理解解析的数据
- [示例](/zh/examples/) - 查看实际的错误处理
