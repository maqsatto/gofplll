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

	ctx := context.Background()
	result, err := client.Reduce(ctx, m, gofplll.ReduceOptions{
		Algorithm: gofplll.AlgLLL,
		Delta:     "0.99",
		Timeout:   10 * time.Second,
	})
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
