package calculator

import (
	"errors"
	"math"
	"testing"
)

func TestEvaluateWithVariables(t *testing.T) {
	// Assignment returns the value and stores it.
	val, env, err := EvaluateWith("x = 5", nil)
	if err != nil {
		t.Fatalf("assignment error: %v", err)
	}
	if val != 5 {
		t.Errorf("x = 5 returned %v, want 5", val)
	}
	if env["x"] != 5 {
		t.Errorf("env[x] = %v, want 5", env["x"])
	}
	if env["ans"] != 5 {
		t.Errorf("env[ans] = %v, want 5 after assignment", env["ans"])
	}

	// Reference the stored variable in a later expression.
	val, env, err = EvaluateWith("x * 2 + 1", env)
	if err != nil {
		t.Fatalf("reference error: %v", err)
	}
	if val != 11 {
		t.Errorf("x * 2 + 1 = %v, want 11", val)
	}
	if env["ans"] != 11 {
		t.Errorf("env[ans] = %v, want 11", env["ans"])
	}

	// "ans" carries the previous result.
	val, _, err = EvaluateWith("ans + 9", env)
	if err != nil {
		t.Fatalf("ans error: %v", err)
	}
	if val != 20 {
		t.Errorf("ans + 9 = %v, want 20", val)
	}
}

func TestEvaluateWithReassignmentAndExprRHS(t *testing.T) {
	env := map[string]float64{}
	if _, env2, err := EvaluateWith("r = 2", env); err != nil {
		t.Fatal(err)
	} else {
		env = env2
	}
	// RHS can be a full expression that references existing variables and funcs.
	val, env, err := EvaluateWith("area = pi * r ^ 2", env)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if math.Abs(val-math.Pi*4) > 1e-9 {
		t.Errorf("area = %v, want %v", val, math.Pi*4)
	}
	if math.Abs(env["area"]-math.Pi*4) > 1e-9 {
		t.Errorf("env[area] = %v, want %v", env["area"], math.Pi*4)
	}
}

func TestEvaluateWithErrors(t *testing.T) {
	tests := []struct {
		name string
		expr string
		env  map[string]float64
	}{
		{"undefined variable", "y + 1", nil},
		{"assign to constant", "pi = 3", nil},
		{"assign to function", "sqrt = 3", nil},
		{"trailing tokens after assignment", "x = 5 6", nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, _, err := EvaluateWith(tt.expr, tt.env); err == nil {
				t.Errorf("EvaluateWith(%q) expected error, got nil", tt.expr)
			}
		})
	}
}

func TestEvaluateWithDoesNotMutateInput(t *testing.T) {
	original := map[string]float64{"a": 1}
	if _, _, err := EvaluateWith("b = 2", original); err != nil {
		t.Fatal(err)
	}
	if _, exists := original["b"]; exists {
		t.Error("EvaluateWith mutated the caller's map")
	}
}

func TestAnsDefaultsToZero(t *testing.T) {
	// "ans" is usable before any prior calculation, defaulting to 0.
	val, _, err := EvaluateWith("ans + 7", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != 7 {
		t.Errorf("ans + 7 = %v, want 7", val)
	}
}

func TestNonFiniteResultsAreRejected(t *testing.T) {
	tests := []struct {
		name string
		expr string
	}{
		{"sqrt of negative", "sqrt(-4)"},
		{"ln of zero", "ln(0)"},
		{"overflow", "2 ^ 100000"},
		{"assignment to non-finite", "z = sqrt(-1)"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := EvaluateWith(tt.expr, nil)
			if !errors.Is(err, ErrNotFinite) {
				t.Errorf("EvaluateWith(%q) error = %v, want ErrNotFinite", tt.expr, err)
			}
		})
	}
}
