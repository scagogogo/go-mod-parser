# Go Mod Parser

[![Go Reference](https://pkg.go.dev/badge/github.com/scagogogo/go-mod-parser.svg)](https://pkg.go.dev/github.com/scagogogo/go-mod-parser)
[![Go Report Card](https://goreportcard.com/badge/github.com/scagogogo/go-mod-parser)](https://goreportcard.com/report/github.com/scagogogo/go-mod-parser)
[![License](https://img.shields.io/github/license/scagogogo/go-mod-parser)](LICENSE)
[![Tests](https://github.com/scagogogo/go-mod-parser/actions/workflows/go-test.yml/badge.svg)](https://github.com/scagogogo/go-mod-parser/actions/workflows/go-test.yml)
[![Documentation](https://img.shields.io/badge/docs-online-blue.svg)](https://scagogogo.github.io/go-mod-parser/)

Go Mod Parser is a comprehensive and easy-to-use library for parsing `go.mod` files. It converts go.mod files into structured Go objects, making dependency management and module analysis easier. Whether you're building dependency analysis tools, module management systems, or need to check project dependencies in CI/CD pipelines, this library provides reliable support.

## 📖 Documentation

**[📚 Complete Documentation](https://scagogogo.github.io/go-mod-parser/)** - Visit our comprehensive documentation website

**Languages:**
- [🇺🇸 English Documentation](https://scagogogo.github.io/go-mod-parser/)
- [🇨🇳 中文文档](https://scagogogo.github.io/go-mod-parser/zh/)

## Features

- ✅ **Complete Directive Support** - Parse all go.mod directives: `module`, `go`, `require`, `replace`, `exclude`, and `retract`
- 🧩 **Structured Data** - Convert go.mod files into easy-to-use Go structs
- 🔍 **Auto Discovery** - Automatically locate go.mod files in project and parent directories
- 🔄 **Dependency Analysis** - Rich helper functions for analyzing module dependencies
- 📝 **Comment Support** - Properly handle `// indirect` markers and other comments
- 🧪 **Well Tested** - Comprehensive unit test coverage ensuring parsing accuracy
- 📚 **Rich Examples** - Multiple practical examples for quick start

## Installation

```bash
go get github.com/scagogogo/go-mod-parser
```

## Quick Start

### Parse a go.mod File

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/go-mod-parser/pkg"
)

func main() {
    // Parse a go.mod file from path
    mod, err := pkg.ParseGoModFile("path/to/go.mod")
    if err != nil {
        log.Fatalf("Failed to parse go.mod: %v", err)
    }
    
    // Access parsed data
    fmt.Printf("Module: %s\n", mod.Name)
    fmt.Printf("Go Version: %s\n", mod.GoVersion)
    
    // List all dependencies
    fmt.Println("Dependencies:")
    for _, req := range mod.Requires {
        indirect := ""
        if req.Indirect {
            indirect = " // indirect"
        }
        fmt.Printf("- %s %s%s\n", req.Path, req.Version, indirect)
    }
}
```

### Auto-discover and Parse

```go
// Find and parse go.mod in current directory or parent directories
mod, err := pkg.FindAndParseGoModInCurrentDir()
if err != nil {
    log.Fatalf("Failed to find and parse go.mod: %v", err)
}

fmt.Printf("Found and parsed module: %s\n", mod.Name)
```

### Parse go.mod Content

```go
content := `module github.com/example/module

go 1.21

require github.com/stretchr/testify v1.8.4
`

// Parse go.mod content
mod, err := pkg.ParseGoModContent(content)
if err != nil {
    log.Fatalf("Failed to parse go.mod content: %v", err)
}

fmt.Printf("Module: %s\n", mod.Name)
```

## Main Features

### 1. Complete go.mod File Structure Parsing

Parse all standard directives in go.mod files:

- **module** - Module declaration
- **go** - Go version requirement
- **require** - Dependency declarations (including indirect markers)
- **replace** - Replacement rules
- **exclude** - Exclusion rules
- **retract** - Version retractions (supports single versions and version ranges)

### 2. Rich Helper Functions

```go
// Check specific dependencies
if pkg.HasRequire(mod, "github.com/stretchr/testify") {
    req := pkg.GetRequire(mod, "github.com/stretchr/testify")
    fmt.Printf("Dependency version: %s (indirect: %v)\n", req.Version, req.Indirect)
}

// Check replacement rules
if pkg.HasReplace(mod, "github.com/old/pkg") {
    rep := pkg.GetReplace(mod, "github.com/old/pkg")
    fmt.Printf("Replace: %s => %s %s\n", rep.Old.Path, rep.New.Path, rep.New.Version)
}

// Check exclusion rules
if pkg.HasExclude(mod, "github.com/problematic/pkg", "v1.0.0") {
    fmt.Println("This version is excluded")
}

// Check version retractions
if pkg.HasRetract(mod, "v1.0.0") {
    fmt.Println("This version has been retracted")
}
```

### 3. Complete API

See [Documentation](https://scagogogo.github.io/go-mod-parser/) for detailed API reference.

| Function | Description |
|----------|-------------|
| `ParseGoModFile(path)` | Parse go.mod file from path |
| `ParseGoModContent(content)` | Parse go.mod content string |
| `FindAndParseGoModFile(dir)` | Find and parse go.mod in directory and parent directories |
| `FindAndParseGoModInCurrentDir()` | Find and parse go.mod in current directory and parent directories |
| `HasRequire(mod, path)` | Check if module has specific dependency |
| `GetRequire(mod, path)` | Get specific dependency of module |
| `HasReplace(mod, path)` | Check if module has specific replacement rule |
| `GetReplace(mod, path)` | Get specific replacement rule of module |
| `HasExclude(mod, path, version)` | Check if module has specific exclusion rule |
| `HasRetract(mod, version)` | Check if module has specific retracted version |

## Examples

The project includes multiple complete examples demonstrating different usage scenarios:

- [00_simple_parser](examples/00_simple_parser) - Simple command-line tool example
- [01_basic_parsing](examples/01_basic_parsing) - Basic parsing example
- [02_find_and_parse](examples/02_find_and_parse) - Find and parse example
- [03_check_dependencies](examples/03_check_dependencies) - Dependency checking example
- [04_replaces_and_excludes](examples/04_replaces_and_excludes) - Replace and exclude rules example
- [05_retract_versions](examples/05_retract_versions) - Version retraction example
- [06_programmatic_api](examples/06_programmatic_api) - Programmatic API example

For detailed explanations, see [examples/README.md](examples/README.md).

## Project Structure

```
pkg/
├── api.go             # Main public API
├── module/            # Module data structure definitions
├── parser/            # go.mod file parsing logic
└── utils/             # Utility functions
```

## Use Cases

- Build dependency analysis tools
- Module version management systems
- Dependency checking in CI/CD pipelines
- Go project build tools
- Module relationship visualization
- Dependency update recommendation systems

## Testing

The library has comprehensive test coverage (96.1%):

```bash
# Run tests
go test -v ./...

# Run tests with coverage
go test -cover ./...
```

## Contributing

Contributions are welcome! Please submit Issues and Pull Requests to improve this project. Make sure to run tests and maintain code style consistency before submitting.

```bash
# Run tests
go test -v ./...

# Run examples
cd examples/01_basic_parsing
go run main.go ../../go.mod
```

## License

This project is open source under the [MIT License](LICENSE).

## Reference Documentation

Here are official reference documents about Go modules and go.mod file format:

1. [Go Modules Reference](https://go.dev/ref/mod) - Authoritative reference for Go module system
2. [Go Modules Wiki](https://github.com/golang/go/wiki/Modules) - More technical details and examples
3. [Go Command Documentation](https://go.dev/doc/modules/gomod-ref) - Detailed go.mod file format reference
4. [Go Modules: retract directive](https://go.dev/doc/modules/version-numbers#retract) - retract directive explanation
5. [Go Language Specification](https://go.dev/ref/spec) - Official Go language specification
