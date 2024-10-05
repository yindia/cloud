// Code generated by mockery v2.46.0. DO NOT EDIT.

package mocks

import (
	interfaces "task/server/repository/interface"

	mock "github.com/stretchr/testify/mock"
)

// TaskManagmentInterface is an autogenerated mock type for the TaskManagmentInterface type
type TaskManagmentInterface struct {
	mock.Mock
}

type TaskManagmentInterface_Expecter struct {
	mock *mock.Mock
}

func (_m *TaskManagmentInterface) EXPECT() *TaskManagmentInterface_Expecter {
	return &TaskManagmentInterface_Expecter{mock: &_m.Mock}
}

// TaskHistoryRepo provides a mock function with given fields:
func (_m *TaskManagmentInterface) TaskHistoryRepo() interfaces.TaskHistoryRepo {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for TaskHistoryRepo")
	}

	var r0 interfaces.TaskHistoryRepo
	if rf, ok := ret.Get(0).(func() interfaces.TaskHistoryRepo); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interfaces.TaskHistoryRepo)
		}
	}

	return r0
}

// TaskManagmentInterface_TaskHistoryRepo_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'TaskHistoryRepo'
type TaskManagmentInterface_TaskHistoryRepo_Call struct {
	*mock.Call
}

// TaskHistoryRepo is a helper method to define mock.On call
func (_e *TaskManagmentInterface_Expecter) TaskHistoryRepo() *TaskManagmentInterface_TaskHistoryRepo_Call {
	return &TaskManagmentInterface_TaskHistoryRepo_Call{Call: _e.mock.On("TaskHistoryRepo")}
}

func (_c *TaskManagmentInterface_TaskHistoryRepo_Call) Run(run func()) *TaskManagmentInterface_TaskHistoryRepo_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *TaskManagmentInterface_TaskHistoryRepo_Call) Return(_a0 interfaces.TaskHistoryRepo) *TaskManagmentInterface_TaskHistoryRepo_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *TaskManagmentInterface_TaskHistoryRepo_Call) RunAndReturn(run func() interfaces.TaskHistoryRepo) *TaskManagmentInterface_TaskHistoryRepo_Call {
	_c.Call.Return(run)
	return _c
}

// TaskRepo provides a mock function with given fields:
func (_m *TaskManagmentInterface) TaskRepo() interfaces.TaskRepo {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for TaskRepo")
	}

	var r0 interfaces.TaskRepo
	if rf, ok := ret.Get(0).(func() interfaces.TaskRepo); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interfaces.TaskRepo)
		}
	}

	return r0
}

// TaskManagmentInterface_TaskRepo_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'TaskRepo'
type TaskManagmentInterface_TaskRepo_Call struct {
	*mock.Call
}

// TaskRepo is a helper method to define mock.On call
func (_e *TaskManagmentInterface_Expecter) TaskRepo() *TaskManagmentInterface_TaskRepo_Call {
	return &TaskManagmentInterface_TaskRepo_Call{Call: _e.mock.On("TaskRepo")}
}

func (_c *TaskManagmentInterface_TaskRepo_Call) Run(run func()) *TaskManagmentInterface_TaskRepo_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *TaskManagmentInterface_TaskRepo_Call) Return(_a0 interfaces.TaskRepo) *TaskManagmentInterface_TaskRepo_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *TaskManagmentInterface_TaskRepo_Call) RunAndReturn(run func() interfaces.TaskRepo) *TaskManagmentInterface_TaskRepo_Call {
	_c.Call.Return(run)
	return _c
}

// NewTaskManagmentInterface creates a new instance of TaskManagmentInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTaskManagmentInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *TaskManagmentInterface {
	mock := &TaskManagmentInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
