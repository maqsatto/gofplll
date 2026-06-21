//go:build linux

package runner

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
)

// LinuxRunner executes commands via exec.CommandContext.
type LinuxRunner struct{}

func (r *LinuxRunner) Run(ctx context.Context, binaryPath string, args []string, stdin []byte) ([]byte, []byte, int, error) {
	cmd := exec.CommandContext(ctx, binaryPath, args...)
	cmd.Stdin = bytes.NewReader(stdin)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	exitCode := 0
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return stdout.Bytes(), stderr.Bytes(), -1, fmt.Errorf("runner: timed out: %w", ctx.Err())
		}
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			return stdout.Bytes(), stderr.Bytes(), 0, fmt.Errorf("runner: exec failed: %w", err)
		}
	}

	return stdout.Bytes(), stderr.Bytes(), exitCode, nil
}
