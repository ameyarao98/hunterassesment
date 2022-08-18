import os

import jwt
from sanic import Sanic, exceptions, response
from sanic.request import Request
from sanic_cors import CORS

app = Sanic("gotta-chat")
CORS(app)


@app.listener("before_server_start")
async def get_auth_public_key(app):
    app.ctx.auth_public_key = os.getenv("AUTH_PUBLIC_KEY")


async def verify_authorization(request):
    try:
        claims = jwt.decode(
            request.headers.get("Creds"), app.ctx.auth_public_key, algorithms="RS256"
        )
    except jwt.DecodeError:
        raise exceptions.Unauthorized("Invalid creds")
    request.ctx.user_id = claims["id"]


app.register_middleware(verify_authorization, "request")


@app.get("/healthcheck")
async def healthcheck(request: Request):
    return response.text("controller")
