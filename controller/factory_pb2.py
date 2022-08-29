# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: factory.proto
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import message as _message
from google.protobuf import reflection as _reflection
from google.protobuf import symbol_database as _symbol_database

# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(
    b'\n\rfactory.proto"\x98\x01\n\x0b\x46\x61\x63toryData\x12\x15\n\rresource_name\x18\x01 \x01(\t\x12\x15\n\rfactory_level\x18\x02 \x01(\x05\x12\x1d\n\x15production_per_second\x18\x03 \x01(\x05\x12"\n\x15next_upgrade_duration\x18\x04 \x01(\x05H\x00\x88\x01\x01\x42\x18\n\x16_next_upgrade_duration"\x17\n\x15GetFactoryDataRequest"=\n\x16GetFactoryDataResponse\x12#\n\rfactory_datas\x18\x01 \x03(\x0b\x32\x0c.FactoryData"$\n\x11\x43reateUserRequest\x12\x0f\n\x07user_id\x18\x01 \x01(\x05"%\n\x12\x43reateUserResponse\x12\x0f\n\x07\x63reated\x18\x01 \x01(\x08"?\n\x15UpgradeFactoryRequest\x12\x0f\n\x07user_id\x18\x01 \x01(\x05\x12\x15\n\rresource_name\x18\x02 \x01(\t"*\n\x16UpgradeFactoryResponse\x12\x10\n\x08upgraded\x18\x01 \x01(\x08\x32\xcc\x01\n\x07\x46\x61\x63tory\x12\x43\n\x0eGetFactoryData\x12\x16.GetFactoryDataRequest\x1a\x17.GetFactoryDataResponse"\x00\x12\x37\n\nCreateUser\x12\x12.CreateUserRequest\x1a\x13.CreateUserResponse"\x00\x12\x43\n\x0eUpgradeFactory\x12\x16.UpgradeFactoryRequest\x1a\x17.UpgradeFactoryResponse"\x00\x42\x06Z\x04./pbb\x06proto3'
)


_FACTORYDATA = DESCRIPTOR.message_types_by_name["FactoryData"]
_GETFACTORYDATAREQUEST = DESCRIPTOR.message_types_by_name["GetFactoryDataRequest"]
_GETFACTORYDATARESPONSE = DESCRIPTOR.message_types_by_name["GetFactoryDataResponse"]
_CREATEUSERREQUEST = DESCRIPTOR.message_types_by_name["CreateUserRequest"]
_CREATEUSERRESPONSE = DESCRIPTOR.message_types_by_name["CreateUserResponse"]
_UPGRADEFACTORYREQUEST = DESCRIPTOR.message_types_by_name["UpgradeFactoryRequest"]
_UPGRADEFACTORYRESPONSE = DESCRIPTOR.message_types_by_name["UpgradeFactoryResponse"]
FactoryData = _reflection.GeneratedProtocolMessageType(
    "FactoryData",
    (_message.Message,),
    {
        "DESCRIPTOR": _FACTORYDATA,
        "__module__": "factory_pb2"
        # @@protoc_insertion_point(class_scope:FactoryData)
    },
)
_sym_db.RegisterMessage(FactoryData)

GetFactoryDataRequest = _reflection.GeneratedProtocolMessageType(
    "GetFactoryDataRequest",
    (_message.Message,),
    {
        "DESCRIPTOR": _GETFACTORYDATAREQUEST,
        "__module__": "factory_pb2"
        # @@protoc_insertion_point(class_scope:GetFactoryDataRequest)
    },
)
_sym_db.RegisterMessage(GetFactoryDataRequest)

GetFactoryDataResponse = _reflection.GeneratedProtocolMessageType(
    "GetFactoryDataResponse",
    (_message.Message,),
    {
        "DESCRIPTOR": _GETFACTORYDATARESPONSE,
        "__module__": "factory_pb2"
        # @@protoc_insertion_point(class_scope:GetFactoryDataResponse)
    },
)
_sym_db.RegisterMessage(GetFactoryDataResponse)

CreateUserRequest = _reflection.GeneratedProtocolMessageType(
    "CreateUserRequest",
    (_message.Message,),
    {
        "DESCRIPTOR": _CREATEUSERREQUEST,
        "__module__": "factory_pb2"
        # @@protoc_insertion_point(class_scope:CreateUserRequest)
    },
)
_sym_db.RegisterMessage(CreateUserRequest)

CreateUserResponse = _reflection.GeneratedProtocolMessageType(
    "CreateUserResponse",
    (_message.Message,),
    {
        "DESCRIPTOR": _CREATEUSERRESPONSE,
        "__module__": "factory_pb2"
        # @@protoc_insertion_point(class_scope:CreateUserResponse)
    },
)
_sym_db.RegisterMessage(CreateUserResponse)

UpgradeFactoryRequest = _reflection.GeneratedProtocolMessageType(
    "UpgradeFactoryRequest",
    (_message.Message,),
    {
        "DESCRIPTOR": _UPGRADEFACTORYREQUEST,
        "__module__": "factory_pb2"
        # @@protoc_insertion_point(class_scope:UpgradeFactoryRequest)
    },
)
_sym_db.RegisterMessage(UpgradeFactoryRequest)

UpgradeFactoryResponse = _reflection.GeneratedProtocolMessageType(
    "UpgradeFactoryResponse",
    (_message.Message,),
    {
        "DESCRIPTOR": _UPGRADEFACTORYRESPONSE,
        "__module__": "factory_pb2"
        # @@protoc_insertion_point(class_scope:UpgradeFactoryResponse)
    },
)
_sym_db.RegisterMessage(UpgradeFactoryResponse)

_FACTORY = DESCRIPTOR.services_by_name["Factory"]
if _descriptor._USE_C_DESCRIPTORS == False:

    DESCRIPTOR._options = None
    DESCRIPTOR._serialized_options = b"Z\004./pb"
    _FACTORYDATA._serialized_start = 18
    _FACTORYDATA._serialized_end = 170
    _GETFACTORYDATAREQUEST._serialized_start = 172
    _GETFACTORYDATAREQUEST._serialized_end = 195
    _GETFACTORYDATARESPONSE._serialized_start = 197
    _GETFACTORYDATARESPONSE._serialized_end = 258
    _CREATEUSERREQUEST._serialized_start = 260
    _CREATEUSERREQUEST._serialized_end = 296
    _CREATEUSERRESPONSE._serialized_start = 298
    _CREATEUSERRESPONSE._serialized_end = 335
    _UPGRADEFACTORYREQUEST._serialized_start = 337
    _UPGRADEFACTORYREQUEST._serialized_end = 400
    _UPGRADEFACTORYRESPONSE._serialized_start = 402
    _UPGRADEFACTORYRESPONSE._serialized_end = 444
    _FACTORY._serialized_start = 447
    _FACTORY._serialized_end = 651
# @@protoc_insertion_point(module_scope)