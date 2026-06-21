# Error Model

`gofplll` uses sentinel errors for all failure modes. Check errors with `errors.Is`.

## Binary Not Found

```go
ErrBinaryNotFound
```

Returned when `fplll` is not found in `$PATH`. Fix: install `fplll`.

## Invalid Matrix

```go
ErrInvalidMatrix
```

Returned when the input matrix is empty, has 0 columns, has rows of different lengths, or contains nil coefficients.

## Invalid Options

```go
ErrInvalidOptions
```

Returned for unknown algorithm, BKZ block size < 2 or > row count, unparseable delta/eta, unknown float type, or negative precision.

## Subprocess Failure

```go
ErrSubprocessFailed
```

Returned when `exec.CommandContext` fails for a reason other than timeout.

## Timeout

```go
ErrTimeout
```

Returned when the subprocess exceeds the timeout. The `ReduceResult` still contains partial stdout/stderr captured before the process was killed.

```go
result, err := client.Reduce(ctx, m,
    gofplll.WithAlgorithm(gofplll.AlgLLL),
    gofplll.WithTimeout(2*time.Second),
)
if errors.Is(err, gofplll.ErrTimeout) {
    fmt.Println("partial stdout:", result.Stdout)
}
```

## Parser Failure

```go
ErrParserFailed
```

Returned when the fplll output cannot be parsed as a matrix.
