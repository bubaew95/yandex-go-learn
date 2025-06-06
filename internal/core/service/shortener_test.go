package service

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/bubaew95/yandex-go-learn/config"
	"github.com/bubaew95/yandex-go-learn/internal/core/model"
	"github.com/stretchr/testify/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateURL(t *testing.T) {
	repo := NewMockShortenerRepository(t)
	service := NewShortenerService(repo, config.Config{
		BaseURL: "http://short.url",
	})

	repo.On("GetURLByID", mock.Anything, mock.Anything).Return("", errors.New("not found")).Once()
	repo.On("SetURL", mock.Anything, mock.Anything, "https://www.yandex.ru").Return(nil).Once()

	url, err := service.GenerateURL(context.Background(), "https://www.yandex.ru", 10)
	require.NoError(t, err)

	assert.Contains(t, url, "http://short.url/")
	repo.AssertExpectations(t)
}

func TestRandStringBytes(t *testing.T) {
	repo := NewMockShortenerRepository(t)
	service := NewShortenerService(repo, config.Config{})

	randomID := service.RandStringBytes(8)

	assert.NotEmpty(t, randomID)
	assert.Len(t, randomID, 8)
}

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

func TestGetURLByOriginalURL(t *testing.T) {
	repo := NewMockShortenerRepository(t)
	service := NewShortenerService(repo, config.Config{
		BaseURL: "http://short.url",
	})

	type mockStruct struct {
		url string
		ok  bool
	}

	tests := []struct {
		name string
		mock mockStruct
		want string
	}{
		{
			name: "success",
			mock: mockStruct{
				url: "wetwet",
				ok:  true,
			},
			want: "http://short.url/wetwet",
		},
		{
			name: "success",
			mock: mockStruct{
				url: "wetwet",
				ok:  false,
			},
			want: "http://short.url/wetwet",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo.On("GetURLByOriginalURL", mock.Anything, tt.mock.url).Return(tt.mock.url, tt.mock.ok).Once()

			url, ok := service.GetURLByOriginalURL(context.Background(), tt.mock.url)
			fmt.Println(url, ok)

			assert.Equal(t, tt.mock.ok, ok)
			if ok {
				assert.Equal(t, tt.want, url)
			}
		})
	}
}

func TestGenerateResponseUrl(t *testing.T) {
	repo := NewMockShortenerRepository(t)
	service := NewShortenerService(repo, config.Config{
		BaseURL: "http://short.url",
	})

	url := service.generateResponseURL("SXhhC3")

	assert.Equal(t, "http://short.url/SXhhC3", url)
}

func TestInsertURL(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		response []model.ShortenerURLResponse
	}{
		{
			name: "Success",
			err:  nil,
		},
		{
			name: "Error",
			err:  errors.New("Mock error"),
		},
	}

	repo := NewMockShortenerRepository(t)
	service := NewShortenerService(repo, config.Config{})

	data := []model.ShortenerURLMapping{
		{CorrelationID: "test-id-1", OriginalURL: "http://example.com"},
		{CorrelationID: "test-id-2", OriginalURL: "http://site.com"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo.On("InsertURLs", mock.Anything, data).Return(tt.err).Once()

			items, err := service.InsertURLs(context.Background(), data)
			if tt.err != nil {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Len(t, items, len(data))
			}
		})
	}
}

func TestShortenerService_GetURLSByUserID(t *testing.T) {
	t.Parallel()

	type args struct {
		userID string
	}
	type repoReturn struct {
		items map[string]string
		err   error
	}
	type want struct {
		result []model.ShortenerURLSForUserResponse
		err    bool
	}

	tests := []struct {
		name       string
		args       args
		repoReturn repoReturn
		want       want
	}{
		{
			name: "success - returns list",
			args: args{userID: "user123"},
			repoReturn: repoReturn{
				items: map[string]string{
					"id1": "http://example.com",
					"id2": "http://yandex.ru",
				},
				err: nil,
			},
			want: want{
				result: []model.ShortenerURLSForUserResponse{
					{ShortURL: "http://short.url/id1", OriginalURL: "http://example.com"},
					{ShortURL: "http://short.url/id2", OriginalURL: "http://yandex.ru"},
				},
				err: false,
			},
		},
		{
			name: "repo error",
			args: args{userID: "user123"},
			repoReturn: repoReturn{
				items: nil,
				err:   errors.New("db error"),
			},
			want: want{
				result: nil,
				err:    true,
			},
		},
		{
			name: "empty result",
			args: args{userID: "user456"},
			repoReturn: repoReturn{
				items: map[string]string{},
				err:   nil,
			},
			want: want{
				result: []model.ShortenerURLSForUserResponse{},
				err:    false,
			},
		},
	}

	for _, tt := range tests {
		tt := tt // захват переменной
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockRepo := NewMockShortenerRepository(t)
			svc := NewShortenerService(mockRepo, config.Config{
				BaseURL: "http://short.url",
			})

			mockRepo.On("GetURLSByUserID", mock.Anything, tt.args.userID).
				Return(tt.repoReturn.items, tt.repoReturn.err).
				Once()

			got, err := svc.GetURLSByUserID(context.Background(), tt.args.userID)

			if tt.want.err {
				require.Error(t, err)
				assert.Nil(t, got)
			} else {
				require.NoError(t, err)
				assert.ElementsMatch(t, tt.want.result, got)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestShortenerService_ScheduleAndRunDeletion(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg := &sync.WaitGroup{}
	mockRepo := NewMockShortenerRepository(t)

	deleteChan := make(chan model.URLToDelete, 10)
	svc := NewShortenerService(mockRepo, config.Config{})
	svcWithChan := svc
	svcWithChan.deleteChan = deleteChan

	itemsToDelete := []model.URLToDelete{
		{ShortLink: "link1", UserID: "user1"},
		{ShortLink: "link2", UserID: "user1"},
	}

	mockRepo.On("DeleteUserURLS", mock.Anything, mock.MatchedBy(func(batch []model.URLToDelete) bool {
		return len(batch) == 2 &&
			batch[0].ShortLink == "link1" &&
			batch[1].ShortLink == "link2"
	})).Return(nil).Once()

	go svcWithChan.Run(ctx, wg)

	svcWithChan.ScheduleURLDeletion(ctx, itemsToDelete)

	time.Sleep(6 * time.Second)

	cancel()
	wg.Wait()

	mockRepo.AssertExpectations(t)
}

func TestShortenerService_Stats(t *testing.T) {
	repo := NewMockShortenerRepository(t)
	service := NewShortenerService(repo, config.Config{})

	expected := model.StatsRespose{
		Users: 3,
		URLs:  15,
	}

	repo.On("Stats", mock.Anything).Return(expected, nil).Once()

	result, err := service.Stats(context.Background())
	require.NoError(t, err)
	assert.Equal(t, expected, result)

	repo.AssertExpectations(t)
}
