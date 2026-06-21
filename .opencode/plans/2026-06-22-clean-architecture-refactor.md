# gofplll Clean Architecture Refactor — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Refactor gofplll from a flat package into a clean architecture with internal packages, functional options, a Reducer interface, and simplified API.

**Architecture:** Root package exposes `Reducer` interface, `Client`, `Matrix`, enums, and functional option constructors. Internal packages hold encoding, validation, and subprocess runner logic. Platform-specific code isolated via build tags in `internal/runner/`.

**Tech Stack:** Go 1.22+, standard library only. No external dependencies.

## Global Constraints

- Linux-only library with build tags (`linux` / `!linux`)
- Matrix coefficients must use `math/big.Int` — no `int64` or `float64`
- No shell execution — use `exec.CommandContext` directly
- No cgo, no Python, no external Go dependencies
- No HNP/ECDSA/Bitcoin/attack logic
- Unit tests pass without fplll installed; integration tests skip gracefully

## File Map

| File | Responsibility |
|------|---------------|
| `gofplll.go` | `Reducer` interface, `Client` struct, `New`, `NewDefault`, `Reduce` |
| `matrix.go` | `Matrix` type, `NewMatrix`, `Set`, `Get`, `Rows`, `Cols` helpers |
| `algorithm.go` | `Algorithm`, `FloatType`, `OutputFormat` enums with `String()` |
| `errors.go` | Sentinel errors |
| `options.go` | `ClientOption`, `ReduceOption` interfaces and `With*` constructors |
| `result.go` | `ReduceResult` struct |
| `internal/encoding/codec.go` | `Encode`, `Decode` implementations |
| `internal/encoding/codec_test.go` | Unit tests for encode/decode |
| `internal/validation/validate.go` | `ValidateMatrix`, `ValidateOptions` |
| `internal/validation/validate_test.go` | Unit tests for validation |
| `internal/runner/runner.go` | `Runner` interface, `ArgsBuilder` |
| `internal/runner/linux.go` | `LinuxRunner` (exec.CommandContext) |
| `internal/runner/unsupported.go` | `UnsupportedRunner` |
| `cmd/gofplll-smoke/main.go` | Smoke CLI |
| `examples/lll_basic/main.go` | LLL example |
| `examples/bkz_basic/main.go` | BKZ example |

---

### Task 1: Create internal/encoding package

**Files:**
- Create: `internal/encoding/codec.go`
- Create: `internal/encoding/codec_test.go`

**Interfaces:**
- Consumes: `Matrix` type (will be `[][]*big.Int` — defined in root, imported here)
- Produces: `Encode(m Matrix) ([]byte, error)`, `Decode(s string) (Matrix, error)`

- [ ] **Step 1: Write the failing tests**

Create `internal/encoding/codec_test.go`:

```go
package encoding

import (
	"math/big"
	"testing"
)

func TestEncode2x2(t *testing.T) {
	m := Matrix{
		{big.NewInt(10), big.NewInt(11)},
		{big.NewInt(11), big.NewInt(12)},
	}
	got, err := Encode(m)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "[[10 11]\n [11 12]]"
	if string(got) != want {
		t.Errorf("Encode() = %q, want %q", string(got), want)
	}
}

func TestEncodeNegative(t *testing.T) {
	m := Matrix{
		{big.NewInt(-1), big.NewInt(2)},
		{big.NewInt(3), big.NewInt(-4)},
	}
	got, err := Encode(m)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "[[-1 2]\n [3 -4]]"
	if string(got) != want {
		t.Errorf("Encode() = %q, want %q", string(got), want)
	}
}

func TestEncodeBigInt(t *testing.T) {
	bigVal := new(big.Int)
	bigVal.SetString("123456789012345678901234567890", 10)
	m := Matrix{
		{bigVal, big.NewInt(1)},
		{big.NewInt(1), bigVal},
	}
	got, err := Encode(m)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "[[123456789012345678901234567890 1]\n [1 123456789012345678901234567890]]"
	if string(got) != want {
		t.Errorf("Encode() = %q, want %q", string(got), want)
	}
}

func TestDecode2x2(t *testing.T) {
	input := "[[10 11]\n [11 12]]"
	got, err := Decode(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 || len(got[0]) != 2 {
		t.Fatalf("got dimensions %dx%d, want 2x2", len(got), len(got[0]))
	}
	if got[0][0].Int64() != 10 || got[0][1].Int64() != 11 {
		t.Errorf("row 0 = [%s %s], want [10 11]", got[0][0], got[0][1])
	}
	if got[1][0].Int64() != 11 || got[1][1].Int64() != 12 {
		t.Errorf("row 1 = [%s %s], want [11 12]", got[1][0], got[1][1])
	}
}

func TestDecodeNegative(t *testing.T) {
	input := "[[-1 2]\n [3 -4]]"
	got, err := Decode(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got[0][0].Int64() != -1 || got[1][1].Int64() != -4 {
		t.Errorf("unexpected values: %v", got)
	}
}

func TestDecodeBigInt(t *testing.T) {
	input := "[[123456789012345678901234567890 1]\n [1 123456789012345678901234567890]]"
	got, err := Decode(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	bigVal := new(big.Int)
	bigVal.SetString("123456789012345678901234567890", 10)
	if got[0][0].Cmp(bigVal) != 0 {
		t.Errorf("got[0][0] = %s, want %s", got[0][0], bigVal)
	}
}

func TestDecodeExtraWhitespace(t *testing.T) {
	input := "  [[  10   11  ]\n   [  11   12  ]]  "
	got, err := Decode(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got[0][0].Int64() != 10 {
		t.Errorf("got[0][0] = %s, want 10", got[0][0])
	}
}

func TestDecodeEmpty(t *testing.T) {
	_, err := Decode("")
	if err == nil {
		t.Error("expected error for empty input")
	}
}

func TestEncodeEmptyMatrix(t *testing.T) {
	_, err := Encode(Matrix{})
	if err == nil {
		t.Error("expected error for empty matrix")
	}
}
```

Note: `Matrix` is `[][]*big.Int` — define it locally in this test file since the encoding package will define it or import it. Actually, encoding needs to know the Matrix type. Since encoding is internal and Matrix is in root, we have a circular dependency. **Solution:** Define `Matrix` as a type alias in the encoding package too, or pass `[][]*big.Int` directly. Better: encoding accepts `[][]*big.Int` directly, root converts.

Actually, let's simplify: `internal/encoding` defines no Matrix type. `Encode` takes `[][]*big.Int` and `Decode` returns `[][]*big.Int`. The root package's `Matrix` type is `[][]*big.Int`, so conversion is implicit.

Rewrite tests with `[][]*big.Int`:

```go
package encoding

import (
	"math/big"
	"testing"
)

type Matrix = [][]*big.Int

func TestEncode2x2(t *testing.T) {
	m := Matrix{
		{big.NewInt(10), big.NewInt(11)},
		{big.NewInt(11), big.NewInt(12)},
	}
	got, err := Encode(m)
	// ... same as above
}

// All tests use Matrix = [][]*big.Int
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./internal/encoding/ -v`
Expected: FAIL — package does not exist yet

- [ ] **Step 3: Create the encoding package with Encode and Decode**

Create `internal/encoding/codec.go`:

```go
package encoding

import (
	"fmt"
	"math/big"
	"strings"
)

// Matrix is [][]*big.Int.
type Matrix = [][]*big.Int

var errEmptyMatrix = fmt.Errorf("encoding: empty matrix")
var errParserFailed = fmt.Errorf("encoding: failed to parse output")

func Encode(m Matrix) ([]byte, error) {
	if len(m) == 0 {
		return nil, errEmptyMatrix
	}
	cols := len(m[0])
	if cols == 0 {
		return nil, errEmptyMatrix
	}
	for i, row := range m {
		if len(row) != cols {
			return nil, fmt.Errorf("encoding: row %d has %d columns, expected %d", i, len(row), cols)
		}
		for j, val := range row {
			if val == nil {
				return nil, fmt.Errorf("encoding: nil at [%d][%d]", i, j)
			}
		}
	}

	var b strings.Builder
	b.WriteString("[")
	for i, row := range m {
		if i > 0 {
			b.WriteString("\n ")
		}
		b.WriteString("[")
		for j, val := range row {
			if j > 0 {
				b.WriteString(" ")
			}
			b.WriteString(val.String())
		}
		b.WriteString("]")
	}
	b.WriteString("]")
	return []byte(b.String()), nil
}

func Decode(out string) (Matrix, error) {
	out = strings.TrimSpace(out)
	if len(out) == 0 {
		return nil, fmt.Errorf("%w: empty output", errParserFailed)
	}

	outerOpen := strings.Index(out, "[[")
	if outerOpen == -1 {
		return nil, fmt.Errorf("%w: missing outer brackets", errParserFailed)
	}
	start := out[outerOpen:]

	var rows Matrix
	depth := 0
	rowStart := -1

	for i := 0; i < len(start); i++ {
		switch start[i] {
		case '[':
			depth++
			if depth == 2 {
				rowStart = i + 1
			}
		case ']':
			depth--
			if depth == 1 && rowStart >= 0 {
				rowStr := strings.TrimSpace(start[rowStart:i])
				if rowStr == "" {
					rowStart = -1
					continue
				}
				fields := strings.Fields(rowStr)
				var row []*big.Int
				for _, f := range fields {
					val := new(big.Int)
					if _, ok := val.SetString(f, 10); !ok {
						return nil, fmt.Errorf("%w: cannot parse %q", errParserFailed, f)
					}
					row = append(row, val)
				}
				rows = append(rows, row)
				rowStart = -1
			}
		}
	}

	if len(rows) == 0 {
		return nil, fmt.Errorf("%w: no rows found", errParserFailed)
	}
	return rows, nil
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./internal/encoding/ -v`
Expected: all PASS

- [ ] **Step 5: Commit**

```bash
git add internal/encoding/
git commit -m "refactor: add internal/encoding package with Encode/Decode"
```

---

### Task 2: Create internal/validation package

**Files:**
- Create: `internal/validation/validate.go`
- Create: `internal/validation/validate_test.go`

**Interfaces:**
- Consumes: `Matrix` type (`[][]*big.Int`)
- Produces: `ValidateMatrix(m Matrix) error`, `ValidateOptions(m Matrix, algo string, blocksize int, delta, eta, floattype string, precision int) error`

- [ ] **Step 1: Write the failing tests**

Create `internal/validation/validate_test.go`:

```go
package validation

import "math/big"

type Matrix = [][]*big.Int

func TestValidateMatrixEmpty(t *testing.T) {
	m := Matrix{}
	if err := ValidateMatrix(m); err == nil {
		t.Error("expected error for empty matrix")
	}
}

func TestValidateMatrixRagged(t *testing.T) {
	m := Matrix{
		{big.NewInt(1), big.NewInt(2)},
		{big.NewInt(3)},
	}
	if err := ValidateMatrix(m); err == nil {
		t.Error("expected error for ragged matrix")
	}
}

func TestValidateMatrixNilCell(t *testing.T) {
	m := Matrix{
		{big.NewInt(1), nil},
		{big.NewInt(3), big.NewInt(4)},
	}
	if err := ValidateMatrix(m); err == nil {
		t.Error("expected error for nil cell")
	}
}

func TestValidateBKZBlockSizeTooSmall(t *testing.T) {
	m := Matrix{{big.NewInt(1), big.NewInt(2)}, {big.NewInt(3), big.NewInt(4)}}
	if err := ValidateOptions(m, "bkz", 1, "", "", "", 0); err == nil {
		t.Error("expected error for block size < 2")
	}
}

func TestValidateBKZBlockSizeTooLarge(t *testing.T) {
	m := Matrix{{big.NewInt(1), big.NewInt(2)}, {big.NewInt(3), big.NewInt(4)}}
	if err := ValidateOptions(m, "bkz", 5, "", "", "", 0); err == nil {
		t.Error("expected error for block size > rows")
	}
}

func TestValidateBKZBlockSizeValid(t *testing.T) {
	m := Matrix{{big.NewInt(1), big.NewInt(2)}, {big.NewInt(3), big.NewInt(4)}}
	if err := ValidateOptions(m, "bkz", 2, "", "", "", 0); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateUnknownAlgorithm(t *testing.T) {
	m := Matrix{{big.NewInt(1), big.NewInt(2)}}
	if err := ValidateOptions(m, "xyz", 0, "", "", "", 0); err == nil {
		t.Error("expected error for unknown algorithm")
	}
}

func TestValidateNegativePrecision(t *testing.T) {
	m := Matrix{{big.NewInt(1), big.NewInt(2)}}
	if err := ValidateOptions(m, "lll", 0, "", "", "", -1); err == nil {
		t.Error("expected error for negative precision")
	}
}

func TestValidateValidLLL(t *testing.T) {
	m := Matrix{{big.NewInt(1), big.NewInt(2)}, {big.NewInt(3), big.NewInt(4)}}
	if err := ValidateOptions(m, "lll", 0, "0.99", "0.51", "mpfr", 128); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./internal/validation/ -v`
Expected: FAIL — package does not exist

- [ ] **Step 3: Create the validation package**

Create `internal/validation/validate.go`:

```go
package validation

import (
	"fmt"
	"math/big"
)

type Matrix = [][]*big.Int

var (
	errInvalidMatrix  = fmt.Errorf("validation: invalid matrix")
	errInvalidOptions = fmt.Errorf("validation: invalid options")
)

var validAlgorithms = map[string]bool{
	"lll": true, "bkz": true, "hkz": true, "svp": true, "cvp": true,
}

var validFloatTypes = map[string]bool{
	"mpfr": true, "double": true, "longdouble": true, "dd": true, "qd": true,
}

func ValidateMatrix(m Matrix) error {
	if len(m) == 0 {
		return fmt.Errorf("%w: matrix must not be empty", errInvalidMatrix)
	}
	cols := len(m[0])
	if cols == 0 {
		return fmt.Errorf("%w: matrix must have at least one column", errInvalidMatrix)
	}
	for i, row := range m {
		if len(row) != cols {
			return fmt.Errorf("%w: row %d has %d columns, expected %d", errInvalidMatrix, i, len(row), cols)
		}
		for j, val := range row {
			if val == nil {
				return fmt.Errorf("%w: nil coefficient at [%d][%d]", errInvalidMatrix, i, j)
			}
		}
	}
	return nil
}

func ValidateOptions(m Matrix, algo string, blockSize int, delta, eta, floatType string, precision int) error {
	if err := ValidateMatrix(m); err != nil {
		return err
	}

	if algo != "" && !validAlgorithms[algo] {
		return fmt.Errorf("%w: unknown algorithm %q", errInvalidOptions, algo)
	}

	if algo == "bkz" {
		if blockSize < 2 {
			return fmt.Errorf("%w: BKZ block size must be >= 2, got %d", errInvalidOptions, blockSize)
		}
		if blockSize > len(m) {
			return fmt.Errorf("%w: BKZ block size %d exceeds number of rows %d", errInvalidOptions, blockSize, len(m))
		}
	}

	if delta != "" {
		if _, ok := new(big.Int).SetString(delta, 10); !ok {
			if _, ok := new(big.Float).SetString(delta); !ok {
				return fmt.Errorf("%w: invalid delta %q", errInvalidOptions, delta)
			}
		}
	}

	if eta != "" {
		if _, ok := new(big.Int).SetString(eta, 10); !ok {
			if _, ok := new(big.Float).SetString(eta); !ok {
				return fmt.Errorf("%w: invalid eta %q", errInvalidOptions, eta)
			}
		}
	}

	if floatType != "" && !validFloatTypes[floatType] {
		return fmt.Errorf("%w: unknown float type %q", errInvalidOptions, floatType)
	}

	if precision < 0 {
		return fmt.Errorf("%w: precision must be >= 0, got %d", errInvalidOptions, precision)
	}

	return nil
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./internal/validation/ -v`
Expected: all PASS

- [ ] **Step 5: Commit**

```bash
git add internal/validation/
git commit -m "refactor: add internal/validation package"
```

---

### Task 3: Create internal/runner package

**Files:**
- Create: `internal/runner/runner.go`
- Create: `internal/runner/linux.go`
- Create: `internal/runner/unsupported.go`

**Interfaces:**
- Consumes: none (standalone)
- Produces: `Runner` interface, `LinuxRunner`, `UnsupportedRunner`, `BuildArgs(...) []string`

- [ ] **Step 1: Create the runner interface**

Create `internal/runner/runner.go`:

```go
package runner

import "context"

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
```

- [ ] **Step 2: Create the Linux runner**

Create `internal/runner/linux.go`:

```go
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
```

- [ ] **Step 3: Create the unsupported runner**

Create `internal/runner/unsupported.go`:

```go
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
```

- [ ] **Step 4: Verify compilation**

Run: `go build ./internal/runner/`
Expected: compiles without errors

- [ ] **Step 5: Commit**

```bash
git add internal/runner/
git commit -m "refactor: add internal/runner package with Linux/unsupported builds"
```

---

### Task 4: Create root package types (errors, algorithm, matrix, result, options)

**Files:**
- Create: `errors.go` (replace existing)
- Create: `algorithm.go`
- Create: `matrix.go` (replace existing)
- Create: `result.go` (replace existing)
- Create: `options.go` (replace existing)

**Interfaces:**
- Consumes: nothing (foundational types)
- Produces: all types used by later tasks

- [ ] **Step 1: Rewrite errors.go**

Replace `errors.go`:

```go
package gofplll

import "errors"

var (
	ErrUnsupportedOS   = errors.New("gofplll: unsupported OS: only Linux is supported")
	ErrBinaryNotFound  = errors.New("gofplll: fplll binary not found")
	ErrInvalidMatrix   = errors.New("gofplll: invalid matrix")
	ErrInvalidOptions  = errors.New("gofplll: invalid options")
	ErrParserFailed    = errors.New("gofplll: failed to parse fplll output")
	ErrSubprocessFailed = errors.New("gofplll: subprocess execution failed")
	ErrTimeout         = errors.New("gofplll: subprocess timed out")
)
```

- [ ] **Step 2: Create algorithm.go**

Create `algorithm.go`:

```go
package gofplll

// Algorithm specifies the lattice reduction algorithm.
type Algorithm string

const (
	AlgLLL Algorithm = "lll"
	AlgBKZ Algorithm = "bkz"
	AlgHKZ Algorithm = "hkz"
	AlgSVP Algorithm = "svp"
	AlgCVP Algorithm = "cvp"
)

// FloatType specifies the floating-point backend.
type FloatType string

const (
	FloatMPFR       FloatType = "mpfr"
	FloatDouble     FloatType = "double"
	FloatLongDouble FloatType = "longdouble"
	FloatDD         FloatType = "dd"
	FloatQD         FloatType = "qd"
)

// OutputFormat specifies what fplll outputs.
type OutputFormat string

const (
	OutputBasis  OutputFormat = "b"
	OutputStatus OutputFormat = "t"
	OutputSVP    OutputFormat = "s"
	OutputCVP    OutputFormat = "c"
)
```

- [ ] **Step 3: Rewrite matrix.go**

Replace `matrix.go`:

```go
package gofplll

import "math/big"

// Matrix is a 2D slice of arbitrary-precision integers.
type Matrix [][]*big.Int

// NewMatrix creates a Matrix with the given dimensions, filled with zeros.
func NewMatrix(rows, cols int) Matrix {
	m := make(Matrix, rows)
	for i := range m {
		m[i] = make([]*big.Int, cols)
		for j := range m[i] {
			m[i][j] = new(big.Int)
		}
	}
	return m
}

// Rows returns the number of rows.
func (m Matrix) Rows() int { return len(m) }

// Cols returns the number of columns (0 if empty).
func (m Matrix) Cols() int {
	if len(m) == 0 {
		return 0
	}
	return len(m[0])
}
```

- [ ] **Step 4: Rewrite result.go**

Replace `result.go`:

```go
package gofplll

import (
	"math/big"
	"time"
)

// ReduceResult holds the output of a reduction call.
type ReduceResult struct {
	Algorithm Algorithm
	Matrix    Matrix
	Vector    []*big.Int
	Stdout    string
	Stderr    string
	ExitCode  int
	Runtime   time.Duration
	Rows      int
	Cols      int
}
```

- [ ] **Step 5: Rewrite options.go**

Replace `options.go`:

```go
package gofplll

import "time"

// ClientOption configures a Client.
type ClientOption interface {
	applyClient(*clientConfig)
}

type clientConfig struct {
	binaryPath string
	workDir    string
}

type withBinaryPath struct{ path string }

func (w withBinaryPath) applyClient(c *clientConfig) { c.binaryPath = w.path }

// WithBinaryPath sets the path to the fplll binary.
func WithBinaryPath(path string) ClientOption { return withBinaryPath{path: path} }

type withWorkDir struct{ dir string }

func (w withWorkDir) applyClient(c *clientConfig) { c.workDir = w.dir }

// WithWorkDir sets the working directory for the subprocess.
func WithWorkDir(dir string) ClientOption { return withWorkDir{dir: dir} }

// ReduceOption configures a single Reduce call.
type ReduceOption interface {
	applyReduce(*reduceConfig)
}

type reduceConfig struct {
	algorithm    Algorithm
	delta        string
	eta          string
	blockSize    int
	maxLoops     int
	maxTime      time.Duration
	autoAbort    bool
	noLLL        bool
	floatType    FloatType
	precision    int
	outputFormat OutputFormat
	verbose      bool
	timeout      time.Duration
}

type withAlgorithm struct{ v Algorithm }

func (w withAlgorithm) applyReduce(c *reduceConfig) { c.algorithm = w.v }

// WithAlgorithm sets the reduction algorithm.
func WithAlgorithm(a Algorithm) ReduceOption { return withAlgorithm{v: a} }

type withDelta struct{ v string }

func (w withDelta) applyReduce(c *reduceConfig) { c.delta = w.v }

// WithDelta sets the LLL delta parameter.
func WithDelta(d string) ReduceOption { return withDelta{v: d} }

type withEta struct{ v string }

func (w withEta) applyReduce(c *reduceConfig) { c.eta = w.v }

// WithEta sets the LLL eta parameter.
func WithEta(e string) ReduceOption { return withEta{v: e} }

type withBlockSize struct{ v int }

func (w withBlockSize) applyReduce(c *reduceConfig) { c.blockSize = w.v }

// WithBlockSize sets the BKZ block size.
func WithBlockSize(n int) ReduceOption { return withBlockSize{v: n} }

type withMaxLoops struct{ v int }

func (w withMaxLoops) applyReduce(c *reduceConfig) { c.maxLoops = w.v }

// WithMaxLoops sets the BKZ max loop count.
func WithMaxLoops(n int) ReduceOption { return withMaxLoops{v: n} }

type withMaxTime struct{ v time.Duration }

func (w withMaxTime) applyReduce(c *reduceConfig) { c.maxTime = w.v }

// WithMaxTime sets the BKZ max time limit.
func WithMaxTime(d time.Duration) ReduceOption { return withMaxTime{v: d} }

type withAutoAbort struct{}

func (w withAutoAbort) applyReduce(c *reduceConfig) { c.autoAbort = true }

// WithAutoAbort enables BKZ auto-abort.
func WithAutoAbort() ReduceOption { return withAutoAbort{} }

type withNoLLL struct{}

func (w withNoLLL) applyReduce(c *reduceConfig) { c.noLLL = true }

// WithNoLLL disables LLL pre-processing.
func WithNoLLL() ReduceOption { return withNoLLL{} }

type withFloatType struct{ v FloatType }

func (w withFloatType) applyReduce(c *reduceConfig) { c.floatType = w.v }

// WithFloatType sets the floating-point backend.
func WithFloatType(f FloatType) ReduceOption { return withFloatType{v: f} }

type withPrecision struct{ v int }

func (w withPrecision) applyReduce(c *reduceConfig) { c.precision = w.v }

// WithPrecision sets the MPFR precision in bits.
func WithPrecision(p int) ReduceOption { return withPrecision{v: p} }

type withOutputFormat struct{ v OutputFormat }

func (w withOutputFormat) applyReduce(c *reduceConfig) { c.outputFormat = w.v }

// WithOutputFormat sets what fplll returns.
func WithOutputFormat(f OutputFormat) ReduceOption { return withOutputFormat{v: f} }

type withVerbose struct{}

func (w withVerbose) applyReduce(c *reduceConfig) { c.verbose = true }

// WithVerbose enables verbose fplll output.
func WithVerbose() ReduceOption { return withVerbose{} }

type withTimeout struct{ v time.Duration }

func (w withTimeout) applyReduce(c *reduceConfig) { c.timeout = w.v }

// WithTimeout sets the subprocess timeout.
func WithTimeout(d time.Duration) ReduceOption { return withTimeout{v: d} }

func applyReduceOptions(opts []ReduceOption) reduceConfig {
	var cfg reduceConfig
	for _, o := range opts {
		o.applyReduce(&cfg)
	}
	return cfg
}
```

- [ ] **Step 6: Verify compilation**

Run: `go build ./...`
Expected: compiles (existing fplll_linux.go will have conflicts — see Task 5)

- [ ] **Step 7: Commit**

```bash
git add errors.go algorithm.go matrix.go result.go options.go
git commit -m "refactor: rewrite root package types with functional options"
```

---

### Task 5: Create gofplll.go (Reducer interface + Client)

**Files:**
- Create: `gofplll.go`
- Delete: `fplll_linux.go`, `fplll_unsupported.go`, `codec.go`, `validate.go`

**Interfaces:**
- Consumes: `internal/encoding`, `internal/validation`, `internal/runner`, all root types
- Produces: `Reducer` interface, `Client`, `New`, `NewDefault`, `Reduce`

- [ ] **Step 1: Delete old files**

```bash
rm -f fplll_linux.go fplll_unsupported.go codec.go validate.go
```

- [ ] **Step 2: Create gofplll.go**

Create `gofplll.go`:

```go
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

type clientConfig struct {
	binaryPath string
	workDir    string
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
```

- [ ] **Step 3: Delete old test files**

```bash
rm -f codec_test.go validate_test.go fplll_integration_test.go
```

- [ ] **Step 4: Create new integration test**

Create `gofplll_test.go`:

```go
package gofplll

import (
	"context"
	"math/big"
	"os/exec"
	"testing"
	"time"
)

func TestNewDefaultNotFound(t *testing.T) {
	// Temporarily rename PATH to ensure fplll is not found
	t.Setenv("PATH", "")
	_, err := NewDefault()
	if err == nil {
		t.Error("expected error when fplll not found")
	}
}

func TestReduceLLLWithFPLLL(t *testing.T) {
	if _, err := exec.LookPath("fplll"); err != nil {
		t.Skip("fplll not installed")
	}

	client, err := NewDefault()
	if err != nil {
		t.Fatalf("NewDefault: %v", err)
	}

	m := Matrix{
		{big.NewInt(10), big.NewInt(11)},
		{big.NewInt(11), big.NewInt(12)},
	}

	result, err := client.Reduce(context.Background(), m,
		WithAlgorithm(AlgLLL),
		WithTimeout(5*time.Second),
	)
	if err != nil {
		t.Fatalf("Reduce: %v", err)
	}
	if result.ExitCode != 0 {
		t.Errorf("expected exit code 0, got %d", result.ExitCode)
	}
	if result.Matrix == nil {
		t.Error("expected parsed matrix")
	}
	if result.Rows == 0 || result.Cols == 0 {
		t.Errorf("expected non-zero dimensions, got %dx%d", result.Rows, result.Cols)
	}
	t.Logf("LLL result: %dx%d, runtime=%v", result.Rows, result.Cols, result.Runtime)
}

func TestReduceTimeout(t *testing.T) {
	if _, err := exec.LookPath("fplll"); err != nil {
		t.Skip("fplll not installed")
	}

	client, err := NewDefault()
	if err != nil {
		t.Fatalf("NewDefault: %v", err)
	}

	m := Matrix{
		{big.NewInt(1), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0)},
		{big.NewInt(0), big.NewInt(1), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0)},
		{big.NewInt(0), big.NewInt(0), big.NewInt(1), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0)},
		{big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(1), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0)},
		{big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(1), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0)},
		{big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(1), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0)},
		{big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(1), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0)},
		{big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(1), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0)},
		{big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(1), big.NewInt(0), big.NewInt(0), big.NewInt(0)},
		{big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(1), big.NewInt(0), big.NewInt(0)},
		{big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(1), big.NewInt(0)},
		{big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(1)},
	}

	_, err = client.Reduce(context.Background(), m,
		WithAlgorithm(AlgLLL),
		WithTimeout(1*time.Nanosecond),
	)
	if err == nil {
		t.Error("expected timeout error")
	}
	if !errors.Is(err, ErrTimeout) {
		t.Errorf("expected ErrTimeout, got: %v", err)
	}
}
```

- [ ] **Step 5: Update smoke CLI and examples**

Update `cmd/gofplll-smoke/main.go`:

```go
package main

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/maqsatto/gofplll"
)

func main() {
	client, err := gofplll.NewDefault()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("binary: fplll (via $PATH)")

	m := gofplll.Matrix{
		{big.NewInt(10), big.NewInt(11)},
		{big.NewInt(11), big.NewInt(12)},
	}

	result, err := client.Reduce(context.Background(), m,
		gofplll.WithAlgorithm(gofplll.AlgLLL),
		gofplll.WithTimeout(5*time.Second),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reduce error: %v\n", err)
		if result != nil {
			fmt.Println("stdout:", result.Stdout)
			fmt.Println("stderr:", result.Stderr)
		}
		os.Exit(1)
	}

	fmt.Println("raw stdout:", result.Stdout)
	fmt.Println("raw stderr:", result.Stderr)
	fmt.Println("parsed matrix:")
	for _, row := range result.Matrix {
		for j, val := range row {
			if j > 0 {
				fmt.Print(" ")
			}
			fmt.Print(val)
		}
		fmt.Println()
	}
	fmt.Printf("runtime: %v\n", result.Runtime)
	fmt.Printf("exit code: %d\n", result.ExitCode)
}
```

Update `examples/lll_basic/main.go`:

```go
package main

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/maqsatto/gofplll"
)

func main() {
	client, err := gofplll.NewDefault()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	m := gofplll.Matrix{
		{big.NewInt(1), big.NewInt(1), big.NewInt(1)},
		{big.NewInt(1), big.NewInt(2), big.NewInt(3)},
		{big.NewInt(1), big.NewInt(3), big.NewInt(5)},
	}

	result, err := client.Reduce(context.Background(), m,
		gofplll.WithAlgorithm(gofplll.AlgLLL),
		gofplll.WithDelta("0.99"),
		gofplll.WithTimeout(10*time.Second),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Reduced matrix:")
	for _, row := range result.Matrix {
		for j, val := range row {
			if j > 0 {
				fmt.Print(" ")
			}
			fmt.Print(val)
		}
		fmt.Println()
	}
	fmt.Printf("Runtime: %v\n", result.Runtime)
}
```

Update `examples/bkz_basic/main.go`:

```go
package main

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/maqsatto/gofplll"
)

func main() {
	client, err := gofplll.NewDefault()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	m := gofplll.Matrix{
		{big.NewInt(1), big.NewInt(1), big.NewInt(1)},
		{big.NewInt(1), big.NewInt(2), big.NewInt(3)},
		{big.NewInt(1), big.NewInt(3), big.NewInt(5)},
	}

	result, err := client.Reduce(context.Background(), m,
		gofplll.WithAlgorithm(gofplll.AlgBKZ),
		gofplll.WithBlockSize(2),
		gofplll.WithTimeout(30*time.Second),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("BKZ reduced matrix:")
	for _, row := range result.Matrix {
		for j, val := range row {
			if j > 0 {
				fmt.Print(" ")
			}
			fmt.Print(val)
		}
		fmt.Println()
	}
	fmt.Printf("Runtime: %v\n", result.Runtime)
}
```

- [ ] **Step 6: Run full validation**

```bash
go fmt ./...
go vet ./...
go test ./...
```

Expected: all pass, integration tests skip if fplll not installed

- [ ] **Step 7: Commit**

```bash
git add -A
git commit -m "refactor: complete clean architecture with Reducer interface and functional options"
```

---

### Task 6: Update documentation and Makefile

**Files:**
- Modify: `Makefile`
- Modify: `README.md`
- Modify: `docs/API.md`
- Modify: `docs/ERROR_MODEL.md`
- Modify: `docs/FPLLL_FORMAT.md`
- Modify: `docs/LINUX_ONLY.md`

- [ ] **Step 1: Update Makefile**

```makefile
.PHONY: test test-integration smoke fmt vet

fmt:
	go fmt ./...

vet:
	go vet ./...

test:
	go test ./...

test-integration:
	go test ./... -run Integration

smoke:
	go run ./cmd/gofplll-smoke
```

- [ ] **Step 2: Update README.md**

Write the new README with functional options examples, cleaner API, and all required sections.

- [ ] **Step 3: Update docs/API.md**

Document the new Reducer interface, functional options, Client, and internal packages.

- [ ] **Step 4: Update docs/ERROR_MODEL.md**

Document the new error wrapping behavior.

- [ ] **Step 5: Commit**

```bash
git add -A
git commit -m "docs: update documentation for clean architecture refactor"
```

---

### Task 7: Final validation

- [ ] **Step 1: Run all acceptance criteria**

```bash
go fmt ./...
go vet ./...
go test ./...
go run ./cmd/gofplll-smoke
```

Expected: all pass

- [ ] **Step 2: Verify no leftover files**

Check that old files (`codec.go`, `validate.go`, `fplll_linux.go`, `fplll_unsupported.go`, old test files) are removed.

- [ ] **Step 3: Verify file tree matches spec**

```
gofplll/
├─ go.mod
├─ gofplll.go
├─ matrix.go
├─ algorithm.go
├─ errors.go
├─ options.go
├─ result.go
├─ gofplll_test.go
├─ internal/
│  ├─ encoding/
│  │  ├─ codec.go
│  │  └─ codec_test.go
│  ├─ validation/
│  │  ├─ validate.go
│  │  └─ validate_test.go
│  └─ runner/
│     ├─ runner.go
│     ├─ linux.go
│     └─ unsupported.go
├─ cmd/gofplll-smoke/main.go
├─ examples/
├─ testdata/
├─ docs/
├─ Makefile
├─ README.md
└─ LICENSE
```
