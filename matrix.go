package gofplll

import "math/big"

// Matrix is a 2D slice of arbitrary-precision integers.
type Matrix [][]*big.Int

// NewMatrix creates a Matrix with the given dimensions, filled with zeros.
func NewMatrix(rows, cols int) Matrix {
	m := make(Matrix, rows)
	for i := range m {
		m[i] = make([]*big.Int, cols)
		for j := range m[i] {
			m[i][j] = new(big.Int)
		}
	}
	return m
}

// Rows returns the number of rows.
func (m Matrix) Rows() int { return len(m) }

// Cols returns the number of columns (0 if empty).
func (m Matrix) Cols() int {
	if len(m) == 0 {
		return 0
	}
	return len(m[0])
}
