// Package calculator provides arithmetic operations for the calculator
// application. It can evaluate full infix expressions (e.g. "2 + 3 * (4 - 1)")
// as well as run individual two-operand operations. It is intentionally free
// of any HTTP or transport concerns so it can be unit-tested and reused
// independently.
package calculator

import "errors"

// ErrDivideByZero is returned when a division (or modulo) by zero is attempted.
var ErrDivideByZero = errors.New("division by zero")

// ErrUnknownOperation is returned when an unsupported operation is requested.
var ErrUnknownOperation = errors.New("unknown operation")

// Add returns a + b.
func Add(a, b float64) float64 { return a + b }

// Subtract returns a - b.
func Subtract(a, b float64) float64 { return a - b }

// Multiply returns a * b.
func Multiply(a, b float64) float64 { return a * b }

// Divide returns a / b, or ErrDivideByZero if b is zero.
func Divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, ErrDivideByZero
	}
	return a / b, nil
}

// Compute dispatches to the appropriate operation based on op.
// Supported ops: "add", "subtract", "multiply", "divide".
func Compute(op string, a, b float64) (float64, error) {
	switch op {
	case "add":
		return Add(a, b), nil
	case "subtract":
		return Subtract(a, b), nil
	case "multiply":
		return Multiply(a, b), nil
	case "divide":
		return Divide(a, b)
	default:
		return 0, ErrUnknownOperation
	}
}

// Evaluate parses and evaluates a full arithmetic expression and returns the
// result. It supports +, -, *, /, % (modulo), ^ (exponent), unary minus/plus,
// parentheses, and decimal numbers, with standard operator precedence.
func Evaluate(expr string) (float64, error) {
	expr, err := normalize(expr)
	if err != nil {
		return 0, err
	}
	tokens, err := tokenize(expr)
	if err != nil {
		return 0, err
	}
	p := &parser{tokens: tokens}
	result, err := p.parseExpression()
	if err != nil {
		return 0, err
	}
	if !p.atEnd() {
		return 0, errors.New("unexpected token: " + p.peek().value)
	}
	return result, nil
}
