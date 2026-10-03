"""Process entrypoint.

Loads config, constructs dependencies, and exposes `app` for uvicorn.
URL validation, code generation, and storage live in their own modules.
"""

from app import config
from app.api import create_app
from app.shortener import Service
from app.store import MemoryStore

cfg = config.load()
app = create_app(Service(MemoryStore()), cfg.base_url)
