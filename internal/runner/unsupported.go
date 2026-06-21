//go:build !linux

package runner

import (
	"context"
	"errors"
)

var errUnsupported = errors.New("runner: unsupported OS: only Linux is supported")

// UnsupportedRunner returns an error on non-Linux platforms.
type UnsupportedRunner struct{}

func (r *UnsupportedRunner) Run(ctx context.Context, binaryPath string, args []string, stdin []byte) ([]byte, []byte, int, error) {
	return nil, nil, -1, errUnsupported
}
