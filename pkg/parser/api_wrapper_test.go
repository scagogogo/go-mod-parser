package parser

import (
	"io"
	"strings"
	"testing"

	"github.com/scagogogo/go-mod-parser/pkg/module"
	"github.com/stretchr/testify/assert"
)

func TestHandleSingleLine(t *testing.T) {
	tests := []struct {
		name           string
		line           string
		expectedModule *module.Module
		expectError    bool
	}{
		{
			name: "module line",
			line: "module github.com/example/module",
			expectedModule: &module.Module{
				Name: "github.com/example/module",
			},
			expectError: false,
		},
		{
			name: "go version line",
			line: "go 1.21",
			expectedModule: &module.Module{
				GoVersion: "1.21",
			},
			expectError: false,
		},
		{
			name: "require line",
			line: "require github.com/stretchr/testify v1.8.4",
			expectedModule: &module.Module{
				Requires: []*module.Require{
					{
						Path:     "github.com/stretchr/testify",
						Version:  "v1.8.4",
						Indirect: false,
					},
				},
			},
			expectError: false,
		},
		{
			name: "replace line",
			line: "replace github.com/old/module => github.com/new/module v1.0.0",
			expectedModule: &module.Module{
				Replaces: []*module.Replace{
					{
						Old: &module.ReplaceItem{
							Path: "github.com/old/module",
						},
						New: &module.ReplaceItem{
							Path:    "github.com/new/module",
							Version: "v1.0.0",
						},
					},
				},
			},
			expectError: false,
		},
		{
			name: "exclude line",
			line: "exclude github.com/bad/module v1.0.0",
			expectedModule: &module.Module{
				Excludes: []*module.Exclude{
					{
						Path:    "github.com/bad/module",
						Version: "v1.0.0",
					},
				},
			},
			expectError: false,
		},
		{
			name: "retract line",
			line: "retract v1.0.0 // security vulnerability",
			expectedModule: &module.Module{
				Retracts: []*module.Retract{
					{
						Version:   "v1.0.0",
						Rationale: "security vulnerability",
					},
				},
			},
			expectError: false,
		},
		{
			name:        "invalid line",
			line:        "invalid content",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mod := &module.Module{}
			err := handleSingleLine(mod, tt.line)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				if tt.expectedModule.Name != "" {
					assert.Equal(t, tt.expectedModule.Name, mod.Name)
				}

				if tt.expectedModule.GoVersion != "" {
					assert.Equal(t, tt.expectedModule.GoVersion, mod.GoVersion)
				}

				if len(tt.expectedModule.Requires) > 0 {
					assert.Equal(t, len(tt.expectedModule.Requires), len(mod.Requires))
					assert.Equal(t, tt.expectedModule.Requires[0].Path, mod.Requires[0].Path)
					assert.Equal(t, tt.expectedModule.Requires[0].Version, mod.Requires[0].Version)
					assert.Equal(t, tt.expectedModule.Requires[0].Indirect, mod.Requires[0].Indirect)
				}

				if len(tt.expectedModule.Replaces) > 0 {
					assert.Equal(t, len(tt.expectedModule.Replaces), len(mod.Replaces))
					assert.Equal(t, tt.expectedModule.Replaces[0].Old.Path, mod.Replaces[0].Old.Path)
					assert.Equal(t, tt.expectedModule.Replaces[0].New.Path, mod.Replaces[0].New.Path)
					assert.Equal(t, tt.expectedModule.Replaces[0].New.Version, mod.Replaces[0].New.Version)
				}

				if len(tt.expectedModule.Excludes) > 0 {
					assert.Equal(t, len(tt.expectedModule.Excludes), len(mod.Excludes))
					assert.Equal(t, tt.expectedModule.Excludes[0].Path, mod.Excludes[0].Path)
					assert.Equal(t, tt.expectedModule.Excludes[0].Version, mod.Excludes[0].Version)
				}

				if len(tt.expectedModule.Retracts) > 0 {
					assert.Equal(t, len(tt.expectedModule.Retracts), len(mod.Retracts))
					assert.Equal(t, tt.expectedModule.Retracts[0].Version, mod.Retracts[0].Version)
					assert.Equal(t, tt.expectedModule.Retracts[0].Rationale, mod.Retracts[0].Rationale)
				}
			}
		})
	}
}

func TestHandleBlockLine(t *testing.T) {
	tests := []struct {
		name           string
		blockType      string
		line           string
		expectedModule *module.Module
		expectError    bool
	}{
		{
			name:      "require block line",
			blockType: "require",
			line:      "github.com/stretchr/testify v1.8.4",
			expectedModule: &module.Module{
				Requires: []*module.Require{
					{
						Path:     "github.com/stretchr/testify",
						Version:  "v1.8.4",
						Indirect: false,
					},
				},
			},
			expectError: false,
		},
		{
			name:      "replace block line",
			blockType: "replace",
			line:      "github.com/old/module => github.com/new/module v1.0.0",
			expectedModule: &module.Module{
				Replaces: []*module.Replace{
					{
						Old: &module.ReplaceItem{
							Path: "github.com/old/module",
						},
						New: &module.ReplaceItem{
							Path:    "github.com/new/module",
							Version: "v1.0.0",
						},
					},
				},
			},
			expectError: false,
		},
		{
			name:      "exclude block line",
			blockType: "exclude",
			line:      "github.com/bad/module v1.0.0",
			expectedModule: &module.Module{
				Excludes: []*module.Exclude{
					{
						Path:    "github.com/bad/module",
						Version: "v1.0.0",
					},
				},
			},
			expectError: false,
		},
		{
			name:      "retract block line",
			blockType: "retract",
			line:      "v1.0.0 // security vulnerability",
			expectedModule: &module.Module{
				Retracts: []*module.Retract{
					{
						Version:   "v1.0.0",
						Rationale: "security vulnerability",
					},
				},
			},
			expectError: false,
		},
		{
			name:        "unknown block type",
			blockType:   "unknown",
			line:        "some content",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mod := &module.Module{}
			err := handleBlockLine(mod, tt.blockType, tt.line)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				if len(tt.expectedModule.Requires) > 0 {
					assert.Equal(t, len(tt.expectedModule.Requires), len(mod.Requires))
					assert.Equal(t, tt.expectedModule.Requires[0].Path, mod.Requires[0].Path)
					assert.Equal(t, tt.expectedModule.Requires[0].Version, mod.Requires[0].Version)
					assert.Equal(t, tt.expectedModule.Requires[0].Indirect, mod.Requires[0].Indirect)
				}

				if len(tt.expectedModule.Replaces) > 0 {
					assert.Equal(t, len(tt.expectedModule.Replaces), len(mod.Replaces))
					assert.Equal(t, tt.expectedModule.Replaces[0].Old.Path, mod.Replaces[0].Old.Path)
					assert.Equal(t, tt.expectedModule.Replaces[0].New.Path, mod.Replaces[0].New.Path)
					assert.Equal(t, tt.expectedModule.Replaces[0].New.Version, mod.Replaces[0].New.Version)
				}

				if len(tt.expectedModule.Excludes) > 0 {
					assert.Equal(t, len(tt.expectedModule.Excludes), len(mod.Excludes))
					assert.Equal(t, tt.expectedModule.Excludes[0].Path, mod.Excludes[0].Path)
					assert.Equal(t, tt.expectedModule.Excludes[0].Version, mod.Excludes[0].Version)
				}

				if len(tt.expectedModule.Retracts) > 0 {
					assert.Equal(t, len(tt.expectedModule.Retracts), len(mod.Retracts))
					assert.Equal(t, tt.expectedModule.Retracts[0].Version, mod.Retracts[0].Version)
					assert.Equal(t, tt.expectedModule.Retracts[0].Rationale, mod.Retracts[0].Rationale)
				}
			}
		})
	}
}

func TestParseFromStringSimple(t *testing.T) {
	// 这个测试已经在parser_test.go中有很多，这里简单测试一下API是否正常
	content := `module github.com/example/module

go 1.21
`
	mod, err := ParseFromString(content)
	assert.NoError(t, err)
	assert.Equal(t, "github.com/example/module", mod.Name)
	assert.Equal(t, "1.21", mod.GoVersion)
}

// 模拟ParseFromReader失败的io.Reader
type errorReader struct{}

func (r errorReader) Read(p []byte) (n int, err error) {
	return 0, io.ErrUnexpectedEOF
}

func TestParseFromReader_Error(t *testing.T) {
	r := errorReader{}
	_, err := ParseFromReader(r)
	assert.Error(t, err)
}

func TestParseFromReader_Success(t *testing.T) {
	r := strings.NewReader(`module github.com/example/module

go 1.21
`)
	mod, err := ParseFromReader(r)
	assert.NoError(t, err)
	assert.Equal(t, "github.com/example/module", mod.Name)
	assert.Equal(t, "1.21", mod.GoVersion)
}
