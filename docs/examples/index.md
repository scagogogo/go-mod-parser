# Examples

This section provides comprehensive examples demonstrating how to use Go Mod Parser in various scenarios. Each example includes complete, runnable code with explanations.

## Overview

The examples are organized by complexity and use case:

- **[Basic Parsing](/examples/basic-parsing)** - Simple parsing operations
- **[File Discovery](/examples/file-discovery)** - Auto-discovering go.mod files
- **[Dependency Analysis](/examples/dependency-analysis)** - Analyzing dependencies and relationships
- **[Advanced Usage](/examples/advanced-usage)** - Complex scenarios and best practices

## Quick Examples

### Parse a go.mod File

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/go-mod-parser/pkg"
)

func main() {
    mod, err := pkg.ParseGoModFile("go.mod")
    if err != nil {
        log.Fatalf("Failed to parse: %v", err)
    }
    
    fmt.Printf("Module: %s\n", mod.Name)
    fmt.Printf("Go Version: %s\n", mod.GoVersion)
    fmt.Printf("Dependencies: %d\n", len(mod.Requires))
}
```

### Check for Specific Dependencies

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
        log.Fatalf("Failed to find go.mod: %v", err)
    }
    
    // Check for popular frameworks
    frameworks := []string{
        "github.com/gin-gonic/gin",
        "github.com/gorilla/mux", 
        "github.com/labstack/echo/v4",
        "github.com/gofiber/fiber/v2",
    }
    
    fmt.Println("Framework Detection:")
    for _, framework := range frameworks {
        if pkg.HasRequire(mod, framework) {
            req := pkg.GetRequire(mod, framework)
            fmt.Printf("✓ %s %s\n", framework, req.Version)
        } else {
            fmt.Printf("✗ %s\n", framework)
        }
    }
}
```

### Analyze Replace Directives

```go
package main

import (
    "fmt"
    "log"
    "strings"
    
    "github.com/scagogogo/go-mod-parser/pkg"
)

func main() {
    mod, err := pkg.FindAndParseGoModInCurrentDir()
    if err != nil {
        log.Fatalf("Failed to parse: %v", err)
    }
    
    if len(mod.Replaces) == 0 {
        fmt.Println("No replace directives found")
        return
    }
    
    fmt.Printf("Found %d replace directive(s):\n\n", len(mod.Replaces))
    
    for i, rep := range mod.Replaces {
        fmt.Printf("%d. %s", i+1, rep.Old.Path)
        if rep.Old.Version != "" {
            fmt.Printf(" %s", rep.Old.Version)
        }
        fmt.Printf(" => %s", rep.New.Path)
        if rep.New.Version != "" {
            fmt.Printf(" %s", rep.New.Version)
        }
        
        // Determine replacement type
        if strings.HasPrefix(rep.New.Path, "./") || strings.HasPrefix(rep.New.Path, "../") {
            fmt.Printf(" (local path)")
        } else {
            fmt.Printf(" (module)")
        }
        fmt.Println()
    }
}
```

## Example Projects

### CLI Tool for go.mod Analysis

```go
package main

import (
    "flag"
    "fmt"
    "log"
    "os"
    "path/filepath"
    "strings"
    
    "github.com/scagogogo/go-mod-parser/pkg"
)

func main() {
    var (
        path = flag.String("path", "", "Path to go.mod file or directory")
        verbose = flag.Bool("verbose", false, "Verbose output")
        checkSecurity = flag.Bool("security", false, "Check for security issues")
    )
    flag.Parse()
    
    var mod *module.Module
    var err error
    
    if *path != "" {
        if strings.HasSuffix(*path, "go.mod") {
            mod, err = pkg.ParseGoModFile(*path)
        } else {
            mod, err = pkg.FindAndParseGoModFile(*path)
        }
    } else {
        mod, err = pkg.FindAndParseGoModInCurrentDir()
    }
    
    if err != nil {
        log.Fatalf("Error: %v", err)
    }
    
    // Basic information
    fmt.Printf("📦 Module: %s\n", mod.Name)
    fmt.Printf("🐹 Go Version: %s\n", mod.GoVersion)
    fmt.Printf("📋 Dependencies: %d\n", len(mod.Requires))
    
    if *verbose {
        printVerboseInfo(mod)
    }
    
    if *checkSecurity {
        checkSecurityIssues(mod)
    }
}

func printVerboseInfo(mod *module.Module) {
    // Print dependencies
    if len(mod.Requires) > 0 {
        fmt.Println("\n📋 Dependencies:")
        for _, req := range mod.Requires {
            fmt.Printf("  %s %s", req.Path, req.Version)
            if req.Indirect {
                fmt.Printf(" (indirect)")
            }
            fmt.Println()
        }
    }
    
    // Print replacements
    if len(mod.Replaces) > 0 {
        fmt.Println("\n🔄 Replacements:")
        for _, rep := range mod.Replaces {
            fmt.Printf("  %s => %s", rep.Old.Path, rep.New.Path)
            if rep.New.Version != "" {
                fmt.Printf(" %s", rep.New.Version)
            }
            fmt.Println()
        }
    }
    
    // Print exclusions
    if len(mod.Excludes) > 0 {
        fmt.Println("\n🚫 Exclusions:")
        for _, exc := range mod.Excludes {
            fmt.Printf("  %s %s\n", exc.Path, exc.Version)
        }
    }
    
    // Print retractions
    if len(mod.Retracts) > 0 {
        fmt.Println("\n⚠️  Retractions:")
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
}

func checkSecurityIssues(mod *module.Module) {
    fmt.Println("\n🔒 Security Analysis:")
    
    issues := 0
    
    // Check for retracted versions in use
    for _, req := range mod.Requires {
        if pkg.HasRetract(mod, req.Version) {
            fmt.Printf("⚠️  Using retracted version: %s %s\n", req.Path, req.Version)
            issues++
        }
    }
    
    // Check for local path replacements
    for _, rep := range mod.Replaces {
        if strings.HasPrefix(rep.New.Path, "./") || strings.HasPrefix(rep.New.Path, "../") {
            fmt.Printf("🔍 Local replacement detected: %s => %s\n", rep.Old.Path, rep.New.Path)
            issues++
        }
    }
    
    if issues == 0 {
        fmt.Println("✅ No security issues detected")
    } else {
        fmt.Printf("Found %d potential security issue(s)\n", issues)
    }
}
```

### Dependency Comparison Tool

```go
package main

import (
    "fmt"
    "log"
    "os"
    
    "github.com/scagogogo/go-mod-parser/pkg"
)

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
    
    compareDependencies(mod1, mod2)
}

func compareDependencies(mod1, mod2 *module.Module) {
    fmt.Printf("Comparing %s vs %s\n\n", mod1.Name, mod2.Name)
    
    // Create maps for easier lookup
    deps1 := make(map[string]string)
    deps2 := make(map[string]string)
    
    for _, req := range mod1.Requires {
        deps1[req.Path] = req.Version
    }
    
    for _, req := range mod2.Requires {
        deps2[req.Path] = req.Version
    }
    
    // Find common dependencies
    fmt.Println("📋 Common Dependencies:")
    for path, version1 := range deps1 {
        if version2, exists := deps2[path]; exists {
            if version1 == version2 {
                fmt.Printf("  ✓ %s %s\n", path, version1)
            } else {
                fmt.Printf("  ⚠️  %s: %s vs %s\n", path, version1, version2)
            }
        }
    }
    
    // Find dependencies only in mod1
    fmt.Println("\n📦 Only in " + mod1.Name + ":")
    for path, version := range deps1 {
        if _, exists := deps2[path]; !exists {
            fmt.Printf("  + %s %s\n", path, version)
        }
    }
    
    // Find dependencies only in mod2
    fmt.Println("\n📦 Only in " + mod2.Name + ":")
    for path, version := range deps2 {
        if _, exists := deps1[path]; !exists {
            fmt.Printf("  + %s %s\n", path, version)
        }
    }
}
```

## Running the Examples

1. **Save the code** to a `.go` file
2. **Initialize a Go module** (if not already done):
   ```bash
   go mod init example
   ```
3. **Add the dependency**:
   ```bash
   go get github.com/scagogogo/go-mod-parser
   ```
4. **Run the example**:
   ```bash
   go run main.go
   ```

## Next Steps

Explore the detailed example categories:

- **[Basic Parsing](/examples/basic-parsing)** - Start here for simple use cases
- **[File Discovery](/examples/file-discovery)** - Learn about auto-discovery features  
- **[Dependency Analysis](/examples/dependency-analysis)** - Advanced dependency analysis
- **[Advanced Usage](/examples/advanced-usage)** - Complex scenarios and patterns

Each section includes multiple examples with full source code and explanations.
