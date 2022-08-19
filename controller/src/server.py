import asyncio
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
class FactoryData:
    resource_name: str
    factory_level: int
    production_per_second: int
    next_upgrade_duration: int | None
    # upgrade_cost: dict



@strawberry.type
class UserResourceData:
    resource_name: str
    factory_level: int
    amount: int
    time_until_upgrade_complete: int | None


@strawberry.type
class Query:
    @strawberry.field
    async def factory_data(self, info) -> typing.List[FactoryData]:
        return []


@strawberry.type
class Mutation:
    @strawberry.mutation
    async def upgrade_factory(self, info, resource_name: str) -> bool:
        return True


@strawberry.type
class Subscription:
    @strawberry.subscription
    async def user_resource(self, info) -> typing.AsyncGenerator[UserResourceData, None]:
        x = 0
        while True:
            yield [
                UserResourceData(
                    resource_name="asd",
                    factory_level=info.context.user_id,
                    amount=x,
                    time_until_upgrade_complete=None,
                )
            ]
            x += 1
            await asyncio.sleep(1)


class ControllerGraphQLView(GraphQLView):
    async def get_context(self, request: Request) -> typing.Dict:
        return request.ctx


app.add_route(
    ControllerGraphQLView.as_view(
        schema=strawberry.Schema(
            query=Query, mutation=Mutation, subscription=Subscription
        )
    ),
    "/graphql",
)
