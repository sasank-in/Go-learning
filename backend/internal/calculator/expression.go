package calculator

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

// tokenType enumerates the kinds of lexical tokens in an expression.
type tokenType int

const (
	tokNumber tokenType = iota
	tokPlus
	tokMinus
	tokStar
	tokSlash
	tokPercent
	tokCaret
	tokLParen
	tokRParen
	tokComma
	tokIdent
	tokEOF
)

type token struct {
	typ   tokenType
	value string
	num   float64
}

// tokenize converts the raw expression string into a slice of tokens,
// ignoring whitespace. It returns an error for any unrecognized character.
func tokenize(expr string) ([]token, error) {
	var tokens []token
	runes := []rune(expr)

	for i := 0; i < len(runes); {
		c := runes[i]

		switch {
		case c == ' ' || c == '\t' || c == '\n' || c == '\r':
			i++
		case c == '+':
			tokens = append(tokens, token{typ: tokPlus, value: "+"})
			i++
		case c == '-':
			tokens = append(tokens, token{typ: tokMinus, value: "-"})
			i++
		case c == '*':
			tokens = append(tokens, token{typ: tokStar, value: "*"})
			i++
		case c == '/':
			tokens = append(tokens, token{typ: tokSlash, value: "/"})
			i++
		case c == '%':
			tokens = append(tokens, token{typ: tokPercent, value: "%"})
			i++
		case c == '^':
			tokens = append(tokens, token{typ: tokCaret, value: "^"})
			i++
		case c == '(':
			tokens = append(tokens, token{typ: tokLParen, value: "("})
			i++
		case c == ')':
			tokens = append(tokens, token{typ: tokRParen, value: ")"})
			i++
		case c == ',':
			tokens = append(tokens, token{typ: tokComma, value: ","})
			i++
		case isDigit(c) || c == '.':
			start := i
			seenDot := false
			for i < len(runes) && (isDigit(runes[i]) || runes[i] == '.') {
				if runes[i] == '.' {
					if seenDot {
						return nil, fmt.Errorf("invalid number: %q", string(runes[start:i+1]))
					}
					seenDot = true
				}
				i++
			}
			// Optional scientific-notation exponent: e / E followed by an
			// optional sign and one or more digits (e.g. 1.5e3, 2E-4).
			if i < len(runes) && (runes[i] == 'e' || runes[i] == 'E') {
				j := i + 1
				if j < len(runes) && (runes[j] == '+' || runes[j] == '-') {
					j++
				}
				if j < len(runes) && isDigit(runes[j]) {
					for j < len(runes) && isDigit(runes[j]) {
						j++
					}
					i = j
				}
			}
			lit := string(runes[start:i])
			val, err := strconv.ParseFloat(lit, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid number: %q", lit)
			}
			tokens = append(tokens, token{typ: tokNumber, value: lit, num: val})
		case isLetter(c):
			start := i
			for i < len(runes) && (isLetter(runes[i]) || isDigit(runes[i])) {
				i++
			}
			tokens = append(tokens, token{typ: tokIdent, value: string(runes[start:i])})
		default:
			return nil, fmt.Errorf("unexpected character: %q", string(c))
		}
	}

	tokens = append(tokens, token{typ: tokEOF, value: ""})
	return tokens, nil
}

func isDigit(c rune) bool { return c >= '0' && c <= '9' }

func isLetter(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

// constants maps recognized constant names (case-insensitive) to their values.
var constants = map[string]float64{
	"pi":  math.Pi,
	"e":   math.E,
	"tau": 2 * math.Pi,
}

// functions maps recognized function names (case-insensitive) to their
// single-argument implementations.
var functions = map[string]func(float64) float64{
	"sqrt":  math.Sqrt,
	"cbrt":  math.Cbrt,
	"abs":   math.Abs,
	"floor": math.Floor,
	"ceil":  math.Ceil,
	"round": math.Round,
	"trunc": math.Trunc,
	"sign":  signum,
	"sin":   math.Sin,
	"cos":   math.Cos,
	"tan":   math.Tan,
	"asin":  math.Asin,
	"acos":  math.Acos,
	"atan":  math.Atan,
	"sinh":  math.Sinh,
	"cosh":  math.Cosh,
	"tanh":  math.Tanh,
	"exp":   math.Exp,
	"ln":    math.Log,
	"log10": math.Log10,
	"log2":  math.Log2,
	"deg":   func(r float64) float64 { return r * 180 / math.Pi },
	"rad":   func(d float64) float64 { return d * math.Pi / 180 },
}

// variadicFunctions maps function names (case-insensitive) to implementations
// that accept one or more arguments. These are dispatched separately from the
// single-argument table above.
var variadicFunctions = map[string]func([]float64) (float64, error){
	"pow":   fnPow,
	"atan2": fnAtan2,
	"hypot": fnHypot,
	"max":   fnMax,
	"min":   fnMin,
	"sum":   fnSum,
	"avg":   fnAvg,
	"log":   fnLog, // log(x) -> base 10; log(x, base) -> arbitrary base
}

func signum(x float64) float64 {
	switch {
	case x > 0:
		return 1
	case x < 0:
		return -1
	default:
		return 0
	}
}

func fnPow(a []float64) (float64, error) {
	if len(a) != 2 {
		return 0, fmt.Errorf("pow expects 2 arguments, got %d", len(a))
	}
	return math.Pow(a[0], a[1]), nil
}

func fnAtan2(a []float64) (float64, error) {
	if len(a) != 2 {
		return 0, fmt.Errorf("atan2 expects 2 arguments, got %d", len(a))
	}
	return math.Atan2(a[0], a[1]), nil
}

func fnHypot(a []float64) (float64, error) {
	if len(a) != 2 {
		return 0, fmt.Errorf("hypot expects 2 arguments, got %d", len(a))
	}
	return math.Hypot(a[0], a[1]), nil
}

func fnLog(a []float64) (float64, error) {
	switch len(a) {
	case 1:
		return math.Log10(a[0]), nil
	case 2:
		if a[1] <= 0 || a[1] == 1 {
			return 0, errors.New("log base must be positive and not equal to 1")
		}
		return math.Log(a[0]) / math.Log(a[1]), nil
	default:
		return 0, fmt.Errorf("log expects 1 or 2 arguments, got %d", len(a))
	}
}

func fnMax(a []float64) (float64, error) {
	if len(a) == 0 {
		return 0, errors.New("max expects at least 1 argument")
	}
	m := a[0]
	for _, v := range a[1:] {
		m = math.Max(m, v)
	}
	return m, nil
}

func fnMin(a []float64) (float64, error) {
	if len(a) == 0 {
		return 0, errors.New("min expects at least 1 argument")
	}
	m := a[0]
	for _, v := range a[1:] {
		m = math.Min(m, v)
	}
	return m, nil
}

func fnSum(a []float64) (float64, error) {
	var s float64
	for _, v := range a {
		s += v
	}
	return s, nil
}

func fnAvg(a []float64) (float64, error) {
	if len(a) == 0 {
		return 0, errors.New("avg expects at least 1 argument")
	}
	s, _ := fnSum(a)
	return s / float64(len(a)), nil
}

// parser is a recursive-descent parser/evaluator over a token slice.
//
// Grammar (lowest to highest precedence):
//
//	expression = term { ("+" | "-") term }
//	term       = factor { ("*" | "/" | "%") factor }
//	factor     = unary { "^" factor }            // right-associative
//	unary      = ("+" | "-") unary | primary
//	primary    = number | "(" expression ")"
type parser struct {
	tokens []token
	pos    int
}

func (p *parser) peek() token { return p.tokens[p.pos] }

func (p *parser) atEnd() bool { return p.peek().typ == tokEOF }

func (p *parser) advance() token {
	t := p.tokens[p.pos]
	if !p.atEnd() {
		p.pos++
	}
	return t
}

func (p *parser) parseExpression() (float64, error) {
	left, err := p.parseTerm()
	if err != nil {
		return 0, err
	}
	for {
		switch p.peek().typ {
		case tokPlus:
			p.advance()
			right, err := p.parseTerm()
			if err != nil {
				return 0, err
			}
			left += right
		case tokMinus:
			p.advance()
			right, err := p.parseTerm()
			if err != nil {
				return 0, err
			}
			left -= right
		default:
			return left, nil
		}
	}
}

func (p *parser) parseTerm() (float64, error) {
	left, err := p.parseFactor()
	if err != nil {
		return 0, err
	}
	for {
		switch p.peek().typ {
		case tokStar:
			p.advance()
			right, err := p.parseFactor()
			if err != nil {
				return 0, err
			}
			left *= right
		case tokSlash:
			p.advance()
			right, err := p.parseFactor()
			if err != nil {
				return 0, err
			}
			if right == 0 {
				return 0, ErrDivideByZero
			}
			left /= right
		case tokPercent:
			p.advance()
			right, err := p.parseFactor()
			if err != nil {
				return 0, err
			}
			if right == 0 {
				return 0, ErrDivideByZero
			}
			left = math.Mod(left, right)
		default:
			return left, nil
		}
	}
}

func (p *parser) parseFactor() (float64, error) {
	base, err := p.parseUnary()
	if err != nil {
		return 0, err
	}
	if p.peek().typ == tokCaret {
		p.advance()
		// Right-associative: 2^3^2 == 2^(3^2).
		exp, err := p.parseFactor()
		if err != nil {
			return 0, err
		}
		return math.Pow(base, exp), nil
	}
	return base, nil
}

func (p *parser) parseUnary() (float64, error) {
	switch p.peek().typ {
	case tokMinus:
		p.advance()
		v, err := p.parseUnary()
		if err != nil {
			return 0, err
		}
		return -v, nil
	case tokPlus:
		p.advance()
		return p.parseUnary()
	default:
		return p.parsePrimary()
	}
}

func (p *parser) parsePrimary() (float64, error) {
	t := p.peek()
	switch t.typ {
	case tokNumber:
		p.advance()
		return t.num, nil
	case tokLParen:
		p.advance()
		v, err := p.parseExpression()
		if err != nil {
			return 0, err
		}
		if p.peek().typ != tokRParen {
			return 0, errors.New("expected closing ')'")
		}
		p.advance()
		return v, nil
	case tokIdent:
		return p.parseIdent(t)
	case tokEOF:
		return 0, errors.New("unexpected end of expression")
	default:
		return 0, fmt.Errorf("unexpected token: %q", t.value)
	}
}

// parseIdent resolves an identifier as either a constant (e.g. "pi") or a
// function call (e.g. "sqrt(2)"). The name is matched case-insensitively.
func (p *parser) parseIdent(t token) (float64, error) {
	p.advance()
	name := strings.ToLower(t.value)

	// Function call: identifier immediately followed by "(".
	if p.peek().typ == tokLParen {
		args, err := p.parseArgList(t.value)
		if err != nil {
			return 0, err
		}

		// Variadic / multi-argument functions take precedence so names like
		// "log" can accept a flexible arity.
		if fn, ok := variadicFunctions[name]; ok {
			return fn(args)
		}
		if fn, ok := functions[name]; ok {
			if len(args) != 1 {
				return 0, fmt.Errorf("%s expects 1 argument, got %d", name, len(args))
			}
			return fn(args[0]), nil
		}
		return 0, fmt.Errorf("unknown function: %q", t.value)
	}

	// Otherwise it must be a constant.
	if val, ok := constants[name]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("unknown identifier: %q", t.value)
}

// parseArgList parses "( expr { , expr } )" and returns the evaluated
// arguments. The opening "(" is expected to be the current token. An empty
// argument list "()" is an error since all supported functions take at least
// one argument. fnName is used only for error messages.
func (p *parser) parseArgList(fnName string) ([]float64, error) {
	p.advance() // consume "("

	if p.peek().typ == tokRParen {
		return nil, fmt.Errorf("%s requires at least one argument", fnName)
	}

	var args []float64
	for {
		v, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		args = append(args, v)

		switch p.peek().typ {
		case tokComma:
			p.advance()
			continue
		case tokRParen:
			p.advance()
			return args, nil
		default:
			return nil, fmt.Errorf("expected ',' or ')' in arguments to %q", fnName)
		}
	}
}

// normalize is a small helper to trim and validate that an expression is
// non-empty before evaluation.
func normalize(expr string) (string, error) {
	trimmed := strings.TrimSpace(expr)
	if trimmed == "" {
		return "", errors.New("empty expression")
	}
	return trimmed, nil
}
