// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/core/ports/shortener_repository.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	model "github.com/bubaew95/yandex-go-learn/internal/core/model"
	gomock "github.com/golang/mock/gomock"
)

// MockShortenerRepository is a mock of ShortenerRepository interface.
type MockShortenerRepository struct {
	ctrl     *gomock.Controller
	recorder *MockShortenerRepositoryMockRecorder
}

// MockShortenerRepositoryMockRecorder is the mock recorder for MockShortenerRepository.
type MockShortenerRepositoryMockRecorder struct {
	mock *MockShortenerRepository
}

// NewMockShortenerRepository creates a new mock instance.
func NewMockShortenerRepository(ctrl *gomock.Controller) *MockShortenerRepository {
	mock := &MockShortenerRepository{ctrl: ctrl}
	mock.recorder = &MockShortenerRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockShortenerRepository) EXPECT() *MockShortenerRepositoryMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockShortenerRepository) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockShortenerRepositoryMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockShortenerRepository)(nil).Close))
}

// GetAllURL mocks base method.
func (m *MockShortenerRepository) GetAllURL(ctx context.Context) map[string]string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllURL", ctx)
	ret0, _ := ret[0].(map[string]string)
	return ret0
}

// GetAllURL indicates an expected call of GetAllURL.
func (mr *MockShortenerRepositoryMockRecorder) GetAllURL(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllURL", reflect.TypeOf((*MockShortenerRepository)(nil).GetAllURL), ctx)
}

// GetURLByID mocks base method.
func (m *MockShortenerRepository) GetURLByID(ctx context.Context, id string) (string, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetURLByID", ctx, id)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// GetURLByID indicates an expected call of GetURLByID.
func (mr *MockShortenerRepositoryMockRecorder) GetURLByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetURLByID", reflect.TypeOf((*MockShortenerRepository)(nil).GetURLByID), ctx, id)
}

// GetURLByOriginalURL mocks base method.
func (m *MockShortenerRepository) GetURLByOriginalURL(ctx context.Context, originalURL string) (string, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetURLByOriginalURL", ctx, originalURL)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// GetURLByOriginalURL indicates an expected call of GetURLByOriginalURL.
func (mr *MockShortenerRepositoryMockRecorder) GetURLByOriginalURL(ctx, originalURL interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetURLByOriginalURL", reflect.TypeOf((*MockShortenerRepository)(nil).GetURLByOriginalURL), ctx, originalURL)
}

// GetURLSByUserID mocks base method.
func (m *MockShortenerRepository) GetURLSByUserID(ctx context.Context, user_id string) (map[string]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetURLSByUserID", ctx, user_id)
	ret0, _ := ret[0].(map[string]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetURLSByUserID indicates an expected call of GetURLSByUserID.
func (mr *MockShortenerRepositoryMockRecorder) GetURLSByUserID(ctx, user_id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetURLSByUserID", reflect.TypeOf((*MockShortenerRepository)(nil).GetURLSByUserID), ctx, user_id)
}

// InsertURLs mocks base method.
func (m *MockShortenerRepository) InsertURLs(ctx context.Context, urls []model.ShortenerURLMapping) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertURLs", ctx, urls)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertURLs indicates an expected call of InsertURLs.
func (mr *MockShortenerRepositoryMockRecorder) InsertURLs(ctx, urls interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertURLs", reflect.TypeOf((*MockShortenerRepository)(nil).InsertURLs), ctx, urls)
}

// Ping mocks base method.
func (m *MockShortenerRepository) Ping() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ping")
	ret0, _ := ret[0].(error)
	return ret0
}

// Ping indicates an expected call of Ping.
func (mr *MockShortenerRepositoryMockRecorder) Ping() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ping", reflect.TypeOf((*MockShortenerRepository)(nil).Ping))
}

// SetURL mocks base method.
func (m *MockShortenerRepository) SetURL(ctx context.Context, id, url string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetURL", ctx, id, url)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetURL indicates an expected call of SetURL.
func (mr *MockShortenerRepositoryMockRecorder) SetURL(ctx, id, url interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetURL", reflect.TypeOf((*MockShortenerRepository)(nil).SetURL), ctx, id, url)
}
