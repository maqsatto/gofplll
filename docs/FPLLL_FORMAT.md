# fplll Matrix Format

## Input Format

fplll reads matrices from stdin in bracket notation:

```
[[10 11]
 [11 12]]
```

Each row is enclosed in `[...]`. Rows are separated by newlines inside the outer `[...]`.

### Rules

- Outer brackets: `[[ ... ]]`
- Elements separated by spaces
- Rows separated by newlines
- All elements must be integers
- No commas

## Encoding (Go -> fplll)

`EncodeMatrix` serializes a `Matrix` into the format above:

```go
m := gofplll.Matrix{
    {big.NewInt(1), big.NewInt(2)},
    {big.NewInt(3), big.NewInt(4)},
}
data, _ := gofplll.EncodeMatrix(m)
// data == "[[1 2]\n [3 4]]"
```

## Decoding (fplll -> Go)

`DecodeMatrix` parses fplll output back into a `Matrix`:

```go
matrix, err := gofplll.DecodeMatrix("[[1 2]\n [3 4]]")
```

It handles:

- Positive integers: `10`
- Negative integers: `-5`
- Large integers: `12345678901234567890`
- Extra whitespace
- Multi-line output

## fplll Output

fplll prints the reduced basis in the same bracket format by default. With `-of s` it prints a single vector (the shortest). With `-of t` it prints status information.
