// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/toruta39/fosite (interfaces: AuthorizeResponder)

// Package internal is a generated GoMock package.
package internal

import (
	http "net/http"
	url "net/url"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockAuthorizeResponder is a mock of AuthorizeResponder interface
type MockAuthorizeResponder struct {
	ctrl     *gomock.Controller
	recorder *MockAuthorizeResponderMockRecorder
}

// MockAuthorizeResponderMockRecorder is the mock recorder for MockAuthorizeResponder
type MockAuthorizeResponderMockRecorder struct {
	mock *MockAuthorizeResponder
}

// NewMockAuthorizeResponder creates a new mock instance
func NewMockAuthorizeResponder(ctrl *gomock.Controller) *MockAuthorizeResponder {
	mock := &MockAuthorizeResponder{ctrl: ctrl}
	mock.recorder = &MockAuthorizeResponderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAuthorizeResponder) EXPECT() *MockAuthorizeResponderMockRecorder {
	return m.recorder
}

// AddFragment mocks base method
func (m *MockAuthorizeResponder) AddFragment(arg0, arg1 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AddFragment", arg0, arg1)
}

// AddFragment indicates an expected call of AddFragment
func (mr *MockAuthorizeResponderMockRecorder) AddFragment(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddFragment", reflect.TypeOf((*MockAuthorizeResponder)(nil).AddFragment), arg0, arg1)
}

// AddHeader mocks base method
func (m *MockAuthorizeResponder) AddHeader(arg0, arg1 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AddHeader", arg0, arg1)
}

// AddHeader indicates an expected call of AddHeader
func (mr *MockAuthorizeResponderMockRecorder) AddHeader(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddHeader", reflect.TypeOf((*MockAuthorizeResponder)(nil).AddHeader), arg0, arg1)
}

// AddQuery mocks base method
func (m *MockAuthorizeResponder) AddQuery(arg0, arg1 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AddQuery", arg0, arg1)
}

// AddQuery indicates an expected call of AddQuery
func (mr *MockAuthorizeResponderMockRecorder) AddQuery(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddQuery", reflect.TypeOf((*MockAuthorizeResponder)(nil).AddQuery), arg0, arg1)
}

// GetCode mocks base method
func (m *MockAuthorizeResponder) GetCode() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCode")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetCode indicates an expected call of GetCode
func (mr *MockAuthorizeResponderMockRecorder) GetCode() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCode", reflect.TypeOf((*MockAuthorizeResponder)(nil).GetCode))
}

// GetFragment mocks base method
func (m *MockAuthorizeResponder) GetFragment() url.Values {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFragment")
	ret0, _ := ret[0].(url.Values)
	return ret0
}

// GetFragment indicates an expected call of GetFragment
func (mr *MockAuthorizeResponderMockRecorder) GetFragment() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFragment", reflect.TypeOf((*MockAuthorizeResponder)(nil).GetFragment))
}

// GetHeader mocks base method
func (m *MockAuthorizeResponder) GetHeader() http.Header {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetHeader")
	ret0, _ := ret[0].(http.Header)
	return ret0
}

// GetHeader indicates an expected call of GetHeader
func (mr *MockAuthorizeResponderMockRecorder) GetHeader() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetHeader", reflect.TypeOf((*MockAuthorizeResponder)(nil).GetHeader))
}

// GetQuery mocks base method
func (m *MockAuthorizeResponder) GetQuery() url.Values {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetQuery")
	ret0, _ := ret[0].(url.Values)
	return ret0
}

// GetQuery indicates an expected call of GetQuery
func (mr *MockAuthorizeResponderMockRecorder) GetQuery() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetQuery", reflect.TypeOf((*MockAuthorizeResponder)(nil).GetQuery))
}
