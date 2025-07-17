# 高级用法

本节涵盖 Go Mod Parser 的高级模式和复杂用例。

## 自定义分析框架

构建灵活的分析框架：

```go
package main

import (
    "fmt"
    "log"
    "strings"
    
    "github.com/scagogogo/go-mod-parser/pkg"
    "github.com/scagogogo/go-mod-parser/pkg/module"
)

// 可扩展分析的分析器接口
type Analyzer interface {
    Name() string
    Analyze(mod *module.Module) AnalysisResult
}

type AnalysisResult struct {
    Summary string
    Details []string
    Issues  []string
}

// 安全分析器
type SecurityAnalyzer struct{}

func (s SecurityAnalyzer) Name() string {
    return "安全分析"
}

func (s SecurityAnalyzer) Analyze(mod *module.Module) AnalysisResult {
    result := AnalysisResult{}
    
    // 检查撤回版本
    for _, req := range mod.Requires {
        if pkg.HasRetract(mod, req.Version) {
            result.Issues = append(result.Issues, 
                fmt.Sprintf("使用撤回版本: %s %s", req.Path, req.Version))
        }
    }
    
    // 检查本地替换
    for _, rep := range mod.Replaces {
        if strings.HasPrefix(rep.New.Path, "./") || strings.HasPrefix(rep.New.Path, "../") {
            result.Issues = append(result.Issues, 
                fmt.Sprintf("本地替换: %s => %s", rep.Old.Path, rep.New.Path))
        }
    }
    
    if len(result.Issues) == 0 {
        result.Summary = "未发现安全问题"
    } else {
        result.Summary = fmt.Sprintf("发现 %d 个安全问题", len(result.Issues))
    }
    
    return result
}

// 性能分析器
type PerformanceAnalyzer struct{}

func (p PerformanceAnalyzer) Name() string {
    return "性能分析"
}

func (p PerformanceAnalyzer) Analyze(mod *module.Module) AnalysisResult {
    result := AnalysisResult{}
    
    heavyDeps := []string{
        "github.com/docker/docker",
        "k8s.io/kubernetes",
        "github.com/aws/aws-sdk-go",
    }
    
    for _, req := range mod.Requires {
        for _, heavy := range heavyDeps {
            if strings.Contains(req.Path, heavy) {
                result.Details = append(result.Details, 
                    fmt.Sprintf("检测到重型依赖: %s", req.Path))
            }
        }
    }
    
    result.Summary = fmt.Sprintf("分析了 %d 个依赖的性能影响", len(mod.Requires))
    return result
}

// 分析运行器
func runAnalysis(mod *module.Module, analyzers []Analyzer) {
    fmt.Printf("为模块运行分析: %s\n", mod.Name)
    fmt.Println(strings.Repeat("=", 60))
    
    for _, analyzer := range analyzers {
        fmt.Printf("\n🔍 %s\n", analyzer.Name())
        result := analyzer.Analyze(mod)
        
        fmt.Printf("总结: %s\n", result.Summary)
        
        if len(result.Details) > 0 {
            fmt.Println("详情:")
            for _, detail := range result.Details {
                fmt.Printf("  • %s\n", detail)
            }
        }
        
        if len(result.Issues) > 0 {
            fmt.Println("问题:")
            for _, issue := range result.Issues {
                fmt.Printf("  ⚠️  %s\n", issue)
            }
        }
    }
}

func main() {
    mod, err := pkg.FindAndParseGoModInCurrentDir()
    if err != nil {
        log.Fatalf("解析 go.mod 失败: %v", err)
    }
    
    analyzers := []Analyzer{
        SecurityAnalyzer{},
        PerformanceAnalyzer{},
    }
    
    runAnalysis(mod, analyzers)
}
```

## 模块比较工具

比较多个 go.mod 文件：

```go
package main

import (
    "fmt"
    "log"
    "os"
    "strings"
    
    "github.com/scagogogo/go-mod-parser/pkg"
    "github.com/scagogogo/go-mod-parser/pkg/module"
)

type ModuleComparison struct {
    Module1 *module.Module
    Module2 *module.Module
    Path1   string
    Path2   string
}

func (mc *ModuleComparison) CompareDependencies() {
    deps1 := make(map[string]string)
    deps2 := make(map[string]string)
    
    for _, req := range mc.Module1.Requires {
        deps1[req.Path] = req.Version
    }
    
    for _, req := range mc.Module2.Requires {
        deps2[req.Path] = req.Version
    }
    
    fmt.Printf("📋 依赖比较\n")
    fmt.Printf("模块 1: %s (%s)\n", mc.Module1.Name, mc.Path1)
    fmt.Printf("模块 2: %s (%s)\n", mc.Module2.Name, mc.Path2)
    fmt.Println(strings.Repeat("-", 60))
    
    // 共同依赖
    fmt.Println("\n🤝 共同依赖:")
    commonCount := 0
    for path, version1 := range deps1 {
        if version2, exists := deps2[path]; exists {
            commonCount++
            if version1 == version2 {
                fmt.Printf("  ✅ %s: %s (相同)\n", path, version1)
            } else {
                fmt.Printf("  ⚠️  %s: %s vs %s\n", path, version1, version2)
            }
        }
    }
    
    // 模块 1 独有
    fmt.Printf("\n📦 仅在 %s 中:\n", mc.Module1.Name)
    unique1Count := 0
    for path, version := range deps1 {
        if _, exists := deps2[path]; !exists {
            fmt.Printf("  + %s %s\n", path, version)
            unique1Count++
        }
    }
    
    // 模块 2 独有
    fmt.Printf("\n📦 仅在 %s 中:\n", mc.Module2.Name)
    unique2Count := 0
    for path, version := range deps2 {
        if _, exists := deps1[path]; !exists {
            fmt.Printf("  + %s %s\n", path, version)
            unique2Count++
        }
    }
    
    fmt.Printf("\n📊 总结:\n")
    fmt.Printf("  共同依赖: %d\n", commonCount)
    fmt.Printf("  %s 独有: %d\n", mc.Module1.Name, unique1Count)
    fmt.Printf("  %s 独有: %d\n", mc.Module2.Name, unique2Count)
}

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
    
    comparison := &ModuleComparison{
        Module1: mod1,
        Module2: mod2,
        Path1:   os.Args[1],
        Path2:   os.Args[2],
    }
    
    comparison.CompareDependencies()
}
```

## 依赖图构建器

构建依赖关系：

```go
package main

import (
    "fmt"
    "log"
    "strings"
    
    "github.com/scagogogo/go-mod-parser/pkg"
    "github.com/scagogogo/go-mod-parser/pkg/module"
)

type DependencyNode struct {
    Path     string
    Version  string
    Indirect bool
    Children []*DependencyNode
}

type DependencyGraph struct {
    Root  *module.Module
    Nodes map[string]*DependencyNode
}

func NewDependencyGraph(mod *module.Module) *DependencyGraph {
    graph := &DependencyGraph{
        Root:  mod,
        Nodes: make(map[string]*DependencyNode),
    }
    
    // 为所有依赖构建节点
    for _, req := range mod.Requires {
        node := &DependencyNode{
            Path:     req.Path,
            Version:  req.Version,
            Indirect: req.Indirect,
        }
        graph.Nodes[req.Path] = node
    }
    
    return graph
}

func (dg *DependencyGraph) PrintGraph() {
    fmt.Printf("📊 %s 的依赖图\n", dg.Root.Name)
    fmt.Println(strings.Repeat("=", 50))
    
    // 按直接/间接分组
    var direct, indirect []*DependencyNode
    for _, node := range dg.Nodes {
        if node.Indirect {
            indirect = append(indirect, node)
        } else {
            direct = append(direct, node)
        }
    }
    
    fmt.Printf("\n🎯 直接依赖 (%d):\n", len(direct))
    for _, node := range direct {
        fmt.Printf("  ├── %s %s\n", node.Path, node.Version)
    }
    
    fmt.Printf("\n🔗 间接依赖 (%d):\n", len(indirect))
    for _, node := range indirect {
        fmt.Printf("  ├── %s %s\n", node.Path, node.Version)
    }
}

func (dg *DependencyGraph) FindCycles() [][]string {
    // 简化的循环检测（真正的循环需要更复杂的实现）
    var cycles [][]string
    
    // 检查替换中的自引用
    for _, rep := range dg.Root.Replaces {
        if rep.Old.Path == rep.New.Path {
            cycles = append(cycles, []string{rep.Old.Path, rep.New.Path})
        }
    }
    
    return cycles
}

func (dg *DependencyGraph) AnalyzeDepth() map[string]int {
    depths := make(map[string]int)
    
    // 基于路径段的简单深度分析
    for path := range dg.Nodes {
        segments := strings.Split(path, "/")
        depths[path] = len(segments)
    }
    
    return depths
}

func main() {
    mod, err := pkg.FindAndParseGoModInCurrentDir()
    if err != nil {
        log.Fatalf("解析 go.mod 失败: %v", err)
    }
    
    graph := NewDependencyGraph(mod)
    graph.PrintGraph()
    
    // 分析循环
    cycles := graph.FindCycles()
    if len(cycles) > 0 {
        fmt.Printf("\n⚠️  检测到潜在循环:\n")
        for _, cycle := range cycles {
            fmt.Printf("  %s\n", strings.Join(cycle, " -> "))
        }
    } else {
        fmt.Printf("\n✅ 未检测到循环\n")
    }
    
    // 分析深度
    depths := graph.AnalyzeDepth()
    fmt.Printf("\n📏 依赖深度分析:\n")
    for path, depth := range depths {
        if depth > 3 {
            fmt.Printf("  深层依赖: %s (深度: %d)\n", path, depth)
        }
    }
}
```

## 配置驱动的分析

创建可配置的分析工具：

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"
    "os"
    "regexp"
    "strings"
    
    "github.com/scagogogo/go-mod-parser/pkg"
    "github.com/scagogogo/go-mod-parser/pkg/module"
)

type AnalysisConfig struct {
    Rules []Rule `json:"rules"`
}

type Rule struct {
    Name        string `json:"name"`
    Type        string `json:"type"`
    Pattern     string `json:"pattern"`
    Severity    string `json:"severity"`
    Description string `json:"description"`
}

type ConfigurableAnalyzer struct {
    Config *AnalysisConfig
}

func (ca *ConfigurableAnalyzer) LoadConfig(filename string) error {
    data, err := os.ReadFile(filename)
    if err != nil {
        return err
    }
    
    return json.Unmarshal(data, &ca.Config)
}

func (ca *ConfigurableAnalyzer) Analyze(mod *module.Module) {
    fmt.Printf("🔍 %s 的可配置分析\n", mod.Name)
    fmt.Println(strings.Repeat("=", 50))
    
    for _, rule := range ca.Config.Rules {
        ca.applyRule(mod, rule)
    }
}

func (ca *ConfigurableAnalyzer) applyRule(mod *module.Module, rule Rule) {
    fmt.Printf("\n📋 规则: %s (%s)\n", rule.Name, rule.Severity)
    fmt.Printf("描述: %s\n", rule.Description)
    
    pattern, err := regexp.Compile(rule.Pattern)
    if err != nil {
        fmt.Printf("❌ 无效模式: %v\n", err)
        return
    }
    
    matches := 0
    
    switch rule.Type {
    case "dependency":
        for _, req := range mod.Requires {
            if pattern.MatchString(req.Path) {
                fmt.Printf("  匹配: %s %s\n", req.Path, req.Version)
                matches++
            }
        }
    case "replace":
        for _, rep := range mod.Replaces {
            if pattern.MatchString(rep.Old.Path) || pattern.MatchString(rep.New.Path) {
                fmt.Printf("  匹配: %s => %s\n", rep.Old.Path, rep.New.Path)
                matches++
            }
        }
    case "version":
        for _, req := range mod.Requires {
            if pattern.MatchString(req.Version) {
                fmt.Printf("  匹配: %s %s\n", req.Path, req.Version)
                matches++
            }
        }
    }
    
    if matches == 0 {
        fmt.Println("  ✅ 未找到匹配")
    } else {
        fmt.Printf("  找到 %d 个匹配\n", matches)
    }
}

// 示例配置
func createExampleConfig() {
    config := AnalysisConfig{
        Rules: []Rule{
            {
                Name:        "废弃的依赖",
                Type:        "dependency",
                Pattern:     "github\\.com/(golang/dep|Masterminds/glide)",
                Severity:    "warning",
                Description: "检查废弃的依赖管理工具",
            },
            {
                Name:        "预发布版本",
                Type:        "version",
                Pattern:     "v\\d+\\.\\d+\\.\\d+-\\w+",
                Severity:    "info",
                Description: "识别预发布版本",
            },
            {
                Name:        "本地替换",
                Type:        "replace",
                Pattern:     "^\\./",
                Severity:    "warning",
                Description: "检查本地路径替换",
            },
        },
    }
    
    data, _ := json.MarshalIndent(config, "", "  ")
    os.WriteFile("analysis-config.json", data, 0644)
    fmt.Println("创建示例配置: analysis-config.json")
}

func main() {
    if len(os.Args) > 1 && os.Args[1] == "create-config" {
        createExampleConfig()
        return
    }
    
    mod, err := pkg.FindAndParseGoModInCurrentDir()
    if err != nil {
        log.Fatalf("解析 go.mod 失败: %v", err)
    }
    
    analyzer := &ConfigurableAnalyzer{}
    
    configFile := "analysis-config.json"
    if len(os.Args) > 1 {
        configFile = os.Args[1]
    }
    
    if err := analyzer.LoadConfig(configFile); err != nil {
        log.Fatalf("加载配置失败: %v", err)
    }
    
    analyzer.Analyze(mod)
}
```

## 下一步

这些高级模式展示了 Go Mod Parser 构建复杂分析工具的灵活性。你可以：

- 使用自定义分析器扩展分析框架
- 使用比较工具构建 CI/CD 集成
- 使用图构建器创建依赖管理仪表板
- 使用配置驱动分析实现策略执行

更多示例，请参见：
- [基础解析](/zh/examples/basic-parsing) - 基本操作
- [文件发现](/zh/examples/file-discovery) - 自动发现功能
- [依赖分析](/zh/examples/dependency-analysis) - 依赖分析模式
