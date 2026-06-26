# Calculator Frontend

An attractive, corporate-styled scientific calculator UI built with **Vite +
TypeScript** and managed with **Bun**. It talks to the Go backend over HTTP.

## Stack

- Bun (package manager / runtime)
- Vite 6 + TypeScript (strict)
- No UI framework — a single, dependency-free TS module + handcrafted CSS

## Features

- Scientific keypad: `sin cos tan √ ^ ln log π e %`, plus arithmetic
- Full keyboard support (digits, operators, `Enter` =, `Esc` clear, `Backspace`)
- Calculation history (click an entry to reuse it)
- Friendly error messages surfaced from the backend
- Responsive layout, dark corporate theme

## Develop

The backend must be running first (defaults to `http://localhost:8080`):

```bash
# in ../backend
go run ./cmd/server
```

Then start the frontend:

```bash
bun install      # first time only
bun run dev      # http://localhost:5173
```

Vite proxies `/api/*` to the Go backend (see `vite.config.ts`), so no CORS or
URL juggling is needed in development.

## Build

```bash
bun run build    # type-checks, then emits static assets to dist/
bun run preview  # preview the production build
```

## Configuration

| Variable        | Default | Purpose                                |
| --------------- | ------- | -------------------------------------- |
| `VITE_API_BASE` | `/api`  | Backend base URL used by the API client |
