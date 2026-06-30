# Calculator Backend (Go)

A small HTTP JSON API providing calculator operations.

## Layout

```
backend/
  cmd/server/         # main entrypoint
  internal/calculator # core arithmetic logic (no transport concerns)
  internal/handlers   # HTTP handlers / JSON API
```

## Run

```bash
cd backend
go run ./cmd/server
# listens on :8080 (override with PORT env var)
```

## Test

```bash
cd backend
go test ./...
```

## API

### `GET /health`
Returns `{"status":"ok"}`.

### `POST /calculate`

Two request forms are supported. If `expression` is non-empty it takes precedence.

**Expression form (recommended)** — evaluates any arithmetic expression:
```json
{ "expression": "2 + 3 * (4 - 1)" }
```
Supports `+ - * / % ^`, unary minus/plus, parentheses, decimals, and
scientific notation (`1.5e3`), with standard operator precedence. `^`
(exponent) is right-associative.

The postfix `!` operator computes factorial (e.g. `5!` → 120).

Functions (case-insensitive):

- **Single-argument**: `sqrt`, `cbrt`, `abs`, `floor`, `ceil`, `round`,
  `trunc`, `sign`, `fact`, `sin`, `cos`, `tan`, `asin`, `acos`, `atan`,
  `sinh`, `cosh`, `tanh`, `exp`, `ln`, `log10`, `log2`, `deg`
  (radians→degrees), `rad` (degrees→radians)
- **Multi-argument**: `pow(x, y)`, `atan2(y, x)`, `hypot(x, y)`,
  `log(x[, base])`, `max(...)`, `min(...)`, `sum(...)`, `avg(...)`,
  `gcd(...)`, `lcm(...)`, `ncr(n, r)`, `npr(n, r)`

Constants: `pi`, `e`, `tau`.

**Angle mode**: include `"angleMode": "deg"` in the request to make `sin`,
`cos`, `tan` (and inverse trig) work in degrees. Defaults to radians.

Examples: `5!`, `ncr(5, 2)`, `gcd(48, 36)`, `sin(90)` with `angleMode: "deg"`.

**Structured form** — a single two-operand operation:
```json
{ "operation": "add", "a": 2, "b": 3 }
```
Supported operations: `add`, `subtract`, `multiply`, `divide`.

Success response:
```json
{ "result": 11 }
```

Error response (e.g. division by zero, malformed expression, unknown operation):
```json
{ "error": "division by zero" }
```

Examples:
```bash
curl -X POST localhost:8080/calculate \
  -H 'Content-Type: application/json' \
  -d '{"expression":"2 + 3 * (4 - 1)"}'        # -> {"result":11}

curl -X POST localhost:8080/calculate \
  -H 'Content-Type: application/json' \
  -d '{"expression":"2 ^ 10"}'                 # -> {"result":1024}
```
