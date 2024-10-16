// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: cloud/v1/cloud.proto

package cloudv1connect

import (
	connect "connectrpc.com/connect"
	context "context"
	errors "errors"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	http "net/http"
	strings "strings"
	v1 "task/pkg/gen/cloud/v1"
)

// This is a compile-time assertion to ensure that this generated file and the connect package are
// compatible. If you get a compiler error that this constant is not defined, this code was
// generated with a version of connect newer than the one compiled into your binary. You can fix the
// problem by either regenerating this code with an older version of connect or updating the connect
// version compiled into your binary.
const _ = connect.IsAtLeastVersion0_1_0

const (
	// TaskManagementServiceName is the fully-qualified name of the TaskManagementService service.
	TaskManagementServiceName = "cloud.v1.TaskManagementService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// TaskManagementServiceCreateTaskProcedure is the fully-qualified name of the
	// TaskManagementService's CreateTask RPC.
	TaskManagementServiceCreateTaskProcedure = "/cloud.v1.TaskManagementService/CreateTask"
	// TaskManagementServiceGetTaskProcedure is the fully-qualified name of the TaskManagementService's
	// GetTask RPC.
	TaskManagementServiceGetTaskProcedure = "/cloud.v1.TaskManagementService/GetTask"
	// TaskManagementServiceListTasksProcedure is the fully-qualified name of the
	// TaskManagementService's ListTasks RPC.
	TaskManagementServiceListTasksProcedure = "/cloud.v1.TaskManagementService/ListTasks"
	// TaskManagementServiceGetTaskHistoryProcedure is the fully-qualified name of the
	// TaskManagementService's GetTaskHistory RPC.
	TaskManagementServiceGetTaskHistoryProcedure = "/cloud.v1.TaskManagementService/GetTaskHistory"
	// TaskManagementServiceUpdateTaskStatusProcedure is the fully-qualified name of the
	// TaskManagementService's UpdateTaskStatus RPC.
	TaskManagementServiceUpdateTaskStatusProcedure = "/cloud.v1.TaskManagementService/UpdateTaskStatus"
	// TaskManagementServiceGetStatusProcedure is the fully-qualified name of the
	// TaskManagementService's GetStatus RPC.
	TaskManagementServiceGetStatusProcedure = "/cloud.v1.TaskManagementService/GetStatus"
	// TaskManagementServiceHeartbeatProcedure is the fully-qualified name of the
	// TaskManagementService's Heartbeat RPC.
	TaskManagementServiceHeartbeatProcedure = "/cloud.v1.TaskManagementService/Heartbeat"
	// TaskManagementServicePullEventsProcedure is the fully-qualified name of the
	// TaskManagementService's PullEvents RPC.
	TaskManagementServicePullEventsProcedure = "/cloud.v1.TaskManagementService/PullEvents"
)

// TaskManagementServiceClient is a client for the cloud.v1.TaskManagementService service.
type TaskManagementServiceClient interface {
	// Creates a new task based on the provided request.
	// Returns a CreateTaskResponse containing the unique identifier of the created task.
	CreateTask(context.Context, *connect.Request[v1.CreateTaskRequest]) (*connect.Response[v1.CreateTaskResponse], error)
	// Retrieves the current status and details of the specified task.
	// Returns a Task message containing all information about the requested task.
	GetTask(context.Context, *connect.Request[v1.GetTaskRequest]) (*connect.Response[v1.Task], error)
	// Lists tasks currently available in the system, with pagination support.
	// Returns a TaskList containing the requested subset of tasks.
	ListTasks(context.Context, *connect.Request[v1.TaskListRequest]) (*connect.Response[v1.TaskList], error)
	// Retrieves the execution history of the specified task.
	// Returns a GetTaskHistoryResponse containing a list of historical status updates.
	GetTaskHistory(context.Context, *connect.Request[v1.GetTaskHistoryRequest]) (*connect.Response[v1.GetTaskHistoryResponse], error)
	// Updates the status of the specified task.
	// Returns an empty response to confirm the update was processed.
	UpdateTaskStatus(context.Context, *connect.Request[v1.UpdateTaskStatusRequest]) (*connect.Response[emptypb.Empty], error)
	// Retrieves the count of tasks for each status.
	// Returns a GetStatusResponse containing a map of status counts.
	GetStatus(context.Context, *connect.Request[v1.GetStatusRequest]) (*connect.Response[v1.GetStatusResponse], error)
	Heartbeat(context.Context, *connect.Request[v1.HeartbeatRequest]) (*connect.Response[v1.HeartbeatResponse], error)
	PullEvents(context.Context, *connect.Request[v1.PullEventsRequest]) (*connect.ServerStreamForClient[v1.PullEventsResponse], error)
}

// NewTaskManagementServiceClient constructs a client for the cloud.v1.TaskManagementService
// service. By default, it uses the Connect protocol with the binary Protobuf Codec, asks for
// gzipped responses, and sends uncompressed requests. To use the gRPC or gRPC-Web protocols, supply
// the connect.WithGRPC() or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewTaskManagementServiceClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) TaskManagementServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &taskManagementServiceClient{
		createTask: connect.NewClient[v1.CreateTaskRequest, v1.CreateTaskResponse](
			httpClient,
			baseURL+TaskManagementServiceCreateTaskProcedure,
			opts...,
		),
		getTask: connect.NewClient[v1.GetTaskRequest, v1.Task](
			httpClient,
			baseURL+TaskManagementServiceGetTaskProcedure,
			opts...,
		),
		listTasks: connect.NewClient[v1.TaskListRequest, v1.TaskList](
			httpClient,
			baseURL+TaskManagementServiceListTasksProcedure,
			opts...,
		),
		getTaskHistory: connect.NewClient[v1.GetTaskHistoryRequest, v1.GetTaskHistoryResponse](
			httpClient,
			baseURL+TaskManagementServiceGetTaskHistoryProcedure,
			opts...,
		),
		updateTaskStatus: connect.NewClient[v1.UpdateTaskStatusRequest, emptypb.Empty](
			httpClient,
			baseURL+TaskManagementServiceUpdateTaskStatusProcedure,
			opts...,
		),
		getStatus: connect.NewClient[v1.GetStatusRequest, v1.GetStatusResponse](
			httpClient,
			baseURL+TaskManagementServiceGetStatusProcedure,
			opts...,
		),
		heartbeat: connect.NewClient[v1.HeartbeatRequest, v1.HeartbeatResponse](
			httpClient,
			baseURL+TaskManagementServiceHeartbeatProcedure,
			opts...,
		),
		pullEvents: connect.NewClient[v1.PullEventsRequest, v1.PullEventsResponse](
			httpClient,
			baseURL+TaskManagementServicePullEventsProcedure,
			opts...,
		),
	}
}

// taskManagementServiceClient implements TaskManagementServiceClient.
type taskManagementServiceClient struct {
	createTask       *connect.Client[v1.CreateTaskRequest, v1.CreateTaskResponse]
	getTask          *connect.Client[v1.GetTaskRequest, v1.Task]
	listTasks        *connect.Client[v1.TaskListRequest, v1.TaskList]
	getTaskHistory   *connect.Client[v1.GetTaskHistoryRequest, v1.GetTaskHistoryResponse]
	updateTaskStatus *connect.Client[v1.UpdateTaskStatusRequest, emptypb.Empty]
	getStatus        *connect.Client[v1.GetStatusRequest, v1.GetStatusResponse]
	heartbeat        *connect.Client[v1.HeartbeatRequest, v1.HeartbeatResponse]
	pullEvents       *connect.Client[v1.PullEventsRequest, v1.PullEventsResponse]
}

// CreateTask calls cloud.v1.TaskManagementService.CreateTask.
func (c *taskManagementServiceClient) CreateTask(ctx context.Context, req *connect.Request[v1.CreateTaskRequest]) (*connect.Response[v1.CreateTaskResponse], error) {
	return c.createTask.CallUnary(ctx, req)
}

// GetTask calls cloud.v1.TaskManagementService.GetTask.
func (c *taskManagementServiceClient) GetTask(ctx context.Context, req *connect.Request[v1.GetTaskRequest]) (*connect.Response[v1.Task], error) {
	return c.getTask.CallUnary(ctx, req)
}

// ListTasks calls cloud.v1.TaskManagementService.ListTasks.
func (c *taskManagementServiceClient) ListTasks(ctx context.Context, req *connect.Request[v1.TaskListRequest]) (*connect.Response[v1.TaskList], error) {
	return c.listTasks.CallUnary(ctx, req)
}

// GetTaskHistory calls cloud.v1.TaskManagementService.GetTaskHistory.
func (c *taskManagementServiceClient) GetTaskHistory(ctx context.Context, req *connect.Request[v1.GetTaskHistoryRequest]) (*connect.Response[v1.GetTaskHistoryResponse], error) {
	return c.getTaskHistory.CallUnary(ctx, req)
}

// UpdateTaskStatus calls cloud.v1.TaskManagementService.UpdateTaskStatus.
func (c *taskManagementServiceClient) UpdateTaskStatus(ctx context.Context, req *connect.Request[v1.UpdateTaskStatusRequest]) (*connect.Response[emptypb.Empty], error) {
	return c.updateTaskStatus.CallUnary(ctx, req)
}

// GetStatus calls cloud.v1.TaskManagementService.GetStatus.
func (c *taskManagementServiceClient) GetStatus(ctx context.Context, req *connect.Request[v1.GetStatusRequest]) (*connect.Response[v1.GetStatusResponse], error) {
	return c.getStatus.CallUnary(ctx, req)
}

// Heartbeat calls cloud.v1.TaskManagementService.Heartbeat.
func (c *taskManagementServiceClient) Heartbeat(ctx context.Context, req *connect.Request[v1.HeartbeatRequest]) (*connect.Response[v1.HeartbeatResponse], error) {
	return c.heartbeat.CallUnary(ctx, req)
}

// PullEvents calls cloud.v1.TaskManagementService.PullEvents.
func (c *taskManagementServiceClient) PullEvents(ctx context.Context, req *connect.Request[v1.PullEventsRequest]) (*connect.ServerStreamForClient[v1.PullEventsResponse], error) {
	return c.pullEvents.CallServerStream(ctx, req)
}

// TaskManagementServiceHandler is an implementation of the cloud.v1.TaskManagementService service.
type TaskManagementServiceHandler interface {
	// Creates a new task based on the provided request.
	// Returns a CreateTaskResponse containing the unique identifier of the created task.
	CreateTask(context.Context, *connect.Request[v1.CreateTaskRequest]) (*connect.Response[v1.CreateTaskResponse], error)
	// Retrieves the current status and details of the specified task.
	// Returns a Task message containing all information about the requested task.
	GetTask(context.Context, *connect.Request[v1.GetTaskRequest]) (*connect.Response[v1.Task], error)
	// Lists tasks currently available in the system, with pagination support.
	// Returns a TaskList containing the requested subset of tasks.
	ListTasks(context.Context, *connect.Request[v1.TaskListRequest]) (*connect.Response[v1.TaskList], error)
	// Retrieves the execution history of the specified task.
	// Returns a GetTaskHistoryResponse containing a list of historical status updates.
	GetTaskHistory(context.Context, *connect.Request[v1.GetTaskHistoryRequest]) (*connect.Response[v1.GetTaskHistoryResponse], error)
	// Updates the status of the specified task.
	// Returns an empty response to confirm the update was processed.
	UpdateTaskStatus(context.Context, *connect.Request[v1.UpdateTaskStatusRequest]) (*connect.Response[emptypb.Empty], error)
	// Retrieves the count of tasks for each status.
	// Returns a GetStatusResponse containing a map of status counts.
	GetStatus(context.Context, *connect.Request[v1.GetStatusRequest]) (*connect.Response[v1.GetStatusResponse], error)
	Heartbeat(context.Context, *connect.Request[v1.HeartbeatRequest]) (*connect.Response[v1.HeartbeatResponse], error)
	PullEvents(context.Context, *connect.Request[v1.PullEventsRequest], *connect.ServerStream[v1.PullEventsResponse]) error
}

// NewTaskManagementServiceHandler builds an HTTP handler from the service implementation. It
// returns the path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewTaskManagementServiceHandler(svc TaskManagementServiceHandler, opts ...connect.HandlerOption) (string, http.Handler) {
	taskManagementServiceCreateTaskHandler := connect.NewUnaryHandler(
		TaskManagementServiceCreateTaskProcedure,
		svc.CreateTask,
		opts...,
	)
	taskManagementServiceGetTaskHandler := connect.NewUnaryHandler(
		TaskManagementServiceGetTaskProcedure,
		svc.GetTask,
		opts...,
	)
	taskManagementServiceListTasksHandler := connect.NewUnaryHandler(
		TaskManagementServiceListTasksProcedure,
		svc.ListTasks,
		opts...,
	)
	taskManagementServiceGetTaskHistoryHandler := connect.NewUnaryHandler(
		TaskManagementServiceGetTaskHistoryProcedure,
		svc.GetTaskHistory,
		opts...,
	)
	taskManagementServiceUpdateTaskStatusHandler := connect.NewUnaryHandler(
		TaskManagementServiceUpdateTaskStatusProcedure,
		svc.UpdateTaskStatus,
		opts...,
	)
	taskManagementServiceGetStatusHandler := connect.NewUnaryHandler(
		TaskManagementServiceGetStatusProcedure,
		svc.GetStatus,
		opts...,
	)
	taskManagementServiceHeartbeatHandler := connect.NewUnaryHandler(
		TaskManagementServiceHeartbeatProcedure,
		svc.Heartbeat,
		opts...,
	)
	taskManagementServicePullEventsHandler := connect.NewServerStreamHandler(
		TaskManagementServicePullEventsProcedure,
		svc.PullEvents,
		opts...,
	)
	return "/cloud.v1.TaskManagementService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case TaskManagementServiceCreateTaskProcedure:
			taskManagementServiceCreateTaskHandler.ServeHTTP(w, r)
		case TaskManagementServiceGetTaskProcedure:
			taskManagementServiceGetTaskHandler.ServeHTTP(w, r)
		case TaskManagementServiceListTasksProcedure:
			taskManagementServiceListTasksHandler.ServeHTTP(w, r)
		case TaskManagementServiceGetTaskHistoryProcedure:
			taskManagementServiceGetTaskHistoryHandler.ServeHTTP(w, r)
		case TaskManagementServiceUpdateTaskStatusProcedure:
			taskManagementServiceUpdateTaskStatusHandler.ServeHTTP(w, r)
		case TaskManagementServiceGetStatusProcedure:
			taskManagementServiceGetStatusHandler.ServeHTTP(w, r)
		case TaskManagementServiceHeartbeatProcedure:
			taskManagementServiceHeartbeatHandler.ServeHTTP(w, r)
		case TaskManagementServicePullEventsProcedure:
			taskManagementServicePullEventsHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedTaskManagementServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedTaskManagementServiceHandler struct{}

func (UnimplementedTaskManagementServiceHandler) CreateTask(context.Context, *connect.Request[v1.CreateTaskRequest]) (*connect.Response[v1.CreateTaskResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloud.v1.TaskManagementService.CreateTask is not implemented"))
}

func (UnimplementedTaskManagementServiceHandler) GetTask(context.Context, *connect.Request[v1.GetTaskRequest]) (*connect.Response[v1.Task], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloud.v1.TaskManagementService.GetTask is not implemented"))
}

func (UnimplementedTaskManagementServiceHandler) ListTasks(context.Context, *connect.Request[v1.TaskListRequest]) (*connect.Response[v1.TaskList], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloud.v1.TaskManagementService.ListTasks is not implemented"))
}

func (UnimplementedTaskManagementServiceHandler) GetTaskHistory(context.Context, *connect.Request[v1.GetTaskHistoryRequest]) (*connect.Response[v1.GetTaskHistoryResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloud.v1.TaskManagementService.GetTaskHistory is not implemented"))
}

func (UnimplementedTaskManagementServiceHandler) UpdateTaskStatus(context.Context, *connect.Request[v1.UpdateTaskStatusRequest]) (*connect.Response[emptypb.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloud.v1.TaskManagementService.UpdateTaskStatus is not implemented"))
}

func (UnimplementedTaskManagementServiceHandler) GetStatus(context.Context, *connect.Request[v1.GetStatusRequest]) (*connect.Response[v1.GetStatusResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloud.v1.TaskManagementService.GetStatus is not implemented"))
}

func (UnimplementedTaskManagementServiceHandler) Heartbeat(context.Context, *connect.Request[v1.HeartbeatRequest]) (*connect.Response[v1.HeartbeatResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloud.v1.TaskManagementService.Heartbeat is not implemented"))
}

func (UnimplementedTaskManagementServiceHandler) PullEvents(context.Context, *connect.Request[v1.PullEventsRequest], *connect.ServerStream[v1.PullEventsResponse]) error {
	return connect.NewError(connect.CodeUnimplemented, errors.New("cloud.v1.TaskManagementService.PullEvents is not implemented"))
}
