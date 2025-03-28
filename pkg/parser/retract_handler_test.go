package parser

import (
	"testing"

	"github.com/scagogogo/go-mod-parser/pkg/module"
	"github.com/stretchr/testify/assert"
)

func TestParseRetractSingleLine(t *testing.T) {
	tests := []struct {
		name              string
		line              string
		expectHandled     bool
		expectVersion     string
		expectVersionLow  string
		expectVersionHigh string
		expectRationale   string
		expectError       bool
	}{
		{
			name:              "valid retract single version",
			line:              "retract v1.0.0",
			expectHandled:     true,
			expectVersion:     "v1.0.0",
			expectVersionLow:  "",
			expectVersionHigh: "",
			expectRationale:   "",
			expectError:       false,
		},
		{
			name:              "valid retract single version with rationale",
			line:              "retract v1.0.0 // security vulnerability",
			expectHandled:     true,
			expectVersion:     "v1.0.0",
			expectVersionLow:  "",
			expectVersionHigh: "",
			expectRationale:   "security vulnerability",
			expectError:       false,
		},
		{
			name:              "valid retract version range",
			line:              "retract [v0.5.0, v0.9.9]",
			expectHandled:     true,
			expectVersion:     "",
			expectVersionLow:  "v0.5.0",
			expectVersionHigh: "v0.9.9",
			expectRationale:   "",
			expectError:       false,
		},
		{
			name:              "valid retract version range with rationale",
			line:              "retract [v0.5.0, v0.9.9] // experimental versions",
			expectHandled:     true,
			expectVersion:     "",
			expectVersionLow:  "v0.5.0",
			expectVersionHigh: "v0.9.9",
			expectRationale:   "experimental versions",
			expectError:       false,
		},
		{
			name:          "not a retract line",
			line:          "module github.com/example/module",
			expectHandled: false,
			expectError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mod := &module.Module{}
			handled, err := parseRetractSingleLine(mod, tt.line)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expectHandled, handled)

			if tt.expectHandled {
				assert.Equal(t, 1, len(mod.Retracts))
				assert.Equal(t, tt.expectVersion, mod.Retracts[0].Version)
				assert.Equal(t, tt.expectVersionLow, mod.Retracts[0].VersionLow)
				assert.Equal(t, tt.expectVersionHigh, mod.Retracts[0].VersionHigh)
				assert.Equal(t, tt.expectRationale, mod.Retracts[0].Rationale)
			} else {
				assert.Equal(t, 0, len(mod.Retracts))
			}
		})
	}
}

func TestParseRetractBlockLine(t *testing.T) {
	tests := []struct {
		name              string
		line              string
		expectVersion     string
		expectVersionLow  string
		expectVersionHigh string
		expectRationale   string
		expectError       bool
	}{
		{
			name:              "valid retract single version in block",
			line:              "v1.0.0",
			expectVersion:     "v1.0.0",
			expectVersionLow:  "",
			expectVersionHigh: "",
			expectRationale:   "",
			expectError:       false,
		},
		{
			name:              "valid retract single version with rationale in block",
			line:              "v1.0.0 // security vulnerability",
			expectVersion:     "v1.0.0",
			expectVersionLow:  "",
			expectVersionHigh: "",
			expectRationale:   "security vulnerability",
			expectError:       false,
		},
		{
			name:              "valid retract version range in block",
			line:              "[v0.5.0, v0.9.9]",
			expectVersion:     "",
			expectVersionLow:  "v0.5.0",
			expectVersionHigh: "v0.9.9",
			expectRationale:   "",
			expectError:       false,
		},
		{
			name:              "valid retract version range with rationale in block",
			line:              "[v0.5.0, v0.9.9] // experimental versions",
			expectVersion:     "",
			expectVersionLow:  "v0.5.0",
			expectVersionHigh: "v0.9.9",
			expectRationale:   "experimental versions",
			expectError:       false,
		},
		{
			name:        "comment line",
			line:        "// comment line only",
			expectError: false,
		},
		{
			name:        "empty line",
			line:        "",
			expectError: false,
		},
		{
			name:        "invalid retract format",
			line:        "invalid format",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mod := &module.Module{}
			err := parseRetractBlockLine(mod, tt.line)

			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, 0, len(mod.Retracts))
			} else {
				assert.NoError(t, err)

				// 跳过注释行和空行的结果检查
				if tt.line == "" || (len(tt.line) > 1 && tt.line[0:2] == "//") {
					assert.Equal(t, 0, len(mod.Retracts))
					return
				}

				assert.Equal(t, 1, len(mod.Retracts))
				assert.Equal(t, tt.expectVersion, mod.Retracts[0].Version)
				assert.Equal(t, tt.expectVersionLow, mod.Retracts[0].VersionLow)
				assert.Equal(t, tt.expectVersionHigh, mod.Retracts[0].VersionHigh)
				assert.Equal(t, tt.expectRationale, mod.Retracts[0].Rationale)
			}
		})
	}
}
