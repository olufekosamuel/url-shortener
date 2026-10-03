"""Process settings from the environment.

12-factor style: the same code runs locally and in production by changing
HOST, PORT, and BASE_URL, not by changing code.
"""

import os
from dataclasses import dataclass


@dataclass(frozen=True)
class Config:
    # Interface to listen on, e.g. "127.0.0.1" or "0.0.0.0".
    host: str
    # TCP port to listen on.
    port: int
    # How short links are presented, e.g. "http://localhost:8080".
    base_url: str


def load() -> Config:
    return Config(
        host=os.environ.get("HOST") or "127.0.0.1",
        port=int(os.environ.get("PORT") or "8080"),
        base_url=os.environ.get("BASE_URL") or "http://localhost:8080",
    )
