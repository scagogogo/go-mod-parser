package pkg

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/scagogogo/go-mod-parser/pkg/module"
	"github.com/stretchr/testify/assert"
)

func TestParseGoModContent(t *testing.T) {
	content := `module github.com/example/module

go 1.21

require github.com/stretchr/testify v1.8.4
`
	mod, err := ParseGoModContent(content)
	if err != nil {
		t.Fatalf("ParseGoModContent failed: %v", err)
	}

	if mod.Name != "github.com/example/module" {
		t.Errorf("Expected module name 'github.com/example/module', got %q", mod.Name)
	}

	if mod.GoVersion != "1.21" {
		t.Errorf("Expected Go version '1.21', got %q", mod.GoVersion)
	}

	if len(mod.Requires) != 1 {
		t.Fatalf("Expected 1 require, got %d", len(mod.Requires))
	}

	if mod.Requires[0].Path != "github.com/stretchr/testify" {
		t.Errorf("Expected require path 'github.com/stretchr/testify', got %q", mod.Requires[0].Path)
	}

	if mod.Requires[0].Version != "v1.8.4" {
		t.Errorf("Expected require version 'v1.8.4', got %q", mod.Requires[0].Version)
	}
}

func TestParseGoModFile(t *testing.T) {
	// 创建临时go.mod文件
	tempDir := t.TempDir()
	goModPath := filepath.Join(tempDir, "go.mod")
	content := `module github.com/example/module

go 1.21

require github.com/stretchr/testify v1.8.4
`
	err := os.WriteFile(goModPath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to write test go.mod file: %v", err)
	}

	// 测试解析文件
	mod, err := ParseGoModFile(goModPath)
	if err != nil {
		t.Fatalf("ParseGoModFile failed: %v", err)
	}

	if mod.Name != "github.com/example/module" {
		t.Errorf("Expected module name 'github.com/example/module', got %q", mod.Name)
	}
}

func TestFindAndParseGoModFile(t *testing.T) {
	// 创建临时目录结构
	tempDir := t.TempDir()
	subDirPath := filepath.Join(tempDir, "level1", "level2")
	err := os.MkdirAll(subDirPath, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory structure: %v", err)
	}

	// 在顶层创建go.mod文件
	goModPath := filepath.Join(tempDir, "go.mod")
	content := `module github.com/example/module

go 1.21
`
	err = os.WriteFile(goModPath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to write test go.mod file: %v", err)
	}

	// 从子目录查找并解析
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(oldDir) // 恢复原目录

	err = os.Chdir(subDirPath)
	if err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	mod, err := FindAndParseGoModFile(".")
	if err != nil {
		t.Fatalf("FindAndParseGoModFile failed: %v", err)
	}

	if mod.Name != "github.com/example/module" {
		t.Errorf("Expected module name 'github.com/example/module', got %q", mod.Name)
	}

	// 测试无法找到go.mod的情况
	emptyDir := filepath.Join(tempDir, "empty")
	err = os.Mkdir(emptyDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create empty test directory: %v", err)
	}

	// 创建一个完全隔离的目录，确保不会从父目录找到go.mod
	isolatedDir := filepath.Join(os.TempDir(), "isolated_test_dir_"+filepath.Base(t.Name()))
	defer os.RemoveAll(isolatedDir) // 确保测试后清理

	err = os.MkdirAll(isolatedDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create isolated test directory: %v", err)
	}

	err = os.Chdir(isolatedDir)
	if err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	_, err = FindAndParseGoModFile(".")
	if err == nil {
		t.Error("Expected error when no go.mod file found, got nil")
	}
}

func TestHelperFunctions(t *testing.T) {
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
	assert.True(t, HasRequire(mod, "github.com/stretchr/testify"))
	assert.True(t, HasRequire(mod, "github.com/pkg/errors"))
	assert.False(t, HasRequire(mod, "github.com/nonexistent/pkg"))

	// 测试 GetRequire
	req := GetRequire(mod, "github.com/stretchr/testify")
	assert.NotNil(t, req)
	assert.Equal(t, "v1.8.4", req.Version)
	assert.False(t, req.Indirect)

	req = GetRequire(mod, "github.com/pkg/errors")
	assert.NotNil(t, req)
	assert.Equal(t, "v0.9.1", req.Version)
	assert.True(t, req.Indirect)

	req = GetRequire(mod, "github.com/nonexistent/pkg")
	assert.Nil(t, req)

	// 测试 HasReplace
	assert.True(t, HasReplace(mod, "github.com/old/pkg"))
	assert.False(t, HasReplace(mod, "github.com/nonexistent/pkg"))

	// 测试 GetReplace
	rep := GetReplace(mod, "github.com/old/pkg")
	assert.NotNil(t, rep)
	assert.Equal(t, "v1.0.0", rep.Old.Version)
	assert.Equal(t, "github.com/new/pkg", rep.New.Path)
	assert.Equal(t, "v2.0.0", rep.New.Version)

	rep = GetReplace(mod, "github.com/nonexistent/pkg")
	assert.Nil(t, rep)

	// 测试 HasExclude
	assert.True(t, HasExclude(mod, "github.com/bad/pkg", "v0.1.0"))
	assert.False(t, HasExclude(mod, "github.com/bad/pkg", "v0.2.0"))
	assert.False(t, HasExclude(mod, "github.com/nonexistent/pkg", "v1.0.0"))

	// 测试 HasRetract
	assert.True(t, HasRetract(mod, "v1.0.0"))
	assert.True(t, HasRetract(mod, "v0.5.0")) // 应该在范围内
	assert.False(t, HasRetract(mod, "v2.0.0"))
}
