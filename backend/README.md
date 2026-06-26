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

Functions (case-insensitive):

- **Single-argument**: `sqrt`, `cbrt`, `abs`, `floor`, `ceil`, `round`,
  `trunc`, `sign`, `sin`, `cos`, `tan`, `asin`, `acos`, `atan`, `sinh`,
  `cosh`, `tanh`, `exp`, `ln`, `log10`, `log2`, `deg` (radians→degrees),
  `rad` (degrees→radians)
- **Multi-argument**: `pow(x, y)`, `atan2(y, x)`, `hypot(x, y)`,
  `log(x[, base])`, `max(...)`, `min(...)`, `sum(...)`, `avg(...)`

Constants: `pi`, `e`, `tau`.

Examples: `sqrt(16) + 2 * pi`, `pow(2, 10)`, `max(1, 9, 4)`, `log(8, 2)`.

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
