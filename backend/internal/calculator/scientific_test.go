package calculator

import (
	"math"
	"testing"
)

func TestFactorialAndCombinatorics(t *testing.T) {
	tests := []struct {
		expr string
		want float64
	}{
		{"5!", 120},
		{"0!", 1},
		{"fact(6)", 720},
		{"3!^2", 36},        // (3!)^2 = 36, postfix binds tighter than ^
		{"2 + 3!", 8},       // 2 + 6
		{"gcd(12, 18)", 6},
		{"gcd(48, 36, 60)", 12},
		{"lcm(4, 6)", 12},
		{"lcm(3, 4, 5)", 60},
		{"ncr(5, 2)", 10},
		{"npr(5, 2)", 20},
		{"ncr(10, 0)", 1},
	}
	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			got, err := Evaluate(tt.expr)
			if err != nil {
				t.Fatalf("Evaluate(%q) error: %v", tt.expr, err)
			}
			if math.Abs(got-tt.want) > 1e-9 {
				t.Errorf("Evaluate(%q) = %v, want %v", tt.expr, got, tt.want)
			}
		})
	}
}

func TestFactorialErrors(t *testing.T) {
	for _, expr := range []string{"(-1)!", "2.5!", "171!"} {
		if _, err := Evaluate(expr); err == nil {
			t.Errorf("Evaluate(%q) expected error, got nil", expr)
		}
	}
}

func TestAngleModes(t *testing.T) {
	deg := func(e string) float64 {
		v, _, err := EvaluateWithOptions(e, Options{Angle: Degrees})
		if err != nil {
			t.Fatalf("EvaluateWithOptions(%q) error: %v", e, err)
		}
		return v
	}

	if got := deg("sin(90)"); math.Abs(got-1) > 1e-9 {
		t.Errorf("sin(90 deg) = %v, want 1", got)
	}
	if got := deg("cos(180)"); math.Abs(got-(-1)) > 1e-9 {
		t.Errorf("cos(180 deg) = %v, want -1", got)
	}
	if got := deg("asin(1)"); math.Abs(got-90) > 1e-9 {
		t.Errorf("asin(1) deg = %v, want 90", got)
	}

	// Radians (default) unchanged.
	if got, _ := Evaluate("sin(0)"); got != 0 {
		t.Errorf("sin(0 rad) = %v, want 0", got)
	}
	if got, _, _ := EvaluateWithOptions("sin(pi/2)", Options{Angle: Radians}); math.Abs(got-1) > 1e-9 {
		t.Errorf("sin(pi/2 rad) = %v, want 1", got)
	}
}
