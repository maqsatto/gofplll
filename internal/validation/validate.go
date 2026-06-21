package validation

import (
	"fmt"
	"math/big"
)

// Matrix is a 2D slice of arbitrary-precision integers.
type Matrix = [][]*big.Int

var (
	errInvalidMatrix  = fmt.Errorf("validation: invalid matrix")
	errInvalidOptions = fmt.Errorf("validation: invalid options")
)

var validAlgorithms = map[string]bool{
	"lll": true, "bkz": true, "hkz": true, "svp": true, "cvp": true,
}

var validFloatTypes = map[string]bool{
	"mpfr": true, "double": true, "longdouble": true, "dd": true, "qd": true,
}

func ValidateMatrix(m Matrix) error {
	if len(m) == 0 {
		return fmt.Errorf("%w: matrix must not be empty", errInvalidMatrix)
	}
	cols := len(m[0])
	if cols == 0 {
		return fmt.Errorf("%w: matrix must have at least one column", errInvalidMatrix)
	}
	for i, row := range m {
		if len(row) != cols {
			return fmt.Errorf("%w: row %d has %d columns, expected %d", errInvalidMatrix, i, len(row), cols)
		}
		for j, val := range row {
			if val == nil {
				return fmt.Errorf("%w: nil coefficient at [%d][%d]", errInvalidMatrix, i, j)
			}
		}
	}
	return nil
}

func ValidateOptions(m Matrix, algo string, blockSize int, delta, eta, floatType string, precision int) error {
	if err := ValidateMatrix(m); err != nil {
		return err
	}

	if algo != "" && !validAlgorithms[algo] {
		return fmt.Errorf("%w: unknown algorithm %q", errInvalidOptions, algo)
	}

	if algo == "bkz" {
		if blockSize < 2 {
			return fmt.Errorf("%w: BKZ block size must be >= 2, got %d", errInvalidOptions, blockSize)
		}
		if blockSize > len(m) {
			return fmt.Errorf("%w: BKZ block size %d exceeds number of rows %d", errInvalidOptions, blockSize, len(m))
		}
	}

	if delta != "" {
		if _, ok := new(big.Int).SetString(delta, 10); !ok {
			if _, ok := new(big.Float).SetString(delta); !ok {
				return fmt.Errorf("%w: invalid delta %q", errInvalidOptions, delta)
			}
		}
	}

	if eta != "" {
		if _, ok := new(big.Int).SetString(eta, 10); !ok {
			if _, ok := new(big.Float).SetString(eta); !ok {
				return fmt.Errorf("%w: invalid eta %q", errInvalidOptions, eta)
			}
		}
	}

	if floatType != "" && !validFloatTypes[floatType] {
		return fmt.Errorf("%w: unknown float type %q", errInvalidOptions, floatType)
	}

	if precision < 0 {
		return fmt.Errorf("%w: precision must be >= 0, got %d", errInvalidOptions, precision)
	}

	return nil
}
