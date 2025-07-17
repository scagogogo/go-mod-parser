# 核心函数

核心函数提供了从各种来源解析 go.mod 文件的主要接口。

## ParseGoModFile

从文件路径解析 go.mod 文件。

```go
func ParseGoModFile(path string) (*module.Module, error)
```

### 参数

- `path` (string): go.mod 文件的路径

### 返回值

- `*module.Module`: 解析的模块数据
- `error`: 解析失败时的错误

### 示例

```go
mod, err := pkg.ParseGoModFile("/path/to/go.mod")
if err != nil {
    log.Fatalf("解析失败: %v", err)
}

fmt.Printf("模块: %s\n", mod.Name)
fmt.Printf("Go 版本: %s\n", mod.GoVersion)
```

### 错误情况

- 文件不存在
- 权限被拒绝
- 无效的 go.mod 语法
- I/O 错误

---

## ParseGoModContent

从字符串解析 go.mod 内容。

```go
func ParseGoModContent(content string) (*module.Module, error)
```

### 参数

- `content` (string): go.mod 文件内容字符串

### 返回值

- `*module.Module`: 解析的模块数据
- `error`: 解析失败时的错误

### 示例

```go
content := `module github.com/example/project

go 1.21

require (
    github.com/gin-gonic/gin v1.9.1
    github.com/stretchr/testify v1.8.4 // indirect
)

replace github.com/old/pkg => github.com/new/pkg v1.0.0
exclude github.com/bad/pkg v1.0.0
retract v1.0.1 // 安全问题
`

mod, err := pkg.ParseGoModContent(content)
if err != nil {
    log.Fatalf("解析失败: %v", err)
}

fmt.Printf("模块: %s\n", mod.Name)
```

### 错误情况

- 无效的 go.mod 语法
- 格式错误的指令
- 不支持的 go.mod 特性

---

## FindAndParseGoModFile

从指定目录开始查找并解析 go.mod 文件，向上搜索父目录。

```go
func FindAndParseGoModFile(dir string) (*module.Module, error)
```

### 参数

- `dir` (string): 搜索的起始目录

### 返回值

- `*module.Module`: 解析的模块数据
- `error`: 未找到 go.mod 或解析失败时的错误

### 示例

```go
// 从特定目录开始搜索
mod, err := pkg.FindAndParseGoModFile("/path/to/project/subdir")
if err != nil {
    log.Fatalf("查找 go.mod 失败: %v", err)
}

fmt.Printf("找到模块: %s\n", mod.Name)
```

### 行为

1. 在指定目录中开始搜索
2. 在当前目录中查找 `go.mod` 文件
3. 如果未找到，移动到父目录
4. 继续直到找到 go.mod 或到达根目录
5. 如果未找到 go.mod 文件则返回错误

### 错误情况

- 在目录树中未找到 go.mod 文件
- 权限被拒绝
- 无效的 go.mod 语法
- I/O 错误

---

## FindAndParseGoModInCurrentDir

从当前工作目录开始查找并解析 go.mod 文件。

```go
func FindAndParseGoModInCurrentDir() (*module.Module, error)
```

### 参数

无

### 返回值

- `*module.Module`: 解析的模块数据
- `error`: 未找到 go.mod 或解析失败时的错误

### 示例

```go
// 从当前目录开始搜索
mod, err := pkg.FindAndParseGoModInCurrentDir()
if err != nil {
    log.Fatalf("查找 go.mod 失败: %v", err)
}

fmt.Printf("当前项目模块: %s\n", mod.Name)
```

### 行为

这等同于调用 `FindAndParseGoModFile("")` 传入空字符串，使用当前工作目录作为起始点。

### 错误情况

- 在当前目录树中未找到 go.mod 文件
- 无法确定当前工作目录
- 权限被拒绝
- 无效的 go.mod 语法
- I/O 错误

---

## 高级用法

### 处理不同输入源

```go
package main

import (
    "fmt"
    "log"
    "os"
    
    "github.com/scagogogo/go-mod-parser/pkg"
)

func parseFromDifferentSources() {
    // 方法 1: 从文件路径解析
    if mod, err := pkg.ParseGoModFile("go.mod"); err == nil {
        fmt.Printf("从文件: %s\n", mod.Name)
    }
    
    // 方法 2: 从内容解析
    content, _ := os.ReadFile("go.mod")
    if mod, err := pkg.ParseGoModContent(string(content)); err == nil {
        fmt.Printf("从内容: %s\n", mod.Name)
    }
    
    // 方法 3: 从当前目录自动发现
    if mod, err := pkg.FindAndParseGoModInCurrentDir(); err == nil {
        fmt.Printf("自动发现: %s\n", mod.Name)
    }
    
    // 方法 4: 从特定目录自动发现
    if mod, err := pkg.FindAndParseGoModFile("/path/to/project"); err == nil {
        fmt.Printf("在项目中找到: %s\n", mod.Name)
    }
}
```

### 错误处理模式

```go
func robustParsing() {
    // 尝试多种方法并提供回退
    var mod *module.Module
    var err error
    
    // 首先尝试当前目录
    mod, err = pkg.FindAndParseGoModInCurrentDir()
    if err != nil {
        // 回退到特定文件
        mod, err = pkg.ParseGoModFile("go.mod")
        if err != nil {
            log.Fatalf("无法解析 go.mod: %v", err)
        }
    }
    
    fmt.Printf("成功解析: %s\n", mod.Name)
}
```

### 性能考虑

```go
func efficientParsing() {
    // 对于重复解析相同内容，使用 ParseGoModContent
    content := `module example.com/project
go 1.21
require github.com/gin-gonic/gin v1.9.1`
    
    // 这避免了文件 I/O 开销
    mod, err := pkg.ParseGoModContent(content)
    if err != nil {
        log.Fatalf("解析错误: %v", err)
    }
    
    // 对于基于文件的解析，如果需要多次解析同一文件，
    // 考虑缓存结果
}
```

## 相关函数

- [辅助函数](/zh/api/helper-functions) - 用于分析解析数据的函数
- [数据结构](/zh/api/data-structures) - Module 类型的详细信息
- [错误处理](/zh/api/error-handling) - 错误类型和处理策略
