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
Supports `+ - * / % ^`, unary minus/plus, parentheses, and decimals, with
standard operator precedence. `^` (exponent) is right-associative.

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
