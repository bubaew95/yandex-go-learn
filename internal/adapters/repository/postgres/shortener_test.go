package postgres

import (
	"context"
	"errors"
	"testing"

	"github.com/bubaew95/yandex-go-learn/internal/adapters/repository/postgres/mock"
	"github.com/bubaew95/yandex-go-learn/internal/core/model"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetURLByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockShortenerRepository(ctrl)
	ctx := context.Background()

	value := "http://noriba.ru"
	m.EXPECT().
		GetURLByID(ctx, "SXhhC3").
		Return(value, nil)

	url, err := m.GetURLByID(ctx, "SXhhC3")
	require.NoError(t, err)
	assert.Equal(t, url, value)
}

func TestInsertURLError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockShortenerRepository(ctrl)
	ctx := context.Background()

	data := []model.ShortenerURLMapping{
		{CorrelationID: "test-id", OriginalURL: "http://example.com"},
	}

	mockError := errors.New("Mock Error")
	m.EXPECT().
		InsertURLs(ctx, data).
		Return(mockError)

	err := m.InsertURLs(ctx, data)
	assert.Error(t, err)
	assert.Equal(t, err, mockError)
}

func TestInsertURLSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockShortenerRepository(ctrl)
	ctx := context.Background()

	data := []model.ShortenerURLMapping{
		{CorrelationID: "test-id-1", OriginalURL: "http://example.com"},
		{CorrelationID: "test-id-2", OriginalURL: "http://site.com"},
	}

	m.EXPECT().
		InsertURLs(ctx, data).
		Return(nil)

	err := m.InsertURLs(ctx, data)
	assert.NoError(t, err)
}
