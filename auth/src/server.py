import os

import jwt
from databases import Database
from sanic import Sanic, exceptions, response
from sanic.request import Request
from sanic_cors import CORS

app = Sanic("gotta-chat")
CORS(app)


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


@app.listener("before_server_stop")
async def close_db(app):
    await app.ctx.db.disconnect()


@app.get("/healthcheck")
async def healthcheck(_: Request):
    return response.text("auth")


@app.post("/signup")
async def singup(request: Request):
    await app.ctx.db.execute(
        query="""INSERT INTO "user"(username, password) VALUES (:username, :password)""",
        values={
            "username": request.json["username"],
            "password": request.json["password"],
        },
    )
    return response.empty()


@app.post("/auth")
async def singup(request: Request):
    user = await app.ctx.db.fetch_one(
        query="""SELECT (id, username) FROM "user" WHERE username = :username and password = :password""",
        values={
            "username": request.json["username"],
            "password": request.json["password"],
        },
    )
    if user is None:
        raise exceptions.Unauthorized("user not found")
    return response.text(jwt.encode(dict(user._mapping), "key", algorithm="HS256"))
