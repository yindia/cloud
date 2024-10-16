
from typing import TypedDict, List, Dict
from enum import Enum
from functools import wraps
from modules.cloud.v1.cloud_pb2 import Workflow

class WorkflowType(Enum):
    PYTHON = "python"

class WorkflowMetadata(TypedDict):
    author: str
    version: str

class WorkflowSpec(TypedDict):
    name: str
    type: WorkflowType
    description: str
    dependencies: List[str]
    metadata: WorkflowMetadata
    base_image: str
    entrypoint: str
    args: List[str]
    env: Dict[str, str]
