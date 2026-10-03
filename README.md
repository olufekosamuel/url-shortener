# url-shortener

A URL shortener in Python with FastAPI.

Create a short code for a long URL, then redirect `GET /<code>` to the original. Redirects are the hot path, so storage sits behind a `Store` protocol and can later move from memory to Postgres or Redis without changing HTTP.

## Layout

| Path | Role |
| --- | --- |
| `app/main.py` | Entrypoint. Loads config and wires modules into the FastAPI `app`. |
| `app/__main__.py` | `python -m app` runs uvicorn on `HOST:PORT`. |
| `app/config.py` | `HOST`, `PORT` and `BASE_URL` from the environment. |
| `app/api.py` | HTTP: routes, status codes, JSON, redirects. |
| `app/shortener.py` | Validate URLs, generate Base62 codes, read/write `Store`. |
| `app/store.py` | `Store` protocol and `MemoryStore` (a dict). |
| `tests/` | pytest suites for each layer. |

```
HTTP request → api → shortener.Service → store.Store
```

## Run

Requires Python 3.10+ and [uv](https://docs.astral.sh/uv/).

```bash
uv sync
uv run pytest
uv run python -m app
```

The server listens on `127.0.0.1:8080` by default and shuts down cleanly on Ctrl+C. Interactive API docs are at `/docs`.

## API

Create a short link:

```bash
curl -X POST localhost:8080/v1/urls -H 'Content-Type: application/json' \
  -d '{"url":"https://go.dev/doc"}'
# 201 {"code":"qXjGFQc","short_url":"http://localhost:8080/qXjGFQc"}
```

Invalid JSON, unknown fields, or a non-http(s) URL return `400` with `{"error": "..."}`.

Follow it:

```bash
curl -i localhost:8080/qXjGFQc
# 302 location: https://go.dev/doc
```

Unknown codes return `404`. Redirects use `302`, not `301`, so browsers do not cache them permanently.

| Env | Default | Meaning |
| --- | --- | --- |
| `HOST` | `127.0.0.1` | Interface to listen on (`0.0.0.0` in a container) |
| `PORT` | `8080` | Listen port |
| `BASE_URL` | `http://localhost:8080` | Prefix for `short_url` in responses |

## Not implemented yet

- Durable database (Postgres) or cache (Redis)
