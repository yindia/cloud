package route

import (
	"context"
	"errors"
	"log"
	"os"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/emptypb"

	cloudv1 "task/pkg/gen/cloud/v1"
	"task/pkg/plugins/email"
	"task/server/repository/model/task"
	"task/server/route/mocks"
)

func TestCreateTask(t *testing.T) {
	mockHandler := mocks.NewTaskManagementHandler(t)

	req := &cloudv1.CreateTaskRequest{
		Name: "Test Task",
		Type: email.PLUGIN_NAME,
		Payload: &cloudv1.Payload{
			Parameters: map[string]string{"key": "value"},
		},
	}

	expectedResp := &cloudv1.CreateTaskResponse{Id: 1}

	mockHandler.EXPECT().CreateTask(mock.Anything, req).Return(expectedResp, nil)

	resp, err := mockHandler.CreateTask(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, expectedResp, resp)
}

func TestGetTask(t *testing.T) {
	mockHandler := mocks.NewTaskManagementHandler(t)

	t.Run("Successful GetTask", func(t *testing.T) {
		req := &cloudv1.GetTaskRequest{Id: 1}
		expectedResp := &cloudv1.Task{
			Id:          1,
			Name:        "Test Task",
			Type:        email.PLUGIN_NAME,
			Status:      cloudv1.TaskStatusEnum_QUEUED,
			Description: "Test task description",
			CreatedAt:   time.Now().Format(time.RFC3339),
		}

		mockHandler.EXPECT().GetTask(mock.Anything, req).Return(expectedResp, nil)

		resp, err := mockHandler.GetTask(context.Background(), req)

		assert.NoError(t, err)
		assert.True(t, cmp.Equal(expectedResp, resp, protocmp.Transform()))
	})

	t.Run("GetTask Not Found", func(t *testing.T) {
		req := &cloudv1.GetTaskRequest{Id: 999}
		mockHandler.EXPECT().GetTask(mock.Anything, req).Return(nil, errors.New("task not found"))

		resp, err := mockHandler.GetTask(context.Background(), req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "task not found")
	})
}

func TestListTasks(t *testing.T) {
	mockHandler := mocks.NewTaskManagementHandler(t)

	t.Run("Successful ListTasks", func(t *testing.T) {
		req := &cloudv1.TaskListRequest{}
		expectedResp := &cloudv1.TaskList{
			Tasks: []*cloudv1.Task{
				{Id: 1, Name: "Task 1"},
				{Id: 2, Name: "Task 2"},
			},
		}

		mockHandler.EXPECT().ListTasks(mock.Anything, req).Return(expectedResp, nil)

		resp, err := mockHandler.ListTasks(context.Background(), req)

		assert.NoError(t, err)
		assert.Equal(t, expectedResp, resp)
		assert.Len(t, resp.Tasks, 2)
	})

}

func TestGetTaskHistory(t *testing.T) {
	mockHandler := mocks.NewTaskManagementHandler(t)

	t.Run("Successful GetTaskHistory", func(t *testing.T) {
		req := &cloudv1.GetTaskHistoryRequest{Id: 1}
		now := time.Now()
		expectedResp := &cloudv1.GetTaskHistoryResponse{
			History: []*cloudv1.TaskHistory{
				{Id: 1, Details: "Created", Status: cloudv1.TaskStatusEnum_QUEUED, CreatedAt: now.Add(-2 * time.Hour).Format(time.RFC3339)},
				{Id: 2, Details: "Started", Status: cloudv1.TaskStatusEnum_RUNNING, CreatedAt: now.Add(-1 * time.Hour).Format(time.RFC3339)},
				{Id: 3, Details: "Completed", Status: cloudv1.TaskStatusEnum_SUCCEEDED, CreatedAt: now.Format(time.RFC3339)},
			},
		}

		mockHandler.EXPECT().GetTaskHistory(mock.Anything, req).Return(expectedResp, nil)

		resp, err := mockHandler.GetTaskHistory(context.Background(), req)

		assert.NoError(t, err)
		assert.True(t, cmp.Equal(expectedResp, resp, protocmp.Transform()))
		assert.Len(t, resp.History, 3)
	})

	t.Run("Empty GetTaskHistory", func(t *testing.T) {
		req := &cloudv1.GetTaskHistoryRequest{Id: 2}
		expectedResp := &cloudv1.GetTaskHistoryResponse{History: []*cloudv1.TaskHistory{}}

		mockHandler.EXPECT().GetTaskHistory(mock.Anything, req).Return(expectedResp, nil)

		resp, err := mockHandler.GetTaskHistory(context.Background(), req)

		assert.NoError(t, err)
		assert.Equal(t, expectedResp, resp)
		assert.Empty(t, resp.History)
	})

	t.Run("GetTaskHistory Error", func(t *testing.T) {
		req := &cloudv1.GetTaskHistoryRequest{Id: 999}
		mockHandler.EXPECT().GetTaskHistory(mock.Anything, req).Return(nil, errors.New("task not found"))

		resp, err := mockHandler.GetTaskHistory(context.Background(), req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "task not found")
	})
}

func TestUpdateTaskStatus(t *testing.T) {
	mockHandler := mocks.NewTaskManagementHandler(t)

	t.Run("Successful UpdateTaskStatus", func(t *testing.T) {
		req := &cloudv1.UpdateTaskStatusRequest{
			Id:     1,
			Status: cloudv1.TaskStatusEnum_RUNNING,
		}
		expectedResp := &emptypb.Empty{}

		mockHandler.EXPECT().UpdateTaskStatus(mock.Anything, req).Return(expectedResp, nil)

		resp, err := mockHandler.UpdateTaskStatus(context.Background(), req)

		assert.NoError(t, err)
		assert.Equal(t, expectedResp, resp)
	})

	t.Run("UpdateTaskStatus Not Found", func(t *testing.T) {
		req := &cloudv1.UpdateTaskStatusRequest{
			Id:     999,
			Status: cloudv1.TaskStatusEnum_SUCCEEDED,
		}
		mockHandler.EXPECT().UpdateTaskStatus(mock.Anything, req).Return(nil, errors.New("task not found"))

		resp, err := mockHandler.UpdateTaskStatus(context.Background(), req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "task not found")
	})

	t.Run("UpdateTaskStatus Invalid Status", func(t *testing.T) {
		req := &cloudv1.UpdateTaskStatusRequest{
			Id:     1,
			Status: cloudv1.TaskStatusEnum_UNKNOWN,
		}
		mockHandler.EXPECT().UpdateTaskStatus(mock.Anything, req).Return(nil, errors.New("invalid status"))

		resp, err := mockHandler.UpdateTaskStatus(context.Background(), req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "invalid status")
	})

	t.Run("UpdateTaskStatus Transition Error", func(t *testing.T) {
		req := &cloudv1.UpdateTaskStatusRequest{
			Id:     1,
			Status: cloudv1.TaskStatusEnum_QUEUED,
		}
		mockHandler.EXPECT().UpdateTaskStatus(mock.Anything, req).Return(nil, errors.New("invalid status transition"))

		resp, err := mockHandler.UpdateTaskStatus(context.Background(), req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "invalid status transition")
	})
}

func TestConvertTaskToProto(t *testing.T) {
	mockServer := &TaskServer{
		logger: log.New(os.Stdout, "TestConvertTaskToProto: ", log.LstdFlags),
	}

	t.Run("Successful conversion", func(t *testing.T) {
		taskModel := &task.Task{

			Name:        "Test Task",
			Description: "Test Description",
			Status:      int(cloudv1.TaskStatusEnum_RUNNING),
			Priority:    2,
			Retries:     3,
			Payload:     `{"key":"value"}`,
			Type:        "email",
		}

		protoTask := mockServer.convertTaskToProto(taskModel)

		assert.Equal(t, int32(taskModel.ID), protoTask.Id)
		assert.Equal(t, taskModel.Name, protoTask.Name)
		assert.Equal(t, taskModel.Description, protoTask.Description)
		assert.Equal(t, cloudv1.TaskStatusEnum(taskModel.Status), protoTask.Status)
		assert.Equal(t, int32(taskModel.Priority), protoTask.Priority)
		assert.Equal(t, int32(taskModel.Retries), protoTask.Retries)
		assert.Equal(t, taskModel.Type, protoTask.Type)
		assert.Equal(t, map[string]string{"key": "value"}, protoTask.Payload.Parameters)
	})

	t.Run("Invalid JSON payload", func(t *testing.T) {
		taskModel := &task.Task{

			Name:    "Invalid Payload Task",
			Payload: `{"key": invalid}`,
		}

		protoTask := mockServer.convertTaskToProto(taskModel)

		assert.Equal(t, int32(taskModel.ID), protoTask.Id)
		assert.Equal(t, taskModel.Name, protoTask.Name)
		assert.Nil(t, protoTask.Payload.Parameters)
	})
}
