import click
import os
import ast
import logging
from typing import List, Dict, Union, Any
from cloud.v1.cloud_pb2 import Task  # Adjusted import
import tarfile
import tempfile
import json
import fnmatch

from typing import TypedDict, List, Dict
from enum import Enum


class TaskType(Enum):
    PYTHON = "python"


class TaskSpec(TypedDict):
    name: str
    type: TaskType
    description: str
    dependencies: List[str]
    metadata: Dict[str, str]
    base_image: str
    entrypoint: str
    args: List[str]
    env: Dict[str, str]


logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')
logger = logging.getLogger(__name__)


def parse_decorator_config(file_path: str) -> Dict[str, List[Dict[str, Union[str, Dict[str, Any]]]]]:
    """
    Parse Python files for task and workflow decorators and extract their configurations.

    Args:
        file_path (str): Path to the Python file to parse.

    Returns:
        Dict[str, List[Dict[str, Union[str, Dict]]]]: Parsed configurations for tasks and workflows.
    """
    with open(file_path, 'r') as file:
        tree = ast.parse(file.read())

    configs = {
        'tasks': []
    }

    def get_value(node: ast.AST) -> Any:
        """Recursively extract values from AST nodes."""
        if isinstance(node, ast.Constant):
            return node.value
        elif isinstance(node, ast.Name):
            return node.id
        elif isinstance(node, ast.Attribute):
            return f"{get_value(node.value)}.{node.attr}"
        elif isinstance(node, ast.Dict):
            return {get_value(k): get_value(v) for k, v in zip(node.keys, node.values)}
        elif isinstance(node, ast.List):
            return [get_value(elt) for elt in node.elts]
        elif isinstance(node, ast.Call):
            func = get_value(node.func)
            args = [get_value(arg) for arg in node.args]
            kwargs = {arg.arg: get_value(arg.value) for arg in node.keywords}
            return f"{func}({', '.join(map(str, args))}{', ' if args and kwargs else ''}{', '.join(f'{k}={v}' for k, v in kwargs.items())})"
        elif isinstance(node, ast.BoolOp):
            op = 'and' if isinstance(node.op, ast.And) else 'or'
            return f" {op} ".join(get_value(value) for value in node.values)
        else:
            return str(node)  # Fallback for other types

    def get_function_info(node: ast.FunctionDef) -> Dict[str, Any]:
        """Extract function argument and return type information."""
        args = []
        for arg in node.args.args:
            arg_info = {"name": arg.arg}
            if arg.annotation:
                arg_info["type"] = get_value(arg.annotation)
            args.append(arg_info)
        
        returns = None
        if node.returns:
            returns = get_value(node.returns)
        
        return {"args": args, "returns": returns}

    for node in ast.walk(tree):
        if isinstance(node, ast.FunctionDef):
            for decorator in node.decorator_list:
                if isinstance(decorator, ast.Call) and decorator.func.id in ['task', 'Workflow']:
                    config: Union[TaskSpec] = {}
                    for kw in decorator.args[0].keywords:
                        config[kw.arg] = get_value(kw.value)
                    
                    function_info = get_function_info(node)
                    
                    item = {
                        'name': node.name,
                        'config': config,
                        'input': function_info['args'],
                        'output': function_info['returns']
                    }
                    
                    if decorator.func.id == 'task':
                        configs['tasks'].append(item)

    return configs

@click.group(invoke_without_command=True)
@click.argument('dir', type=click.Path(exists=True), nargs=1, required=False, default='.')
@click.option('--image', '-i', help='Specify a custom image')
@click.option('--verbose', '-v', is_flag=True, help='Enable verbose output')
def cli(dir: str, image: str, verbose: bool) -> None:
    """
    Find Python files, read their content, and search for task and workflow decorators.

    Args:
        dir (str): Directory to search for Python files (default: current directory).
        verbose (bool): Enable verbose output.
    """
    if verbose:
        logger.setLevel(logging.DEBUG)
    
    logger.info(f"Searching for Python files in directory: {dir}")
    python_files: List[str] = find_python_files(dir)
    logger.info(f"Found {len(python_files)} Python files")
    
    for file_path in python_files:
        process_file(file_path, image)

def find_python_files(directory: str) -> List[str]:
    """
    Find all Python files in the given directory and its subdirectories.

    Args:
        directory (str): Directory to search for Python files.

    Returns:
        List[str]: List of paths to Python files.
    """
    python_files: List[str] = []
    for root, _, files in os.walk(directory):
        for file in files:
            if file.endswith('.py'):
                python_files.append(os.path.join(root, file))
    return python_files

def process_file(file_path: str, custom_image: str = None) -> None:
    """
    Read the file content, search for task and workflow decorators, and store protos in a .tgz file.

    Args:
        file_path (str): Path to the Python file to process.
        custom_image (str): Custom image to use instead of base_image if provided.
    """
    logger.info(f"Processing file: {file_path}")
    parsed_configs = parse_decorator_config(file_path)
    
    with tempfile.TemporaryDirectory() as temp_dir:
        task_protos = []
     
        if parsed_configs['tasks']:
            logger.info(f"Found {len(parsed_configs['tasks'])} tasks in {file_path}")
            for task in parsed_configs['tasks']:
                logger.debug(f"Task: {task['name']}")
                logger.debug(f"  Config: {task['config']}")
                logger.debug(f"  Input: {task['input']}")
                logger.debug(f"  Output: {task['output']}")
                
                # task_proto = Task(
                #     name=task['config']["name"],
                #     type=task['config']["type"],
                #     description=task['config']["description"],
                #     dependencies=task['config']["dependencies"],
                #     metadata=task['config']["metadata"],
                #     base_image=custom_image if custom_image else task['config']["base_image"],
                #     entrypoint=task['config']["entrypoint"],
                #     args=task['config']["args"],
                #     env=task['config']["env"]
                # )
                # task_protos.append(task_proto)
        

        
        # Save protos in binary format
        for i, task_proto in enumerate(task_protos):
            with open(os.path.join(temp_dir, f'task_{i}.pb'), 'wb') as f:
                f.write(task_proto.SerializeToString())
        
        # Compress the entire code
        code_tar_path = os.path.join(temp_dir, 'code.tar.gz')
        with tarfile.open(code_tar_path, "w:gz") as code_tar:
            for root, dirs, files in os.walk(os.path.dirname(file_path)):
                # Read .gitignore patterns
                gitignore_patterns = []
                gitignore_path = os.path.join(root, '.gitignore')
                if os.path.exists(gitignore_path):
                    with open(gitignore_path, 'r') as gitignore_file:
                        gitignore_patterns = gitignore_file.read().splitlines()
                
                for file in files:
                    file_path = os.path.join(root, file)
                    relative_path = os.path.relpath(file_path, os.path.dirname(file_path))
                    
                    # Check if the file should be ignored
                    if not any(fnmatch.fnmatch(relative_path, pattern) for pattern in gitignore_patterns):
                        code_tar.add(file_path, arcname=relative_path)
        
        # Create package.tgz file
        output_filename = 'package.tgz'
        with tarfile.open(output_filename, "w:gz") as tar:
            tar.add(temp_dir, arcname="")
    
    logger.info(f"Stored task and workflow protos in {output_filename}")

if __name__ == '__main__':
    cli()
