import asyncio
import os
import typing

import jwt
import strawberry
from sanic import Blueprint, Sanic, exceptions, response
from sanic.request import Request
from sanic_cors import CORS
from strawberry.sanic.views import GraphQLView

app = Sanic("gotta-chat")
CORS(app)


@app.listener("before_server_start")
async def get_auth_public_key(app):
    app.ctx.auth_public_key = os.getenv("AUTH_PUBLIC_KEY")


@app.get("/healthcheck")
async def healthcheck(request: Request):
    return response.text("controller")


graphql_blueprint = Blueprint("graphql")


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
    async def enter_game(self, info) -> bool:
        return True

    @strawberry.mutation
    async def upgrade_factory(self, info, resource_name: str) -> bool:
        return True


@strawberry.type
class Subscription:
    @strawberry.subscription
    async def user_resources(
        self, info
    ) -> typing.AsyncGenerator[typing.List[UserResourceData], None]:
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

    @strawberry.subscription
    async def count(self, target: int = 100) -> typing.AsyncGenerator[int, None]:
        for i in range(target):
            yield i
            await asyncio.sleep(0.5)


class ControllerGraphQLView(GraphQLView):
    async def get_context(self, request: Request) -> typing.Dict:
        return request.ctx


graphql_blueprint.add_route(
    ControllerGraphQLView.as_view(
        schema=strawberry.Schema(
            query=Query, mutation=Mutation, subscription=Subscription
        )
    ),
    "/graphql",
)


@graphql_blueprint.middleware("request")
async def verify_authorization(request):
    try:
        claims = jwt.decode(
            request.headers.get("Creds"), app.ctx.auth_public_key, algorithms="RS256"
        )
    except jwt.DecodeError:
        raise exceptions.Unauthorized("Invalid creds")
    request.ctx.user_id = claims["id"]


app.blueprint(graphql_blueprint)
