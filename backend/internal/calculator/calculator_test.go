package calculator

import (
	"errors"
	"testing"
)

func TestCompute(t *testing.T) {
	tests := []struct {
		name    string
		op      string
		a, b    float64
		want    float64
		wantErr error
	}{
		{"add", "add", 2, 3, 5, nil},
		{"subtract", "subtract", 5, 3, 2, nil},
		{"multiply", "multiply", 4, 3, 12, nil},
		{"divide", "divide", 10, 2, 5, nil},
		{"divide by zero", "divide", 1, 0, 0, ErrDivideByZero},
		{"unknown op", "modulo", 1, 1, 0, ErrUnknownOperation},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Compute(tt.op, tt.a, tt.b)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Compute(%q) error = %v, want %v", tt.op, err, tt.wantErr)
			}
			if err == nil && got != tt.want {
				t.Errorf("Compute(%q, %v, %v) = %v, want %v", tt.op, tt.a, tt.b, got, tt.want)
			}
		})
	}
}
