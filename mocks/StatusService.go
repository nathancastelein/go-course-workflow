// Code generated by mockery v2.45.1. DO NOT EDIT.

package mocks

import (
	pokemon "github.com/nathancastelein/go-course-workflows/pokemon"
	mock "github.com/stretchr/testify/mock"
)

// StatusService is an autogenerated mock type for the StatusService type
type StatusService struct {
	mock.Mock
}

// Paralyze provides a mock function with given fields: _a0
func (_m *StatusService) Paralyze(_a0 *pokemon.Pokemon) error {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for Paralyze")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*pokemon.Pokemon) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewStatusService creates a new instance of StatusService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewStatusService(t interface {
	mock.TestingT
	Cleanup(func())
}) *StatusService {
	mock := &StatusService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}