import os
import typing

import jwt
import strawberry
from sanic import Sanic, exceptions, response
from sanic.request import Request
from sanic_cors import CORS
from strawberry.sanic.views import GraphQLView

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


@strawberry.type
class UserResourceData:
    resource_name: str
    factory_level: int
    amount: int
    time_until_upgrade_complete: int | None


@strawberry.type
class Query:
    @strawberry.field
    async def dashboard(self, info) -> typing.List[UserResourceData]:
        return [
            UserResourceData(
                resource_name="asd",
                factory_level=info.context.user_id,
                amount=1,
                time_until_upgrade_complete=None,
            )
        ]


@strawberry.type
class Mutation:
    @strawberry.mutation
    async def upgrade_factory(self, info, resource_name: str) -> None:
        ...


class ControllerGraphQLView(GraphQLView):
    async def get_context(self, request: Request) -> typing.Dict:
        return request.ctx


app.add_route(
    ControllerGraphQLView.as_view(
        schema=strawberry.Schema(query=Query, mutation=Mutation)
    ),
    "/graphql",
)
