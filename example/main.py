from typing import TypedDict, List
from enum import Enum
from pytask import task, TaskType, TaskSpec

@task(TaskSpec(
    name="hello_world",
    type=TaskType.PYTHON,
    description="This is a test task",
    dependencies="This is a test task",
    metadata={  
        "author": "Yuvraj Singh",
        "version": "1.0.0",
    },
    base_image="python:3.12",
    entrypoint="python",
    args=[],
    env={}
))
def hello_world(a : int, b : int) -> bool:
    print("Hello, World!")
    return True


@task(TaskSpec(
    name="hello_world",
    type=TaskType.PYTHON,
    description="This is a test task",
    dependencies="This is a test task",
    metadata={  
        "author": "Yuvraj Singh",
        "version": "1.0.0",
    },
    base_image="python:3.12",
    entrypoint="python",
    args=[],
    env={}
))
def hello_world_2(a : int, b : int) -> bool:
    print("Hello, World!")
    return True


