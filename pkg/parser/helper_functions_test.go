package parser_test

import (
	"testing"

	"github.com/scagogogo/go-mod-parser/pkg/module"
	"github.com/scagogogo/go-mod-parser/pkg/parser"
	"github.com/stretchr/testify/assert"
)

func TestHelperFunctions(t *testing.T) {
	// 创建一个测试模块
	mod := &module.Module{
		Name:      "example.com/test",
		GoVersion: "1.18",
		Requires: []*module.Require{
			{Path: "github.com/test/req1", Version: "v1.0.0"},
			{Path: "github.com/test/req2", Version: "v2.0.0", Indirect: true},
		},
		Replaces: []*module.Replace{
			{
				Old: &module.ReplaceItem{Path: "github.com/test/old1"},
				New: &module.ReplaceItem{Path: "github.com/test/new1", Version: "v1.0.0"},
			},
		},
		Excludes: []*module.Exclude{
			{Path: "github.com/test/exclude1", Version: "v1.0.0"},
		},
		Retracts: []*module.Retract{
			{Version: "v1.0.0", Rationale: "buggy version"},
			{VersionLow: "v2.0.0", VersionHigh: "v2.9.9", Rationale: "security issue"},
		},
	}

	// 测试 HasRequire 和 GetRequire
	assert.True(t, parser.HasRequire(mod, "github.com/test/req1"))
	assert.False(t, parser.HasRequire(mod, "github.com/test/nonexistent"))

	req := parser.GetRequire(mod, "github.com/test/req2")
	assert.NotNil(t, req)
	assert.Equal(t, "v2.0.0", req.Version)
	assert.True(t, req.Indirect)

	assert.Nil(t, parser.GetRequire(mod, "github.com/test/nonexistent"))

	// 测试 HasReplace 和 GetReplace
	assert.True(t, parser.HasReplace(mod, "github.com/test/old1"))
	assert.False(t, parser.HasReplace(mod, "github.com/test/nonexistent"))

	rep := parser.GetReplace(mod, "github.com/test/old1")
	assert.NotNil(t, rep)
	assert.Equal(t, "github.com/test/new1", rep.New.Path)

	assert.Nil(t, parser.GetReplace(mod, "github.com/test/nonexistent"))

	// 测试 HasExclude
	assert.True(t, parser.HasExclude(mod, "github.com/test/exclude1", "v1.0.0"))
	assert.False(t, parser.HasExclude(mod, "github.com/test/exclude1", "v2.0.0"))
	assert.False(t, parser.HasExclude(mod, "github.com/test/nonexistent", "v1.0.0"))

	// 测试 HasRetract
	assert.True(t, parser.HasRetract(mod, "v1.0.0"))
	assert.True(t, parser.HasRetract(mod, "v2.5.0")) // 在范围 [v2.0.0, v2.9.9] 内
	assert.False(t, parser.HasRetract(mod, "v3.0.0"))
}
