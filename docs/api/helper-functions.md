# Helper Functions

Helper functions provide convenient ways to analyze and query parsed go.mod data. These functions make it easy to check for specific dependencies, replacements, exclusions, and retractions.

## Dependency Functions

### HasRequire

Check if a module has a specific dependency.

```go
func HasRequire(mod *module.Module, path string) bool
```

#### Parameters

- `mod` (*module.Module): The parsed module
- `path` (string): The dependency path to check

#### Returns

- `bool`: True if the dependency exists, false otherwise

#### Example

```go
if pkg.HasRequire(mod, "github.com/gin-gonic/gin") {
    fmt.Println("Project uses Gin framework")
} else {
    fmt.Println("Gin framework not found")
}
```

---

### GetRequire

Get detailed information about a specific dependency.

```go
func GetRequire(mod *module.Module, path string) *module.Require
```

#### Parameters

- `mod` (*module.Module): The parsed module
- `path` (string): The dependency path to retrieve

#### Returns

- `*module.Require`: The dependency details, or nil if not found

#### Example

```go
req := pkg.GetRequire(mod, "github.com/gin-gonic/gin")
if req != nil {
    fmt.Printf("Version: %s\n", req.Version)
    fmt.Printf("Indirect: %v\n", req.Indirect)
} else {
    fmt.Println("Dependency not found")
}
```

#### Safe Usage Pattern

```go
// Always check if dependency exists first
if pkg.HasRequire(mod, "github.com/gin-gonic/gin") {
    req := pkg.GetRequire(mod, "github.com/gin-gonic/gin")
    // req is guaranteed to be non-nil here
    fmt.Printf("Found Gin version: %s\n", req.Version)
}

// Or handle nil case
if req := pkg.GetRequire(mod, "github.com/gin-gonic/gin"); req != nil {
    fmt.Printf("Found Gin version: %s\n", req.Version)
}
```

---

## Replace Functions

### HasReplace

Check if a module has a replacement for a specific path.

```go
func HasReplace(mod *module.Module, path string) bool
```

#### Parameters

- `mod` (*module.Module): The parsed module
- `path` (string): The original module path to check

#### Returns

- `bool`: True if a replacement exists, false otherwise

#### Example

```go
if pkg.HasReplace(mod, "github.com/old/package") {
    fmt.Println("Package has been replaced")
}
```

---

### GetReplace

Get replacement details for a specific module path.

```go
func GetReplace(mod *module.Module, path string) *module.Replace
```

#### Parameters

- `mod` (*module.Module): The parsed module
- `path` (string): The original module path

#### Returns

- `*module.Replace`: The replacement details, or nil if not found

#### Example

```go
if rep := pkg.GetReplace(mod, "github.com/old/package"); rep != nil {
    fmt.Printf("Replaced with: %s", rep.New.Path)
    if rep.New.Version != "" {
        fmt.Printf(" %s", rep.New.Version)
    }
    fmt.Println()
}
```

#### Analyzing Replacement Types

```go
rep := pkg.GetReplace(mod, "github.com/old/package")
if rep != nil {
    if strings.HasPrefix(rep.New.Path, "./") || strings.HasPrefix(rep.New.Path, "../") {
        fmt.Println("Local path replacement")
    } else {
        fmt.Println("Module replacement")
        if rep.New.Version != "" {
            fmt.Printf("Replacement version: %s\n", rep.New.Version)
        }
    }
}
```

---

## Exclude Functions

### HasExclude

Check if a specific module version is excluded.

```go
func HasExclude(mod *module.Module, path, version string) bool
```

#### Parameters

- `mod` (*module.Module): The parsed module
- `path` (string): The module path
- `version` (string): The version to check

#### Returns

- `bool`: True if the version is excluded, false otherwise

#### Example

```go
if pkg.HasExclude(mod, "github.com/problematic/pkg", "v1.0.0") {
    fmt.Println("Version v1.0.0 is excluded")
}
```

#### Checking Multiple Versions

```go
problematicVersions := []string{"v1.0.0", "v1.0.1", "v1.0.2"}
for _, version := range problematicVersions {
    if pkg.HasExclude(mod, "github.com/problematic/pkg", version) {
        fmt.Printf("Version %s is excluded\n", version)
    }
}
```

---

## Retract Functions

### HasRetract

Check if a specific version has been retracted.

```go
func HasRetract(mod *module.Module, version string) bool
```

#### Parameters

- `mod` (*module.Module): The parsed module
- `version` (string): The version to check

#### Returns

- `bool`: True if the version is retracted, false otherwise

#### Example

```go
if pkg.HasRetract(mod, "v1.0.1") {
    fmt.Println("Version v1.0.1 has been retracted")
}
```

#### Note on Version Ranges

This function checks both single version retractions and version ranges. If a version falls within a retracted range, it will return true.

```go
// For retract [v1.0.0, v1.0.5]
fmt.Println(pkg.HasRetract(mod, "v1.0.2")) // true
fmt.Println(pkg.HasRetract(mod, "v1.0.6")) // false
```

---

## Advanced Usage Patterns

### Comprehensive Dependency Analysis

```go
func analyzeDependencies(mod *module.Module) {
    fmt.Printf("Analyzing %d dependencies:\n", len(mod.Requires))
    
    for _, req := range mod.Requires {
        fmt.Printf("\n📦 %s %s", req.Path, req.Version)
        
        // Check if it's indirect
        if req.Indirect {
            fmt.Printf(" (indirect)")
        }
        
        // Check if it has a replacement
        if pkg.HasReplace(mod, req.Path) {
            rep := pkg.GetReplace(mod, req.Path)
            fmt.Printf("\n   🔄 Replaced with: %s", rep.New.Path)
            if rep.New.Version != "" {
                fmt.Printf(" %s", rep.New.Version)
            }
        }
        
        // Check if any version is excluded
        excluded := false
        for _, exc := range mod.Excludes {
            if exc.Path == req.Path {
                fmt.Printf("\n   🚫 Version %s is excluded", exc.Version)
                excluded = true
            }
        }
        
        fmt.Println()
    }
}
```

### Validation Functions

```go
func validateModule(mod *module.Module) []string {
    var issues []string
    
    // Check for dependencies that are also replaced
    for _, req := range mod.Requires {
        if pkg.HasReplace(mod, req.Path) {
            rep := pkg.GetReplace(mod, req.Path)
            issues = append(issues, 
                fmt.Sprintf("Dependency %s is replaced with %s", 
                    req.Path, rep.New.Path))
        }
    }
    
    // Check for retracted versions in use
    for _, req := range mod.Requires {
        if pkg.HasRetract(mod, req.Version) {
            issues = append(issues, 
                fmt.Sprintf("Using retracted version %s of %s", 
                    req.Version, req.Path))
        }
    }
    
    return issues
}
```

### Dependency Filtering

```go
func filterDependencies(mod *module.Module, filter func(*module.Require) bool) []*module.Require {
    var filtered []*module.Require
    for _, req := range mod.Requires {
        if filter(req) {
            filtered = append(filtered, req)
        }
    }
    return filtered
}

// Usage examples
directDeps := filterDependencies(mod, func(req *module.Require) bool {
    return !req.Indirect
})

testDeps := filterDependencies(mod, func(req *module.Require) bool {
    return strings.Contains(req.Path, "test")
})

replacedDeps := filterDependencies(mod, func(req *module.Require) bool {
    return pkg.HasReplace(mod, req.Path)
})
```

### Security Analysis

```go
func checkSecurity(mod *module.Module) {
    fmt.Println("Security Analysis:")
    
    // Check for retracted versions
    retractedCount := 0
    for _, req := range mod.Requires {
        if pkg.HasRetract(mod, req.Version) {
            fmt.Printf("⚠️  Using retracted version: %s %s\n", 
                req.Path, req.Version)
            retractedCount++
        }
    }
    
    // Check for local replacements (potential security risk)
    localReplacements := 0
    for _, rep := range mod.Replaces {
        if strings.HasPrefix(rep.New.Path, "./") || 
           strings.HasPrefix(rep.New.Path, "../") {
            fmt.Printf("🔍 Local replacement: %s => %s\n", 
                rep.Old.Path, rep.New.Path)
            localReplacements++
        }
    }
    
    fmt.Printf("\nSummary: %d retracted versions, %d local replacements\n", 
        retractedCount, localReplacements)
}
```

## Performance Tips

1. **Cache Results**: If you're calling these functions multiple times with the same parameters, consider caching the results.

2. **Batch Checks**: Instead of calling `HasRequire` followed by `GetRequire`, just call `GetRequire` and check for nil.

3. **Early Returns**: Use early returns in your analysis functions to avoid unnecessary processing.

```go
// Efficient pattern
if req := pkg.GetRequire(mod, path); req != nil {
    // Process req
    return req.Version
}

// Less efficient
if pkg.HasRequire(mod, path) {
    req := pkg.GetRequire(mod, path)
    return req.Version
}
```

## Related Documentation

- [Core Functions](/api/core-functions) - Functions for parsing go.mod files
- [Data Structures](/api/data-structures) - Details about the data types
- [Examples](/examples/) - Practical usage examples
