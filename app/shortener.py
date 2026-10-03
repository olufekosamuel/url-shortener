"""Validate URLs, generate short codes, and read/write the store."""

import secrets
import string
from collections.abc import Callable
from urllib.parse import urlsplit

from app.store import ConflictError, NotFoundError, Store

ALPHABET = string.digits + string.ascii_lowercase + string.ascii_uppercase
CODE_LENGTH = 7
MAX_ATTEMPTS = 8


class InvalidURLError(ValueError):
    def __init__(self) -> None:
        super().__init__("url must be an absolute http or https URL")


class CodeCollisionError(Exception):
    def __init__(self) -> None:
        super().__init__("could not allocate a unique short code")


def random_code(n: int = CODE_LENGTH) -> str:
    """Return n random characters from the Base62 alphabet.

    Base62 stays URL-safe (no +, /, or punctuation). Seven characters give
    62^7 possible codes, so collisions are rare; save still rejects
    duplicates and shorten retries. secrets.choice picks uniformly, so no
    character is favoured.
    """
    if n <= 0:
        raise ValueError("code length must be positive")
    return "".join(secrets.choice(ALPHABET) for _ in range(n))


class Service:
    def __init__(
        self,
        store: Store,
        generate: Callable[[], str] = random_code,
    ) -> None:
        self._store = store
        self._generate = generate

    async def shorten(self, raw_url: str) -> str:
        """Validate raw_url, allocate a unique code, and store the mapping."""
        long_url = normalize_url(raw_url)
        for _ in range(MAX_ATTEMPTS):
            code = self._generate()
            try:
                await self._store.save(code, long_url)
            except ConflictError:
                continue
            return code
        raise CodeCollisionError()

    async def resolve(self, code: str) -> str:
        """Return the long URL for code."""
        if not code:
            raise NotFoundError()
        return await self._store.get(code)


def normalize_url(raw: str) -> str:
    try:
        parts = urlsplit(raw.strip())
    except ValueError:
        raise InvalidURLError() from None
    if parts.scheme not in ("http", "https") or not parts.netloc:
        raise InvalidURLError()
    return parts.geturl()
