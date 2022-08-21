import os

import jwt
from databases import Database
from sanic import Sanic, exceptions, request, response
from sanic_cors import CORS

app = Sanic("gotta-chat")
CORS(app)

PRIVATE_KEY = """
-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----
"""


@app.listener("before_server_start")
async def setup_db(app):
    app.ctx.db = Database(os.getenv("POSTGRES_DSN"))
    await app.ctx.db.connect()
    await app.ctx.db.execute(
        query="""
        CREATE TABLE IF NOT EXISTS "user"(
            id SERIAL PRIMARY KEY,
            username VARCHAR(100) UNIQUE, 
            password VARCHAR(256)
        )
            """
    )


@app.get("/healthcheck")
async def healthcheck(_: request.Request):
    return response.text("auth")


@app.post("/signup")
async def singup(request: request.Request):
    await app.ctx.db.execute(
        query="""INSERT INTO "user"(username, password) VALUES (:username, :password)""",
        values={
            "username": request.json["username"],
            "password": str(hash(request.json["password"])),
        },
    )
    return response.empty()


@app.post("/auth")
async def auth(request: request.Request):
    user = await app.ctx.db.fetch_one(
        query="""SELECT (id, username) FROM "user" WHERE username=:username and password=:password""",
        values={
            "username": request.json["username"],
            "password": str(hash(request.json["password"])),
        },
    )
    if user is None:
        raise exceptions.Unauthorized("user not found")
    return response.text(
        jwt.encode(
            {
                "id": user.row[0],
                "username": user.row[1],
            },
            PRIVATE_KEY,
            algorithm="RS256",
        )
    )
