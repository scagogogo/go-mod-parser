package parser

import "regexp"

// moduleRegexp 匹配module声明
var moduleRegexp = regexp.MustCompile(`^module\s+([^\s]+)$`)

// goRegexp 匹配go版本声明
var goRegexp = regexp.MustCompile(`^go\s+([^\s]+)$`)

// singleRequireRegexp 匹配单行require声明
var singleRequireRegexp = regexp.MustCompile(`^require\s+([^\s]+)\s+([^\s]+)(.*)$`)

// singleReplaceRegexp 匹配单行replace声明
var singleReplaceRegexp = regexp.MustCompile(`^replace\s+([^\s]+)\s+=>\s+([^\s]+)\s+([^\s]+)$`)

// singleExcludeRegexp 匹配单行exclude声明
var singleExcludeRegexp = regexp.MustCompile(`^exclude\s+([^\s]+)\s+([^\s]+)$`)

// singleRetractVersionRegexp 匹配单行retract声明（单个版本）
var singleRetractVersionRegexp = regexp.MustCompile(`^retract\s+([^\s\[\]]+)(.*)$`)

// singleRetractVersionRangeRegexp 匹配单行retract声明（版本范围）
var singleRetractVersionRangeRegexp = regexp.MustCompile(`^retract\s+\[\s*([^\s,]+)\s*,\s*([^\s,\]]+)\s*\](.*)$`)

// indirectCommentRegexp 匹配indirect注释
var indirectCommentRegexp = regexp.MustCompile(`//\s*indirect`)

// rationaleRegexp 匹配retract理由
var rationaleRegexp = regexp.MustCompile(`//\s*(.+)`)
