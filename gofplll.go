package gofplll

import (
	"context"
	"fmt"
	"math/big"
	"os/exec"
	"strings"
	"time"

	"github.com/maqsatto/gofplll/internal/encoding"
	"github.com/maqsatto/gofplll/internal/runner"
	"github.com/maqsatto/gofplll/internal/validation"
)

// Reducer is the interface for lattice reduction backends.
type Reducer interface {
	Reduce(ctx context.Context, matrix Matrix, opts ...ReduceOption) (*ReduceResult, error)
}


// Client is a Reducer backed by the fplll subprocess.
type Client struct {
	binaryPath string
	workDir    string
	r          runner.Runner
}

// New creates a Client with the given binary path.
// If path is empty, fplll is looked up in $PATH.
func New(path string, opts ...ClientOption) (*Client, error) {
	cfg := clientConfig{binaryPath: path}
	for _, o := range opts {
		o.applyClient(&cfg)
	}

	if cfg.binaryPath == "" {
		p, err := exec.LookPath("fplll")
		if err != nil {
			return nil, ErrBinaryNotFound
		}
		cfg.binaryPath = p
	}

	return &Client{
		binaryPath: cfg.binaryPath,
		workDir:    cfg.workDir,
		r:          &runner.LinuxRunner{},
	}, nil
}

// NewDefault creates a Client using the system fplll on $PATH.
func NewDefault(opts ...ClientOption) (*Client, error) {
	return New("", opts...)
}

// Reduce runs fplll on the given matrix with the provided options.
func (c *Client) Reduce(ctx context.Context, matrix Matrix, opts ...ReduceOption) (*ReduceResult, error) {
	cfg := applyReduceOptions(opts)

	algo := string(cfg.algorithm)
	if algo == "" {
		algo = "lll"
	}

	if err := validation.ValidateOptions(
		[][]*big.Int(matrix),
		algo,
		cfg.blockSize,
		cfg.delta,
		cfg.eta,
		string(cfg.floatType),
		cfg.precision,
	); err != nil {
		return nil, wrapValidationErr(err)
	}

	input, err := encoding.Encode([][]*big.Int(matrix))
	if err != nil {
		return nil, wrapValidationErr(err)
	}

	if cfg.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, cfg.timeout)
		defer cancel()
	}

	maxTimeSecs := 0
	if cfg.maxTime > 0 {
		maxTimeSecs = int(cfg.maxTime.Seconds())
		if maxTimeSecs < 1 {
			maxTimeSecs = 1
		}
	}

	args := runner.BuildArgs(runner.ReduceConfig{
		Algorithm:    algo,
		Delta:        cfg.delta,
		Eta:          cfg.eta,
		BlockSize:    cfg.blockSize,
		MaxLoops:     cfg.maxLoops,
		MaxTimeSecs:  maxTimeSecs,
		AutoAbort:    cfg.autoAbort,
		NoLLL:        cfg.noLLL,
		FloatType:    string(cfg.floatType),
		Precision:    cfg.precision,
		OutputFormat: string(cfg.outputFormat),
		Verbose:      cfg.verbose,
	})

	start := time.Now()
	stdout, stderr, exitCode, err := c.r.Run(ctx, c.binaryPath, args, input)
	elapsed := time.Since(start)

	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return &ReduceResult{
				Algorithm: Algorithm(algo),
				Stdout:    string(stdout),
				Stderr:    string(stderr),
				ExitCode:  -1,
				Runtime:   elapsed,
			}, fmt.Errorf("%w: %v", ErrTimeout, ctx.Err())
		}
		return nil, fmt.Errorf("%w: %v", ErrSubprocessFailed, err)
	}

	result := &ReduceResult{
		Algorithm: Algorithm(algo),
		Stdout:    string(stdout),
		Stderr:    string(stderr),
		ExitCode:  exitCode,
		Runtime:   elapsed,
	}

	out := strings.TrimSpace(string(stdout))
	if cfg.outputFormat == OutputStatus || cfg.outputFormat == OutputSVP || cfg.outputFormat == OutputCVP {
		result.Vector = parseVector(out)
	} else {
		parsed, err := encoding.Decode(out)
		if err == nil {
			result.Matrix = Matrix(parsed)
			result.Rows = len(parsed)
			if result.Rows > 0 {
				result.Cols = len(parsed[0])
			}
		}
	}

	return result, nil
}

func parseVector(s string) []*big.Int {
	fields := strings.Fields(s)
	var vec []*big.Int
	for _, f := range fields {
		v := new(big.Int)
		if _, ok := v.SetString(f, 10); ok {
			vec = append(vec, v)
		}
	}
	return vec
}

func wrapValidationErr(err error) error {
	switch {
	case strings.Contains(err.Error(), "invalid matrix"):
		return fmt.Errorf("%w: %v", ErrInvalidMatrix, err)
	case strings.Contains(err.Error(), "invalid options"):
		return fmt.Errorf("%w: %v", ErrInvalidOptions, err)
	default:
		return err
	}
}
