// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/core/ports/ishortener_repository.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	model "github.com/bubaew95/yandex-go-learn/internal/core/model"
	gomock "github.com/golang/mock/gomock"
)

// MockShortenerRepositoryInterface is a mock of ShortenerRepositoryInterface interface.
type MockShortenerRepositoryInterface struct {
	ctrl     *gomock.Controller
	recorder *MockShortenerRepositoryInterfaceMockRecorder
}

// MockShortenerRepositoryInterfaceMockRecorder is the mock recorder for MockShortenerRepositoryInterface.
type MockShortenerRepositoryInterfaceMockRecorder struct {
	mock *MockShortenerRepositoryInterface
}

// NewMockShortenerRepositoryInterface creates a new mock instance.
func NewMockShortenerRepositoryInterface(ctrl *gomock.Controller) *MockShortenerRepositoryInterface {
	mock := &MockShortenerRepositoryInterface{ctrl: ctrl}
	mock.recorder = &MockShortenerRepositoryInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockShortenerRepositoryInterface) EXPECT() *MockShortenerRepositoryInterfaceMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockShortenerRepositoryInterface) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockShortenerRepositoryInterfaceMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockShortenerRepositoryInterface)(nil).Close))
}

// GetAllURL mocks base method.
func (m *MockShortenerRepositoryInterface) GetAllURL(ctx context.Context) map[string]string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllURL", ctx)
	ret0, _ := ret[0].(map[string]string)
	return ret0
}

// GetAllURL indicates an expected call of GetAllURL.
func (mr *MockShortenerRepositoryInterfaceMockRecorder) GetAllURL(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllURL", reflect.TypeOf((*MockShortenerRepositoryInterface)(nil).GetAllURL), ctx)
}

// GetURLByID mocks base method.
func (m *MockShortenerRepositoryInterface) GetURLByID(ctx context.Context, id string) (string, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetURLByID", ctx, id)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// GetURLByID indicates an expected call of GetURLByID.
func (mr *MockShortenerRepositoryInterfaceMockRecorder) GetURLByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetURLByID", reflect.TypeOf((*MockShortenerRepositoryInterface)(nil).GetURLByID), ctx, id)
}

// InsertURLs mocks base method.
func (m *MockShortenerRepositoryInterface) InsertURLs(ctx context.Context, urls []model.ShortenerURLMapping) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertURLs", ctx, urls)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertURLs indicates an expected call of InsertURLs.
func (mr *MockShortenerRepositoryInterfaceMockRecorder) InsertURLs(ctx, urls interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertURLs", reflect.TypeOf((*MockShortenerRepositoryInterface)(nil).InsertURLs), ctx, urls)
}

// Ping mocks base method.
func (m *MockShortenerRepositoryInterface) Ping() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ping")
	ret0, _ := ret[0].(error)
	return ret0
}

// Ping indicates an expected call of Ping.
func (mr *MockShortenerRepositoryInterfaceMockRecorder) Ping() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ping", reflect.TypeOf((*MockShortenerRepositoryInterface)(nil).Ping))
}

// SetURL mocks base method.
func (m *MockShortenerRepositoryInterface) SetURL(ctx context.Context, id, url string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetURL", ctx, id, url)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetURL indicates an expected call of SetURL.
func (mr *MockShortenerRepositoryInterfaceMockRecorder) SetURL(ctx, id, url interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetURL", reflect.TypeOf((*MockShortenerRepositoryInterface)(nil).SetURL), ctx, id, url)
}
