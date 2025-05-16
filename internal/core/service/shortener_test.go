package service

import (
	"context"
	"errors"
	"testing"

	"github.com/bubaew95/yandex-go-learn/config"
	"github.com/bubaew95/yandex-go-learn/internal/core/model"
	"github.com/stretchr/testify/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetURLByID(t *testing.T) {
	repo := NewMockShortenerRepository(t)
	service := NewShortenerService(repo, config.Config{})

	link := "https://example.com"

	repo.On("GetURLByID", mock.Anything, "SXhhC3").
		Return(link, nil)

	url, err := service.GetURLByID(context.Background(), "SXhhC3")
	require.NoError(t, err)
	assert.Equal(t, link, url)
}

func TestInsertURLError(t *testing.T) {
	repo := NewMockShortenerRepository(t)
	service := NewShortenerService(repo, config.Config{})

	data := []model.ShortenerURLMapping{
		{CorrelationID: "test-id", OriginalURL: "http://example.com"},
	}

	mockError := errors.New("Mock Error")

	repo.On("InsertURLs", mock.Anything, data).Return(mockError)

	items, err := service.InsertURLs(context.Background(), data)
	assert.Error(t, err)

	assert.Equal(t, err, mockError)
	assert.Nil(t, items)
}

func TestInsertURLSuccess(t *testing.T) {
	repo := NewMockShortenerRepository(t)
	service := NewShortenerService(repo, config.Config{})

	data := []model.ShortenerURLMapping{
		{CorrelationID: "test-id-1", OriginalURL: "http://example.com"},
		{CorrelationID: "test-id-2", OriginalURL: "http://site.com"},
	}

	repo.On("InsertURLs", mock.Anything, data).Return(nil)

	items, err := service.InsertURLs(context.Background(), data)
	require.NoError(t, err)

	assert.Equal(t, data[0].CorrelationID, items[0].CorrelationID)
}
