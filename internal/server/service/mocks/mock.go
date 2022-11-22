// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/dragun-igor/messenger/internal/server/service (interfaces: Repository)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	model "github.com/dragun-igor/messenger/internal/pkg/model"
	gomock "github.com/golang/mock/gomock"
)

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// CheckLoginExists mocks base method.
func (m *MockRepository) CheckLoginExists(arg0 context.Context, arg1 model.AuthData) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckLoginExists", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckLoginExists indicates an expected call of CheckLoginExists.
func (mr *MockRepositoryMockRecorder) CheckLoginExists(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckLoginExists", reflect.TypeOf((*MockRepository)(nil).CheckLoginExists), arg0, arg1)
}

// CheckNameExists mocks base method.
func (m *MockRepository) CheckNameExists(arg0 context.Context, arg1 model.AuthData) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckNameExists", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckNameExists indicates an expected call of CheckNameExists.
func (mr *MockRepositoryMockRecorder) CheckNameExists(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckNameExists", reflect.TypeOf((*MockRepository)(nil).CheckNameExists), arg0, arg1)
}

// CreateUser mocks base method.
func (m *MockRepository) CreateUser(arg0 context.Context, arg1 model.AuthData) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockRepositoryMockRecorder) CreateUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockRepository)(nil).CreateUser), arg0, arg1)
}

// GetUser mocks base method.
func (m *MockRepository) GetUser(arg0 context.Context, arg1 model.AuthData) (model.AuthData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", arg0, arg1)
	ret0, _ := ret[0].(model.AuthData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser.
func (mr *MockRepositoryMockRecorder) GetUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockRepository)(nil).GetUser), arg0, arg1)
}

// InsertMessage mocks base method.
func (m *MockRepository) InsertMessage(arg0 context.Context, arg1 model.Message) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertMessage", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertMessage indicates an expected call of InsertMessage.
func (mr *MockRepositoryMockRecorder) InsertMessage(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertMessage", reflect.TypeOf((*MockRepository)(nil).InsertMessage), arg0, arg1)
}
