# gofplll Clean Architecture Refactor

## Goal

Refactor gofplll from a flat single-package structure into a clean, idiomatic Go library with internal packages, functional options, a Reducer interface, and a simplified public API surface.

## Architecture

```
gofplll/
├─ go.mod
├─ gofplll.go            # Reducer interface, Client, New, NewDefault
├─ matrix.go             # Matrix type + helpers
├─ algorithm.go          # Algorithm, FloatType, OutputFormat (with String/Parse)
├─ errors.go             # Sentinel errors
├─ options.go            # ClientOption + ReduceOption functional option types
├─ result.go             # ReduceResult
├─ internal/
│  ├─ encoding/
│  │  ├─ codec.go        # Encode/Decode implementations
│  │  └─ codec_test.go
│  ├─ validation/
│  │  ├─ validate.go     # ValidateMatrix/ValidateOptions implementations
│  │  └─ validate_test.go
│  └─ runner/
│     ├─ runner.go       # Runner interface + ArgsBuilder
│     ├─ linux.go        # Linux subprocess (exec.CommandContext)
│     └─ unsupported.go  # Non-Linux stub
├─ cmd/gofplll-smoke/
├─ examples/
├─ testdata/
└─ docs/
```

## Public API

### Core Types

```go
type Matrix [][]*big.Int

type Algorithm string   // LLL, BKZ, HKZ, SVP, CVP
type FloatType string   // MPFR, Double, LongDouble, DD, QD
type OutputFormat string // Basis, Status, SVP, CVP
```

### Reducer Interface

```go
type Reducer interface {
    Reduce(ctx context.Context, matrix Matrix, opts ...ReduceOption) (*ReduceResult, error)
}
```

### Client

```go
type Client struct { binaryPath string }  // private field

func New(path string, opts ...ClientOption) (*Client, error)
func NewDefault(opts ...ClientOption) (*Client, error)
func (c *Client) Reduce(ctx context.Context, matrix Matrix, opts ...ReduceOption) (*ReduceResult, error)
```

### Functional Options

```go
// Client construction
type ClientOption interface { applyClient(*clientConfig) }
func WithBinaryPath(path string) ClientOption
func WithWorkDir(dir string) ClientOption

// Reduce call
type ReduceOption interface { applyReduce(*reduceConfig) }
func WithAlgorithm(a Algorithm) ReduceOption
func WithDelta(d string) ReduceOption
func WithEta(e string) ReduceOption
func WithBlockSize(n int) ReduceOption
func WithMaxLoops(n int) ReduceOption
func WithMaxTime(d time.Duration) ReduceOption
func WithAutoAbort() ReduceOption
func WithNoLLL() ReduceOption
func WithFloatType(f FloatType) ReduceOption
func WithPrecision(p int) ReduceOption
func WithOutputFormat(f OutputFormat) ReduceOption
func WithVerbose() ReduceOption
func WithTimeout(d time.Duration) ReduceOption
```

### Result

```go
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

## Internal Packages

### internal/encoding

- `Encode(m Matrix) ([]byte, error)` — serialize to fplll bracket format
- `Decode(s string) (Matrix, error)` — parse fplll bracket format

### internal/validation

- `ValidateMatrix(m Matrix) error` — non-empty, rectangular, no nils
- `ValidateOptions(matrix Matrix, cfg *reduceConfig) error` — algorithm, BKZ params, delta/eta, precision

### internal/runner

```go
type Runner interface {
    Run(ctx context.Context, binaryPath string, args []string, stdin []byte) (stdout, stderr []byte, exitCode int, err error)
}
```

- `LinuxRunner` — exec.CommandContext implementation
- `UnsupportedRunner` — returns ErrUnsupportedOS

## Key Design Decisions

1. **Functional options over config structs** — cleaner API, no zero-value confusion, extensible without breaking changes
2. **Reducer interface** — enables mocking for tests, swappable backends in future
3. **Private binaryPath** — prevents mutation after construction
4. **Internal packages** — encoding, validation, runner hidden from importers
5. **WorkDir on ClientOption** — was dead field on ReduceOptions, now properly on client construction
6. **Validation in internal/validation** — single validation path, no redundant calls

## Migration Notes

- `ReduceOptions{Algorithm: "lll", BlockSize: 10}` → `WithAlgorithm("lll"), WithBlockSize(10)`
- `Client.BinaryPath` → `client.BinaryPath()` (getter method) or just internal
- `EncodeMatrix`/`DecodeMatrix` → removed from public API (internal only)
- `ValidateMatrix`/`ValidateOptions` → removed from public API (internal only)
