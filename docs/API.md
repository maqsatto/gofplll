# API Reference

## Types

### `Matrix`

```go
type Matrix [][]*big.Int
```

A 2D matrix of arbitrary-precision integers.

### `Algorithm`

```go
type Algorithm string
```

Supported reduction algorithms:

| Constant   | Value  | Description          |
|------------|--------|----------------------|
| `AlgLLL`   | `"lll"` | LLL reduction       |
| `AlgBKZ`   | `"bkz"` | BKZ block reduction |
| `AlgHKZ`   | `"hkz"` | HKZ reduction       |
| `AlgSVP`   | `"svp"` | Shortest vector     |
| `AlgCVP`   | `"cvp"` | Closest vector      |

### `FloatType`

```go
type FloatType string
```

Floating-point backends:

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

| Constant        | Value | Description         |
|-----------------|-------|---------------------|
| `OutputBasis`   | `"b"` | Output basis        |
| `OutputStatus`  | `"t"` | Output status       |
| `OutputSVP`     | `"s"` | Output SVP vector   |
| `OutputCVP`     | `"c"` | Output CVP vector   |

### `ReduceOptions`

```go
type ReduceOptions struct {
    Algorithm    Algorithm
    Delta        string
    Eta          string
    BlockSize    int
    MaxLoops     int
    MaxTime      time.Duration
    AutoAbort    bool
    NoLLL        bool
    FloatType    FloatType
    Precision    int
    OutputFormat OutputFormat
    Verbose      bool
    Timeout      time.Duration
    WorkDir      string
}
```

### `ReduceResult`

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

## Client

### `New`

```go
func New(binaryPath string) (*Client, error)
```

Creates a new Client. If `binaryPath` is empty, searches for `fplll` in `$PATH`.

### `NewDefault`

```go
func NewDefault() (*Client, error)
```

Equivalent to `New("")`.

### `Reduce`

```go
func (c *Client) Reduce(ctx context.Context, matrix Matrix, opts ReduceOptions) (*ReduceResult, error)
```

Runs fplll on the given matrix. Returns the reduced matrix, raw stdout/stderr, exit code, and runtime. Supports timeout via `context.Context` or `opts.Timeout`.

## Codec

### `EncodeMatrix`

```go
func EncodeMatrix(m Matrix) ([]byte, error)
```

Serializes a Matrix into fplll bracket format.

### `DecodeMatrix`

```go
func DecodeMatrix(out string) (Matrix, error)
```

Parses fplll bracket format output into a Matrix.

## Validation

### `ValidateMatrix`

```go
func ValidateMatrix(m Matrix) error
```

Checks matrix is non-empty, rectangular, and has no nil cells.

### `ValidateOptions`

```go
func ValidateOptions(matrix Matrix, opts ReduceOptions) error
```

Validates options against the given matrix. Enforces BKZ block size constraints.

## Errors

| Error                  | Description                              |
|------------------------|------------------------------------------|
| `ErrUnsupportedOS`     | Non-Linux build                          |
| `ErrBinaryNotFound`    | `fplll` not found in `$PATH`            |
| `ErrInvalidMatrix`     | Empty, ragged, or nil-coefficient matrix |
| `ErrInvalidOptions`    | Invalid algorithm, block size, or delta  |
| `ErrParserFailed`      | Cannot parse fplll output                |
| `ErrSubprocessFailed`  | Subprocess execution failed              |
| `ErrTimeout`           | Subprocess exceeded timeout              |
