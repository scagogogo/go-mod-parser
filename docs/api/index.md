# API Reference

Go Mod Parser provides a comprehensive API for parsing and analyzing go.mod files. The library is organized into several packages, each serving a specific purpose.

## Package Overview

### Main Package (`pkg`)

The main package provides the primary API that most users will interact with. It offers high-level functions for parsing go.mod files and analyzing their contents.

```go
import "github.com/scagogogo/go-mod-parser/pkg"
```

**Key Functions:**
- `ParseGoModFile(path string)` - Parse a go.mod file from disk
- `ParseGoModContent(content string)` - Parse go.mod content from string
- `FindAndParseGoModFile(dir string)` - Auto-discover and parse go.mod files
- `HasRequire(mod, path)` - Check if a dependency exists
- `GetRequire(mod, path)` - Get dependency details

### Module Package (`pkg/module`)

Defines the data structures used to represent parsed go.mod content.

```go
import "github.com/scagogogo/go-mod-parser/pkg/module"
```

**Key Types:**
- `Module` - Represents a complete go.mod file
- `Require` - Represents a dependency
- `Replace` - Represents a replace directive
- `Exclude` - Represents an exclude directive
- `Retract` - Represents a retract directive

## API Categories

### [Core Functions](/api/core-functions)
Primary parsing functions for different input sources:
- File-based parsing
- String-based parsing  
- Auto-discovery parsing

### [Data Structures](/api/data-structures)
Complete reference for all data types:
- Module structure
- Dependency types
- Directive representations

### [Helper Functions](/api/helper-functions)
Utility functions for analyzing parsed data:
- Dependency checking
- Replace directive analysis
- Exclude and retract validation

### [Error Handling](/api/error-handling)
Error types and handling patterns:
- Parse errors
- File system errors
- Validation errors

## Quick Reference

| Function | Description | Returns |
|----------|-------------|---------|
| `ParseGoModFile(path)` | Parse go.mod file from path | `(*Module, error)` |
| `ParseGoModContent(content)` | Parse go.mod from string | `(*Module, error)` |
| `FindAndParseGoModFile(dir)` | Find and parse go.mod in directory | `(*Module, error)` |
| `FindAndParseGoModInCurrentDir()` | Find and parse go.mod in current dir | `(*Module, error)` |
| `HasRequire(mod, path)` | Check if dependency exists | `bool` |
| `GetRequire(mod, path)` | Get dependency details | `*Require` |
| `HasReplace(mod, path)` | Check if replace exists | `bool` |
| `GetReplace(mod, path)` | Get replace details | `*Replace` |
| `HasExclude(mod, path, version)` | Check if version excluded | `bool` |
| `HasRetract(mod, version)` | Check if version retracted | `bool` |

## Usage Patterns

### Basic Parsing

```go
// Parse from file
mod, err := pkg.ParseGoModFile("go.mod")

// Parse from string
mod, err := pkg.ParseGoModContent(content)

// Auto-discover
mod, err := pkg.FindAndParseGoModInCurrentDir()
```

### Dependency Analysis

```go
// Check dependencies
if pkg.HasRequire(mod, "github.com/gin-gonic/gin") {
    req := pkg.GetRequire(mod, "github.com/gin-gonic/gin")
    fmt.Printf("Version: %s, Indirect: %v\n", req.Version, req.Indirect)
}
```

### Replace Directive Analysis

```go
// Check replacements
if pkg.HasReplace(mod, "github.com/old/pkg") {
    rep := pkg.GetReplace(mod, "github.com/old/pkg")
    fmt.Printf("Replaced with: %s %s\n", rep.New.Path, rep.New.Version)
}
```

## Error Handling

All parsing functions return an error as the second return value. Always check for errors:

```go
mod, err := pkg.ParseGoModFile("go.mod")
if err != nil {
    // Handle error appropriately
    log.Fatalf("Failed to parse go.mod: %v", err)
}
```

Common error scenarios:
- File not found
- Invalid go.mod syntax
- Permission errors
- Malformed directives

## Thread Safety

The library is designed to be thread-safe for read operations. The returned `Module` struct and its fields are safe to read concurrently from multiple goroutines. However, if you need to modify the parsed data, you should implement your own synchronization.

## Performance Considerations

- Parsing is generally fast, but file I/O can be the bottleneck
- Consider caching parsed results for frequently accessed files
- Use `ParseGoModContent` for in-memory content to avoid file system overhead
- The auto-discovery functions may traverse multiple directories

## Next Steps

- [Core Functions](/api/core-functions) - Detailed function documentation
- [Data Structures](/api/data-structures) - Complete type reference
- [Helper Functions](/api/helper-functions) - Analysis utilities
- [Examples](/examples/) - Practical usage examples
