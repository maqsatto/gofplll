package gofplll

import (
	"context"
	"errors"
	"math/big"
	"os/exec"
	"testing"
	"time"
)

func TestNewDefaultNotFound(t *testing.T) {
	t.Setenv("PATH", "")
	_, err := NewDefault()
	if err == nil {
		t.Error("expected error when fplll not found")
	}
}

func TestReduceLLLWithFPLLL(t *testing.T) {
	if _, err := exec.LookPath("fplll"); err != nil {
		t.Skip("fplll not installed")
	}

	client, err := NewDefault()
	if err != nil {
		t.Fatalf("NewDefault: %v", err)
	}

	m := Matrix{
		{big.NewInt(10), big.NewInt(11)},
		{big.NewInt(11), big.NewInt(12)},
	}

	result, err := client.Reduce(context.Background(), m,
		WithAlgorithm(AlgLLL),
		WithTimeout(5*time.Second),
	)
	if err != nil {
		t.Fatalf("Reduce: %v", err)
	}
	if result.ExitCode != 0 {
		t.Errorf("expected exit code 0, got %d", result.ExitCode)
	}
	if result.Matrix == nil {
		t.Error("expected parsed matrix")
	}
	if result.Rows == 0 || result.Cols == 0 {
		t.Errorf("expected non-zero dimensions, got %dx%d", result.Rows, result.Cols)
	}
	t.Logf("LLL result: %dx%d, runtime=%v", result.Rows, result.Cols, result.Runtime)
}

func TestReduceTimeout(t *testing.T) {
	if _, err := exec.LookPath("fplll"); err != nil {
		t.Skip("fplll not installed")
	}

	client, err := NewDefault()
	if err != nil {
		t.Fatalf("NewDefault: %v", err)
	}

	m := Matrix{
		{big.NewInt(1), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0)},
		{big.NewInt(0), big.NewInt(1), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0)},
		{big.NewInt(0), big.NewInt(0), big.NewInt(1), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0)},
		{big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(1), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0)},
		{big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(1), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0)},
		{big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(1), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0)},
		{big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(1), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0)},
		{big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(1), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0)},
		{big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(1), big.NewInt(0), big.NewInt(0), big.NewInt(0)},
		{big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(1), big.NewInt(0), big.NewInt(0)},
		{big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(1), big.NewInt(0)},
		{big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(1)},
	}

	_, err = client.Reduce(context.Background(), m,
		WithAlgorithm(AlgLLL),
		WithTimeout(1*time.Nanosecond),
	)
	if err == nil {
		t.Error("expected timeout error")
	}
	if !errors.Is(err, ErrTimeout) {
		t.Errorf("expected ErrTimeout, got: %v", err)
	}
}
