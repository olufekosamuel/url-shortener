import asyncio

import pytest

from app.shortener import ALPHABET, InvalidURLError, Service, random_code
from app.store import MemoryStore, NotFoundError


def test_shorten_and_resolve():
    svc = Service(MemoryStore(), generate=lambda: "abc1234")

    code = asyncio.run(svc.shorten("https://example.com/a"))
    assert code == "abc1234"
    assert asyncio.run(svc.resolve(code)) == "https://example.com/a"


@pytest.mark.parametrize(
    "raw", ["", "example.com", "javascript:alert(1)", "ftp://example.com", "http://[::1"]
)
def test_shorten_rejects_invalid_url(raw):
    with pytest.raises(InvalidURLError):
        asyncio.run(Service(MemoryStore()).shorten(raw))


def test_shorten_retries_on_conflict():
    st = MemoryStore()
    asyncio.run(st.save("aaaaaaa", "https://example.com/old"))
    codes = iter(["aaaaaaa", "bbbbbbb"])

    svc = Service(st, generate=lambda: next(codes))
    assert asyncio.run(svc.shorten("https://example.com/new")) == "bbbbbbb"


def test_resolve_missing():
    with pytest.raises(NotFoundError):
        asyncio.run(Service(MemoryStore()).resolve("missing"))


def test_random_code():
    got = random_code(7)
    assert len(got) == 7
    assert all(c in ALPHABET for c in got)
