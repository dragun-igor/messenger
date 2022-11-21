// Code generated by MockGen. DO NOT EDIT.
// Source: internal/server/resources/interface.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	model "github.com/dragun-igor/messenger/internal/server/model"
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
func (m *MockRepository) CheckLoginExists(ctx context.Context, user model.User) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckLoginExists", ctx, user)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckLoginExists indicates an expected call of CheckLoginExists.
func (mr *MockRepositoryMockRecorder) CheckLoginExists(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckLoginExists", reflect.TypeOf((*MockRepository)(nil).CheckLoginExists), ctx, user)
}

// CheckNameExists mocks base method.
func (m *MockRepository) CheckNameExists(ctx context.Context, user model.User) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckNameExists", ctx, user)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckNameExists indicates an expected call of CheckNameExists.
func (mr *MockRepositoryMockRecorder) CheckNameExists(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckNameExists", reflect.TypeOf((*MockRepository)(nil).CheckNameExists), ctx, user)
}

// Close mocks base method.
func (m *MockRepository) Close(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockRepositoryMockRecorder) Close(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockRepository)(nil).Close), ctx)
}

// CreateUser mocks base method.
func (m *MockRepository) CreateUser(ctx context.Context, user model.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", ctx, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockRepositoryMockRecorder) CreateUser(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockRepository)(nil).CreateUser), ctx, user)
}

// InsertMessage mocks base method.
func (m *MockRepository) InsertMessage(ctx context.Context, message model.Message) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertMessage", ctx, message)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertMessage indicates an expected call of InsertMessage.
func (mr *MockRepositoryMockRecorder) InsertMessage(ctx, message interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertMessage", reflect.TypeOf((*MockRepository)(nil).InsertMessage), ctx, message)
}

// LogIn mocks base method.
func (m *MockRepository) LogIn(ctx context.Context, user model.User) (string, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LogIn", ctx, user)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// LogIn indicates an expected call of LogIn.
func (mr *MockRepositoryMockRecorder) LogIn(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LogIn", reflect.TypeOf((*MockRepository)(nil).LogIn), ctx, user)
}
