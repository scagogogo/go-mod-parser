package pkg

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseGoModFile(t *testing.T) {
	// 创建临时go.mod文件
	tempDir := t.TempDir()
	goModPath := filepath.Join(tempDir, "go.mod")
	content := `module github.com/example/test

go 1.21

require (
	github.com/stretchr/testify v1.8.4
	github.com/example/dep v1.0.0 // indirect
)

replace github.com/old/pkg => github.com/new/pkg v1.0.0

exclude github.com/bad/pkg v1.0.0

retract v1.0.1 // security issue
`
	err := os.WriteFile(goModPath, []byte(content), 0644)
	require.NoError(t, err)

	// 测试解析文件
	mod, err := ParseGoModFile(goModPath)
	require.NoError(t, err)
	assert.Equal(t, "github.com/example/test", mod.Name)
	assert.Equal(t, "1.21", mod.GoVersion)
	assert.Len(t, mod.Requires, 2)
	assert.Len(t, mod.Replaces, 1)
	assert.Len(t, mod.Excludes, 1)
	assert.Len(t, mod.Retracts, 1)
}

func TestParseGoModFile_NotFound(t *testing.T) {
	// 测试文件不存在的情况
	_, err := ParseGoModFile("/nonexistent/go.mod")
	assert.Error(t, err)
}

func TestParseGoModContent(t *testing.T) {
	content := `module github.com/example/test

go 1.21

require github.com/stretchr/testify v1.8.4
`
	mod, err := ParseGoModContent(content)
	require.NoError(t, err)
	assert.Equal(t, "github.com/example/test", mod.Name)
	assert.Equal(t, "1.21", mod.GoVersion)
	assert.Len(t, mod.Requires, 1)
	assert.Equal(t, "github.com/stretchr/testify", mod.Requires[0].Path)
}

func TestParseGoModContent_Invalid(t *testing.T) {
	// 测试无效内容
	_, err := ParseGoModContent("invalid go.mod content")
	assert.Error(t, err)
}

func TestFindAndParseGoModFile(t *testing.T) {
	// 创建临时目录结构
	tempDir := t.TempDir()
	subDir := filepath.Join(tempDir, "subdir")
	err := os.Mkdir(subDir, 0755)
	require.NoError(t, err)

	// 在父目录创建go.mod文件
	goModPath := filepath.Join(tempDir, "go.mod")
	content := `module github.com/example/test

go 1.21
`
	err = os.WriteFile(goModPath, []byte(content), 0644)
	require.NoError(t, err)

	// 从子目录查找并解析
	mod, err := FindAndParseGoModFile(subDir)
	require.NoError(t, err)
	assert.Equal(t, "github.com/example/test", mod.Name)
}

func TestFindAndParseGoModFile_NotFound(t *testing.T) {
	// 测试在没有go.mod的目录中查找
	tempDir := t.TempDir()
	_, err := FindAndParseGoModFile(tempDir)
	assert.Error(t, err)
}

func TestFindAndParseGoModInCurrentDir(t *testing.T) {
	// 保存当前工作目录
	originalWd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(originalWd)

	// 创建临时目录并切换到该目录
	tempDir := t.TempDir()
	err = os.Chdir(tempDir)
	require.NoError(t, err)

	// 在当前目录创建go.mod文件
	content := `module github.com/example/test

go 1.21
`
	err = os.WriteFile("go.mod", []byte(content), 0644)
	require.NoError(t, err)

	// 测试在当前目录查找并解析
	mod, err := FindAndParseGoModInCurrentDir()
	require.NoError(t, err)
	assert.Equal(t, "github.com/example/test", mod.Name)
}

func TestHelperFunctions(t *testing.T) {
	content := `module github.com/example/test

go 1.21

require (
	github.com/stretchr/testify v1.8.4
	github.com/example/dep v1.0.0 // indirect
)

replace github.com/old/pkg => github.com/new/pkg v1.0.0

exclude github.com/bad/pkg v1.0.0

retract v1.0.1 // security issue
`
	mod, err := ParseGoModContent(content)
	require.NoError(t, err)

	// 测试HasRequire和GetRequire
	assert.True(t, HasRequire(mod, "github.com/stretchr/testify"))
	assert.False(t, HasRequire(mod, "github.com/nonexistent/pkg"))

	req := GetRequire(mod, "github.com/stretchr/testify")
	require.NotNil(t, req)
	assert.Equal(t, "v1.8.4", req.Version)
	assert.False(t, req.Indirect)

	req = GetRequire(mod, "github.com/example/dep")
	require.NotNil(t, req)
	assert.True(t, req.Indirect)

	// 测试HasReplace和GetReplace
	assert.True(t, HasReplace(mod, "github.com/old/pkg"))
	assert.False(t, HasReplace(mod, "github.com/nonexistent/pkg"))

	rep := GetReplace(mod, "github.com/old/pkg")
	require.NotNil(t, rep)
	assert.Equal(t, "github.com/new/pkg", rep.New.Path)
	assert.Equal(t, "v1.0.0", rep.New.Version)

	// 测试HasExclude
	assert.True(t, HasExclude(mod, "github.com/bad/pkg", "v1.0.0"))
	assert.False(t, HasExclude(mod, "github.com/bad/pkg", "v2.0.0"))
	assert.False(t, HasExclude(mod, "github.com/nonexistent/pkg", "v1.0.0"))

	// 测试HasRetract
	assert.True(t, HasRetract(mod, "v1.0.1"))
	assert.False(t, HasRetract(mod, "v1.0.0"))
}
