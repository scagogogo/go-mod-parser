# Data Structures

This page documents all the data structures used by Go Mod Parser to represent parsed go.mod file content.

## Module

The `Module` struct represents a complete go.mod file and contains all parsed directives.

```go
type Module struct {
    Name      string      // Module name from module directive
    GoVersion string      // Go version from go directive
    Requires  []*Require  // Dependencies from require directives
    Replaces  []*Replace  // Replace directives
    Excludes  []*Exclude  // Exclude directives
    Retracts  []*Retract  // Retract directives
}
```

### Fields

#### Name (string)
The module name as declared in the `module` directive.

**Example:**
```go
// For: module github.com/example/project
fmt.Println(mod.Name) // Output: github.com/example/project
```

#### GoVersion (string)
The Go version requirement from the `go` directive.

**Example:**
```go
// For: go 1.21
fmt.Println(mod.GoVersion) // Output: 1.21
```

#### Requires ([]*Require)
Slice of all dependencies declared in `require` directives.

#### Replaces ([]*Replace)
Slice of all module replacements from `replace` directives.

#### Excludes ([]*Exclude)
Slice of all excluded module versions from `exclude` directives.

#### Retracts ([]*Retract)
Slice of all retracted versions from `retract` directives.

---

## Require

Represents a single dependency declaration.

```go
type Require struct {
    Path     string // Module path
    Version  string // Module version
    Indirect bool   // Whether this is an indirect dependency
}
```

### Fields

#### Path (string)
The module path of the dependency.

**Example:**
```go
// For: require github.com/gin-gonic/gin v1.9.1
fmt.Println(req.Path) // Output: github.com/gin-gonic/gin
```

#### Version (string)
The version constraint for the dependency.

**Example:**
```go
// For: require github.com/gin-gonic/gin v1.9.1
fmt.Println(req.Version) // Output: v1.9.1
```

#### Indirect (bool)
Indicates whether this is an indirect dependency (marked with `// indirect` comment).

**Example:**
```go
// For: require github.com/example/pkg v1.0.0 // indirect
fmt.Println(req.Indirect) // Output: true

// For: require github.com/example/pkg v1.0.0
fmt.Println(req.Indirect) // Output: false
```

### Usage Example

```go
for _, req := range mod.Requires {
    status := "direct"
    if req.Indirect {
        status = "indirect"
    }
    fmt.Printf("%s %s (%s)\n", req.Path, req.Version, status)
}
```

---

## Replace

Represents a module replacement directive.

```go
type Replace struct {
    Old *ReplaceItem // Original module being replaced
    New *ReplaceItem // Replacement module
}
```

### Fields

#### Old (*ReplaceItem)
The original module that is being replaced.

#### New (*ReplaceItem)
The replacement module.

### Usage Example

```go
for _, rep := range mod.Replaces {
    fmt.Printf("Replace %s", rep.Old.Path)
    if rep.Old.Version != "" {
        fmt.Printf(" %s", rep.Old.Version)
    }
    fmt.Printf(" => %s", rep.New.Path)
    if rep.New.Version != "" {
        fmt.Printf(" %s", rep.New.Version)
    }
    fmt.Println()
}
```

---

## ReplaceItem

Represents a module reference in a replace directive.

```go
type ReplaceItem struct {
    Path    string // Module path
    Version string // Module version (may be empty)
}
```

### Fields

#### Path (string)
The module path.

#### Version (string)
The module version. May be empty for local path replacements.

### Examples

```go
// For: replace github.com/old/pkg => github.com/new/pkg v1.0.0
// Old: Path="github.com/old/pkg", Version=""
// New: Path="github.com/new/pkg", Version="v1.0.0"

// For: replace github.com/old/pkg v1.0.0 => ./local/path
// Old: Path="github.com/old/pkg", Version="v1.0.0"  
// New: Path="./local/path", Version=""
```

---

## Exclude

Represents an excluded module version.

```go
type Exclude struct {
    Path    string // Module path
    Version string // Excluded version
}
```

### Fields

#### Path (string)
The module path to exclude.

#### Version (string)
The specific version to exclude.

### Usage Example

```go
for _, exc := range mod.Excludes {
    fmt.Printf("Excluded: %s %s\n", exc.Path, exc.Version)
}

// Check if a specific version is excluded
for _, exc := range mod.Excludes {
    if exc.Path == "github.com/problematic/pkg" && exc.Version == "v1.0.0" {
        fmt.Println("Version v1.0.0 is excluded")
    }
}
```

---

## Retract

Represents a retracted version or version range.

```go
type Retract struct {
    Version     string // Single retracted version
    VersionLow  string // Low end of version range
    VersionHigh string // High end of version range  
    Rationale   string // Reason for retraction
}
```

### Fields

#### Version (string)
For single version retractions, this contains the version. Empty for version ranges.

#### VersionLow (string)
For version range retractions, this is the lower bound. Empty for single versions.

#### VersionHigh (string)
For version range retractions, this is the upper bound. Empty for single versions.

#### Rationale (string)
The reason for retraction, extracted from comments. May be empty.

### Usage Examples

```go
for _, ret := range mod.Retracts {
    if ret.Version != "" {
        // Single version retraction
        fmt.Printf("Retracted: %s", ret.Version)
    } else {
        // Version range retraction
        fmt.Printf("Retracted: [%s, %s]", ret.VersionLow, ret.VersionHigh)
    }
    
    if ret.Rationale != "" {
        fmt.Printf(" (%s)", ret.Rationale)
    }
    fmt.Println()
}
```

### Retraction Examples

```go
// Single version: retract v1.0.1 // security issue
// Version="v1.0.1", VersionLow="", VersionHigh="", Rationale="security issue"

// Version range: retract [v1.0.0, v1.0.5] // broken builds  
// Version="", VersionLow="v1.0.0", VersionHigh="v1.0.5", Rationale="broken builds"
```

---

## Complete Example

Here's a complete example showing how to work with all data structures:

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

replace (
    github.com/old/pkg => github.com/new/pkg v1.0.0
    github.com/local/pkg => ./local/path
)

exclude (
    github.com/bad/pkg v1.0.0
    github.com/another/bad v2.0.0
)

retract (
    v1.0.1 // security vulnerability
    [v1.0.2, v1.0.5] // broken builds
)
`

    mod, err := pkg.ParseGoModContent(content)
    if err != nil {
        log.Fatalf("Parse error: %v", err)
    }
    
    // Access basic info
    fmt.Printf("Module: %s\n", mod.Name)
    fmt.Printf("Go Version: %s\n", mod.GoVersion)
    
    // Process dependencies
    fmt.Printf("\nDependencies (%d):\n", len(mod.Requires))
    for _, req := range mod.Requires {
        fmt.Printf("  %s %s", req.Path, req.Version)
        if req.Indirect {
            fmt.Printf(" (indirect)")
        }
        fmt.Println()
    }
    
    // Process replacements
    fmt.Printf("\nReplacements (%d):\n", len(mod.Replaces))
    for _, rep := range mod.Replaces {
        fmt.Printf("  %s => %s", rep.Old.Path, rep.New.Path)
        if rep.New.Version != "" {
            fmt.Printf(" %s", rep.New.Version)
        }
        fmt.Println()
    }
    
    // Process exclusions
    fmt.Printf("\nExclusions (%d):\n", len(mod.Excludes))
    for _, exc := range mod.Excludes {
        fmt.Printf("  %s %s\n", exc.Path, exc.Version)
    }
    
    // Process retractions
    fmt.Printf("\nRetractions (%d):\n", len(mod.Retracts))
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
```
