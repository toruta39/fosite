// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/toruta39/fosite/handler/pkce (interfaces: PKCERequestStorage)

// Package internal is a generated GoMock package.
package internal

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"

	fosite "github.com/toruta39/fosite"
)

// MockPKCERequestStorage is a mock of PKCERequestStorage interface
type MockPKCERequestStorage struct {
	ctrl     *gomock.Controller
	recorder *MockPKCERequestStorageMockRecorder
}

// MockPKCERequestStorageMockRecorder is the mock recorder for MockPKCERequestStorage
type MockPKCERequestStorageMockRecorder struct {
	mock *MockPKCERequestStorage
}

// NewMockPKCERequestStorage creates a new mock instance
func NewMockPKCERequestStorage(ctrl *gomock.Controller) *MockPKCERequestStorage {
	mock := &MockPKCERequestStorage{ctrl: ctrl}
	mock.recorder = &MockPKCERequestStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockPKCERequestStorage) EXPECT() *MockPKCERequestStorageMockRecorder {
	return m.recorder
}

// CreatePKCERequestSession mocks base method
func (m *MockPKCERequestStorage) CreatePKCERequestSession(arg0 context.Context, arg1 string, arg2 fosite.Requester) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreatePKCERequestSession", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreatePKCERequestSession indicates an expected call of CreatePKCERequestSession
func (mr *MockPKCERequestStorageMockRecorder) CreatePKCERequestSession(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePKCERequestSession", reflect.TypeOf((*MockPKCERequestStorage)(nil).CreatePKCERequestSession), arg0, arg1, arg2)
}

// DeletePKCERequestSession mocks base method
func (m *MockPKCERequestStorage) DeletePKCERequestSession(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeletePKCERequestSession", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeletePKCERequestSession indicates an expected call of DeletePKCERequestSession
func (mr *MockPKCERequestStorageMockRecorder) DeletePKCERequestSession(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeletePKCERequestSession", reflect.TypeOf((*MockPKCERequestStorage)(nil).DeletePKCERequestSession), arg0, arg1)
}

// GetPKCERequestSession mocks base method
func (m *MockPKCERequestStorage) GetPKCERequestSession(arg0 context.Context, arg1 string, arg2 fosite.Session) (fosite.Requester, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPKCERequestSession", arg0, arg1, arg2)
	ret0, _ := ret[0].(fosite.Requester)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPKCERequestSession indicates an expected call of GetPKCERequestSession
func (mr *MockPKCERequestStorageMockRecorder) GetPKCERequestSession(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPKCERequestSession", reflect.TypeOf((*MockPKCERequestStorage)(nil).GetPKCERequestSession), arg0, arg1, arg2)
}
