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
