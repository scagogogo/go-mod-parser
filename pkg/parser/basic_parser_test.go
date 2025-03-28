package parser

import (
	"testing"

	"github.com/scagogogo/go-mod-parser/pkg/module"
	"github.com/stretchr/testify/assert"
)

func TestParseModuleName(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		expected string
		handled  bool
	}{
		{
			name:     "valid module declaration",
			line:     "module github.com/example/module",
			expected: "github.com/example/module",
			handled:  true,
		},
		{
			name:     "invalid module declaration - wrong format",
			line:     "modulegithub.com/example/module",
			expected: "",
			handled:  false,
		},
		{
			name:     "invalid module declaration - extra tokens",
			line:     "module github.com/example/module extra",
			expected: "",
			handled:  false,
		},
		{
			name:     "not a module declaration",
			line:     "go 1.21",
			expected: "",
			handled:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mod := &module.Module{}
			handled, err := parseModuleName(mod, tt.line)
			assert.NoError(t, err)
			assert.Equal(t, tt.handled, handled)
			assert.Equal(t, tt.expected, mod.Name)
		})
	}
}

func TestParseGoVersion(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		expected string
		handled  bool
	}{
		{
			name:     "valid go version",
			line:     "go 1.21",
			expected: "1.21",
			handled:  true,
		},
		{
			name:     "valid go version with patch",
			line:     "go 1.21.0",
			expected: "1.21.0",
			handled:  true,
		},
		{
			name:     "invalid go version - wrong format",
			line:     "go1.21",
			expected: "",
			handled:  false,
		},
		{
			name:     "invalid go version - extra tokens",
			line:     "go 1.21 extra",
			expected: "",
			handled:  false,
		},
		{
			name:     "not a go version declaration",
			line:     "module github.com/example/module",
			expected: "",
			handled:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mod := &module.Module{}
			handled, err := parseGoVersion(mod, tt.line)
			assert.NoError(t, err)
			assert.Equal(t, tt.handled, handled)
			assert.Equal(t, tt.expected, mod.GoVersion)
		})
	}
}

func TestIsIndirect(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		expected bool
	}{
		{
			name:     "indirect comment",
			line:     "github.com/example/module v1.0.0 // indirect",
			expected: true,
		},
		{
			name:     "indirect comment with extra text",
			line:     "github.com/example/module v1.0.0 // indirect comment",
			expected: true,
		},
		{
			name:     "not indirect comment",
			line:     "github.com/example/module v1.0.0 // comment",
			expected: false,
		},
		{
			name:     "no comment",
			line:     "github.com/example/module v1.0.0",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isIndirect(tt.line)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBlockStartsWith(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		keyword  string
		expected bool
	}{
		{
			name:     "line starts with keyword and ends with (",
			line:     "require (",
			keyword:  "require",
			expected: true,
		},
		{
			name:     "line with spaces starts with keyword and ends with (",
			line:     "  require  (",
			keyword:  "require",
			expected: true,
		},
		{
			name:     "line starts with keyword but doesn't end with (",
			line:     "require github.com/example/module v1.0.0",
			keyword:  "require",
			expected: false,
		},
		{
			name:     "line doesn't start with keyword",
			line:     "module github.com/example/module",
			keyword:  "require",
			expected: false,
		},
		{
			name:     "empty line",
			line:     "",
			keyword:  "require",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := blockStartsWith(tt.line, tt.keyword)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsBlockEnd(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		expected bool
	}{
		{
			name:     "just closing bracket",
			line:     ")",
			expected: true,
		},
		{
			name:     "closing bracket with spaces",
			line:     "  )  ",
			expected: true,
		},
		{
			name:     "closing bracket with text",
			line:     ") // comment",
			expected: false,
		},
		{
			name:     "not a closing bracket",
			line:     "github.com/example/module v1.0.0",
			expected: false,
		},
		{
			name:     "empty line",
			line:     "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isBlockEnd(tt.line)
			assert.Equal(t, tt.expected, result)
		})
	}
}
