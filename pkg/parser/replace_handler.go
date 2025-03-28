package parser

import (
	"strings"

	"github.com/scagogogo/go-mod-parser/pkg/module"
)

// parseReplaceSingleLine 解析单行replace语句
func parseReplaceSingleLine(mod *module.Module, line string) (bool, error) {
	if matches := singleReplaceRegexp.FindStringSubmatch(line); len(matches) == 4 {
		mod.Replaces = append(mod.Replaces, &module.Replace{
			Old: &module.ReplaceItem{
				Path: matches[1],
			},
			New: &module.ReplaceItem{
				Path:    matches[2],
				Version: matches[3],
			},
		})
		return true, nil
	}
	return false, nil
}

// parseReplaceBlockLine 解析replace块内的语句
func parseReplaceBlockLine(mod *module.Module, line string) error {
	parts := strings.Split(line, "=>")
	if len(parts) != 2 {
		return ErrInvalidReplace
	}

	oldParts := strings.Fields(strings.TrimSpace(parts[0]))
	newParts := strings.Fields(strings.TrimSpace(parts[1]))

	oldPath := oldParts[0]
	var oldVersion string
	if len(oldParts) > 1 {
		oldVersion = oldParts[1]
	}

	if len(newParts) < 1 {
		return ErrInvalidReplace
	}

	newPath := newParts[0]
	var newVersion string
	if len(newParts) > 1 {
		newVersion = newParts[1]
	}

	mod.Replaces = append(mod.Replaces, &module.Replace{
		Old: &module.ReplaceItem{
			Path:    oldPath,
			Version: oldVersion,
		},
		New: &module.ReplaceItem{
			Path:    newPath,
			Version: newVersion,
		},
	})
	return nil
}
