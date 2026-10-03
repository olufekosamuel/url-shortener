import asyncio

import pytest

from app.store import ConflictError, MemoryStore, NotFoundError


def test_save_and_get():
    s = MemoryStore()
    asyncio.run(s.save("abc123", "https://example.com"))
    assert asyncio.run(s.get("abc123")) == "https://example.com"


def test_missing_code():
    with pytest.raises(NotFoundError):
        asyncio.run(MemoryStore().get("missing"))


def test_duplicate_code():
    s = MemoryStore()
    asyncio.run(s.save("abc123", "https://example.com/a"))
    with pytest.raises(ConflictError):
        asyncio.run(s.save("abc123", "https://example.com/b"))
    assert asyncio.run(s.get("abc123")) == "https://example.com/a"
