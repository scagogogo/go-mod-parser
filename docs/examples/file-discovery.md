# File Discovery

Go Mod Parser provides powerful auto-discovery features to automatically locate go.mod files in your project structure.

## Auto-Discovery from Current Directory

The simplest way to find and parse a go.mod file:

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/go-mod-parser/pkg"
)

func main() {
    // Find and parse go.mod starting from current directory
    mod, err := pkg.FindAndParseGoModInCurrentDir()
    if err != nil {
        log.Fatalf("Failed to find go.mod: %v", err)
    }
    
    fmt.Printf("Found module: %s\n", mod.Name)
}
```

## Auto-Discovery from Specific Directory

Start the search from a specific directory:

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/go-mod-parser/pkg"
)

func main() {
    // Find and parse go.mod starting from a specific directory
    projectDir := "/path/to/project/subdirectory"
    mod, err := pkg.FindAndParseGoModFile(projectDir)
    if err != nil {
        log.Fatalf("Failed to find go.mod in %s: %v", projectDir, err)
    }
    
    fmt.Printf("Found module: %s\n", mod.Name)
}
```

## How Auto-Discovery Works

The auto-discovery process:

1. Starts in the specified directory (or current directory)
2. Looks for a `go.mod` file in the current directory
3. If not found, moves to the parent directory
4. Repeats until a go.mod file is found or the root directory is reached
5. Returns an error if no go.mod file is found

## Robust Discovery with Fallbacks

```go
package main

import (
    "fmt"
    "log"
    "os"
    "path/filepath"
    
    "github.com/scagogogo/go-mod-parser/pkg"
)

func findGoModWithFallbacks() (*module.Module, error) {
    // Try current directory first
    mod, err := pkg.FindAndParseGoModInCurrentDir()
    if err == nil {
        fmt.Println("Found go.mod in current directory tree")
        return mod, nil
    }
    
    // Try specific paths as fallbacks
    fallbackPaths := []string{
        ".",
        "..",
        "../..",
        filepath.Join(os.Getenv("HOME"), "go", "src", "myproject"),
    }
    
    for _, path := range fallbackPaths {
        fmt.Printf("Trying fallback path: %s\n", path)
        mod, err := pkg.FindAndParseGoModFile(path)
        if err == nil {
            fmt.Printf("Found go.mod at: %s\n", path)
            return mod, nil
        }
    }
    
    return nil, fmt.Errorf("no go.mod file found in any location")
}

func main() {
    mod, err := findGoModWithFallbacks()
    if err != nil {
        log.Fatalf("Discovery failed: %v", err)
    }
    
    fmt.Printf("Successfully found module: %s\n", mod.Name)
}
```

## Working with Monorepos

In monorepo scenarios, you might have multiple go.mod files:

```go
package main

import (
    "fmt"
    "log"
    "os"
    "path/filepath"
    
    "github.com/scagogogo/go-mod-parser/pkg"
)

func findAllGoModFiles(rootDir string) (map[string]*module.Module, error) {
    modules := make(map[string]*module.Module)
    
    err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        
        if info.Name() == "go.mod" {
            mod, parseErr := pkg.ParseGoModFile(path)
            if parseErr != nil {
                fmt.Printf("Warning: Failed to parse %s: %v\n", path, parseErr)
                return nil // Continue walking
            }
            
            modules[path] = mod
            fmt.Printf("Found module: %s at %s\n", mod.Name, path)
        }
        
        return nil
    })
    
    return modules, err
}

func main() {
    rootDir := "." // or specify your monorepo root
    
    modules, err := findAllGoModFiles(rootDir)
    if err != nil {
        log.Fatalf("Error walking directory: %v", err)
    }
    
    fmt.Printf("\nFound %d modules:\n", len(modules))
    for path, mod := range modules {
        fmt.Printf("- %s: %s\n", mod.Name, path)
    }
}
```

## Discovery with Validation

Validate discovered go.mod files:

```go
package main

import (
    "fmt"
    "log"
    "strings"
    
    "github.com/scagogogo/go-mod-parser/pkg"
)

func validateModule(mod *module.Module) []string {
    var issues []string
    
    if mod.Name == "" {
        issues = append(issues, "module name is empty")
    }
    
    if mod.GoVersion == "" {
        issues = append(issues, "go version not specified")
    }
    
    if !strings.Contains(mod.Name, ".") {
        issues = append(issues, "module name should contain domain")
    }
    
    return issues
}

func discoverAndValidate() {
    mod, err := pkg.FindAndParseGoModInCurrentDir()
    if err != nil {
        log.Fatalf("Discovery failed: %v", err)
    }
    
    fmt.Printf("Discovered module: %s\n", mod.Name)
    
    issues := validateModule(mod)
    if len(issues) > 0 {
        fmt.Println("\nValidation issues:")
        for _, issue := range issues {
            fmt.Printf("- %s\n", issue)
        }
    } else {
        fmt.Println("✅ Module validation passed")
    }
}

func main() {
    discoverAndValidate()
}
```

## Performance Considerations

For better performance in large directory trees:

```go
package main

import (
    "fmt"
    "log"
    "os"
    "path/filepath"
    "time"
    
    "github.com/scagogogo/go-mod-parser/pkg"
)

func timedDiscovery(startDir string) {
    start := time.Now()
    
    mod, err := pkg.FindAndParseGoModFile(startDir)
    
    duration := time.Since(start)
    
    if err != nil {
        log.Printf("Discovery failed in %v: %v", duration, err)
        return
    }
    
    fmt.Printf("Found %s in %v\n", mod.Name, duration)
}

func efficientDiscovery(startDir string, maxDepth int) (*module.Module, error) {
    currentDir := startDir
    depth := 0
    
    for depth < maxDepth {
        goModPath := filepath.Join(currentDir, "go.mod")
        
        if _, err := os.Stat(goModPath); err == nil {
            return pkg.ParseGoModFile(goModPath)
        }
        
        parentDir := filepath.Dir(currentDir)
        if parentDir == currentDir {
            // Reached root directory
            break
        }
        
        currentDir = parentDir
        depth++
    }
    
    return nil, fmt.Errorf("no go.mod found within %d levels", maxDepth)
}

func main() {
    // Standard discovery
    fmt.Println("Standard discovery:")
    timedDiscovery(".")
    
    // Limited depth discovery
    fmt.Println("\nLimited depth discovery:")
    mod, err := efficientDiscovery(".", 5)
    if err != nil {
        log.Printf("Limited discovery failed: %v", err)
    } else {
        fmt.Printf("Found with limited search: %s\n", mod.Name)
    }
}
```

## Next Steps

- [Dependency Analysis](/examples/dependency-analysis) - Analyze discovered modules
- [Advanced Usage](/examples/advanced-usage) - Complex discovery patterns
- [Basic Parsing](/examples/basic-parsing) - Learn basic parsing operations
