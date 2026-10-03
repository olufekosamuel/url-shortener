"""HTTP: routing, JSON, status codes, redirects.

Routes:
    POST /v1/urls   create a short link
    GET  /{code}    302 to the original URL

This module stays thin. It calls shortener.Service and does not import a
database driver.
"""

import logging

from fastapi import FastAPI, Request
from fastapi.exceptions import RequestValidationError
from fastapi.responses import JSONResponse, RedirectResponse
from pydantic import BaseModel, ConfigDict

from app.shortener import InvalidURLError, Service
from app.store import NotFoundError

log = logging.getLogger(__name__)


class CreateRequest(BaseModel):
    # Reject unknown fields so a typo like {"link": ...} is a 400, not a
    # silently missing url.
    model_config = ConfigDict(extra="forbid")

    url: str


class CreateResponse(BaseModel):
    code: str
    short_url: str


def create_app(service: Service, base_url: str) -> FastAPI:
    base_url = base_url.rstrip("/")
    app = FastAPI(title="url-shortener")

    # FastAPI answers bad bodies with 422 and a detailed list; keep the
    # simpler 400 {"error": ...} contract for every client error.
    @app.exception_handler(RequestValidationError)
    async def bad_request(_: Request, __: RequestValidationError) -> JSONResponse:
        return JSONResponse({"error": "invalid JSON body"}, status_code=400)

    @app.exception_handler(InvalidURLError)
    async def invalid_url(_: Request, exc: InvalidURLError) -> JSONResponse:
        return JSONResponse({"error": str(exc)}, status_code=400)

    @app.exception_handler(NotFoundError)
    async def not_found(_: Request, exc: NotFoundError) -> JSONResponse:
        return JSONResponse({"error": str(exc)}, status_code=404)

    @app.post("/v1/urls", status_code=201, response_model=CreateResponse)
    async def create(req: CreateRequest) -> CreateResponse:
        code = await service.shorten(req.url)
        return CreateResponse(code=code, short_url=f"{base_url}/{code}")

    # 302, not 301: browsers cache 301s indefinitely, which would hide later
    # changes to a link and skip any future click counting.
    @app.get("/{code}")
    async def redirect(code: str) -> RedirectResponse:
        long_url = await service.resolve(code)
        return RedirectResponse(long_url, status_code=302)

    return app
