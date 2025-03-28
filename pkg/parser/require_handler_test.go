package parser

import (
	"testing"

	"github.com/scagogogo/go-mod-parser/pkg/module"
	"github.com/stretchr/testify/assert"
)

func TestParseRequireSingleLine(t *testing.T) {
	tests := []struct {
		name           string
		line           string
		expectHandled  bool
		expectPath     string
		expectVersion  string
		expectIndirect bool
		expectError    bool
	}{
		{
			name:           "valid require",
			line:           "require github.com/example/module v1.0.0",
			expectHandled:  true,
			expectPath:     "github.com/example/module",
			expectVersion:  "v1.0.0",
			expectIndirect: false,
			expectError:    false,
		},
		{
			name:           "valid require with indirect",
			line:           "require github.com/example/module v1.0.0 // indirect",
			expectHandled:  true,
			expectPath:     "github.com/example/module",
			expectVersion:  "v1.0.0",
			expectIndirect: true,
			expectError:    false,
		},
		{
			name:           "valid require with comments",
			line:           "require github.com/example/module v1.0.0 // some comment",
			expectHandled:  true,
			expectPath:     "github.com/example/module",
			expectVersion:  "v1.0.0",
			expectIndirect: false,
			expectError:    false,
		},
		{
			name:          "not a require line",
			line:          "module github.com/example/module",
			expectHandled: false,
			expectError:   false,
		},
		{
			name:          "invalid require format",
			line:          "require github.com/example/module",
			expectHandled: false,
			expectError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mod := &module.Module{}
			handled, err := parseRequireSingleLine(mod, tt.line)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expectHandled, handled)

			if tt.expectHandled {
				assert.Equal(t, 1, len(mod.Requires))
				assert.Equal(t, tt.expectPath, mod.Requires[0].Path)
				assert.Equal(t, tt.expectVersion, mod.Requires[0].Version)
				assert.Equal(t, tt.expectIndirect, mod.Requires[0].Indirect)
			} else {
				assert.Equal(t, 0, len(mod.Requires))
			}
		})
	}
}

func TestParseRequireBlockLine(t *testing.T) {
	tests := []struct {
		name           string
		line           string
		expectPath     string
		expectVersion  string
		expectIndirect bool
		expectError    bool
	}{
		{
			name:           "valid require in block",
			line:           "github.com/example/module v1.0.0",
			expectPath:     "github.com/example/module",
			expectVersion:  "v1.0.0",
			expectIndirect: false,
			expectError:    false,
		},
		{
			name:           "valid require with indirect in block",
			line:           "github.com/example/module v1.0.0 // indirect",
			expectPath:     "github.com/example/module",
			expectVersion:  "v1.0.0",
			expectIndirect: true,
			expectError:    false,
		},
		{
			name:           "valid require with comments in block",
			line:           "github.com/example/module v1.0.0 // some comment",
			expectPath:     "github.com/example/module",
			expectVersion:  "v1.0.0",
			expectIndirect: false,
			expectError:    false,
		},
		{
			name:        "invalid require format - missing version",
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
			err := parseRequireBlockLine(mod, tt.line)

			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, 0, len(mod.Requires))
			} else {
				assert.NoError(t, err)
				assert.Equal(t, 1, len(mod.Requires))
				assert.Equal(t, tt.expectPath, mod.Requires[0].Path)
				assert.Equal(t, tt.expectVersion, mod.Requires[0].Version)
				assert.Equal(t, tt.expectIndirect, mod.Requires[0].Indirect)
			}
		})
	}
}
