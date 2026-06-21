# Error Model

`gofplll` uses sentinel errors for all failure modes. You can check errors with `errors.Is`.

## Binary Not Found

```go
ErrBinaryNotFound
```

Returned when `fplll` is not found in `$PATH`. Fix: install `fplll`.

## Invalid Matrix

```go
ErrInvalidMatrix
```

Returned when the input matrix:

- Is empty (0 rows)
- Has 0 columns
- Has rows of different lengths (ragged)
- Contains nil coefficients

## Invalid Options

```go
ErrInvalidOptions
```

Returned when:

- An unknown algorithm is specified
- BKZ block size is less than 2
- BKZ block size exceeds the number of matrix rows
- Delta or Eta values cannot be parsed
- Precision is negative

## Subprocess Failure

```go
ErrSubprocessFailed
```

Returned when `exec.CommandContext` fails for a reason other than timeout (e.g., binary permissions, OS-level errors).

## Timeout

```go
ErrTimeout
```

Returned when the subprocess exceeds the timeout. The `ReduceResult` still contains partial stdout/stderr captured before the process was killed.

```go
result, err := client.Reduce(ctx, m, opts)
if errors.Is(err, gofplll.ErrTimeout) {
    fmt.Println("partial stdout:", result.Stdout)
}
```

## Parser Failure

```go
ErrParserFailed
```

Returned when the fplll output cannot be parsed as a matrix. This can happen if:

- fplll output format changed
- Output contains non-numeric values
- Output brackets are malformed

## Error Hierarchy

All errors are package-level `var` values. Wrap them with `fmt.Errorf` when needed:

```go
if errors.Is(err, gofplll.ErrInvalidMatrix) {
    // handle invalid matrix
}
```
