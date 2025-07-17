# Advanced Usage

This section covers advanced patterns and complex use cases for Go Mod Parser.

## Custom Analysis Framework

Build a flexible analysis framework:

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/go-mod-parser/pkg"
    "github.com/scagogogo/go-mod-parser/pkg/module"
)

// Analyzer interface for extensible analysis
type Analyzer interface {
    Name() string
    Analyze(mod *module.Module) AnalysisResult
}

type AnalysisResult struct {
    Summary string
    Details []string
    Issues  []string
}

// Security analyzer
type SecurityAnalyzer struct{}

func (s SecurityAnalyzer) Name() string {
    return "Security Analysis"
}

func (s SecurityAnalyzer) Analyze(mod *module.Module) AnalysisResult {
    result := AnalysisResult{}
    
    // Check for retracted versions
    for _, req := range mod.Requires {
        if pkg.HasRetract(mod, req.Version) {
            result.Issues = append(result.Issues, 
                fmt.Sprintf("Using retracted version: %s %s", req.Path, req.Version))
        }
    }
    
    // Check for local replacements
    for _, rep := range mod.Replaces {
        if strings.HasPrefix(rep.New.Path, "./") || strings.HasPrefix(rep.New.Path, "../") {
            result.Issues = append(result.Issues, 
                fmt.Sprintf("Local replacement: %s => %s", rep.Old.Path, rep.New.Path))
        }
    }
    
    if len(result.Issues) == 0 {
        result.Summary = "No security issues found"
    } else {
        result.Summary = fmt.Sprintf("%d security issues found", len(result.Issues))
    }
    
    return result
}

// Performance analyzer
type PerformanceAnalyzer struct{}

func (p PerformanceAnalyzer) Name() string {
    return "Performance Analysis"
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
                    fmt.Sprintf("Heavy dependency detected: %s", req.Path))
            }
        }
    }
    
    result.Summary = fmt.Sprintf("Analyzed %d dependencies for performance impact", len(mod.Requires))
    return result
}

// Analysis runner
func runAnalysis(mod *module.Module, analyzers []Analyzer) {
    fmt.Printf("Running analysis for module: %s\n", mod.Name)
    fmt.Println(strings.Repeat("=", 60))
    
    for _, analyzer := range analyzers {
        fmt.Printf("\n🔍 %s\n", analyzer.Name())
        result := analyzer.Analyze(mod)
        
        fmt.Printf("Summary: %s\n", result.Summary)
        
        if len(result.Details) > 0 {
            fmt.Println("Details:")
            for _, detail := range result.Details {
                fmt.Printf("  • %s\n", detail)
            }
        }
        
        if len(result.Issues) > 0 {
            fmt.Println("Issues:")
            for _, issue := range result.Issues {
                fmt.Printf("  ⚠️  %s\n", issue)
            }
        }
    }
}

func main() {
    mod, err := pkg.FindAndParseGoModInCurrentDir()
    if err != nil {
        log.Fatalf("Failed to parse go.mod: %v", err)
    }
    
    analyzers := []Analyzer{
        SecurityAnalyzer{},
        PerformanceAnalyzer{},
    }
    
    runAnalysis(mod, analyzers)
}
```

## Module Comparison Tool

Compare multiple go.mod files:

```go
package main

import (
    "fmt"
    "log"
    "os"
    
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
    
    fmt.Printf("📋 Dependency Comparison\n")
    fmt.Printf("Module 1: %s (%s)\n", mc.Module1.Name, mc.Path1)
    fmt.Printf("Module 2: %s (%s)\n", mc.Module2.Name, mc.Path2)
    fmt.Println(strings.Repeat("-", 60))
    
    // Common dependencies
    fmt.Println("\n🤝 Common Dependencies:")
    commonCount := 0
    for path, version1 := range deps1 {
        if version2, exists := deps2[path]; exists {
            commonCount++
            if version1 == version2 {
                fmt.Printf("  ✅ %s: %s (same)\n", path, version1)
            } else {
                fmt.Printf("  ⚠️  %s: %s vs %s\n", path, version1, version2)
            }
        }
    }
    
    // Unique to module 1
    fmt.Printf("\n📦 Only in %s:\n", mc.Module1.Name)
    unique1Count := 0
    for path, version := range deps1 {
        if _, exists := deps2[path]; !exists {
            fmt.Printf("  + %s %s\n", path, version)
            unique1Count++
        }
    }
    
    // Unique to module 2
    fmt.Printf("\n📦 Only in %s:\n", mc.Module2.Name)
    unique2Count := 0
    for path, version := range deps2 {
        if _, exists := deps1[path]; !exists {
            fmt.Printf("  + %s %s\n", path, version)
            unique2Count++
        }
    }
    
    fmt.Printf("\n📊 Summary:\n")
    fmt.Printf("  Common dependencies: %d\n", commonCount)
    fmt.Printf("  Unique to %s: %d\n", mc.Module1.Name, unique1Count)
    fmt.Printf("  Unique to %s: %d\n", mc.Module2.Name, unique2Count)
}

func main() {
    if len(os.Args) != 3 {
        fmt.Println("Usage: compare <go.mod1> <go.mod2>")
        os.Exit(1)
    }
    
    mod1, err := pkg.ParseGoModFile(os.Args[1])
    if err != nil {
        log.Fatalf("Failed to parse %s: %v", os.Args[1], err)
    }
    
    mod2, err := pkg.ParseGoModFile(os.Args[2])
    if err != nil {
        log.Fatalf("Failed to parse %s: %v", os.Args[2], err)
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

## Dependency Graph Builder

Build dependency relationships:

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
    
    // Build nodes for all dependencies
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
    fmt.Printf("📊 Dependency Graph for %s\n", dg.Root.Name)
    fmt.Println(strings.Repeat("=", 50))
    
    // Group by direct/indirect
    var direct, indirect []*DependencyNode
    for _, node := range dg.Nodes {
        if node.Indirect {
            indirect = append(indirect, node)
        } else {
            direct = append(direct, node)
        }
    }
    
    fmt.Printf("\n🎯 Direct Dependencies (%d):\n", len(direct))
    for _, node := range direct {
        fmt.Printf("  ├── %s %s\n", node.Path, node.Version)
    }
    
    fmt.Printf("\n🔗 Indirect Dependencies (%d):\n", len(indirect))
    for _, node := range indirect {
        fmt.Printf("  ├── %s %s\n", node.Path, node.Version)
    }
}

func (dg *DependencyGraph) FindCycles() [][]string {
    // Simplified cycle detection (would need more complex implementation for real cycles)
    var cycles [][]string
    
    // Check for self-references in replacements
    for _, rep := range dg.Root.Replaces {
        if rep.Old.Path == rep.New.Path {
            cycles = append(cycles, []string{rep.Old.Path, rep.New.Path})
        }
    }
    
    return cycles
}

func (dg *DependencyGraph) AnalyzeDepth() map[string]int {
    depths := make(map[string]int)
    
    // Simple depth analysis based on path segments
    for path := range dg.Nodes {
        segments := strings.Split(path, "/")
        depths[path] = len(segments)
    }
    
    return depths
}

func main() {
    mod, err := pkg.FindAndParseGoModInCurrentDir()
    if err != nil {
        log.Fatalf("Failed to parse go.mod: %v", err)
    }
    
    graph := NewDependencyGraph(mod)
    graph.PrintGraph()
    
    // Analyze cycles
    cycles := graph.FindCycles()
    if len(cycles) > 0 {
        fmt.Printf("\n⚠️  Potential Cycles Detected:\n")
        for _, cycle := range cycles {
            fmt.Printf("  %s\n", strings.Join(cycle, " -> "))
        }
    } else {
        fmt.Printf("\n✅ No cycles detected\n")
    }
    
    // Analyze depths
    depths := graph.AnalyzeDepth()
    fmt.Printf("\n📏 Dependency Depth Analysis:\n")
    for path, depth := range depths {
        if depth > 3 {
            fmt.Printf("  Deep dependency: %s (depth: %d)\n", path, depth)
        }
    }
}
```

## Configuration-Driven Analysis

Create configurable analysis tools:

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"
    "os"
    "regexp"
    
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
    fmt.Printf("🔍 Configurable Analysis for %s\n", mod.Name)
    fmt.Println(strings.Repeat("=", 50))
    
    for _, rule := range ca.Config.Rules {
        ca.applyRule(mod, rule)
    }
}

func (ca *ConfigurableAnalyzer) applyRule(mod *module.Module, rule Rule) {
    fmt.Printf("\n📋 Rule: %s (%s)\n", rule.Name, rule.Severity)
    fmt.Printf("Description: %s\n", rule.Description)
    
    pattern, err := regexp.Compile(rule.Pattern)
    if err != nil {
        fmt.Printf("❌ Invalid pattern: %v\n", err)
        return
    }
    
    matches := 0
    
    switch rule.Type {
    case "dependency":
        for _, req := range mod.Requires {
            if pattern.MatchString(req.Path) {
                fmt.Printf("  Match: %s %s\n", req.Path, req.Version)
                matches++
            }
        }
    case "replace":
        for _, rep := range mod.Replaces {
            if pattern.MatchString(rep.Old.Path) || pattern.MatchString(rep.New.Path) {
                fmt.Printf("  Match: %s => %s\n", rep.Old.Path, rep.New.Path)
                matches++
            }
        }
    case "version":
        for _, req := range mod.Requires {
            if pattern.MatchString(req.Version) {
                fmt.Printf("  Match: %s %s\n", req.Path, req.Version)
                matches++
            }
        }
    }
    
    if matches == 0 {
        fmt.Println("  ✅ No matches found")
    } else {
        fmt.Printf("  Found %d matches\n", matches)
    }
}

// Example configuration
func createExampleConfig() {
    config := AnalysisConfig{
        Rules: []Rule{
            {
                Name:        "Deprecated Dependencies",
                Type:        "dependency",
                Pattern:     "github\\.com/(golang/dep|Masterminds/glide)",
                Severity:    "warning",
                Description: "Check for deprecated dependency management tools",
            },
            {
                Name:        "Pre-release Versions",
                Type:        "version",
                Pattern:     "v\\d+\\.\\d+\\.\\d+-\\w+",
                Severity:    "info",
                Description: "Identify pre-release versions",
            },
            {
                Name:        "Local Replacements",
                Type:        "replace",
                Pattern:     "^\\./",
                Severity:    "warning",
                Description: "Check for local path replacements",
            },
        },
    }
    
    data, _ := json.MarshalIndent(config, "", "  ")
    os.WriteFile("analysis-config.json", data, 0644)
    fmt.Println("Created example configuration: analysis-config.json")
}

func main() {
    if len(os.Args) > 1 && os.Args[1] == "create-config" {
        createExampleConfig()
        return
    }
    
    mod, err := pkg.FindAndParseGoModInCurrentDir()
    if err != nil {
        log.Fatalf("Failed to parse go.mod: %v", err)
    }
    
    analyzer := &ConfigurableAnalyzer{}
    
    configFile := "analysis-config.json"
    if len(os.Args) > 1 {
        configFile = os.Args[1]
    }
    
    if err := analyzer.LoadConfig(configFile); err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }
    
    analyzer.Analyze(mod)
}
```

## Next Steps

These advanced patterns demonstrate the flexibility of Go Mod Parser for building sophisticated analysis tools. You can:

- Extend the analysis framework with custom analyzers
- Build CI/CD integrations using the comparison tools
- Create dependency management dashboards with the graph builder
- Implement policy enforcement with configuration-driven analysis

For more examples, see:
- [Basic Parsing](/examples/basic-parsing) - Fundamental operations
- [File Discovery](/examples/file-discovery) - Auto-discovery features
- [Dependency Analysis](/examples/dependency-analysis) - Dependency analysis patterns
