# url-shortener

A URL shortener in Go.

Create a short code for a long URL, then redirect `GET /<code>` to the original. Redirects are the hot path, so storage sits behind a `Store` interface and can later move from memory to Postgres or Redis without changing HTTP.

## Layout

| Path | Role |
| --- | --- |
| `cmd/server` | Binary. `main` loads config and wires packages. |
| `internal/` | Private application code. Other modules cannot import it. |
| `internal/config` | `ADDR` and `BASE_URL` from the environment. |
| `internal/handler` | HTTP: routes, status codes, JSON. Not implemented yet. |
| `internal/shortener` | Validate URLs, generate Base62 codes, read/write `Store`. |
| `internal/store` | `Store` interface and `MemoryStore` (map + `RWMutex`). |

```
HTTP request → handler → shortener.Service → store.Store
```

## Run

Requires Go 1.21+.

```bash
go test ./...
go run ./cmd/server
```

The process logs config and exits. The HTTP server is not started yet.

```bash
go build -o bin/server ./cmd/server
```

## Not implemented yet

- HTTP server
- Durable database (Postgres) or cache (Redis)

## Module

```
github.com/olufekosamuel/url-shortener
```
