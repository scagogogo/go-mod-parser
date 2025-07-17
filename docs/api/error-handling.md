# Error Handling

Go Mod Parser provides comprehensive error handling to help you diagnose and handle various failure scenarios when parsing go.mod files.

## Error Types

### File System Errors

These errors occur when there are issues accessing files or directories.

#### File Not Found

```go
mod, err := pkg.ParseGoModFile("/nonexistent/go.mod")
if err != nil {
    // Error: open /nonexistent/go.mod: no such file or directory
    fmt.Printf("Error: %v\n", err)
}
```

#### Permission Denied

```go
mod, err := pkg.ParseGoModFile("/root/go.mod")
if err != nil {
    // Error: open /root/go.mod: permission denied
    fmt.Printf("Error: %v\n", err)
}
```

#### Directory Not Found (Auto-discovery)

```go
mod, err := pkg.FindAndParseGoModFile("/nonexistent/directory")
if err != nil {
    // Error: go.mod file not found
    fmt.Printf("Error: %v\n", err)
}
```

### Parse Errors

These errors occur when the go.mod file content is malformed or contains invalid syntax.

#### Invalid Module Declaration

```go
content := `invalid module declaration`
mod, err := pkg.ParseGoModContent(content)
if err != nil {
    // Error: line 1: unrecognized line format: invalid module declaration
    fmt.Printf("Parse error: %v\n", err)
}
```

#### Invalid Require Format

```go
content := `module example.com/test

require invalid-require-format
`
mod, err := pkg.ParseGoModContent(content)
if err != nil {
    // Error: line 3: invalid require declaration
    fmt.Printf("Parse error: %v\n", err)
}
```

#### Invalid Replace Format

```go
content := `module example.com/test

replace invalid => format
`
mod, err := pkg.ParseGoModContent(content)
if err != nil {
    // Error: line 3: invalid replace declaration
    fmt.Printf("Parse error: %v\n", err)
}
```

### Block Parsing Errors

Errors that occur within block statements (require, replace, exclude, retract blocks).

```go
content := `module example.com/test

require (
    github.com/example/pkg
)
`
mod, err := pkg.ParseGoModContent(content)
if err != nil {
    // Error: line 4: invalid require declaration
    fmt.Printf("Block parse error: %v\n", err)
}
```

## Error Handling Patterns

### Basic Error Checking

Always check for errors when calling parsing functions:

```go
mod, err := pkg.ParseGoModFile("go.mod")
if err != nil {
    log.Fatalf("Failed to parse go.mod: %v", err)
}
// Continue with mod...
```

### Graceful Error Handling

```go
func parseGoModSafely(path string) (*module.Module, error) {
    mod, err := pkg.ParseGoModFile(path)
    if err != nil {
        // Log the error but don't crash
        log.Printf("Warning: Failed to parse %s: %v", path, err)
        return nil, err
    }
    return mod, nil
}
```

### Error Type Detection

```go
func handleParseError(err error) {
    errStr := err.Error()
    
    switch {
    case strings.Contains(errStr, "no such file"):
        fmt.Println("File not found - check the path")
    case strings.Contains(errStr, "permission denied"):
        fmt.Println("Permission denied - check file permissions")
    case strings.Contains(errStr, "go.mod file not found"):
        fmt.Println("No go.mod found in directory tree")
    case strings.Contains(errStr, "line"):
        fmt.Println("Parse error - check go.mod syntax")
    default:
        fmt.Printf("Unknown error: %v\n", err)
    }
}
```

### Fallback Strategies

```go
func parseWithFallbacks(primaryPath, fallbackPath string) (*module.Module, error) {
    // Try primary path first
    mod, err := pkg.ParseGoModFile(primaryPath)
    if err == nil {
        return mod, nil
    }
    
    log.Printf("Primary parse failed: %v", err)
    
    // Try fallback path
    mod, err = pkg.ParseGoModFile(fallbackPath)
    if err == nil {
        log.Printf("Successfully parsed fallback: %s", fallbackPath)
        return mod, nil
    }
    
    // Try auto-discovery as last resort
    mod, err = pkg.FindAndParseGoModInCurrentDir()
    if err == nil {
        log.Println("Successfully auto-discovered go.mod")
        return mod, nil
    }
    
    return nil, fmt.Errorf("all parsing attempts failed")
}
```

### Validation After Parsing

```go
func validateParsedModule(mod *module.Module) error {
    if mod.Name == "" {
        return fmt.Errorf("module name is empty")
    }
    
    if mod.GoVersion == "" {
        return fmt.Errorf("go version is not specified")
    }
    
    // Check for duplicate dependencies
    seen := make(map[string]bool)
    for _, req := range mod.Requires {
        if seen[req.Path] {
            return fmt.Errorf("duplicate dependency: %s", req.Path)
        }
        seen[req.Path] = true
    }
    
    return nil
}

func parseAndValidate(path string) (*module.Module, error) {
    mod, err := pkg.ParseGoModFile(path)
    if err != nil {
        return nil, fmt.Errorf("parse error: %w", err)
    }
    
    if err := validateParsedModule(mod); err != nil {
        return nil, fmt.Errorf("validation error: %w", err)
    }
    
    return mod, nil
}
```

## Robust Parsing Function

Here's a comprehensive example that handles multiple error scenarios:

```go
package main

import (
    "fmt"
    "log"
    "os"
    "path/filepath"
    "strings"
    
    "github.com/scagogogo/go-mod-parser/pkg"
    "github.com/scagogogo/go-mod-parser/pkg/module"
)

func robustParseGoMod(path string) (*module.Module, error) {
    // Normalize path
    absPath, err := filepath.Abs(path)
    if err != nil {
        return nil, fmt.Errorf("failed to resolve path %s: %w", path, err)
    }
    
    // Check if file exists
    if _, err := os.Stat(absPath); os.IsNotExist(err) {
        return nil, fmt.Errorf("go.mod file not found: %s", absPath)
    } else if err != nil {
        return nil, fmt.Errorf("failed to access file %s: %w", absPath, err)
    }
    
    // Attempt to parse
    mod, err := pkg.ParseGoModFile(absPath)
    if err != nil {
        return nil, fmt.Errorf("failed to parse %s: %w", absPath, err)
    }
    
    // Basic validation
    if mod.Name == "" {
        return nil, fmt.Errorf("parsed module has empty name")
    }
    
    return mod, nil
}

func parseWithRetry(paths []string) (*module.Module, error) {
    var lastErr error
    
    for _, path := range paths {
        mod, err := robustParseGoMod(path)
        if err == nil {
            return mod, nil
        }
        
        log.Printf("Failed to parse %s: %v", path, err)
        lastErr = err
    }
    
    return nil, fmt.Errorf("all parsing attempts failed, last error: %w", lastErr)
}

func main() {
    // Try multiple possible locations
    candidates := []string{
        "go.mod",
        "./go.mod", 
        "../go.mod",
        "../../go.mod",
    }
    
    mod, err := parseWithRetry(candidates)
    if err != nil {
        // Final fallback: auto-discovery
        mod, err = pkg.FindAndParseGoModInCurrentDir()
        if err != nil {
            log.Fatalf("Could not find or parse any go.mod file: %v", err)
        }
        log.Println("Successfully used auto-discovery")
    }
    
    fmt.Printf("Successfully parsed module: %s\n", mod.Name)
}
```

## Error Prevention

### Input Validation

```go
func validateInputs(path string) error {
    if path == "" {
        return fmt.Errorf("path cannot be empty")
    }
    
    if !strings.HasSuffix(path, "go.mod") && !isDirectory(path) {
        return fmt.Errorf("path must be a go.mod file or directory")
    }
    
    return nil
}

func isDirectory(path string) bool {
    info, err := os.Stat(path)
    return err == nil && info.IsDir()
}
```

### Content Validation

```go
func validateGoModContent(content string) error {
    if content == "" {
        return fmt.Errorf("content cannot be empty")
    }
    
    if !strings.Contains(content, "module ") {
        return fmt.Errorf("content does not appear to be a valid go.mod file")
    }
    
    return nil
}

func parseContentSafely(content string) (*module.Module, error) {
    if err := validateGoModContent(content); err != nil {
        return nil, err
    }
    
    return pkg.ParseGoModContent(content)
}
```

## Best Practices

1. **Always Check Errors**: Never ignore errors from parsing functions
2. **Provide Context**: Wrap errors with additional context using `fmt.Errorf`
3. **Log Appropriately**: Use appropriate log levels for different error types
4. **Implement Fallbacks**: Have backup strategies for critical parsing operations
5. **Validate Results**: Check parsed data for consistency and completeness
6. **Handle Edge Cases**: Consider empty files, permission issues, and network problems

## Related Documentation

- [Core Functions](/api/core-functions) - Functions that can return errors
- [Data Structures](/api/data-structures) - Understanding the parsed data
- [Examples](/examples/) - See error handling in practice
