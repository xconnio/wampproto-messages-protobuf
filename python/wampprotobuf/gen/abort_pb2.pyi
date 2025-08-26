from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Optional as _Optional

DESCRIPTOR: _descriptor.FileDescriptor

class Abort(_message.Message):
    __slots__ = ("reason", "payload_serializer_id")
    REASON_FIELD_NUMBER: _ClassVar[int]
    PAYLOAD_SERIALIZER_ID_FIELD_NUMBER: _ClassVar[int]
    reason: str
    payload_serializer_id: int
    def __init__(self, reason: _Optional[str] = ..., payload_serializer_id: _Optional[int] = ...) -> None: ...
