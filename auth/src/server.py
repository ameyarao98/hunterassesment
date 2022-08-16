import os

from databases import Database
from sanic import Sanic, response
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
            username VARCHAR(100) PRIMARY KEY, 
            password VARCHAR(256)
        )
            """
    )


@app.listener("before_server_stop")
async def close_db(app):
    await app.ctx.db.disconnect()


@app.get("/")
async def healthcheck(request: Request):
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
