# Core Functions

The core functions provide the primary interface for parsing go.mod files from various sources.

## ParseGoModFile

Parse a go.mod file from a file path.

```go
func ParseGoModFile(path string) (*module.Module, error)
```

### Parameters

- `path` (string): The file path to the go.mod file

### Returns

- `*module.Module`: Parsed module data
- `error`: Error if parsing fails

### Example

```go
mod, err := pkg.ParseGoModFile("/path/to/go.mod")
if err != nil {
    log.Fatalf("Failed to parse: %v", err)
}

fmt.Printf("Module: %s\n", mod.Name)
fmt.Printf("Go Version: %s\n", mod.GoVersion)
```

### Error Conditions

- File does not exist
- Permission denied
- Invalid go.mod syntax
- I/O errors

---

## ParseGoModContent

Parse go.mod content from a string.

```go
func ParseGoModContent(content string) (*module.Module, error)
```

### Parameters

- `content` (string): The go.mod file content as a string

### Returns

- `*module.Module`: Parsed module data
- `error`: Error if parsing fails

### Example

```go
content := `module github.com/example/project

go 1.21

require (
    github.com/gin-gonic/gin v1.9.1
    github.com/stretchr/testify v1.8.4 // indirect
)

replace github.com/old/pkg => github.com/new/pkg v1.0.0
exclude github.com/bad/pkg v1.0.0
retract v1.0.1 // security issue
`

mod, err := pkg.ParseGoModContent(content)
if err != nil {
    log.Fatalf("Failed to parse: %v", err)
}

fmt.Printf("Module: %s\n", mod.Name)
```

### Error Conditions

- Invalid go.mod syntax
- Malformed directives
- Unsupported go.mod features

---

## FindAndParseGoModFile

Find and parse a go.mod file starting from a specified directory, searching upward through parent directories.

```go
func FindAndParseGoModFile(dir string) (*module.Module, error)
```

### Parameters

- `dir` (string): Starting directory for the search

### Returns

- `*module.Module`: Parsed module data
- `error`: Error if no go.mod found or parsing fails

### Example

```go
// Search starting from a specific directory
mod, err := pkg.FindAndParseGoModFile("/path/to/project/subdir")
if err != nil {
    log.Fatalf("Failed to find go.mod: %v", err)
}

fmt.Printf("Found module: %s\n", mod.Name)
```

### Behavior

1. Starts searching in the specified directory
2. Looks for `go.mod` file in current directory
3. If not found, moves to parent directory
4. Continues until go.mod is found or root directory is reached
5. Returns error if no go.mod file is found

### Error Conditions

- No go.mod file found in directory tree
- Permission denied
- Invalid go.mod syntax
- I/O errors

---

## FindAndParseGoModInCurrentDir

Find and parse a go.mod file starting from the current working directory.

```go
func FindAndParseGoModInCurrentDir() (*module.Module, error)
```

### Parameters

None

### Returns

- `*module.Module`: Parsed module data
- `error`: Error if no go.mod found or parsing fails

### Example

```go
// Search starting from current directory
mod, err := pkg.FindAndParseGoModInCurrentDir()
if err != nil {
    log.Fatalf("Failed to find go.mod: %v", err)
}

fmt.Printf("Current project module: %s\n", mod.Name)
```

### Behavior

This is equivalent to calling `FindAndParseGoModFile("")` with an empty string, which uses the current working directory as the starting point.

### Error Conditions

- No go.mod file found in current directory tree
- Unable to determine current working directory
- Permission denied
- Invalid go.mod syntax
- I/O errors

---

## Advanced Usage

### Handling Different Input Sources

```go
package main

import (
    "fmt"
    "log"
    "os"
    
    "github.com/scagogogo/go-mod-parser/pkg"
)

func parseFromDifferentSources() {
    // Method 1: Parse from file path
    if mod, err := pkg.ParseGoModFile("go.mod"); err == nil {
        fmt.Printf("From file: %s\n", mod.Name)
    }
    
    // Method 2: Parse from content
    content, _ := os.ReadFile("go.mod")
    if mod, err := pkg.ParseGoModContent(string(content)); err == nil {
        fmt.Printf("From content: %s\n", mod.Name)
    }
    
    // Method 3: Auto-discover from current directory
    if mod, err := pkg.FindAndParseGoModInCurrentDir(); err == nil {
        fmt.Printf("Auto-discovered: %s\n", mod.Name)
    }
    
    // Method 4: Auto-discover from specific directory
    if mod, err := pkg.FindAndParseGoModFile("/path/to/project"); err == nil {
        fmt.Printf("Found in project: %s\n", mod.Name)
    }
}
```

### Error Handling Patterns

```go
func robustParsing() {
    // Try multiple methods with fallbacks
    var mod *module.Module
    var err error
    
    // Try current directory first
    mod, err = pkg.FindAndParseGoModInCurrentDir()
    if err != nil {
        // Fallback to specific file
        mod, err = pkg.ParseGoModFile("go.mod")
        if err != nil {
            log.Fatalf("Could not parse go.mod: %v", err)
        }
    }
    
    fmt.Printf("Successfully parsed: %s\n", mod.Name)
}
```

### Performance Considerations

```go
func efficientParsing() {
    // For repeated parsing of the same content, use ParseGoModContent
    content := `module example.com/project
go 1.21
require github.com/gin-gonic/gin v1.9.1`
    
    // This avoids file I/O overhead
    mod, err := pkg.ParseGoModContent(content)
    if err != nil {
        log.Fatalf("Parse error: %v", err)
    }
    
    // For file-based parsing, consider caching results
    // if you need to parse the same file multiple times
}
```

## Related Functions

- [Helper Functions](/api/helper-functions) - Functions for analyzing parsed data
- [Data Structures](/api/data-structures) - Details about the Module type
- [Error Handling](/api/error-handling) - Error types and handling strategies
