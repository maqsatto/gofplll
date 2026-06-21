# gofplll

A Linux-only Go library that wraps the `fplll` command-line binary through a safe subprocess API for integer lattice reduction.

## Scope Boundary

This project is **NOT** a cryptographic attack tool. It provides a generic integer lattice reduction backend API for Go applications.

**What this library does:** `integer matrix -> fplll subprocess -> parsed reduced matrix / result / stdout / stderr / runtime`

**What this library does NOT do:** HNP, ECDSA, Bitcoin, nonce recovery, private key recovery, blockchain logic, or any attack pipeline.

## Requirements

- Linux (WSL2 supported)
- Go 1.22+
- `fplll` binary installed

### Install fplll

**Arch/CachyOS:**

```bash
sudo pacman -S fplll
```

**Debian/Ubuntu:**

```bash
sudo apt install fplll-tools
```

## Install

```bash
go get github.com/maqsatto/gofplll
```

## Quick Start

```go
package main

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/maqsatto/gofplll"
)

func main() {
	client, _ := gofplll.NewDefault()

	m := gofplll.Matrix{
		{big.NewInt(10), big.NewInt(11)},
		{big.NewInt(11), big.NewInt(12)},
	}

	result, _ := client.Reduce(context.Background(), m, gofplll.ReduceOptions{
		Algorithm: gofplll.AlgLLL,
		Timeout:   5 * time.Second,
	})

	fmt.Println(result.Matrix)
}
```

## Basic LLL Example

```go
client, err := gofplll.NewDefault()
if err != nil {
    log.Fatal(err)
}

m := gofplll.Matrix{
    {big.NewInt(1), big.NewInt(1), big.NewInt(1)},
    {big.NewInt(1), big.NewInt(2), big.NewInt(3)},
    {big.NewInt(1), big.NewInt(3), big.NewInt(5)},
}

result, err := client.Reduce(context.Background(), m, gofplll.ReduceOptions{
    Algorithm: gofplll.AlgLLL,
    Delta:     "0.99",
})
```

## Basic BKZ Example

```go
result, err := client.Reduce(context.Background(), m, gofplll.ReduceOptions{
    Algorithm: gofplll.AlgBKZ,
    BlockSize: 2,
    MaxLoops:  10,
})
```

## Error Handling

```go
result, err := client.Reduce(ctx, m, opts)
if err != nil {
    if errors.Is(err, gofplll.ErrBinaryNotFound) {
        fmt.Println("fplll is not installed")
    } else if errors.Is(err, gofplll.ErrInvalidMatrix) {
        fmt.Println("bad matrix input")
    } else if errors.Is(err, gofplll.ErrTimeout) {
        fmt.Println("timed out, partial result:", result.Stdout)
    }
}
```

## Timeout Example

```go
result, err := client.Reduce(ctx, m, gofplll.ReduceOptions{
    Algorithm: gofplll.AlgLLL,
    Timeout:   2 * time.Second,
})
if errors.Is(err, gofplll.ErrTimeout) {
    fmt.Println("partial stdout:", result.Stdout)
}
```

## License

MIT License. See [LICENSE](LICENSE).
