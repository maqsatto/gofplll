package encoding

import (
	"math/big"
	"testing"
)

func TestEncode2x2(t *testing.T) {
	m := Matrix{
		{big.NewInt(10), big.NewInt(11)},
		{big.NewInt(11), big.NewInt(12)},
	}
	got, err := Encode(m)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "[[10 11]\n [11 12]]"
	if string(got) != want {
		t.Errorf("Encode() = %q, want %q", string(got), want)
	}
}

func TestEncodeNegative(t *testing.T) {
	m := Matrix{
		{big.NewInt(-1), big.NewInt(2)},
		{big.NewInt(3), big.NewInt(-4)},
	}
	got, err := Encode(m)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "[[-1 2]\n [3 -4]]"
	if string(got) != want {
		t.Errorf("Encode() = %q, want %q", string(got), want)
	}
}

func TestEncodeBigInt(t *testing.T) {
	bigVal := new(big.Int)
	bigVal.SetString("123456789012345678901234567890", 10)
	m := Matrix{
		{bigVal, big.NewInt(1)},
		{big.NewInt(1), bigVal},
	}
	got, err := Encode(m)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "[[123456789012345678901234567890 1]\n [1 123456789012345678901234567890]]"
	if string(got) != want {
		t.Errorf("Encode() = %q, want %q", string(got), want)
	}
}

func TestDecode2x2(t *testing.T) {
	input := "[[10 11]\n [11 12]]"
	got, err := Decode(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 || len(got[0]) != 2 {
		t.Fatalf("got dimensions %dx%d, want 2x2", len(got), len(got[0]))
	}
	if got[0][0].Int64() != 10 || got[0][1].Int64() != 11 {
		t.Errorf("row 0 = [%s %s], want [10 11]", got[0][0], got[0][1])
	}
	if got[1][0].Int64() != 11 || got[1][1].Int64() != 12 {
		t.Errorf("row 1 = [%s %s], want [11 12]", got[1][0], got[1][1])
	}
}

func TestDecodeNegative(t *testing.T) {
	input := "[[-1 2]\n [3 -4]]"
	got, err := Decode(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got[0][0].Int64() != -1 || got[1][1].Int64() != -4 {
		t.Errorf("unexpected values: %v", got)
	}
}

func TestDecodeBigInt(t *testing.T) {
	input := "[[123456789012345678901234567890 1]\n [1 123456789012345678901234567890]]"
	got, err := Decode(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	bigVal := new(big.Int)
	bigVal.SetString("123456789012345678901234567890", 10)
	if got[0][0].Cmp(bigVal) != 0 {
		t.Errorf("got[0][0] = %s, want %s", got[0][0], bigVal)
	}
}

func TestDecodeExtraWhitespace(t *testing.T) {
	input := "  [[  10   11  ]\n   [  11   12  ]]  "
	got, err := Decode(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got[0][0].Int64() != 10 {
		t.Errorf("got[0][0] = %s, want 10", got[0][0])
	}
}

func TestDecodeEmpty(t *testing.T) {
	_, err := Decode("")
	if err == nil {
		t.Error("expected error for empty input")
	}
}

func TestEncodeEmptyMatrix(t *testing.T) {
	_, err := Encode(Matrix{})
	if err == nil {
		t.Error("expected error for empty matrix")
	}
}
