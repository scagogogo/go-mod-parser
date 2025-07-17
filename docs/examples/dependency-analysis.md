# Dependency Analysis

This section demonstrates how to analyze dependencies, replacements, exclusions, and retractions in go.mod files.

## Basic Dependency Checking

Check if specific dependencies exist:

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/go-mod-parser/pkg"
)

func main() {
    mod, err := pkg.FindAndParseGoModInCurrentDir()
    if err != nil {
        log.Fatalf("Failed to parse go.mod: %v", err)
    }
    
    // Check for specific dependencies
    dependencies := []string{
        "github.com/gin-gonic/gin",
        "github.com/gorilla/mux",
        "github.com/stretchr/testify",
    }
    
    fmt.Println("Dependency Check:")
    for _, dep := range dependencies {
        if pkg.HasRequire(mod, dep) {
            req := pkg.GetRequire(mod, dep)
            fmt.Printf("✅ %s %s", dep, req.Version)
            if req.Indirect {
                fmt.Printf(" (indirect)")
            }
            fmt.Println()
        } else {
            fmt.Printf("❌ %s (not found)\n", dep)
        }
    }
}
```

## Framework Detection

Detect popular Go frameworks:

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/go-mod-parser/pkg"
)

func detectFrameworks(mod *module.Module) {
    frameworks := map[string]string{
        "github.com/gin-gonic/gin":     "Gin Web Framework",
        "github.com/gorilla/mux":       "Gorilla Mux Router",
        "github.com/labstack/echo/v4":  "Echo Web Framework",
        "github.com/gofiber/fiber/v2":  "Fiber Web Framework",
        "github.com/beego/beego/v2":    "Beego Framework",
        "github.com/revel/revel":       "Revel Framework",
    }
    
    fmt.Println("🔍 Framework Detection:")
    found := false
    
    for path, name := range frameworks {
        if pkg.HasRequire(mod, path) {
            req := pkg.GetRequire(mod, path)
            fmt.Printf("  ✅ %s (%s)\n", name, req.Version)
            found = true
        }
    }
    
    if !found {
        fmt.Println("  ❌ No popular web frameworks detected")
    }
}

func main() {
    mod, err := pkg.FindAndParseGoModInCurrentDir()
    if err != nil {
        log.Fatalf("Failed to parse go.mod: %v", err)
    }
    
    detectFrameworks(mod)
}
```

## Dependency Statistics

Analyze dependency patterns:

```go
package main

import (
    "fmt"
    "log"
    "sort"
    "strings"
    
    "github.com/scagogogo/go-mod-parser/pkg"
)

func analyzeDependencyStats(mod *module.Module) {
    direct := 0
    indirect := 0
    domains := make(map[string]int)
    
    for _, req := range mod.Requires {
        if req.Indirect {
            indirect++
        } else {
            direct++
        }
        
        // Extract domain from module path
        parts := strings.Split(req.Path, "/")
        if len(parts) > 0 {
            domain := parts[0]
            domains[domain]++
        }
    }
    
    fmt.Printf("📊 Dependency Statistics:\n")
    fmt.Printf("  Total Dependencies: %d\n", len(mod.Requires))
    fmt.Printf("  Direct Dependencies: %d\n", direct)
    fmt.Printf("  Indirect Dependencies: %d\n", indirect)
    fmt.Printf("  Replace Directives: %d\n", len(mod.Replaces))
    fmt.Printf("  Exclude Directives: %d\n", len(mod.Excludes))
    fmt.Printf("  Retract Directives: %d\n", len(mod.Retracts))
    
    // Top domains
    fmt.Println("\n🌐 Top Dependency Domains:")
    type domainCount struct {
        domain string
        count  int
    }
    
    var sortedDomains []domainCount
    for domain, count := range domains {
        sortedDomains = append(sortedDomains, domainCount{domain, count})
    }
    
    sort.Slice(sortedDomains, func(i, j int) bool {
        return sortedDomains[i].count > sortedDomains[j].count
    })
    
    for i, dc := range sortedDomains {
        if i >= 5 { // Show top 5
            break
        }
        fmt.Printf("  %d. %s (%d dependencies)\n", i+1, dc.domain, dc.count)
    }
}

func main() {
    mod, err := pkg.FindAndParseGoModInCurrentDir()
    if err != nil {
        log.Fatalf("Failed to parse go.mod: %v", err)
    }
    
    analyzeDependencyStats(mod)
}
```

## Replace Directive Analysis

Analyze replacement patterns:

```go
package main

import (
    "fmt"
    "log"
    "strings"
    
    "github.com/scagogogo/go-mod-parser/pkg"
)

func analyzeReplaces(mod *module.Module) {
    if len(mod.Replaces) == 0 {
        fmt.Println("🔄 No replace directives found")
        return
    }
    
    fmt.Printf("🔄 Replace Directive Analysis (%d total):\n", len(mod.Replaces))
    
    localReplaces := 0
    moduleReplaces := 0
    
    for i, rep := range mod.Replaces {
        fmt.Printf("\n%d. %s", i+1, rep.Old.Path)
        if rep.Old.Version != "" {
            fmt.Printf(" %s", rep.Old.Version)
        }
        fmt.Printf(" => %s", rep.New.Path)
        if rep.New.Version != "" {
            fmt.Printf(" %s", rep.New.Version)
        }
        
        if strings.HasPrefix(rep.New.Path, "./") || strings.HasPrefix(rep.New.Path, "../") {
            fmt.Printf(" (local path)")
            localReplaces++
        } else {
            fmt.Printf(" (module)")
            moduleReplaces++
        }
    }
    
    fmt.Printf("\n\nSummary:\n")
    fmt.Printf("  Local Path Replacements: %d\n", localReplaces)
    fmt.Printf("  Module Replacements: %d\n", moduleReplaces)
}

func main() {
    mod, err := pkg.FindAndParseGoModInCurrentDir()
    if err != nil {
        log.Fatalf("Failed to parse go.mod: %v", err)
    }
    
    analyzeReplaces(mod)
}
```

## Security Analysis

Check for potential security issues:

```go
package main

import (
    "fmt"
    "log"
    "strings"
    
    "github.com/scagogogo/go-mod-parser/pkg"
)

func securityAnalysis(mod *module.Module) {
    fmt.Println("🔒 Security Analysis:")
    
    issues := 0
    
    // Check for retracted versions in use
    fmt.Println("\n⚠️  Retracted Version Check:")
    retractedFound := false
    for _, req := range mod.Requires {
        if pkg.HasRetract(mod, req.Version) {
            fmt.Printf("  ❌ Using retracted version: %s %s\n", req.Path, req.Version)
            issues++
            retractedFound = true
        }
    }
    if !retractedFound {
        fmt.Println("  ✅ No retracted versions in use")
    }
    
    // Check for local path replacements
    fmt.Println("\n🔍 Local Path Replacement Check:")
    localFound := false
    for _, rep := range mod.Replaces {
        if strings.HasPrefix(rep.New.Path, "./") || strings.HasPrefix(rep.New.Path, "../") {
            fmt.Printf("  ⚠️  Local replacement: %s => %s\n", rep.Old.Path, rep.New.Path)
            issues++
            localFound = true
        }
    }
    if !localFound {
        fmt.Println("  ✅ No local path replacements")
    }
    
    // Check for development dependencies in production
    fmt.Println("\n🧪 Development Dependency Check:")
    devDeps := []string{"testify", "mock", "test", "debug", "dev"}
    devFound := false
    for _, req := range mod.Requires {
        if !req.Indirect {
            for _, devKeyword := range devDeps {
                if strings.Contains(strings.ToLower(req.Path), devKeyword) {
                    fmt.Printf("  ⚠️  Potential dev dependency as direct: %s\n", req.Path)
                    devFound = true
                    break
                }
            }
        }
    }
    if !devFound {
        fmt.Println("  ✅ No obvious development dependencies as direct")
    }
    
    fmt.Printf("\n📋 Summary: %d potential security issues found\n", issues)
}

func main() {
    mod, err := pkg.FindAndParseGoModInCurrentDir()
    if err != nil {
        log.Fatalf("Failed to parse go.mod: %v", err)
    }
    
    securityAnalysis(mod)
}
```

## Version Analysis

Analyze version patterns:

```go
package main

import (
    "fmt"
    "log"
    "regexp"
    "strings"
    
    "github.com/scagogogo/go-mod-parser/pkg"
)

func analyzeVersions(mod *module.Module) {
    fmt.Println("📋 Version Analysis:")
    
    semverPattern := regexp.MustCompile(`^v\d+\.\d+\.\d+`)
    preReleasePattern := regexp.MustCompile(`-\w+`)
    
    semverCount := 0
    preReleaseCount := 0
    pseudoVersionCount := 0
    
    fmt.Println("\nDependency Versions:")
    for _, req := range mod.Requires {
        version := req.Version
        fmt.Printf("  %s: %s", req.Path, version)
        
        if semverPattern.MatchString(version) {
            semverCount++
            if preReleasePattern.MatchString(version) {
                fmt.Printf(" (pre-release)")
                preReleaseCount++
            } else {
                fmt.Printf(" (stable)")
            }
        } else if strings.Contains(version, "-") && len(version) > 20 {
            fmt.Printf(" (pseudo-version)")
            pseudoVersionCount++
        } else {
            fmt.Printf(" (other)")
        }
        
        if req.Indirect {
            fmt.Printf(" [indirect]")
        }
        fmt.Println()
    }
    
    fmt.Printf("\nVersion Summary:\n")
    fmt.Printf("  Semantic Versions: %d\n", semverCount)
    fmt.Printf("  Pre-release Versions: %d\n", preReleaseCount)
    fmt.Printf("  Pseudo Versions: %d\n", pseudoVersionCount)
    fmt.Printf("  Other Versions: %d\n", len(mod.Requires)-semverCount-pseudoVersionCount)
}

func main() {
    mod, err := pkg.FindAndParseGoModInCurrentDir()
    if err != nil {
        log.Fatalf("Failed to parse go.mod: %v", err)
    }
    
    analyzeVersions(mod)
}
```

## Comprehensive Analysis Tool

A complete analysis tool combining all features:

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/go-mod-parser/pkg"
)

func comprehensiveAnalysis(mod *module.Module) {
    fmt.Printf("📦 Module: %s\n", mod.Name)
    fmt.Printf("🐹 Go Version: %s\n", mod.GoVersion)
    fmt.Println(strings.Repeat("=", 50))
    
    // Basic stats
    analyzeDependencyStats(mod)
    fmt.Println()
    
    // Framework detection
    detectFrameworks(mod)
    fmt.Println()
    
    // Replace analysis
    analyzeReplaces(mod)
    fmt.Println()
    
    // Security analysis
    securityAnalysis(mod)
    fmt.Println()
    
    // Version analysis
    analyzeVersions(mod)
}

func main() {
    mod, err := pkg.FindAndParseGoModInCurrentDir()
    if err != nil {
        log.Fatalf("Failed to parse go.mod: %v", err)
    }
    
    comprehensiveAnalysis(mod)
}
```

## Next Steps

- [Advanced Usage](/examples/advanced-usage) - Complex analysis patterns
- [File Discovery](/examples/file-discovery) - Auto-discovery features
- [Basic Parsing](/examples/basic-parsing) - Basic parsing operations
