package gormimpl

import (
	"context"
	"testing"

	"task/server/repository/mocks"
	"task/server/repository/model/task"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateTask(t *testing.T) {
	mockRepo := mocks.NewTaskRepo(t)

	newTask := task.Task{Name: "Test Task", Status: 2}
	mockRepo.EXPECT().CreateTask(mock.Anything, newTask).Return(newTask, nil)

	createdTask, err := mockRepo.CreateTask(context.Background(), newTask)

	assert.NoError(t, err)
	assert.Equal(t, newTask.Name, createdTask.Name)
}

func TestGetTaskByID(t *testing.T) {
	mockRepo := mocks.NewTaskRepo(t)

	taskToCreate := task.Task{Name: "Test Task", Status: 0}
	mockRepo.EXPECT().CreateTask(mock.Anything, taskToCreate).Return(taskToCreate, nil)
	mockRepo.EXPECT().GetTaskByID(mock.Anything, uint(1)).Return(&taskToCreate, nil)

	mockRepo.CreateTask(context.Background(), taskToCreate)

	retrievedTask, err := mockRepo.GetTaskByID(context.Background(), 1)

	assert.NoError(t, err)
	assert.Equal(t, taskToCreate.Name, retrievedTask.Name)
}

func TestUpdateTaskStatus(t *testing.T) {
	mockRepo := mocks.NewTaskRepo(t)

	taskToCreate := task.Task{Name: "Test Task", Status: 1}
	mockRepo.EXPECT().CreateTask(mock.Anything, taskToCreate).Return(taskToCreate, nil)
	mockRepo.EXPECT().UpdateTaskStatus(mock.Anything, uint(1), 3).Return(nil)
	mockRepo.EXPECT().GetTaskByID(mock.Anything, uint(1)).Return(&task.Task{Status: 3}, nil)

	mockRepo.CreateTask(context.Background(), taskToCreate)

	err := mockRepo.UpdateTaskStatus(context.Background(), 1, 3)

	assert.NoError(t, err)

	updatedTask, _ := mockRepo.GetTaskByID(context.Background(), 1)
	assert.Equal(t, 3, updatedTask.Status)
}
