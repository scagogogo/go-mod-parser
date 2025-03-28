package parser

import (
	"errors"
)

var (
	// ErrInvalidModuleDeclaration 表示无法解析module声明
	ErrInvalidModuleDeclaration = errors.New("invalid module declaration")
	// ErrInvalidGoVersion 表示无法解析go版本声明
	ErrInvalidGoVersion = errors.New("invalid go version declaration")
	// ErrInvalidRequire 表示无法解析require声明
	ErrInvalidRequire = errors.New("invalid require declaration")
	// ErrInvalidReplace 表示无法解析replace声明
	ErrInvalidReplace = errors.New("invalid replace declaration")
	// ErrInvalidExclude 表示无法解析exclude声明
	ErrInvalidExclude = errors.New("invalid exclude declaration")
	// ErrInvalidRetract 表示无法解析retract声明
	ErrInvalidRetract = errors.New("invalid retract declaration")
)
