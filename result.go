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
