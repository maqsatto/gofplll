package validation

import (
	"math/big"
	"testing"
)

func TestValidateMatrixEmpty(t *testing.T) {
	m := Matrix{}
	if err := ValidateMatrix(m); err == nil {
		t.Error("expected error for empty matrix")
	}
}

func TestValidateMatrixRagged(t *testing.T) {
	m := Matrix{
		{big.NewInt(1), big.NewInt(2)},
		{big.NewInt(3)},
	}
	if err := ValidateMatrix(m); err == nil {
		t.Error("expected error for ragged matrix")
	}
}

func TestValidateMatrixNilCell(t *testing.T) {
	m := Matrix{
		{big.NewInt(1), nil},
		{big.NewInt(3), big.NewInt(4)},
	}
	if err := ValidateMatrix(m); err == nil {
		t.Error("expected error for nil cell")
	}
}

func TestValidateBKZBlockSizeTooSmall(t *testing.T) {
	m := Matrix{{big.NewInt(1), big.NewInt(2)}, {big.NewInt(3), big.NewInt(4)}}
	if err := ValidateOptions(m, "bkz", 1, "", "", "", 0); err == nil {
		t.Error("expected error for block size < 2")
	}
}

func TestValidateBKZBlockSizeTooLarge(t *testing.T) {
	m := Matrix{{big.NewInt(1), big.NewInt(2)}, {big.NewInt(3), big.NewInt(4)}}
	if err := ValidateOptions(m, "bkz", 5, "", "", "", 0); err == nil {
		t.Error("expected error for block size > rows")
	}
}

func TestValidateBKZBlockSizeValid(t *testing.T) {
	m := Matrix{{big.NewInt(1), big.NewInt(2)}, {big.NewInt(3), big.NewInt(4)}}
	if err := ValidateOptions(m, "bkz", 2, "", "", "", 0); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateUnknownAlgorithm(t *testing.T) {
	m := Matrix{{big.NewInt(1), big.NewInt(2)}}
	if err := ValidateOptions(m, "xyz", 0, "", "", "", 0); err == nil {
		t.Error("expected error for unknown algorithm")
	}
}

func TestValidateNegativePrecision(t *testing.T) {
	m := Matrix{{big.NewInt(1), big.NewInt(2)}}
	if err := ValidateOptions(m, "lll", 0, "", "", "", -1); err == nil {
		t.Error("expected error for negative precision")
	}
}

func TestValidateValidLLL(t *testing.T) {
	m := Matrix{{big.NewInt(1), big.NewInt(2)}, {big.NewInt(3), big.NewInt(4)}}
	if err := ValidateOptions(m, "lll", 0, "0.99", "0.51", "mpfr", 128); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
