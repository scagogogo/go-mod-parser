package parser

import (
	"strings"

	"github.com/scagogogo/go-mod-parser/pkg/module"
)

// parseRequireSingleLine 解析单行require语句
func parseRequireSingleLine(mod *module.Module, line string) (bool, error) {
	if matches := singleRequireRegexp.FindStringSubmatch(line); len(matches) >= 3 {
		// 检查是否有indirect注释
		indirect := false
		if len(matches) > 3 && matches[3] != "" {
			indirect = indirectCommentRegexp.MatchString(matches[3])
		}

		mod.Requires = append(mod.Requires, &module.Require{
			Path:     matches[1],
			Version:  matches[2],
			Indirect: indirect,
		})
		return true, nil
	}
	return false, nil
}

// parseRequireBlockLine 解析require块内的语句
func parseRequireBlockLine(mod *module.Module, line string) error {
	parts := strings.Fields(line)
	if len(parts) < 2 {
		return ErrInvalidRequire
	}

	// 检查是否有indirect注释
	indirect := false
	for i := 2; i < len(parts); i++ {
		if strings.Contains(parts[i], "//") {
			commentText := strings.Join(parts[i:], " ")
			indirect = indirectCommentRegexp.MatchString(commentText)
			break
		}
	}

	mod.Requires = append(mod.Requires, &module.Require{
		Path:     parts[0],
		Version:  parts[1],
		Indirect: indirect,
	})
	return nil
}
