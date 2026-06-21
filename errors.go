package gofplll

import "errors"

var (
	ErrUnsupportedOS    = errors.New("gofplll: unsupported OS: only Linux is supported")
	ErrBinaryNotFound   = errors.New("gofplll: fplll binary not found")
	ErrInvalidMatrix    = errors.New("gofplll: invalid matrix")
	ErrInvalidOptions   = errors.New("gofplll: invalid options")
	ErrParserFailed     = errors.New("gofplll: failed to parse fplll output")
	ErrSubprocessFailed = errors.New("gofplll: subprocess execution failed")
	ErrTimeout          = errors.New("gofplll: subprocess timed out")
)
