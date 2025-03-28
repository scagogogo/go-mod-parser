package parser

import (
	"testing"

	"github.com/scagogogo/go-mod-parser/pkg/module"
	"github.com/stretchr/testify/assert"
)

func TestParseReplaceSingleLine(t *testing.T) {
	tests := []struct {
		name             string
		line             string
		expectHandled    bool
		expectOldPath    string
		expectOldVersion string
		expectNewPath    string
		expectNewVersion string
		expectError      bool
	}{
		{
			name:             "valid replace",
			line:             "replace github.com/old/module => github.com/new/module v1.0.0",
			expectHandled:    true,
			expectOldPath:    "github.com/old/module",
			expectOldVersion: "",
			expectNewPath:    "github.com/new/module",
			expectNewVersion: "v1.0.0",
			expectError:      false,
		},
		{
			name:          "not a replace line",
			line:          "module github.com/example/module",
			expectHandled: false,
			expectError:   false,
		},
		{
			name:          "invalid replace format",
			line:          "replace github.com/old/module github.com/new/module v1.0.0",
			expectHandled: false,
			expectError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mod := &module.Module{}
			handled, err := parseReplaceSingleLine(mod, tt.line)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expectHandled, handled)

			if tt.expectHandled {
				assert.Equal(t, 1, len(mod.Replaces))
				assert.Equal(t, tt.expectOldPath, mod.Replaces[0].Old.Path)
				assert.Equal(t, tt.expectOldVersion, mod.Replaces[0].Old.Version)
				assert.Equal(t, tt.expectNewPath, mod.Replaces[0].New.Path)
				assert.Equal(t, tt.expectNewVersion, mod.Replaces[0].New.Version)
			} else {
				assert.Equal(t, 0, len(mod.Replaces))
			}
		})
	}
}

func TestParseReplaceBlockLine(t *testing.T) {
	tests := []struct {
		name             string
		line             string
		expectOldPath    string
		expectOldVersion string
		expectNewPath    string
		expectNewVersion string
		expectError      bool
	}{
		{
			name:             "valid replace in block without old version",
			line:             "github.com/old/module => github.com/new/module v1.0.0",
			expectOldPath:    "github.com/old/module",
			expectOldVersion: "",
			expectNewPath:    "github.com/new/module",
			expectNewVersion: "v1.0.0",
			expectError:      false,
		},
		{
			name:             "valid replace in block with old version",
			line:             "github.com/old/module v0.5.0 => github.com/new/module v1.0.0",
			expectOldPath:    "github.com/old/module",
			expectOldVersion: "v0.5.0",
			expectNewPath:    "github.com/new/module",
			expectNewVersion: "v1.0.0",
			expectError:      false,
		},
		{
			name:             "valid replace without new version",
			line:             "github.com/old/module => github.com/new/module",
			expectOldPath:    "github.com/old/module",
			expectOldVersion: "",
			expectNewPath:    "github.com/new/module",
			expectNewVersion: "",
			expectError:      false,
		},
		{
			name:        "invalid replace format - missing =>",
			line:        "github.com/old/module github.com/new/module v1.0.0",
			expectError: true,
		},
		{
			name:        "invalid replace format - empty new path",
			line:        "github.com/old/module => ",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mod := &module.Module{}
			err := parseReplaceBlockLine(mod, tt.line)

			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, 0, len(mod.Replaces))
			} else {
				assert.NoError(t, err)
				assert.Equal(t, 1, len(mod.Replaces))
				assert.Equal(t, tt.expectOldPath, mod.Replaces[0].Old.Path)
				assert.Equal(t, tt.expectOldVersion, mod.Replaces[0].Old.Version)
				assert.Equal(t, tt.expectNewPath, mod.Replaces[0].New.Path)
				assert.Equal(t, tt.expectNewVersion, mod.Replaces[0].New.Version)
			}
		})
	}
}
