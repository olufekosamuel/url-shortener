"""Storage for short-code → long-URL mappings."""

from typing import Protocol


class NotFoundError(Exception):
    """A short code has no mapping."""

    def __init__(self) -> None:
        super().__init__("short url not found")


class ConflictError(Exception):
    """save was asked to use a code that already exists."""

    def __init__(self) -> None:
        super().__init__("short code already exists")


class Store(Protocol):
    """Persists short-code → long-URL mappings.

    HTTP and domain logic depend on this protocol, not on a dict or a
    database. Methods are async so a remote backend (Postgres, Redis) can
    replace MemoryStore without changing callers.
    """

    async def save(self, code: str, long_url: str) -> None: ...

    async def get(self, code: str) -> str: ...


class MemoryStore:
    """Keeps mappings in a dict.

    No lock is needed: the app runs on one asyncio event loop, and neither
    method awaits between checking and writing the dict, so no other request
    can interleave. Data is lost when the process exits.
    """

    def __init__(self) -> None:
        self._urls: dict[str, str] = {}  # code -> long URL

    async def save(self, code: str, long_url: str) -> None:
        if code in self._urls:
            raise ConflictError()
        self._urls[code] = long_url

    async def get(self, code: str) -> str:
        try:
            return self._urls[code]
        except KeyError:
            raise NotFoundError() from None
