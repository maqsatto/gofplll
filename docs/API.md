# API Reference

## Core Types

### `Matrix`

```go
type Matrix [][]*big.Int
```

A 2D matrix of arbitrary-precision integers.

**Methods:**
- `NewMatrix(rows, cols int) Matrix` — create zero-filled matrix
- `Rows() int` — number of rows
- `Cols() int` — number of columns

### `Algorithm`

```go
type Algorithm string
```

| Constant   | Value  |
|------------|--------|
| `AlgLLL`   | `"lll"` |
| `AlgBKZ`   | `"bkz"` |
| `AlgHKZ`   | `"hkz"` |
| `AlgSVP`   | `"svp"` |
| `AlgCVP`   | `"cvp"` |

### `FloatType`

```go
type FloatType string
```

| Constant          | Value          |
|-------------------|----------------|
| `FloatMPFR`       | `"mpfr"`       |
| `FloatDouble`     | `"double"`     |
| `FloatLongDouble` | `"longdouble"` |
| `FloatDD`         | `"dd"`         |
| `FloatQD`         | `"qd"`         |

### `OutputFormat`

```go
type OutputFormat string
```

| Constant        | Value |
|-----------------|-------|
| `OutputBasis`   | `"b"` |
| `OutputStatus`  | `"t"` |
| `OutputSVP`     | `"s"` |
| `OutputCVP`     | `"c"` |

## Reducer Interface

```go
type Reducer interface {
    Reduce(ctx context.Context, matrix Matrix, opts ...ReduceOption) (*ReduceResult, error)
}
```

The primary interface for lattice reduction backends. `Client` implements this.

## Client

### `New`

```go
func New(path string, opts ...ClientOption) (*Client, error)
```

Creates a new Client. If `path` is empty, searches for `fplll` in `$PATH`.

### `NewDefault`

```go
func NewDefault(opts ...ClientOption) (*Client, error)
```

Equivalent to `New("")`.

### `Reduce`

```go
func (c *Client) Reduce(ctx context.Context, matrix Matrix, opts ...ReduceOption) (*ReduceResult, error)
```

Runs fplll on the given matrix. Returns the reduced matrix, raw stdout/stderr, exit code, and runtime.

## Client Options

| Option | Description |
|--------|-------------|
| `WithBinaryPath(path)` | Set fplll binary path |
| `WithWorkDir(dir)` | Set subprocess working directory |

## Reduce Options

| Option | Description |
|--------|-------------|
| `WithAlgorithm(a)` | Set reduction algorithm |
| `WithDelta(d)` | Set LLL delta parameter |
| `WithEta(e)` | Set LLL eta parameter |
| `WithBlockSize(n)` | Set BKZ block size |
| `WithMaxLoops(n)` | Set BKZ max loop count |
| `WithMaxTime(d)` | Set BKZ max time limit |
| `WithAutoAbort()` | Enable BKZ auto-abort |
| `WithNoLLL()` | Disable LLL pre-processing |
| `WithFloatType(f)` | Set floating-point backend |
| `WithPrecision(p)` | Set MPFR precision in bits |
| `WithOutputFormat(f)` | Set output format |
| `WithVerbose()` | Enable verbose output |
| `WithTimeout(d)` | Set subprocess timeout |

## ReduceResult

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

## Errors

| Error | Description |
|-------|-------------|
| `ErrUnsupportedOS` | Non-Linux build |
| `ErrBinaryNotFound` | `fplll` not found in `$PATH` |
| `ErrInvalidMatrix` | Empty, ragged, or nil-coefficient matrix |
| `ErrInvalidOptions` | Invalid algorithm, block size, or delta |
| `ErrParserFailed` | Cannot parse fplll output |
| `ErrSubprocessFailed` | Subprocess execution failed |
| `ErrTimeout` | Subprocess exceeded timeout |
