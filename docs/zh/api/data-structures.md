# 数据结构

本页面记录了 Go Mod Parser 用于表示解析的 go.mod 文件内容的所有数据结构。

## Module

`Module` 结构体表示完整的 go.mod 文件，包含所有解析的指令。

```go
type Module struct {
    Name      string      // 来自 module 指令的模块名称
    GoVersion string      // 来自 go 指令的 Go 版本
    Requires  []*Require  // 来自 require 指令的依赖
    Replaces  []*Replace  // Replace 指令
    Excludes  []*Exclude  // Exclude 指令
    Retracts  []*Retract  // Retract 指令
}
```

### 字段

#### Name (string)
在 `module` 指令中声明的模块名称。

**示例:**
```go
// 对于: module github.com/example/project
fmt.Println(mod.Name) // 输出: github.com/example/project
```

#### GoVersion (string)
来自 `go` 指令的 Go 版本要求。

**示例:**
```go
// 对于: go 1.21
fmt.Println(mod.GoVersion) // 输出: 1.21
```

#### Requires ([]*Require)
在 `require` 指令中声明的所有依赖的切片。

#### Replaces ([]*Replace)
来自 `replace` 指令的所有模块替换的切片。

#### Excludes ([]*Exclude)
来自 `exclude` 指令的所有排除模块版本的切片。

#### Retracts ([]*Retract)
来自 `retract` 指令的所有撤回版本的切片。

---

## Require

表示单个依赖声明。

```go
type Require struct {
    Path     string // 模块路径
    Version  string // 模块版本
    Indirect bool   // 是否为间接依赖
}
```

### 字段

#### Path (string)
依赖的模块路径。

**示例:**
```go
// 对于: require github.com/gin-gonic/gin v1.9.1
fmt.Println(req.Path) // 输出: github.com/gin-gonic/gin
```

#### Version (string)
依赖的版本约束。

**示例:**
```go
// 对于: require github.com/gin-gonic/gin v1.9.1
fmt.Println(req.Version) // 输出: v1.9.1
```

#### Indirect (bool)
指示这是否为间接依赖（用 `// indirect` 注释标记）。

**示例:**
```go
// 对于: require github.com/example/pkg v1.0.0 // indirect
fmt.Println(req.Indirect) // 输出: true

// 对于: require github.com/example/pkg v1.0.0
fmt.Println(req.Indirect) // 输出: false
```

### 使用示例

```go
for _, req := range mod.Requires {
    status := "直接"
    if req.Indirect {
        status = "间接"
    }
    fmt.Printf("%s %s (%s)\n", req.Path, req.Version, status)
}
```

---

## Replace

表示模块替换指令。

```go
type Replace struct {
    Old *ReplaceItem // 被替换的原始模块
    New *ReplaceItem // 替换模块
}
```

### 字段

#### Old (*ReplaceItem)
被替换的原始模块。

#### New (*ReplaceItem)
替换模块。

### 使用示例

```go
for _, rep := range mod.Replaces {
    fmt.Printf("替换 %s", rep.Old.Path)
    if rep.Old.Version != "" {
        fmt.Printf(" %s", rep.Old.Version)
    }
    fmt.Printf(" => %s", rep.New.Path)
    if rep.New.Version != "" {
        fmt.Printf(" %s", rep.New.Version)
    }
    fmt.Println()
}
```

---

## ReplaceItem

表示替换指令中的模块引用。

```go
type ReplaceItem struct {
    Path    string // 模块路径
    Version string // 模块版本（可能为空）
}
```

### 字段

#### Path (string)
模块路径。

#### Version (string)
模块版本。对于本地路径替换可能为空。

### 示例

```go
// 对于: replace github.com/old/pkg => github.com/new/pkg v1.0.0
// Old: Path="github.com/old/pkg", Version=""
// New: Path="github.com/new/pkg", Version="v1.0.0"

// 对于: replace github.com/old/pkg v1.0.0 => ./local/path
// Old: Path="github.com/old/pkg", Version="v1.0.0"  
// New: Path="./local/path", Version=""
```

---

## Exclude

表示排除的模块版本。

```go
type Exclude struct {
    Path    string // 模块路径
    Version string // 排除的版本
}
```

### 字段

#### Path (string)
要排除的模块路径。

#### Version (string)
要排除的特定版本。

### 使用示例

```go
for _, exc := range mod.Excludes {
    fmt.Printf("排除: %s %s\n", exc.Path, exc.Version)
}

// 检查特定版本是否被排除
for _, exc := range mod.Excludes {
    if exc.Path == "github.com/problematic/pkg" && exc.Version == "v1.0.0" {
        fmt.Println("版本 v1.0.0 被排除")
    }
}
```

---

## Retract

表示撤回的版本或版本范围。

```go
type Retract struct {
    Version     string // 单个撤回版本
    VersionLow  string // 版本范围的下限
    VersionHigh string // 版本范围的上限
    Rationale   string // 撤回原因
}
```

### 字段

#### Version (string)
对于单个版本撤回，包含版本。对于版本范围为空。

#### VersionLow (string)
对于版本范围撤回，这是下限。对于单个版本为空。

#### VersionHigh (string)
对于版本范围撤回，这是上限。对于单个版本为空。

#### Rationale (string)
从注释中提取的撤回原因。可能为空。

### 使用示例

```go
for _, ret := range mod.Retracts {
    if ret.Version != "" {
        // 单个版本撤回
        fmt.Printf("撤回: %s", ret.Version)
    } else {
        // 版本范围撤回
        fmt.Printf("撤回: [%s, %s]", ret.VersionLow, ret.VersionHigh)
    }
    
    if ret.Rationale != "" {
        fmt.Printf(" (%s)", ret.Rationale)
    }
    fmt.Println()
}
```

### 撤回示例

```go
// 单个版本: retract v1.0.1 // 安全问题
// Version="v1.0.1", VersionLow="", VersionHigh="", Rationale="安全问题"

// 版本范围: retract [v1.0.0, v1.0.5] // 构建损坏
// Version="", VersionLow="v1.0.0", VersionHigh="v1.0.5", Rationale="构建损坏"
```

---

## 完整示例

这是一个展示如何使用所有数据结构的完整示例：

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

replace (
    github.com/old/pkg => github.com/new/pkg v1.0.0
    github.com/local/pkg => ./local/path
)

exclude (
    github.com/bad/pkg v1.0.0
    github.com/another/bad v2.0.0
)

retract (
    v1.0.1 // 安全漏洞
    [v1.0.2, v1.0.5] // 构建损坏
)
`

    mod, err := pkg.ParseGoModContent(content)
    if err != nil {
        log.Fatalf("解析错误: %v", err)
    }
    
    // 访问基本信息
    fmt.Printf("模块: %s\n", mod.Name)
    fmt.Printf("Go 版本: %s\n", mod.GoVersion)
    
    // 处理依赖
    fmt.Printf("\n依赖 (%d):\n", len(mod.Requires))
    for _, req := range mod.Requires {
        fmt.Printf("  %s %s", req.Path, req.Version)
        if req.Indirect {
            fmt.Printf(" (间接)")
        }
        fmt.Println()
    }
    
    // 处理替换
    fmt.Printf("\n替换 (%d):\n", len(mod.Replaces))
    for _, rep := range mod.Replaces {
        fmt.Printf("  %s => %s", rep.Old.Path, rep.New.Path)
        if rep.New.Version != "" {
            fmt.Printf(" %s", rep.New.Version)
        }
        fmt.Println()
    }
    
    // 处理排除
    fmt.Printf("\n排除 (%d):\n", len(mod.Excludes))
    for _, exc := range mod.Excludes {
        fmt.Printf("  %s %s\n", exc.Path, exc.Version)
    }
    
    // 处理撤回
    fmt.Printf("\n撤回 (%d):\n", len(mod.Retracts))
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
```

## 相关文档

- [核心函数](/zh/api/core-functions) - 返回这些数据结构的函数
- [辅助函数](/zh/api/helper-functions) - 用于分析这些结构的函数
- [示例](/zh/examples/) - 实际使用示例
