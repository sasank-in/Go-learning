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
		{"sqrt", "sqrt(16)", 4},
		{"nested func", "sqrt(sqrt(16))", 2},
		{"func in expr", "2 + sqrt(9) * 2", 8},
		{"pi constant", "pi", math.Pi},
		{"const in expr", "2 * pi", 2 * math.Pi},
		{"e constant", "e", math.E},
		{"abs negative", "abs(-7)", 7},
		{"floor", "floor(3.9)", 3},
		{"ceil", "ceil(3.1)", 4},
		{"log10", "log(1000)", 3},
		{"ln of e", "ln(e)", 1},
		{"func with expr arg", "sqrt(2 ^ 2 + 3 ^ 2 - 4)", 3},
		{"scientific notation", "1.5e3", 1500},
		{"sci negative exp", "2e-2", 0.02},
		{"case insensitive", "SQRT(4) + PI - pi", 2},
		{"pow", "pow(2, 10)", 1024},
		{"hypot", "hypot(3, 4)", 5},
		{"max two", "max(3, 7)", 7},
		{"max many", "max(1, 9, 4, 2)", 9},
		{"min many", "min(5, 2, 8, 1)", 1},
		{"sum", "sum(1, 2, 3, 4)", 10},
		{"avg", "avg(2, 4, 6)", 4},
		{"log base 2", "log(8, 2)", 3},
		{"log default base 10", "log(1000)", 3},
		{"log10 explicit", "log10(100)", 2},
		{"nested multiarg", "max(pow(2, 3), 7)", 8},
		{"sign positive", "sign(42)", 1},
		{"sign negative", "sign(-3)", -1},
		{"trunc", "trunc(3.9)", 3},
		{"tanh zero", "tanh(0)", 0},
		{"deg of pi", "deg(pi)", 180},
		{"rad of 180", "rad(180)", math.Pi},
		{"expr args", "pow(1 + 1, 2 + 1)", 8},
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
		{"unknown function", "foo(2)", nil},
		{"unknown identifier", "x + 1", nil},
		{"function missing parens", "sqrt 4", nil},
		{"function missing arg", "sqrt()", nil},
		{"function unclosed", "sqrt(4", nil},
		{"too many args", "sqrt(4, 9)", nil},
		{"pow too few args", "pow(2)", nil},
		{"empty args", "max()", nil},
		{"trailing comma", "max(1, 2,)", nil},
		{"bad log base", "log(8, 1)", nil},
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
