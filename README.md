# url-shortener

A URL shortener in Go.

Create a short code for a long URL, then redirect `GET /<code>` to the original. Redirects are the hot path, so storage sits behind a `Store` interface and can later move from memory to Postgres or Redis without changing HTTP.

## Layout

| Path | Role |
| --- | --- |
| `cmd/server` | Binary. `main` loads config and wires packages. |
| `internal/` | Private application code. Other modules cannot import it. |
| `internal/config` | `ADDR` and `BASE_URL` from the environment. |
| `internal/handler` | HTTP: routes, status codes, JSON, redirects. |
| `internal/shortener` | Validate URLs, generate Base62 codes, read/write `Store`. |
| `internal/store` | `Store` interface and `MemoryStore` (map + `RWMutex`). |

```
HTTP request → handler → shortener.Service → store.Store
```

## Run

Requires Go 1.22+ (method and wildcard routing in `net/http`).

```bash
go test ./...
go run ./cmd/server
```

The server listens on `ADDR` (default `:8080`) and shuts down cleanly on Ctrl+C.

## API

Create a short link:

```bash
curl -X POST localhost:8080/v1/urls -d '{"url":"https://go.dev/doc"}'
# 201 {"code":"m7KLiqv","short_url":"http://localhost:8080/m7KLiqv"}
```

Invalid JSON or a non-http(s) URL returns `400` with `{"error": "..."}`.

Follow it:

```bash
curl -i localhost:8080/m7KLiqv
# 302 Location: https://go.dev/doc
```

Unknown codes return `404`. Redirects use `302`, not `301`, so browsers do not cache them permanently.

| Env | Default | Meaning |
| --- | --- | --- |
| `ADDR` | `:8080` | Listen address |
| `BASE_URL` | `http://localhost:8080` | Prefix for `short_url` in responses |

```bash
go build -o bin/server ./cmd/server
```

## Not implemented yet

- Durable database (Postgres) or cache (Redis)

## Module

```
github.com/olufekosamuel/url-shortener
```
