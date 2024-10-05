package gormimpl

import (
	"context"
	"testing"

	"task/server/repository/mocks"
	"task/server/repository/model/task"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateTaskHistory(t *testing.T) {
	mockRepo := mocks.NewTaskHistoryRepo(t) // {{ edit_1 }}
	history := task.TaskHistory{TaskID: 1, Status: 3}

	mockRepo.EXPECT().CreateTaskHistory(mock.Anything, history).Return(history, nil) // {{ edit_2 }}

	result, err := mockRepo.CreateTaskHistory(context.Background(), history)

	assert.NoError(t, err)
	assert.Equal(t, history, result)
	mockRepo.AssertExpectations(t)
}

func TestGetTaskHistory(t *testing.T) {
	mockRepo := mocks.NewTaskHistoryRepo(t) // {{ edit_3 }}
	histories := []task.TaskHistory{{TaskID: 1, Status: 3}}

	mockRepo.EXPECT().GetTaskHistory(mock.Anything, uint(1)).Return(histories, nil) // {{ edit_4 }}

	result, err := mockRepo.GetTaskHistory(context.Background(), 1)

	assert.NoError(t, err)
	assert.Equal(t, histories, result)
	assert.Len(t, result, len(histories))
	mockRepo.AssertExpectations(t)
}
