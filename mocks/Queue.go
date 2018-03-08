// Code generated by mockery v1.0.0
package mocks

import mock "github.com/stretchr/testify/mock"
import model "github.com/wayofthepie/task-store/pkg/model"

import uuid "github.com/satori/go.uuid"

// Queue is an autogenerated mock type for the Queue type
type Queue struct {
	mock.Mock
}

// Push provides a mock function with given fields: spec
func (_m *Queue) Push(spec *model.Spec) (*uuid.UUID, error) {
	ret := _m.Called(spec)

	var r0 *uuid.UUID
	if rf, ok := ret.Get(0).(func(*model.Spec) *uuid.UUID); ok {
		r0 = rf(spec)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*uuid.UUID)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*model.Spec) error); ok {
		r1 = rf(spec)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
