package parser

import (
	"strings"

	"github.com/scagogogo/go-mod-parser/pkg/module"
)

// parseRetractSingleLine 解析单行retract语句
func parseRetractSingleLine(mod *module.Module, line string) (bool, error) {
	// 首先尝试匹配单个版本的retract
	if matches := singleRetractVersionRegexp.FindStringSubmatch(line); len(matches) == 3 {
		retract := &module.Retract{
			Version: matches[1],
		}

		// 检查是否有理由说明
		if rationaleMatch := rationaleRegexp.FindStringSubmatch(matches[2]); len(rationaleMatch) == 2 {
			retract.Rationale = strings.TrimSpace(rationaleMatch[1])
		}

		mod.Retracts = append(mod.Retracts, retract)
		return true, nil
	}

	// 然后尝试匹配版本范围的retract
	if matches := singleRetractVersionRangeRegexp.FindStringSubmatch(line); len(matches) == 4 {
		retract := &module.Retract{
			VersionLow:  matches[1],
			VersionHigh: matches[2],
		}

		// 检查是否有理由说明
		if rationaleMatch := rationaleRegexp.FindStringSubmatch(matches[3]); len(rationaleMatch) == 2 {
			retract.Rationale = strings.TrimSpace(rationaleMatch[1])
		}

		mod.Retracts = append(mod.Retracts, retract)
		return true, nil
	}

	return false, nil
}

// parseRetractBlockLine 解析retract块内的语句
func parseRetractBlockLine(mod *module.Module, line string) error {
	line = strings.TrimSpace(line)

	// 跳过空行和纯注释行
	if line == "" || strings.HasPrefix(line, "//") {
		return nil
	}

	// 尝试匹配单个版本
	if strings.Contains(line, "[") && strings.Contains(line, "]") {
		// 解析 [v1.0.0, v1.9.9] 格式
		rangeStart := strings.Index(line, "[")
		rangeEnd := strings.Index(line, "]")
		if rangeStart >= 0 && rangeEnd > rangeStart {
			versionRange := line[rangeStart+1 : rangeEnd]
			versions := strings.Split(versionRange, ",")
			if len(versions) == 2 {
				retract := &module.Retract{
					VersionLow:  strings.TrimSpace(versions[0]),
					VersionHigh: strings.TrimSpace(versions[1]),
				}

				// 获取理由（如果有）
				if rangeEnd+1 < len(line) {
					rest := strings.TrimSpace(line[rangeEnd+1:])
					if strings.HasPrefix(rest, "//") {
						rationaleMatches := rationaleRegexp.FindStringSubmatch(rest)
						if len(rationaleMatches) > 1 {
							retract.Rationale = strings.TrimSpace(rationaleMatches[1])
						}
					}
				}

				mod.Retracts = append(mod.Retracts, retract)
				return nil
			}
		}
		// 版本范围格式错误
		return ErrInvalidRetract
	} else {
		// 如果不是版本范围，则为单个版本
		parts := strings.Fields(line)
		if len(parts) < 1 {
			return ErrInvalidRetract
		}

		// 检查是否为有效版本号 (简单检查是否以v开头)
		if len(parts[0]) < 1 || (parts[0][0] != 'v' && parts[0][0] != 'V') {
			return ErrInvalidRetract
		}

		retract := &module.Retract{
			Version: parts[0],
		}

		// 提取理由（如果有）
		for i := 1; i < len(parts); i++ {
			if strings.HasPrefix(parts[i], "//") {
				commentText := strings.Join(parts[i:], " ")
				rationaleMatches := rationaleRegexp.FindStringSubmatch(commentText)
				if len(rationaleMatches) > 1 {
					retract.Rationale = strings.TrimSpace(rationaleMatches[1])
				}
				break
			}
		}

		mod.Retracts = append(mod.Retracts, retract)
		return nil
	}
}
