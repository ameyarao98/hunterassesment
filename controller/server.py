# import asyncio
import os
# import typing

import grpc
import jwt
import strawberry
from sanic import Blueprint, Sanic, exceptions, request, response
from sanic_cors import CORS
from strawberry.sanic.views import GraphQLView

import factory_pb2
import factory_pb2_grpc

app = Sanic("controller")
CORS(app)


@app.listener("before_server_start")
async def get_auth_public_key(app):
    app.ctx.auth_public_key = os.getenv("AUTH_PUBLIC_KEY")
    app.ctx.grpc_client = factory_pb2_grpc.FactoryStub(
        grpc.insecure_channel("factory:8080")
    )


@app.get("/healthcheck")
async def healthcheck(request: request.Request):
    return response.text("controller")


graphql_blueprint = Blueprint("graphql")


@strawberry.type
class FactoryData:
    resource_name: str
    factory_level: int
    production_per_second: int
    next_upgrade_duration: int | None


@strawberry.type
class UserResourceData:
    resource_name: str
    factory_level: int
    amount: int
    production_rate: int
    time_until_upgrade_complete: int | None


@strawberry.type
class Query:
    @strawberry.field
    async def factory_data(self, info) -> list[FactoryData]:
        return (
            app.ctx.grpc_client.GetFactoryData.future(
                factory_pb2.GetFactoryDataRequest()
            )
            .result()
            .factory_datas
        )

    @strawberry.field
    async def user_resource_data(self, info) -> list[UserResourceData]:
        return (
            app.ctx.grpc_client.GetUserResourceData.future(
                factory_pb2.GetUserResourceDataRequest(user_id=info.context.user_id)
            )
            .result()
            .user_resource_datas
        )


@strawberry.type
class Mutation:
    @strawberry.mutation
    async def enter_game(self, info) -> bool:
        return (
            app.ctx.grpc_client.CreateUser.future(
                factory_pb2.CreateUserRequest(user_id=info.context.user_id)
            )
            .result()
            .created
        )

    @strawberry.mutation
    async def upgrade_factory(self, info, resource_name: str) -> bool:
        return (
            app.ctx.grpc_client.UpgradeFactory.future(
                factory_pb2.UpgradeFactoryRequest(
                    user_id=info.context.user_id, resource_name=resource_name
                )
            )
            .result()
            .upgraded
        )


# @strawberry.type
# class Subscription:
#     @strawberry.subscription
#     async def user_resources(
#         self, info
#     ) -> typing.AsyncGenerator[list[UserResourceData], None]:
#         x = 0
#         while True:
#             yield [
#                 UserResourceData(
#                     resource_name="asd",
#                     factory_level=info.context.user_id,
#                     amount=x,
#                     time_until_upgrade_complete=None,
#                 )
#             ]
#             x += 1
#             await asyncio.sleep(1)


class ControllerGraphQLView(GraphQLView):
    async def get_context(self, request: request.Request) -> dict:
        return request.ctx


graphql_blueprint.add_route(
    ControllerGraphQLView.as_view(
        schema=strawberry.Schema(
            query=Query,
            mutation=Mutation,
            # subscription=Subscription,
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
