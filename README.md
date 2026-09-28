# url-shortener

A URL shortener in Go, built in layers

Structure plus an in-memory `Store`. Packages compile; nothing shortens a URL over HTTP yet.

## What we are building

Two requests, two jobs:

1. **Create** — `POST /v1/urls` with a long URL → generate a short code → store the mapping → return `http://host/<code>`.
2. **Redirect** — `GET /<code>` → look up the long URL → respond with `302` (or `301`).

Reads (redirects) will dominate writes (creates). That is why we keep storage behind an interface: later we can add Postgres for durability and Redis for hot lookups without rewriting HTTP.

## Why this folder layout

Go convention, not personal taste:

| Path | Role |
| --- | --- |
| `cmd/server` | The binary. `main` wires dependencies. No business logic. |
| `internal/` | Private application code. Other Go modules cannot import it. |
| `internal/config` | Env-based settings (`ADDR`, `BASE_URL`). |
| `internal/handler` | HTTP transport: routes, status codes, JSON. |
| `internal/shortener` | Domain logic: validation, code generation, orchestration. |
| `internal/store` | Persistence **port** (`Store`) plus `MemoryStore` (map + `RWMutex`). |

Request flow we will implement:

```
HTTP request → handler → shortener.Service → store.Store → memory / Postgres / Redis
```

Each arrow is a dependency. `handler` never talks to a database. That is the design you can draw on a whiteboard.

## Run the skeleton

Requires Go 1.21+.

```bash
go run ./cmd/server
```

You should see the loaded `addr` and `baseURL`. The HTTP server is not started yet.

```bash
go test ./internal/store
go build -o bin/server ./cmd/server
```

## What is intentionally missing

- No HTTP server
- No code generator (Base62 / hash / counter)
- No durable database (Postgres) or cache (Redis)

Next: the shortener service (validation + code generation), then HTTP.

## Module path

```
github.com/olufekosamuel/url-shortener
```
