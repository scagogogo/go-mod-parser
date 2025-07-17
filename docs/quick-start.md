# Quick Start

This guide will help you get started with Go Mod Parser in just a few minutes.

## Basic Usage

### Parse a go.mod File

The most common use case is parsing an existing go.mod file:

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/go-mod-parser/pkg"
)

func main() {
    // Parse a go.mod file by path
    mod, err := pkg.ParseGoModFile("path/to/go.mod")
    if err != nil {
        log.Fatalf("Failed to parse go.mod: %v", err)
    }
    
    // Access basic information
    fmt.Printf("Module Name: %s\n", mod.Name)
    fmt.Printf("Go Version: %s\n", mod.GoVersion)
    
    // List all dependencies
    fmt.Println("\nDependencies:")
    for _, req := range mod.Requires {
        indirect := ""
        if req.Indirect {
            indirect = " // indirect"
        }
        fmt.Printf("- %s %s%s\n", req.Path, req.Version, indirect)
    }
}
```

### Parse go.mod Content

You can also parse go.mod content directly from a string:

```go
content := `module github.com/example/myproject

go 1.21

require (
    github.com/gin-gonic/gin v1.9.1
    github.com/stretchr/testify v1.8.4 // indirect
)

replace github.com/old/package => github.com/new/package v1.0.0

exclude github.com/problematic/package v1.0.0

retract v1.0.1 // security vulnerability
`

mod, err := pkg.ParseGoModContent(content)
if err != nil {
    log.Fatalf("Failed to parse content: %v", err)
}

fmt.Printf("Parsed module: %s\n", mod.Name)
```

### Auto-discover go.mod Files

The library can automatically find go.mod files in the current directory or parent directories:

```go
// Find and parse go.mod in current directory or parent directories
mod, err := pkg.FindAndParseGoModInCurrentDir()
if err != nil {
    log.Fatalf("Failed to find and parse go.mod: %v", err)
}

fmt.Printf("Found module: %s\n", mod.Name)

// Or specify a starting directory
mod, err = pkg.FindAndParseGoModFile("/path/to/project")
if err != nil {
    log.Fatalf("Failed to find go.mod: %v", err)
}
```

## Working with Dependencies

### Check if a Dependency Exists

```go
// Check if a specific dependency exists
if pkg.HasRequire(mod, "github.com/gin-gonic/gin") {
    fmt.Println("Project uses Gin framework")
    
    // Get detailed information about the dependency
    req := pkg.GetRequire(mod, "github.com/gin-gonic/gin")
    fmt.Printf("Version: %s\n", req.Version)
    fmt.Printf("Indirect: %v\n", req.Indirect)
}
```

### Analyze Replace Directives

```go
// Check for replace directives
if pkg.HasReplace(mod, "github.com/old/package") {
    replace := pkg.GetReplace(mod, "github.com/old/package")
    fmt.Printf("Package %s is replaced with %s %s\n", 
        replace.Old.Path, replace.New.Path, replace.New.Version)
}
```

### Check Excluded Packages

```go
// Check if a specific version is excluded
if pkg.HasExclude(mod, "github.com/problematic/package", "v1.0.0") {
    fmt.Println("Version v1.0.0 of problematic package is excluded")
}
```

### Check Retracted Versions

```go
// Check if a version is retracted
if pkg.HasRetract(mod, "v1.0.1") {
    fmt.Println("Version v1.0.1 has been retracted")
}
```

## Complete Example

Here's a complete example that demonstrates most features:

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/go-mod-parser/pkg"
)

func main() {
    // Parse go.mod file
    mod, err := pkg.FindAndParseGoModInCurrentDir()
    if err != nil {
        log.Fatalf("Error: %v", err)
    }
    
    // Print basic info
    fmt.Printf("📦 Module: %s\n", mod.Name)
    fmt.Printf("🐹 Go Version: %s\n", mod.GoVersion)
    
    // Analyze dependencies
    fmt.Printf("\n📋 Dependencies (%d):\n", len(mod.Requires))
    for _, req := range mod.Requires {
        status := "direct"
        if req.Indirect {
            status = "indirect"
        }
        fmt.Printf("  • %s %s (%s)\n", req.Path, req.Version, status)
    }
    
    // Show replace directives
    if len(mod.Replaces) > 0 {
        fmt.Printf("\n🔄 Replace Directives (%d):\n", len(mod.Replaces))
        for _, rep := range mod.Replaces {
            fmt.Printf("  • %s => %s %s\n", 
                rep.Old.Path, rep.New.Path, rep.New.Version)
        }
    }
    
    // Show excluded packages
    if len(mod.Excludes) > 0 {
        fmt.Printf("\n🚫 Excluded Packages (%d):\n", len(mod.Excludes))
        for _, exc := range mod.Excludes {
            fmt.Printf("  • %s %s\n", exc.Path, exc.Version)
        }
    }
    
    // Show retracted versions
    if len(mod.Retracts) > 0 {
        fmt.Printf("\n⚠️  Retracted Versions (%d):\n", len(mod.Retracts))
        for _, ret := range mod.Retracts {
            if ret.Version != "" {
                fmt.Printf("  • %s", ret.Version)
            } else {
                fmt.Printf("  • [%s, %s]", ret.VersionLow, ret.VersionHigh)
            }
            if ret.Rationale != "" {
                fmt.Printf(" (%s)", ret.Rationale)
            }
            fmt.Println()
        }
    }
}
```

## Next Steps

- Explore the [API Reference](/api/) for detailed documentation
- Check out more [Examples](/examples/) for advanced usage patterns
- Learn about [Data Structures](/api/data-structures) used by the library
