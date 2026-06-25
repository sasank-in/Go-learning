package calculator

import (
	"errors"
	"math"
	"testing"
)

func TestEvaluate(t *testing.T) {
	tests := []struct {
		name string
		expr string
		want float64
	}{
		{"simple add", "2 + 3", 5},
		{"precedence", "2 + 3 * 4", 14},
		{"parentheses", "2 + 3 * (4 - 1)", 11},
		{"nested parens", "((1 + 2) * (3 + 4))", 21},
		{"decimals", "0.1 + 0.2", 0.3},
		{"unary minus", "-5 + 3", -2},
		{"unary in parens", "3 * (-2)", -6},
		{"division", "10 / 4", 2.5},
		{"modulo", "10 % 3", 1},
		{"power", "2 ^ 8", 256},
		{"power right assoc", "2 ^ 3 ^ 2", 512},
		{"power and mult", "2 * 3 ^ 2", 18},
		{"whitespace heavy", "  7   -   2  ", 5},
		{"chained", "1 + 2 + 3 + 4 + 5", 15},
		{"double negative", "--5", 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Evaluate(tt.expr)
			if err != nil {
				t.Fatalf("Evaluate(%q) unexpected error: %v", tt.expr, err)
			}
			if math.Abs(got-tt.want) > 1e-9 {
				t.Errorf("Evaluate(%q) = %v, want %v", tt.expr, got, tt.want)
			}
		})
	}
}

func TestEvaluateErrors(t *testing.T) {
	tests := []struct {
		name    string
		expr    string
		wantErr error // nil means "any error"
	}{
		{"divide by zero", "1 / 0", ErrDivideByZero},
		{"modulo by zero", "5 % 0", ErrDivideByZero},
		{"empty", "", nil},
		{"only operator", "+", nil},
		{"unbalanced paren", "(1 + 2", nil},
		{"trailing operator", "3 +", nil},
		{"bad character", "2 $ 3", nil},
		{"double dot", "1.2.3", nil},
		{"dangling number", "2 3", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Evaluate(tt.expr)
			if err == nil {
				t.Fatalf("Evaluate(%q) expected error, got nil", tt.expr)
			}
			if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Errorf("Evaluate(%q) error = %v, want %v", tt.expr, err, tt.wantErr)
			}
		})
	}
}
