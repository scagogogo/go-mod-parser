package parser

import (
	"testing"

	"github.com/scagogogo/go-mod-parser/pkg/module"
	"github.com/stretchr/testify/assert"
)

func TestParseExcludeSingleLine(t *testing.T) {
	tests := []struct {
		name          string
		line          string
		expectHandled bool
		expectPath    string
		expectVersion string
		expectError   bool
	}{
		{
			name:          "valid exclude",
			line:          "exclude github.com/example/module v1.0.0",
			expectHandled: true,
			expectPath:    "github.com/example/module",
			expectVersion: "v1.0.0",
			expectError:   false,
		},
		{
			name:          "not an exclude line",
			line:          "module github.com/example/module",
			expectHandled: false,
			expectError:   false,
		},
		{
			name:          "invalid exclude format",
			line:          "exclude github.com/example/module",
			expectHandled: false,
			expectError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mod := &module.Module{}
			handled, err := parseExcludeSingleLine(mod, tt.line)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expectHandled, handled)

			if tt.expectHandled {
				assert.Equal(t, 1, len(mod.Excludes))
				assert.Equal(t, tt.expectPath, mod.Excludes[0].Path)
				assert.Equal(t, tt.expectVersion, mod.Excludes[0].Version)
			} else {
				assert.Equal(t, 0, len(mod.Excludes))
			}
		})
	}
}

func TestParseExcludeBlockLine(t *testing.T) {
	tests := []struct {
		name          string
		line          string
		expectPath    string
		expectVersion string
		expectError   bool
	}{
		{
			name:          "valid exclude in block",
			line:          "github.com/example/module v1.0.0",
			expectPath:    "github.com/example/module",
			expectVersion: "v1.0.0",
			expectError:   false,
		},
		{
			name:        "invalid exclude format - missing version",
			line:        "github.com/example/module",
			expectError: true,
		},
		{
			name:        "empty line",
			line:        "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mod := &module.Module{}
			err := parseExcludeBlockLine(mod, tt.line)

			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, 0, len(mod.Excludes))
			} else {
				assert.NoError(t, err)
				assert.Equal(t, 1, len(mod.Excludes))
				assert.Equal(t, tt.expectPath, mod.Excludes[0].Path)
				assert.Equal(t, tt.expectVersion, mod.Excludes[0].Version)
			}
		})
	}
}
