import pytest
from fastapi.testclient import TestClient

from app.api import create_app
from app.shortener import Service
from app.store import MemoryStore


@pytest.fixture
def client():
    app = create_app(Service(MemoryStore()), "http://short.test/")
    return TestClient(app, follow_redirects=False)


def test_create_and_redirect(client):
    resp = client.post("/v1/urls", json={"url": "https://example.com/a"})
    assert resp.status_code == 201, resp.text
    body = resp.json()
    assert body["code"]
    assert body["short_url"] == "http://short.test/" + body["code"]

    resp = client.get("/" + body["code"])
    assert resp.status_code == 302
    assert resp.headers["location"] == "https://example.com/a"


@pytest.mark.parametrize(
    "body",
    [
        "",
        "not json",
        '{"url":"example.com"}',
        '{"url":"ftp://example.com"}',
        '{"link":"https://example.com"}',
    ],
)
def test_create_bad_requests(client, body):
    resp = client.post(
        "/v1/urls", content=body, headers={"Content-Type": "application/json"}
    )
    assert resp.status_code == 400, resp.text
    assert "error" in resp.json()


def test_redirect_unknown_code(client):
    resp = client.get("/missing")
    assert resp.status_code == 404


def test_wrong_method(client):
    # GET /v1/urls does not match GET /{code} (single segment), and
    # /v1/urls only allows POST.
    assert client.get("/v1/urls").status_code == 405
