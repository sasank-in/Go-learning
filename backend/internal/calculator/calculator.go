// Package calculator provides arithmetic operations for the calculator
// application. It can evaluate full infix expressions (e.g. "2 + 3 * (4 - 1)")
// as well as run individual two-operand operations. It is intentionally free
// of any HTTP or transport concerns so it can be unit-tested and reused
// independently.
package calculator

import (
	"errors"
	"math"
)

// ErrDivideByZero is returned when a division (or modulo) by zero is attempted.
var ErrDivideByZero = errors.New("division by zero")

// ErrUnknownOperation is returned when an unsupported operation is requested.
var ErrUnknownOperation = errors.New("unknown operation")

// ErrNotFinite is returned when a computation produces a result that is not a
// finite real number (NaN or ±Inf), e.g. sqrt(-1), ln(0), or an overflow.
var ErrNotFinite = errors.New("result is not a finite number")

// checkFinite returns ErrNotFinite if v is NaN or infinite; otherwise nil.
func checkFinite(v float64) error {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return ErrNotFinite
	}
	return nil
}

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
	result, _, err := EvaluateWith(expr, nil)
	return result, err
}

// AngleMode controls how trigonometric functions interpret their input
// (and how inverse trig functions produce output).
type AngleMode int

const (
	// Radians is the default angle mode.
	Radians AngleMode = iota
	// Degrees interprets/produces trig angles in degrees.
	Degrees
)

// Options configures an evaluation.
type Options struct {
	// Vars is the incoming variable environment (never mutated).
	Vars map[string]float64
	// Angle selects radians (default) or degrees for trig functions.
	Angle AngleMode
}

// EvaluateWith parses and evaluates an expression in the context of a variable
// environment, returning the result and the (possibly updated) environment.
// Trig functions use radians. See EvaluateWithOptions for angle-mode control.
//
// In addition to everything Evaluate supports, it allows:
//   - Variable references: a bare identifier resolves to a previously defined
//     variable (e.g. "x * 2").
//   - Assignment: "name = expression" stores the evaluated result under name
//     and returns it. Names may not shadow built-in constants or functions.
//
// The supplied vars map is never mutated; a new map is returned so callers can
// keep their own copy intact. Passing a nil map is fine (treated as empty).
func EvaluateWith(expr string, vars map[string]float64) (float64, map[string]float64, error) {
	return EvaluateWithOptions(expr, Options{Vars: vars})
}

// EvaluateWithOptions is EvaluateWith with explicit options (e.g. angle mode).
func EvaluateWithOptions(expr string, opts Options) (float64, map[string]float64, error) {
	env := make(map[string]float64, len(opts.Vars)+1)
	for k, v := range opts.Vars {
		env[k] = v
	}

	expr, err := normalize(expr)
	if err != nil {
		return 0, env, err
	}
	tokens, err := tokenize(expr)
	if err != nil {
		return 0, env, err
	}

	p := &parser{tokens: tokens, env: env, angle: opts.Angle}

	// Detect a leading assignment: IDENT "=" ...
	if name, ok := p.assignmentTarget(); ok {
		if isReserved(name) {
			return 0, env, errors.New("cannot assign to reserved name: " + name)
		}
		p.advance() // identifier
		p.advance() // "="
		value, err := p.parseExpression()
		if err != nil {
			return 0, env, err
		}
		if !p.atEnd() {
			return 0, env, errors.New("unexpected token: " + p.peek().value)
		}
		if err := checkFinite(value); err != nil {
			return 0, env, err
		}
		env[name] = value
		env["ans"] = value
		return value, env, nil
	}

	result, err := p.parseExpression()
	if err != nil {
		return 0, env, err
	}
	if !p.atEnd() {
		return 0, env, errors.New("unexpected token: " + p.peek().value)
	}
	if err := checkFinite(result); err != nil {
		return 0, env, err
	}
	env["ans"] = result
	return result, env, nil
}
