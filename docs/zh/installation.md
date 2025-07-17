# 安装

## 系统要求

- Go 1.21 或更高版本
- Git（用于获取模块）

## 通过 Go 模块安装

最简单的安装方式是使用 Go 模块：

```bash
go get github.com/scagogogo/go-mod-parser
```

## 验证安装

创建一个简单的测试文件来验证安装：

```go
// test.go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/go-mod-parser/pkg"
)

func main() {
    // 测试解析简单的 go.mod 内容
    content := `module github.com/example/test

go 1.21

require github.com/stretchr/testify v1.8.4
`
    
    mod, err := pkg.ParseGoModContent(content)
    if err != nil {
        log.Fatalf("解析失败: %v", err)
    }
    
    fmt.Printf("成功解析模块: %s\n", mod.Name)
    fmt.Printf("Go 版本: %s\n", mod.GoVersion)
    fmt.Printf("依赖数量: %d\n", len(mod.Requires))
}
```

运行测试：

```bash
go run test.go
```

预期输出：
```
成功解析模块: github.com/example/test
Go 版本: 1.21
依赖数量: 1
```

## 在项目中导入

在你的 Go 文件中添加导入：

```go
import "github.com/scagogogo/go-mod-parser/pkg"
```

导入特定子包：

```go
import (
    "github.com/scagogogo/go-mod-parser/pkg"
    "github.com/scagogogo/go-mod-parser/pkg/module"
)
```

## 开发环境设置

如果你想要贡献代码或修改库：

1. 克隆仓库：
```bash
git clone https://github.com/scagogogo/go-mod-parser.git
cd go-mod-parser
```

2. 安装依赖：
```bash
go mod download
```

3. 运行测试：
```bash
go test -v ./...
```

4. 运行示例：
```bash
cd examples/01_basic_parsing
go run main.go ../../go.mod
```

## 故障排除

### 模块未找到

如果遇到"模块未找到"错误：

1. 确保使用 Go 1.21 或更高版本：
```bash
go version
```

2. 如果还没有初始化模块：
```bash
go mod init your-project-name
```

3. 尝试显式获取模块：
```bash
go get -u github.com/scagogogo/go-mod-parser
```

### 导入错误

如果遇到导入错误：

1. 检查你的 Go 模块是否正确初始化
2. 确保导入路径正确
3. 运行 `go mod tidy` 清理依赖

### 版本冲突

如果有版本冲突：

1. 检查 go.mod 文件中的冲突版本
2. 使用 `go mod why` 了解依赖链
3. 考虑使用 `replace` 指令解决冲突

### 网络问题

如果在中国大陆遇到网络问题：

1. 配置 Go 模块代理：
```bash
go env -w GOPROXY=https://goproxy.cn,direct
```

2. 或使用其他代理：
```bash
go env -w GOPROXY=https://proxy.golang.org,direct
```

3. 如果需要，配置私有模块：
```bash
go env -w GOPRIVATE=your-private-domain.com
```

### 权限问题

如果遇到权限问题：

1. 确保有写入 `$GOPATH` 或模块缓存的权限
2. 在 Linux/macOS 上，可能需要调整目录权限：
```bash
chmod -R 755 $GOPATH
```

### 代理配置

如果在企业环境中需要代理：

```bash
# 设置 HTTP 代理
export HTTP_PROXY=http://proxy.company.com:8080
export HTTPS_PROXY=http://proxy.company.com:8080

# 设置不使用代理的域名
export NO_PROXY=localhost,127.0.0.1,company.com
```

## 验证完整安装

运行这个更全面的测试来验证所有功能：

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/go-mod-parser/pkg"
)

func main() {
    content := `module github.com/example/comprehensive

go 1.21

require (
    github.com/gin-gonic/gin v1.9.1
    github.com/stretchr/testify v1.8.4 // indirect
)

replace github.com/old/pkg => github.com/new/pkg v1.0.0

exclude github.com/bad/pkg v1.0.0

retract v1.0.1 // 安全漏洞
`
    
    mod, err := pkg.ParseGoModContent(content)
    if err != nil {
        log.Fatalf("解析失败: %v", err)
    }
    
    fmt.Printf("✅ 模块解析成功: %s\n", mod.Name)
    fmt.Printf("✅ Go 版本: %s\n", mod.GoVersion)
    fmt.Printf("✅ 依赖项: %d\n", len(mod.Requires))
    fmt.Printf("✅ 替换规则: %d\n", len(mod.Replaces))
    fmt.Printf("✅ 排除规则: %d\n", len(mod.Excludes))
    fmt.Printf("✅ 撤回版本: %d\n", len(mod.Retracts))
    
    // 测试辅助函数
    if pkg.HasRequire(mod, "github.com/gin-gonic/gin") {
        fmt.Println("✅ 依赖检查功能正常")
    }
    
    if pkg.HasReplace(mod, "github.com/old/pkg") {
        fmt.Println("✅ 替换检查功能正常")
    }
    
    if pkg.HasExclude(mod, "github.com/bad/pkg", "v1.0.0") {
        fmt.Println("✅ 排除检查功能正常")
    }
    
    if pkg.HasRetract(mod, "v1.0.1") {
        fmt.Println("✅ 撤回检查功能正常")
    }
    
    fmt.Println("\n🎉 所有功能验证通过！")
}
```

如果所有测试都通过，说明 Go Mod Parser 已经正确安装并可以使用了。
