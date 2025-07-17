# API 参考

Go Mod Parser 提供了全面的 API 用于解析和分析 go.mod 文件。库分为几个包，每个包都有特定的用途。

## 包概览

### 主包 (`pkg`)

主包提供了大多数用户会使用的主要 API。它提供了解析 go.mod 文件和分析其内容的高级函数。

```go
import "github.com/scagogogo/go-mod-parser/pkg"
```

**核心函数:**
- `ParseGoModFile(path string)` - 从磁盘解析 go.mod 文件
- `ParseGoModContent(content string)` - 从字符串解析 go.mod 内容
- `FindAndParseGoModFile(dir string)` - 自动发现并解析 go.mod 文件
- `HasRequire(mod, path)` - 检查依赖是否存在
- `GetRequire(mod, path)` - 获取依赖详情

### 模块包 (`pkg/module`)

定义用于表示解析的 go.mod 内容的数据结构。

```go
import "github.com/scagogogo/go-mod-parser/pkg/module"
```

**核心类型:**
- `Module` - 表示完整的 go.mod 文件
- `Require` - 表示一个依赖
- `Replace` - 表示一个替换指令
- `Exclude` - 表示一个排除指令
- `Retract` - 表示一个撤回指令

## API 分类

### [核心函数](/zh/api/core-functions)
不同输入源的主要解析函数：
- 基于文件的解析
- 基于字符串的解析
- 自动发现解析

### [数据结构](/zh/api/data-structures)
所有数据类型的完整参考：
- 模块结构
- 依赖类型
- 指令表示

### [辅助函数](/zh/api/helper-functions)
用于分析解析数据的实用函数：
- 依赖检查
- 替换指令分析
- 排除和撤回验证

### [错误处理](/zh/api/error-handling)
错误类型和处理模式：
- 解析错误
- 文件系统错误
- 验证错误

## 快速参考

| 函数 | 描述 | 返回值 |
|------|------|--------|
| `ParseGoModFile(path)` | 从路径解析 go.mod 文件 | `(*Module, error)` |
| `ParseGoModContent(content)` | 从字符串解析 go.mod | `(*Module, error)` |
| `FindAndParseGoModFile(dir)` | 在目录中查找并解析 go.mod | `(*Module, error)` |
| `FindAndParseGoModInCurrentDir()` | 在当前目录查找并解析 go.mod | `(*Module, error)` |
| `HasRequire(mod, path)` | 检查依赖是否存在 | `bool` |
| `GetRequire(mod, path)` | 获取依赖详情 | `*Require` |
| `HasReplace(mod, path)` | 检查替换是否存在 | `bool` |
| `GetReplace(mod, path)` | 获取替换详情 | `*Replace` |
| `HasExclude(mod, path, version)` | 检查版本是否被排除 | `bool` |
| `HasRetract(mod, version)` | 检查版本是否被撤回 | `bool` |

## 使用模式

### 基本解析

```go
// 从文件解析
mod, err := pkg.ParseGoModFile("go.mod")

// 从字符串解析
mod, err := pkg.ParseGoModContent(content)

// 自动发现
mod, err := pkg.FindAndParseGoModInCurrentDir()
```

### 依赖分析

```go
// 检查依赖
if pkg.HasRequire(mod, "github.com/gin-gonic/gin") {
    req := pkg.GetRequire(mod, "github.com/gin-gonic/gin")
    fmt.Printf("版本: %s, 间接: %v\n", req.Version, req.Indirect)
}
```

### 替换指令分析

```go
// 检查替换
if pkg.HasReplace(mod, "github.com/old/pkg") {
    rep := pkg.GetReplace(mod, "github.com/old/pkg")
    fmt.Printf("替换为: %s %s\n", rep.New.Path, rep.New.Version)
}
```

## 错误处理

所有解析函数都返回错误作为第二个返回值。始终检查错误：

```go
mod, err := pkg.ParseGoModFile("go.mod")
if err != nil {
    // 适当处理错误
    log.Fatalf("解析 go.mod 失败: %v", err)
}
```

常见错误场景：
- 文件未找到
- 无效的 go.mod 语法
- 权限错误
- 格式错误的指令

## 线程安全

库设计为读操作线程安全。返回的 `Module` 结构体及其字段可以安全地从多个 goroutine 并发读取。但是，如果需要修改解析的数据，应该实现自己的同步机制。

## 性能考虑

- 解析通常很快，但文件 I/O 可能是瓶颈
- 考虑为频繁访问的文件缓存解析结果
- 使用 `ParseGoModContent` 处理内存中的内容以避免文件系统开销
- 自动发现函数可能遍历多个目录

## 最佳实践

### 错误处理
```go
// 推荐：提供上下文信息
mod, err := pkg.ParseGoModFile(path)
if err != nil {
    return fmt.Errorf("解析 %s 失败: %w", path, err)
}
```

### 依赖检查
```go
// 高效：直接调用 GetRequire 并检查 nil
if req := pkg.GetRequire(mod, path); req != nil {
    // 处理 req
}

// 低效：两次调用
if pkg.HasRequire(mod, path) {
    req := pkg.GetRequire(mod, path)
    // 处理 req
}
```

### 输入验证
```go
func parseGoModSafely(path string) (*module.Module, error) {
    if path == "" {
        return nil, fmt.Errorf("路径不能为空")
    }
    
    return pkg.ParseGoModFile(path)
}
```

## 下一步

- [核心函数](/zh/api/core-functions) - 详细的函数文档
- [数据结构](/zh/api/data-structures) - 完整的类型参考
- [辅助函数](/zh/api/helper-functions) - 分析工具
- [示例](/zh/examples/) - 实际使用示例
