# 辅助函数

辅助函数提供了分析和查询解析的 go.mod 数据的便捷方法。这些函数使检查特定依赖、替换、排除和撤回变得容易。

## 依赖函数

### HasRequire

检查模块是否有特定的依赖。

```go
func HasRequire(mod *module.Module, path string) bool
```

#### 参数

- `mod` (*module.Module): 解析的模块
- `path` (string): 要检查的依赖路径

#### 返回值

- `bool`: 如果依赖存在返回 true，否则返回 false

#### 示例

```go
if pkg.HasRequire(mod, "github.com/gin-gonic/gin") {
    fmt.Println("项目使用 Gin 框架")
} else {
    fmt.Println("未找到 Gin 框架")
}
```

---

### GetRequire

获取特定依赖的详细信息。

```go
func GetRequire(mod *module.Module, path string) *module.Require
```

#### 参数

- `mod` (*module.Module): 解析的模块
- `path` (string): 要获取的依赖路径

#### 返回值

- `*module.Require`: 依赖详情，如果未找到则为 nil

#### 示例

```go
req := pkg.GetRequire(mod, "github.com/gin-gonic/gin")
if req != nil {
    fmt.Printf("版本: %s\n", req.Version)
    fmt.Printf("间接: %v\n", req.Indirect)
} else {
    fmt.Println("未找到依赖")
}
```

#### 安全使用模式

```go
// 总是先检查依赖是否存在
if pkg.HasRequire(mod, "github.com/gin-gonic/gin") {
    req := pkg.GetRequire(mod, "github.com/gin-gonic/gin")
    // 这里 req 保证不为 nil
    fmt.Printf("找到 Gin 版本: %s\n", req.Version)
}

// 或处理 nil 情况
if req := pkg.GetRequire(mod, "github.com/gin-gonic/gin"); req != nil {
    fmt.Printf("找到 Gin 版本: %s\n", req.Version)
}
```

---

## 替换函数

### HasReplace

检查模块是否有特定路径的替换。

```go
func HasReplace(mod *module.Module, path string) bool
```

#### 参数

- `mod` (*module.Module): 解析的模块
- `path` (string): 要检查的原始模块路径

#### 返回值

- `bool`: 如果替换存在返回 true，否则返回 false

#### 示例

```go
if pkg.HasReplace(mod, "github.com/old/package") {
    fmt.Println("包已被替换")
}
```

---

### GetReplace

获取特定模块路径的替换详情。

```go
func GetReplace(mod *module.Module, path string) *module.Replace
```

#### 参数

- `mod` (*module.Module): 解析的模块
- `path` (string): 原始模块路径

#### 返回值

- `*module.Replace`: 替换详情，如果未找到则为 nil

#### 示例

```go
if rep := pkg.GetReplace(mod, "github.com/old/package"); rep != nil {
    fmt.Printf("替换为: %s", rep.New.Path)
    if rep.New.Version != "" {
        fmt.Printf(" %s", rep.New.Version)
    }
    fmt.Println()
}
```

#### 分析替换类型

```go
rep := pkg.GetReplace(mod, "github.com/old/package")
if rep != nil {
    if strings.HasPrefix(rep.New.Path, "./") || strings.HasPrefix(rep.New.Path, "../") {
        fmt.Println("本地路径替换")
    } else {
        fmt.Println("模块替换")
        if rep.New.Version != "" {
            fmt.Printf("替换版本: %s\n", rep.New.Version)
        }
    }
}
```

---

## 排除函数

### HasExclude

检查特定模块版本是否被排除。

```go
func HasExclude(mod *module.Module, path, version string) bool
```

#### 参数

- `mod` (*module.Module): 解析的模块
- `path` (string): 模块路径
- `version` (string): 要检查的版本

#### 返回值

- `bool`: 如果版本被排除返回 true，否则返回 false

#### 示例

```go
if pkg.HasExclude(mod, "github.com/problematic/pkg", "v1.0.0") {
    fmt.Println("版本 v1.0.0 被排除")
}
```

#### 检查多个版本

```go
problematicVersions := []string{"v1.0.0", "v1.0.1", "v1.0.2"}
for _, version := range problematicVersions {
    if pkg.HasExclude(mod, "github.com/problematic/pkg", version) {
        fmt.Printf("版本 %s 被排除\n", version)
    }
}
```

---

## 撤回函数

### HasRetract

检查特定版本是否已被撤回。

```go
func HasRetract(mod *module.Module, version string) bool
```

#### 参数

- `mod` (*module.Module): 解析的模块
- `version` (string): 要检查的版本

#### 返回值

- `bool`: 如果版本被撤回返回 true，否则返回 false

#### 示例

```go
if pkg.HasRetract(mod, "v1.0.1") {
    fmt.Println("版本 v1.0.1 已被撤回")
}
```

#### 版本范围说明

此函数检查单个版本撤回和版本范围。如果版本在撤回范围内，将返回 true。

```go
// 对于 retract [v1.0.0, v1.0.5]
fmt.Println(pkg.HasRetract(mod, "v1.0.2")) // true
fmt.Println(pkg.HasRetract(mod, "v1.0.6")) // false
```

---

## 高级使用模式

### 综合依赖分析

```go
func analyzeDependencies(mod *module.Module) {
    fmt.Printf("分析 %d 个依赖:\n", len(mod.Requires))
    
    for _, req := range mod.Requires {
        fmt.Printf("\n📦 %s %s", req.Path, req.Version)
        
        // 检查是否为间接依赖
        if req.Indirect {
            fmt.Printf(" (间接)")
        }
        
        // 检查是否有替换
        if pkg.HasReplace(mod, req.Path) {
            rep := pkg.GetReplace(mod, req.Path)
            fmt.Printf("\n   🔄 替换为: %s", rep.New.Path)
            if rep.New.Version != "" {
                fmt.Printf(" %s", rep.New.Version)
            }
        }
        
        // 检查是否有版本被排除
        excluded := false
        for _, exc := range mod.Excludes {
            if exc.Path == req.Path {
                fmt.Printf("\n   🚫 版本 %s 被排除", exc.Version)
                excluded = true
            }
        }
        
        fmt.Println()
    }
}
```

### 验证函数

```go
func validateModule(mod *module.Module) []string {
    var issues []string
    
    // 检查同时被依赖和替换的包
    for _, req := range mod.Requires {
        if pkg.HasReplace(mod, req.Path) {
            rep := pkg.GetReplace(mod, req.Path)
            issues = append(issues, 
                fmt.Sprintf("依赖 %s 被替换为 %s", 
                    req.Path, rep.New.Path))
        }
    }
    
    // 检查使用中的撤回版本
    for _, req := range mod.Requires {
        if pkg.HasRetract(mod, req.Version) {
            issues = append(issues, 
                fmt.Sprintf("使用撤回版本 %s 的 %s", 
                    req.Version, req.Path))
        }
    }
    
    return issues
}
```

### 依赖过滤

```go
func filterDependencies(mod *module.Module, filter func(*module.Require) bool) []*module.Require {
    var filtered []*module.Require
    for _, req := range mod.Requires {
        if filter(req) {
            filtered = append(filtered, req)
        }
    }
    return filtered
}

// 使用示例
directDeps := filterDependencies(mod, func(req *module.Require) bool {
    return !req.Indirect
})

testDeps := filterDependencies(mod, func(req *module.Require) bool {
    return strings.Contains(req.Path, "test")
})

replacedDeps := filterDependencies(mod, func(req *module.Require) bool {
    return pkg.HasReplace(mod, req.Path)
})
```

### 安全分析

```go
func checkSecurity(mod *module.Module) {
    fmt.Println("安全分析:")
    
    // 检查撤回版本
    retractedCount := 0
    for _, req := range mod.Requires {
        if pkg.HasRetract(mod, req.Version) {
            fmt.Printf("⚠️  使用撤回版本: %s %s\n", 
                req.Path, req.Version)
            retractedCount++
        }
    }
    
    // 检查本地替换（潜在安全风险）
    localReplacements := 0
    for _, rep := range mod.Replaces {
        if strings.HasPrefix(rep.New.Path, "./") || 
           strings.HasPrefix(rep.New.Path, "../") {
            fmt.Printf("🔍 本地替换: %s => %s\n", 
                rep.Old.Path, rep.New.Path)
            localReplacements++
        }
    }
    
    fmt.Printf("\n总结: %d 个撤回版本，%d 个本地替换\n", 
        retractedCount, localReplacements)
}
```

## 性能提示

1. **缓存结果**: 如果多次使用相同参数调用这些函数，考虑缓存结果。

2. **批量检查**: 不要先调用 `HasRequire` 再调用 `GetRequire`，直接调用 `GetRequire` 并检查 nil。

3. **提前返回**: 在分析函数中使用提前返回避免不必要的处理。

```go
// 高效模式
if req := pkg.GetRequire(mod, path); req != nil {
    // 处理 req
    return req.Version
}

// 低效模式
if pkg.HasRequire(mod, path) {
    req := pkg.GetRequire(mod, path)
    return req.Version
}
```

## 相关文档

- [核心函数](/zh/api/core-functions) - 用于解析 go.mod 文件的函数
- [数据结构](/zh/api/data-structures) - 数据类型的详细信息
- [示例](/zh/examples/) - 实际使用示例
