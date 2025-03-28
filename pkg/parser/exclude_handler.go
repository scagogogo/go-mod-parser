package parser

import (
	"strings"

	"github.com/scagogogo/go-mod-parser/pkg/module"
)

// parseExcludeSingleLine 解析单行exclude语句
func parseExcludeSingleLine(mod *module.Module, line string) (bool, error) {
	if matches := singleExcludeRegexp.FindStringSubmatch(line); len(matches) == 3 {
		mod.Excludes = append(mod.Excludes, &module.Exclude{
			Path:    matches[1],
			Version: matches[2],
		})
		return true, nil
	}
	return false, nil
}

// parseExcludeBlockLine 解析exclude块内的语句
func parseExcludeBlockLine(mod *module.Module, line string) error {
	parts := strings.Fields(line)
	if len(parts) < 2 {
		return ErrInvalidExclude
	}
	mod.Excludes = append(mod.Excludes, &module.Exclude{
		Path:    parts[0],
		Version: parts[1],
	})
	return nil
}
