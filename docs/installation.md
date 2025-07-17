# Installation

## Requirements

- Go 1.21 or later
- Git (for fetching the module)

## Install via Go Modules

The easiest way to install Go Mod Parser is using Go modules:

```bash
go get github.com/scagogogo/go-mod-parser
```

## Verify Installation

Create a simple test file to verify the installation:

```go
// test.go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/go-mod-parser/pkg"
)

func main() {
    // Test parsing a simple go.mod content
    content := `module github.com/example/test

go 1.21

require github.com/stretchr/testify v1.8.4
`
    
    mod, err := pkg.ParseGoModContent(content)
    if err != nil {
        log.Fatalf("Failed to parse: %v", err)
    }
    
    fmt.Printf("Successfully parsed module: %s\n", mod.Name)
    fmt.Printf("Go version: %s\n", mod.GoVersion)
    fmt.Printf("Dependencies: %d\n", len(mod.Requires))
}
```

Run the test:

```bash
go run test.go
```

Expected output:
```
Successfully parsed module: github.com/example/test
Go version: 1.21
Dependencies: 1
```

## Import in Your Project

Add the import to your Go files:

```go
import "github.com/scagogogo/go-mod-parser/pkg"
```

For specific subpackages:

```go
import (
    "github.com/scagogogo/go-mod-parser/pkg"
    "github.com/scagogogo/go-mod-parser/pkg/module"
)
```

## Development Setup

If you want to contribute or modify the library:

1. Clone the repository:
```bash
git clone https://github.com/scagogogo/go-mod-parser.git
cd go-mod-parser
```

2. Install dependencies:
```bash
go mod download
```

3. Run tests:
```bash
go test -v ./...
```

4. Run examples:
```bash
cd examples/01_basic_parsing
go run main.go ../../go.mod
```

## Troubleshooting

### Module Not Found

If you get a "module not found" error:

1. Ensure you're using Go 1.21 or later:
```bash
go version
```

2. Initialize your module if you haven't:
```bash
go mod init your-project-name
```

3. Try fetching the module explicitly:
```bash
go get -u github.com/scagogogo/go-mod-parser
```

### Import Errors

If you encounter import errors:

1. Check your Go module is properly initialized
2. Ensure the import path is correct
3. Run `go mod tidy` to clean up dependencies

### Version Conflicts

If you have version conflicts:

1. Check your go.mod file for conflicting versions
2. Use `go mod why` to understand dependency chains
3. Consider using `replace` directives if needed
