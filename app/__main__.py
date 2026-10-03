"""`python -m app` runs the server on HOST:PORT.

uvicorn handles SIGINT/SIGTERM and lets in-flight requests finish before
exiting.
"""

import uvicorn

from app.main import cfg

if __name__ == "__main__":
    uvicorn.run("app.main:app", host=cfg.host, port=cfg.port, timeout_keep_alive=60)
