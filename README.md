# Calculator Application

A scientific calculator with a **Go** backend (expression-evaluation HTTP API)
and an attractive, corporate-styled **TypeScript** frontend (Vite + Bun).

## Structure

```
calculator-application/
  backend/    # Go HTTP JSON API — expression parser + evaluator
  frontend/   # Vite + TypeScript UI (managed with Bun)
```

## Features

- Full arithmetic expressions with precedence and parentheses: `2 + 3 * (4 - 1)`
- Operators: `+ - * / % ^`, unary minus/plus, decimals, scientific notation
- Scientific functions: trig/hyperbolic, `sqrt`, `ln`, `log`, `exp`, rounding,
  `deg`/`rad`, and more
- Multi-argument functions: `pow`, `hypot`, `atan2`, `max`, `min`, `sum`, `avg`,
  `log(x, base)`
- Constants: `pi`, `e`, `tau`
- Calculation history, keyboard support, friendly error messages

## Quick start

```bash
# Terminal 1 — backend (http://localhost:8080)
cd backend && go run ./cmd/server

# Terminal 2 — frontend (http://localhost:5173)
cd frontend && bun install && bun run dev
```

See [backend/README.md](backend/README.md) and
[frontend/README.md](frontend/README.md) for details.
