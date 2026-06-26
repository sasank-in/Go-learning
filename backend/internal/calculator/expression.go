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
	"sin":   math.Sin,
	"cos":   math.Cos,
	"tan":   math.Tan,
	"asin":  math.Asin,
	"acos":  math.Acos,
	"atan":  math.Atan,
	"exp":   math.Exp,
	"ln":    math.Log,
	"log":   math.Log10,
	"log2":  math.Log2,
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
		fn, ok := functions[name]
		if !ok {
			return 0, fmt.Errorf("unknown function: %q", t.value)
		}
		p.advance() // consume "("
		arg, err := p.parseExpression()
		if err != nil {
			return 0, err
		}
		if p.peek().typ != tokRParen {
			return 0, fmt.Errorf("expected closing ')' after argument to %q", t.value)
		}
		p.advance() // consume ")"
		return fn(arg), nil
	}

	// Otherwise it must be a constant.
	if val, ok := constants[name]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("unknown identifier: %q", t.value)
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
