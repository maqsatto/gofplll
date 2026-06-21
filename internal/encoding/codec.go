package encoding

import (
	"fmt"
	"math/big"
	"strings"
)

// Matrix is a 2D slice of arbitrary-precision integers.
type Matrix = [][]*big.Int

func Encode(m Matrix) ([]byte, error) {
	if len(m) == 0 {
		return nil, fmt.Errorf("encoding: empty matrix")
	}
	cols := len(m[0])
	if cols == 0 {
		return nil, fmt.Errorf("encoding: empty matrix")
	}
	for i, row := range m {
		if len(row) != cols {
			return nil, fmt.Errorf("encoding: row %d has %d columns, expected %d", i, len(row), cols)
		}
		for j, val := range row {
			if val == nil {
				return nil, fmt.Errorf("encoding: nil at [%d][%d]", i, j)
			}
		}
	}

	var b strings.Builder
	b.WriteString("[")
	for i, row := range m {
		if i > 0 {
			b.WriteString("\n ")
		}
		b.WriteString("[")
		for j, val := range row {
			if j > 0 {
				b.WriteString(" ")
			}
			b.WriteString(val.String())
		}
		b.WriteString("]")
	}
	b.WriteString("]")
	return []byte(b.String()), nil
}

func Decode(out string) (Matrix, error) {
	out = strings.TrimSpace(out)
	if len(out) == 0 {
		return nil, fmt.Errorf("encoding: empty output")
	}

	outerOpen := strings.Index(out, "[[")
	if outerOpen == -1 {
		return nil, fmt.Errorf("encoding: missing outer brackets")
	}
	start := out[outerOpen:]

	var rows Matrix
	depth := 0
	rowStart := -1

	for i := 0; i < len(start); i++ {
		switch start[i] {
		case '[':
			depth++
			if depth == 2 {
				rowStart = i + 1
			}
		case ']':
			depth--
			if depth == 1 && rowStart >= 0 {
				rowStr := strings.TrimSpace(start[rowStart:i])
				if rowStr == "" {
					rowStart = -1
					continue
				}
				fields := strings.Fields(rowStr)
				var row []*big.Int
				for _, f := range fields {
					val := new(big.Int)
					if _, ok := val.SetString(f, 10); !ok {
						return nil, fmt.Errorf("encoding: cannot parse %q as integer", f)
					}
					row = append(row, val)
				}
				rows = append(rows, row)
				rowStart = -1
			}
		}
	}

	if len(rows) == 0 {
		return nil, fmt.Errorf("encoding: no rows found")
	}
	return rows, nil
}
