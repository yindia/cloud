package route

import (
	"context"
	v1 "task/pkg/gen/cloud/v1"

	"google.golang.org/protobuf/types/known/emptypb"
)

// TaskManagementHandler defines the methods for task management operations.
//

//go:generate mockery --output=./mocks --case=underscore --all --with-expecter
type TaskManagementHandler interface {
	CreateTask(ctx context.Context, req *v1.CreateTaskRequest) (*v1.CreateTaskResponse, error)
	GetTask(ctx context.Context, req *v1.GetTaskRequest) (*v1.Task, error)
	GetTaskHistory(ctx context.Context, req *v1.GetTaskHistoryRequest) (*v1.GetTaskHistoryResponse, error)
	UpdateTaskStatus(ctx context.Context, req *v1.UpdateTaskStatusRequest) (*emptypb.Empty, error)
	ListTasks(ctx context.Context, req *v1.TaskListRequest) (*v1.TaskList, error) // Updated to match the proto definition
}
