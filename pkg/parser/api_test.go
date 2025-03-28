package parser_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/scagogogo/go-mod-parser/pkg/module"
	"github.com/scagogogo/go-mod-parser/pkg/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFindAndParseGoModFile(t *testing.T) {
	// 创建临时目录结构
	tempDir := t.TempDir()
	subDirPath := filepath.Join(tempDir, "level1", "level2")
	require.NoError(t, os.MkdirAll(subDirPath, 0755))

	// 在顶层创建go.mod文件
	goModPath := filepath.Join(tempDir, "go.mod")
	content := `module example.com/test
go 1.18
require github.com/stretchr/testify v1.8.4
`
	require.NoError(t, os.WriteFile(goModPath, []byte(content), 0644))

	// 保存当前目录
	oldDir, err := os.Getwd()
	require.NoError(t, err)
	defer func() { require.NoError(t, os.Chdir(oldDir)) }()

	// 从子目录查找并解析
	require.NoError(t, os.Chdir(subDirPath))

	mod, err := parser.FindAndParseGoModFile(".")
	require.NoError(t, err)
	require.NotNil(t, mod)

	assert.Equal(t, "example.com/test", mod.Name)
	assert.Equal(t, "1.18", mod.GoVersion)
	assert.Len(t, mod.Requires, 1)
	assert.Equal(t, "github.com/stretchr/testify", mod.Requires[0].Path)
	assert.Equal(t, "v1.8.4", mod.Requires[0].Version)
}

func TestParseGoModContent(t *testing.T) {
	content := `module example.com/test
go 1.18
require github.com/stretchr/testify v1.8.4
`
	mod, err := parser.ParseGoModContent(content)
	require.NoError(t, err)
	require.NotNil(t, mod)

	assert.Equal(t, "example.com/test", mod.Name)
	assert.Equal(t, "1.18", mod.GoVersion)
	assert.Len(t, mod.Requires, 1)
	assert.Equal(t, "github.com/stretchr/testify", mod.Requires[0].Path)
	assert.Equal(t, "v1.8.4", mod.Requires[0].Version)
}

func TestParseGoModFile(t *testing.T) {
	// 创建临时go.mod文件
	tempDir := t.TempDir()
	goModPath := filepath.Join(tempDir, "go.mod")
	content := `module example.com/test
go 1.18
require github.com/stretchr/testify v1.8.4
`
	require.NoError(t, os.WriteFile(goModPath, []byte(content), 0644))

	mod, err := parser.ParseGoModFile(goModPath)
	require.NoError(t, err)
	require.NotNil(t, mod)

	assert.Equal(t, "example.com/test", mod.Name)
	assert.Equal(t, "1.18", mod.GoVersion)
	assert.Len(t, mod.Requires, 1)
	assert.Equal(t, "github.com/stretchr/testify", mod.Requires[0].Path)
	assert.Equal(t, "v1.8.4", mod.Requires[0].Version)
}

func TestHelperFunctionsAPI(t *testing.T) {
	// 创建一个测试模块
	mod := &module.Module{
		Name:      "github.com/example/module",
		GoVersion: "1.21",
		Requires: []*module.Require{
			{Path: "github.com/stretchr/testify", Version: "v1.8.4", Indirect: false},
			{Path: "github.com/pkg/errors", Version: "v0.9.1", Indirect: true},
		},
		Replaces: []*module.Replace{
			{
				Old: &module.ReplaceItem{Path: "github.com/old/pkg", Version: "v1.0.0"},
				New: &module.ReplaceItem{Path: "github.com/new/pkg", Version: "v2.0.0"},
			},
		},
		Excludes: []*module.Exclude{
			{Path: "github.com/bad/pkg", Version: "v0.1.0"},
		},
		Retracts: []*module.Retract{
			{Version: "v1.0.0", Rationale: "Critical bug"},
			{VersionLow: "v0.1.0", VersionHigh: "v0.9.0", Rationale: "Experimental versions"},
		},
	}

	// 测试 HasRequire
	assert.True(t, parser.HasRequire(mod, "github.com/stretchr/testify"))
	assert.True(t, parser.HasRequire(mod, "github.com/pkg/errors"))
	assert.False(t, parser.HasRequire(mod, "github.com/nonexistent/pkg"))

	// 测试 GetRequire
	req := parser.GetRequire(mod, "github.com/stretchr/testify")
	assert.NotNil(t, req)
	assert.Equal(t, "v1.8.4", req.Version)
	assert.False(t, req.Indirect)

	req = parser.GetRequire(mod, "github.com/pkg/errors")
	assert.NotNil(t, req)
	assert.Equal(t, "v0.9.1", req.Version)
	assert.True(t, req.Indirect)

	req = parser.GetRequire(mod, "github.com/nonexistent/pkg")
	assert.Nil(t, req)

	// 测试 HasReplace
	assert.True(t, parser.HasReplace(mod, "github.com/old/pkg"))
	assert.False(t, parser.HasReplace(mod, "github.com/nonexistent/pkg"))

	// 测试 GetReplace
	rep := parser.GetReplace(mod, "github.com/old/pkg")
	assert.NotNil(t, rep)
	assert.Equal(t, "v1.0.0", rep.Old.Version)
	assert.Equal(t, "github.com/new/pkg", rep.New.Path)
	assert.Equal(t, "v2.0.0", rep.New.Version)

	rep = parser.GetReplace(mod, "github.com/nonexistent/pkg")
	assert.Nil(t, rep)

	// 测试 HasExclude
	assert.True(t, parser.HasExclude(mod, "github.com/bad/pkg", "v0.1.0"))
	assert.False(t, parser.HasExclude(mod, "github.com/bad/pkg", "v0.2.0"))
	assert.False(t, parser.HasExclude(mod, "github.com/nonexistent/pkg", "v1.0.0"))

	// 测试 HasRetract
	assert.True(t, parser.HasRetract(mod, "v1.0.0"))
	assert.True(t, parser.HasRetract(mod, "v0.5.0")) // 应该在范围内
	assert.False(t, parser.HasRetract(mod, "v2.0.0"))
}
