---
layout: home

hero:
  name: "Go Mod Parser"
  text: "Comprehensive Go Module Parser"
  tagline: "Parse and analyze go.mod files with ease"
  image:
    src: /logo.svg
    alt: Go Mod Parser
  actions:
    - theme: brand
      text: Get Started
      link: /quick-start
    - theme: alt
      text: API Reference
      link: /api/
    - theme: alt
      text: View on GitHub
      link: https://github.com/scagogogo/go-mod-parser

features:
  - icon: 🧩
    title: Complete Directive Support
    details: Parse all go.mod directives including module, go, require, replace, exclude, and retract
  - icon: 🔍
    title: Auto Discovery
    details: Automatically find and parse go.mod files in project directories and parent directories
  - icon: 📝
    title: Comment Support
    details: Properly handle indirect comments and other annotations in go.mod files
  - icon: 🔄
    title: Dependency Analysis
    details: Rich helper functions for analyzing module dependencies and relationships
  - icon: 🧪
    title: Well Tested
    details: Comprehensive unit test coverage ensuring parsing accuracy and reliability
  - icon: 📚
    title: Rich Examples
    details: Multiple practical examples demonstrating different usage scenarios
---

## Quick Example

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/go-mod-parser/pkg"
)

func main() {
    // Parse a go.mod file
    mod, err := pkg.ParseGoModFile("go.mod")
    if err != nil {
        log.Fatalf("Failed to parse go.mod: %v", err)
    }
    
    // Access parsed data
    fmt.Printf("Module: %s\n", mod.Name)
    fmt.Printf("Go Version: %s\n", mod.GoVersion)
    
    // List dependencies
    for _, req := range mod.Requires {
        fmt.Printf("- %s %s\n", req.Path, req.Version)
    }
}
```

## Installation

```bash
go get github.com/scagogogo/go-mod-parser
```

## Use Cases

- **Dependency Analysis Tools** - Build tools to analyze project dependencies
- **Module Version Management** - Create systems for managing module versions
- **CI/CD Pipeline Integration** - Check dependencies in continuous integration
- **Build Tools** - Integrate into Go project build systems
- **Dependency Visualization** - Create visual representations of module relationships
- **Update Recommendation Systems** - Build tools that suggest dependency updates
