
from typing import TypedDict, List
from enum import Enum
from functools import wraps
from modules.cloud.v1.cloud_pb2 import Task

class TaskType(Enum):
    PYTHON = "python"

class TaskMetadata(TypedDict):
    author: str
    version: str

class TaskSpec(TypedDict):
    name: str
    type: TaskType
    description: str
    dependencies: List[str]
    metadata: TaskMetadata
    base_image: str
    entrypoint: str
