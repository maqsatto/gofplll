package runner

import (
	"context"
	"fmt"
)

// Runner executes an external command and captures output.
type Runner interface {
	Run(ctx context.Context, binaryPath string, args []string, stdin []byte) (stdout, stderr []byte, exitCode int, err error)
}

// ReduceConfig holds the parameters needed to build CLI args.
type ReduceConfig struct {
	Algorithm    string
	Delta        string
	Eta          string
	BlockSize    int
	MaxLoops     int
	MaxTimeSecs  int
	AutoAbort    bool
	NoLLL        bool
	FloatType    string
	Precision    int
	OutputFormat string
	Verbose      bool
}

// BuildArgs converts a ReduceConfig into CLI argument slices.
func BuildArgs(cfg ReduceConfig) []string {
	var args []string

	if cfg.Algorithm != "" {
		args = append(args, "-a", cfg.Algorithm)
	}
	if cfg.Delta != "" {
		args = append(args, "-d", cfg.Delta)
	}
	if cfg.Eta != "" {
		args = append(args, "-e", cfg.Eta)
	}
	if cfg.FloatType != "" {
		args = append(args, "-f", cfg.FloatType)
	}
	if cfg.Precision > 0 {
		args = append(args, "-p", fmt.Sprintf("%d", cfg.Precision))
	}
	if cfg.OutputFormat != "" {
		args = append(args, "-of", cfg.OutputFormat)
	}
	if cfg.Verbose {
		args = append(args, "-v")
	}
	if cfg.NoLLL {
		args = append(args, "-nolll")
	}
	if cfg.Algorithm == "bkz" {
		if cfg.BlockSize > 0 {
			args = append(args, "-b", fmt.Sprintf("%d", cfg.BlockSize))
		}
		if cfg.MaxLoops > 0 {
			args = append(args, "-bkzmaxloops", fmt.Sprintf("%d", cfg.MaxLoops))
		}
		if cfg.MaxTimeSecs > 0 {
			args = append(args, "-bkzmaxtime", fmt.Sprintf("%d", cfg.MaxTimeSecs))
		}
		if cfg.AutoAbort {
			args = append(args, "-bkzautoabort")
		}
	}

	return args
}
