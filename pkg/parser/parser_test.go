package parser

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFromString_ModuleName(t *testing.T) {
	content := `module github.com/example/module

go 1.21
`
	mod, err := ParseFromString(content)
	assert.NoError(t, err)
	assert.Equal(t, "github.com/example/module", mod.Name)
	assert.Equal(t, "1.21", mod.GoVersion)
}

func TestParseFromString_SingleRequire(t *testing.T) {
	content := `module github.com/example/module

go 1.21

require github.com/stretchr/testify v1.8.4
`
	mod, err := ParseFromString(content)
	assert.NoError(t, err)
	assert.Equal(t, "github.com/example/module", mod.Name)
	assert.Equal(t, 1, len(mod.Requires))
	assert.Equal(t, "github.com/stretchr/testify", mod.Requires[0].Path)
	assert.Equal(t, "v1.8.4", mod.Requires[0].Version)
	assert.False(t, mod.Requires[0].Indirect)
}

func TestParseFromString_SingleRequireIndirect(t *testing.T) {
	content := `module github.com/example/module

go 1.21

require github.com/stretchr/testify v1.8.4 // indirect
`
	mod, err := ParseFromString(content)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(mod.Requires))
	assert.Equal(t, "github.com/stretchr/testify", mod.Requires[0].Path)
	assert.Equal(t, "v1.8.4", mod.Requires[0].Version)
	assert.True(t, mod.Requires[0].Indirect)
}

func TestParseFromString_MultipleRequires(t *testing.T) {
	content := `module github.com/example/module

go 1.21

require (
	github.com/stretchr/testify v1.8.4
	github.com/pkg/errors v0.9.1
)
`
	mod, err := ParseFromString(content)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(mod.Requires))
	assert.Equal(t, "github.com/stretchr/testify", mod.Requires[0].Path)
	assert.Equal(t, "v1.8.4", mod.Requires[0].Version)
	assert.Equal(t, "github.com/pkg/errors", mod.Requires[1].Path)
	assert.Equal(t, "v0.9.1", mod.Requires[1].Version)
}

func TestParseFromString_MultipleRequiresWithIndirect(t *testing.T) {
	content := `module github.com/example/module

go 1.21

require (
	github.com/stretchr/testify v1.8.4
	github.com/pkg/errors v0.9.1 // indirect
	golang.org/x/text v0.12.0 // indirect comment with more text
)
`
	mod, err := ParseFromString(content)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(mod.Requires))
	assert.Equal(t, "github.com/stretchr/testify", mod.Requires[0].Path)
	assert.Equal(t, "v1.8.4", mod.Requires[0].Version)
	assert.False(t, mod.Requires[0].Indirect)

	assert.Equal(t, "github.com/pkg/errors", mod.Requires[1].Path)
	assert.Equal(t, "v0.9.1", mod.Requires[1].Version)
	assert.True(t, mod.Requires[1].Indirect)

	assert.Equal(t, "golang.org/x/text", mod.Requires[2].Path)
	assert.Equal(t, "v0.12.0", mod.Requires[2].Version)
	assert.True(t, mod.Requires[2].Indirect)
}

func TestParseFromString_Replace(t *testing.T) {
	content := `module github.com/example/module

go 1.21

replace github.com/old/module => github.com/new/module v1.0.0
`
	mod, err := ParseFromString(content)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(mod.Replaces))
	assert.Equal(t, "github.com/old/module", mod.Replaces[0].Old.Path)
	assert.Equal(t, "", mod.Replaces[0].Old.Version)
	assert.Equal(t, "github.com/new/module", mod.Replaces[0].New.Path)
	assert.Equal(t, "v1.0.0", mod.Replaces[0].New.Version)
}

func TestParseFromString_ReplaceWithVersion(t *testing.T) {
	content := `module github.com/example/module

go 1.21

replace (
	github.com/old/module v1.0.0 => github.com/new/module v2.0.0
	github.com/another/old => github.com/another/new v1.0.0
)
`
	mod, err := ParseFromString(content)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(mod.Replaces))

	assert.Equal(t, "github.com/old/module", mod.Replaces[0].Old.Path)
	assert.Equal(t, "v1.0.0", mod.Replaces[0].Old.Version)
	assert.Equal(t, "github.com/new/module", mod.Replaces[0].New.Path)
	assert.Equal(t, "v2.0.0", mod.Replaces[0].New.Version)

	assert.Equal(t, "github.com/another/old", mod.Replaces[1].Old.Path)
	assert.Equal(t, "", mod.Replaces[1].Old.Version)
	assert.Equal(t, "github.com/another/new", mod.Replaces[1].New.Path)
	assert.Equal(t, "v1.0.0", mod.Replaces[1].New.Version)
}

func TestParseFromString_Exclude(t *testing.T) {
	content := `module github.com/example/module

go 1.21

exclude github.com/bad/module v1.0.0
`
	mod, err := ParseFromString(content)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(mod.Excludes))
	assert.Equal(t, "github.com/bad/module", mod.Excludes[0].Path)
	assert.Equal(t, "v1.0.0", mod.Excludes[0].Version)
}

func TestParseFromString_MultipleExcludes(t *testing.T) {
	content := `module github.com/example/module

go 1.21

exclude (
	github.com/bad/module v1.0.0
	github.com/another/bad v2.0.0
)
`
	mod, err := ParseFromString(content)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(mod.Excludes))
	assert.Equal(t, "github.com/bad/module", mod.Excludes[0].Path)
	assert.Equal(t, "v1.0.0", mod.Excludes[0].Version)
	assert.Equal(t, "github.com/another/bad", mod.Excludes[1].Path)
	assert.Equal(t, "v2.0.0", mod.Excludes[1].Version)
}

func TestParseFromString_SingleRetract(t *testing.T) {
	content := `module github.com/example/module

go 1.21

retract v1.0.0 // Contains critical security vulnerability
`
	mod, err := ParseFromString(content)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(mod.Retracts))
	assert.Equal(t, "v1.0.0", mod.Retracts[0].Version)
	assert.Equal(t, "Contains critical security vulnerability", mod.Retracts[0].Rationale)
}

func TestParseFromString_RetractVersionRange(t *testing.T) {
	content := `module github.com/example/module

go 1.21

retract [v1.0.0, v1.9.9] // Contains critical security vulnerability
`
	mod, err := ParseFromString(content)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(mod.Retracts))
	assert.Equal(t, "", mod.Retracts[0].Version)
	assert.Equal(t, "v1.0.0", mod.Retracts[0].VersionLow)
	assert.Equal(t, "v1.9.9", mod.Retracts[0].VersionHigh)
	assert.Equal(t, "Contains critical security vulnerability", mod.Retracts[0].Rationale)
}

func TestParseFromString_MultipleRetracts(t *testing.T) {
	content := `module github.com/example/module

go 1.21

retract (
	v1.0.0 // Contains critical security vulnerability
	[v0.5.0, v0.9.9] // Experimental versions
)
`
	mod, err := ParseFromString(content)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(mod.Retracts))
	assert.Equal(t, "v1.0.0", mod.Retracts[0].Version)
	assert.Equal(t, "Contains critical security vulnerability", mod.Retracts[0].Rationale)

	assert.Equal(t, "", mod.Retracts[1].Version)
	assert.Equal(t, "v0.5.0", mod.Retracts[1].VersionLow)
	assert.Equal(t, "v0.9.9", mod.Retracts[1].VersionHigh)
	assert.Equal(t, "Experimental versions", mod.Retracts[1].Rationale)
}

func TestParseFromString_ComplexFile(t *testing.T) {
	content := `module github.com/example/module

go 1.21

require (
	github.com/stretchr/testify v1.8.4
	github.com/pkg/errors v0.9.1 // indirect
)

replace (
	github.com/old/module => github.com/new/module v1.0.0
	github.com/another/old v1.0.0 => github.com/another/new v2.0.0
)

exclude (
	github.com/bad/module v1.0.0
	github.com/another/bad v2.0.0
)

retract (
	v1.0.0 // Critical bug
	[v0.9.0, v0.9.9] // Experimental
)
`
	mod, err := ParseFromString(content)
	assert.NoError(t, err)
	assert.Equal(t, "github.com/example/module", mod.Name)
	assert.Equal(t, "1.21", mod.GoVersion)
	assert.Equal(t, 2, len(mod.Requires))
	assert.Equal(t, 2, len(mod.Replaces))
	assert.Equal(t, 2, len(mod.Excludes))
	assert.Equal(t, 2, len(mod.Retracts))

	// 验证第一个require不是indirect
	assert.False(t, mod.Requires[0].Indirect)
	// 验证第二个require是indirect
	assert.True(t, mod.Requires[1].Indirect)

	// 验证第一个retract
	assert.Equal(t, "v1.0.0", mod.Retracts[0].Version)
	assert.Equal(t, "Critical bug", mod.Retracts[0].Rationale)

	// 验证第二个retract
	assert.Equal(t, "v0.9.0", mod.Retracts[1].VersionLow)
	assert.Equal(t, "v0.9.9", mod.Retracts[1].VersionHigh)
	assert.Equal(t, "Experimental", mod.Retracts[1].Rationale)
}

func TestParseFromString_WithComments(t *testing.T) {
	content := `// This is a comment at the top of the file
module github.com/example/module

// Go version comment
go 1.21

// Dependencies
require (
	// A test framework
	github.com/stretchr/testify v1.8.4
	// Error handling
	github.com/pkg/errors v0.9.1
)
`
	mod, err := ParseFromString(content)
	assert.NoError(t, err)
	assert.Equal(t, "github.com/example/module", mod.Name)
	assert.Equal(t, "1.21", mod.GoVersion)
	assert.Equal(t, 2, len(mod.Requires))
}

func TestParseFromReader(t *testing.T) {
	content := `module github.com/example/module

go 1.21

require github.com/stretchr/testify v1.8.4
`
	reader := strings.NewReader(content)
	mod, err := ParseFromReader(reader)
	assert.NoError(t, err)
	assert.Equal(t, "github.com/example/module", mod.Name)
	assert.Equal(t, 1, len(mod.Requires))
}

func TestParseFromString_InvalidContent(t *testing.T) {
	content := `invalid content`
	_, err := ParseFromString(content)
	assert.Error(t, err)
}
