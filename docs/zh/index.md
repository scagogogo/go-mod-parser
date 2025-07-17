---
layout: home

hero:
  name: "Go Mod Parser"
  text: "全面的 Go 模块解析器"
  tagline: "轻松解析和分析 go.mod 文件"
  image:
    src: /logo.svg
    alt: Go Mod Parser
  actions:
    - theme: brand
      text: 开始使用
      link: /zh/quick-start
    - theme: alt
      text: API 参考
      link: /zh/api/
    - theme: alt
      text: 查看 GitHub
      link: https://github.com/scagogogo/go-mod-parser

features:
  - icon: 🧩
    title: 完整指令支持
    details: 解析所有 go.mod 指令，包括 module、go、require、replace、exclude 和 retract
  - icon: 🔍
    title: 自动发现
    details: 自动在项目目录和父目录中查找并解析 go.mod 文件
  - icon: 📝
    title: 注释支持
    details: 正确处理 go.mod 文件中的间接依赖注释和其他注解
  - icon: 🔄
    title: 依赖分析
    details: 提供丰富的辅助函数用于分析模块依赖关系和模块间关系
  - icon: 🧪
    title: 测试完善
    details: 全面的单元测试覆盖，确保解析的准确性和可靠性
  - icon: 📚
    title: 示例丰富
    details: 多个实用示例，展示不同使用场景的最佳实践
---

## 快速示例

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
        log.Fatalf("解析失败: %v", err)
    }
    
    // 访问解析数据
    fmt.Printf("模块: %s\n", mod.Name)
    fmt.Printf("Go 版本: %s\n", mod.GoVersion)
    
    // 列出依赖项
    for _, req := range mod.Requires {
        fmt.Printf("- %s %s\n", req.Path, req.Version)
    }
}
```

## 安装

```bash
go get github.com/scagogogo/go-mod-parser
```

## 应用场景

- **依赖分析工具** - 构建分析项目依赖的工具
- **模块版本管理** - 创建管理模块版本的系统
- **CI/CD 流程集成** - 在持续集成中检查依赖
- **构建工具** - 集成到 Go 项目构建系统中
- **依赖可视化** - 创建模块关系的可视化表示
- **更新推荐系统** - 构建建议依赖更新的工具
