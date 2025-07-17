# Basic Parsing

This section demonstrates the fundamental parsing operations with Go Mod Parser.

## Simple File Parsing

The most basic operation is parsing a go.mod file from disk:

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/go-mod-parser/pkg"
)

func main() {
    // Parse go.mod file
    mod, err := pkg.ParseGoModFile("go.mod")
    if err != nil {
        log.Fatalf("Failed to parse go.mod: %v", err)
    }
    
    // Display basic information
    fmt.Printf("Module: %s\n", mod.Name)
    fmt.Printf("Go Version: %s\n", mod.GoVersion)
    fmt.Printf("Dependencies: %d\n", len(mod.Requires))
}
```

## Content Parsing

Parse go.mod content directly from a string:

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
`
    
    mod, err := pkg.ParseGoModContent(content)
    if err != nil {
        log.Fatalf("Failed to parse content: %v", err)
    }
    
    fmt.Printf("Module: %s\n", mod.Name)
    fmt.Printf("Go Version: %s\n", mod.GoVersion)
    
    for _, req := range mod.Requires {
        fmt.Printf("Dependency: %s %s\n", req.Path, req.Version)
        if req.Indirect {
            fmt.Println("  (indirect)")
        }
    }
}
```

## Error Handling

Always handle errors when parsing:

```go
package main

import (
    "fmt"
    "log"
    "os"
    
    "github.com/scagogogo/go-mod-parser/pkg"
)

func main() {
    mod, err := pkg.ParseGoModFile("go.mod")
    if err != nil {
        if os.IsNotExist(err) {
            log.Fatal("go.mod file not found")
        } else {
            log.Fatalf("Parse error: %v", err)
        }
    }
    
    fmt.Printf("Successfully parsed: %s\n", mod.Name)
}
```

## Complete Example

Here's a complete example that demonstrates basic parsing with comprehensive output:

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/go-mod-parser/pkg"
)

func main() {
    // Parse the go.mod file
    mod, err := pkg.ParseGoModFile("go.mod")
    if err != nil {
        log.Fatalf("Error parsing go.mod: %v", err)
    }
    
    // Print module information
    fmt.Printf("📦 Module: %s\n", mod.Name)
    fmt.Printf("🐹 Go Version: %s\n", mod.GoVersion)
    
    // Print dependencies
    if len(mod.Requires) > 0 {
        fmt.Printf("\n📋 Dependencies (%d):\n", len(mod.Requires))
        for i, req := range mod.Requires {
            fmt.Printf("%d. %s %s", i+1, req.Path, req.Version)
            if req.Indirect {
                fmt.Printf(" (indirect)")
            }
            fmt.Println()
        }
    } else {
        fmt.Println("\n📋 No dependencies found")
    }
    
    // Print replace directives
    if len(mod.Replaces) > 0 {
        fmt.Printf("\n🔄 Replace Directives (%d):\n", len(mod.Replaces))
        for i, rep := range mod.Replaces {
            fmt.Printf("%d. %s => %s", i+1, rep.Old.Path, rep.New.Path)
            if rep.New.Version != "" {
                fmt.Printf(" %s", rep.New.Version)
            }
            fmt.Println()
        }
    }
    
    // Print exclude directives
    if len(mod.Excludes) > 0 {
        fmt.Printf("\n🚫 Exclude Directives (%d):\n", len(mod.Excludes))
        for i, exc := range mod.Excludes {
            fmt.Printf("%d. %s %s\n", i+1, exc.Path, exc.Version)
        }
    }
    
    // Print retract directives
    if len(mod.Retracts) > 0 {
        fmt.Printf("\n⚠️  Retract Directives (%d):\n", len(mod.Retracts))
        for i, ret := range mod.Retracts {
            fmt.Printf("%d. ", i+1)
            if ret.Version != "" {
                fmt.Printf("%s", ret.Version)
            } else {
                fmt.Printf("[%s, %s]", ret.VersionLow, ret.VersionHigh)
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

- [File Discovery](/examples/file-discovery) - Learn about auto-discovery features
- [Dependency Analysis](/examples/dependency-analysis) - Analyze dependencies in detail
- [Advanced Usage](/examples/advanced-usage) - Complex scenarios and patterns
