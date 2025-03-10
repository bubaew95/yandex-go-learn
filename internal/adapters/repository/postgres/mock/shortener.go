// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/core/ports/ports.go

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

// DeleteUserURLS mocks base method.
func (m *MockShortenerRepository) DeleteUserURLS(ctx context.Context, items []model.URLToDelete) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUserURLS", ctx, items)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUserURLS indicates an expected call of DeleteUserURLS.
func (mr *MockShortenerRepositoryMockRecorder) DeleteUserURLS(ctx, items interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUserURLS", reflect.TypeOf((*MockShortenerRepository)(nil).DeleteUserURLS), ctx, items)
}

// GetURLByID mocks base method.
func (m *MockShortenerRepository) GetURLByID(ctx context.Context, id string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetURLByID", ctx, id)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
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
func (m *MockShortenerRepository) GetURLSByUserID(ctx context.Context, userID string) (map[string]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetURLSByUserID", ctx, userID)
	ret0, _ := ret[0].(map[string]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetURLSByUserID indicates an expected call of GetURLSByUserID.
func (mr *MockShortenerRepositoryMockRecorder) GetURLSByUserID(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetURLSByUserID", reflect.TypeOf((*MockShortenerRepository)(nil).GetURLSByUserID), ctx, userID)
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
func (m *MockShortenerRepository) Ping(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ping", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Ping indicates an expected call of Ping.
func (mr *MockShortenerRepositoryMockRecorder) Ping(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ping", reflect.TypeOf((*MockShortenerRepository)(nil).Ping), ctx)
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

// MockShortenerService is a mock of ShortenerService interface.
type MockShortenerService struct {
	ctrl     *gomock.Controller
	recorder *MockShortenerServiceMockRecorder
}

// MockShortenerServiceMockRecorder is the mock recorder for MockShortenerService.
type MockShortenerServiceMockRecorder struct {
	mock *MockShortenerService
}

// NewMockShortenerService creates a new mock instance.
func NewMockShortenerService(ctrl *gomock.Controller) *MockShortenerService {
	mock := &MockShortenerService{ctrl: ctrl}
	mock.recorder = &MockShortenerServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockShortenerService) EXPECT() *MockShortenerServiceMockRecorder {
	return m.recorder
}

// DeleteUserURLS mocks base method.
func (m *MockShortenerService) DeleteUserURLS(ctx context.Context, items []model.URLToDelete) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUserURLS", ctx, items)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUserURLS indicates an expected call of DeleteUserURLS.
func (mr *MockShortenerServiceMockRecorder) DeleteUserURLS(ctx, items interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUserURLS", reflect.TypeOf((*MockShortenerService)(nil).DeleteUserURLS), ctx, items)
}

// GenerateURL mocks base method.
func (m *MockShortenerService) GenerateURL(ctx context.Context, url string, randomStringLength int) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateURL", ctx, url, randomStringLength)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateURL indicates an expected call of GenerateURL.
func (mr *MockShortenerServiceMockRecorder) GenerateURL(ctx, url, randomStringLength interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateURL", reflect.TypeOf((*MockShortenerService)(nil).GenerateURL), ctx, url, randomStringLength)
}

// GetURLByID mocks base method.
func (m *MockShortenerService) GetURLByID(ctx context.Context, id string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetURLByID", ctx, id)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetURLByID indicates an expected call of GetURLByID.
func (mr *MockShortenerServiceMockRecorder) GetURLByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetURLByID", reflect.TypeOf((*MockShortenerService)(nil).GetURLByID), ctx, id)
}

// GetURLByOriginalURL mocks base method.
func (m *MockShortenerService) GetURLByOriginalURL(ctx context.Context, originalURL string) (string, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetURLByOriginalURL", ctx, originalURL)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// GetURLByOriginalURL indicates an expected call of GetURLByOriginalURL.
func (mr *MockShortenerServiceMockRecorder) GetURLByOriginalURL(ctx, originalURL interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetURLByOriginalURL", reflect.TypeOf((*MockShortenerService)(nil).GetURLByOriginalURL), ctx, originalURL)
}

// GetURLSByUserID mocks base method.
func (m *MockShortenerService) GetURLSByUserID(ctx context.Context, userID string) ([]model.ShortenerURLSForUserResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetURLSByUserID", ctx, userID)
	ret0, _ := ret[0].([]model.ShortenerURLSForUserResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetURLSByUserID indicates an expected call of GetURLSByUserID.
func (mr *MockShortenerServiceMockRecorder) GetURLSByUserID(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetURLSByUserID", reflect.TypeOf((*MockShortenerService)(nil).GetURLSByUserID), ctx, userID)
}

// InsertURLs mocks base method.
func (m *MockShortenerService) InsertURLs(ctx context.Context, urls []model.ShortenerURLMapping) ([]model.ShortenerURLResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertURLs", ctx, urls)
	ret0, _ := ret[0].([]model.ShortenerURLResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// InsertURLs indicates an expected call of InsertURLs.
func (mr *MockShortenerServiceMockRecorder) InsertURLs(ctx, urls interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertURLs", reflect.TypeOf((*MockShortenerService)(nil).InsertURLs), ctx, urls)
}

// Ping mocks base method.
func (m *MockShortenerService) Ping(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ping", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Ping indicates an expected call of Ping.
func (mr *MockShortenerServiceMockRecorder) Ping(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ping", reflect.TypeOf((*MockShortenerService)(nil).Ping), ctx)
}

// RandStringBytes mocks base method.
func (m *MockShortenerService) RandStringBytes(n int) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RandStringBytes", n)
	ret0, _ := ret[0].(string)
	return ret0
}

// RandStringBytes indicates an expected call of RandStringBytes.
func (mr *MockShortenerServiceMockRecorder) RandStringBytes(n interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RandStringBytes", reflect.TypeOf((*MockShortenerService)(nil).RandStringBytes), n)
}

// ScheduleURLDeletion mocks base method.
func (m *MockShortenerService) ScheduleURLDeletion(ctx context.Context, items []model.URLToDelete) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ScheduleURLDeletion", ctx, items)
}

// ScheduleURLDeletion indicates an expected call of ScheduleURLDeletion.
func (mr *MockShortenerServiceMockRecorder) ScheduleURLDeletion(ctx, items interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ScheduleURLDeletion", reflect.TypeOf((*MockShortenerService)(nil).ScheduleURLDeletion), ctx, items)
}
