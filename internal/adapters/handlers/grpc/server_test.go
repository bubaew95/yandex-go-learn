package grpc

import (
	"context"
	"errors"
	"github.com/bubaew95/yandex-go-learn/internal/adapters/constants"
	"github.com/bubaew95/yandex-go-learn/internal/core/model"
	pb "github.com/bubaew95/yandex-go-learn/internal/proto"
	"github.com/bubaew95/yandex-go-learn/pkg/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"testing"
)

func TestCreateURL(t *testing.T) {
	t.Parallel()

	type mockData struct {
		shortURL string
		err      error
	}

	tests := []struct {
		name       string
		originURL  string
		statusCode codes.Code
		mockData   mockData
	}{
		{
			name:       "Simple url",
			originURL:  "https://practicum.yandex.ru/",
			statusCode: codes.OK,
			mockData: mockData{
				shortURL: "https://yan.dex/WfgSF3",
			},
		},
		{
			name:      "Url already exists",
			originURL: "https://practicum.yandex.ru/",
			mockData: mockData{
				err:      constants.ErrUniqueIndex,
				shortURL: "https://yan.dex/WfgSF3",
			},
			statusCode: codes.AlreadyExists,
		},
		{
			name:       "Data is empty",
			originURL:  "",
			statusCode: codes.InvalidArgument,
			mockData: mockData{
				err: constants.ErrParamsIsEmpty,
			},
		},
		{
			name:       "Other error",
			originURL:  "https://practicum.yandex.ru/",
			statusCode: codes.Internal,
			mockData: mockData{
				err: errors.New("some error"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewMockShortenerService(t)

			server := &Server{
				service: service,
			}

			service.On("GenerateURL", mock.Anything, tt.originURL, mock.Anything).
				Return(tt.mockData.shortURL, tt.mockData.err)

			if tt.statusCode == codes.AlreadyExists {
				service.On("GetURLByOriginalURL", mock.Anything, tt.originURL).
					Return(tt.mockData.shortURL, true)
			}

			resp, err := server.CreateURL(context.Background(), &pb.ShortenRequest{
				OriginalUrl: tt.originURL,
			})

			if err != nil {
				if e, ok := status.FromError(err); ok {
					assert.Equal(t, tt.statusCode, e.Code())
				}
			} else {
				assert.Equal(t, tt.mockData.shortURL, resp.ShortUrl)
			}
		})
	}
}

func TestGetURL(t *testing.T) {
	t.Parallel()

	service := NewMockShortenerService(t)

	server := &Server{
		service: service,
	}

	type want struct {
		statusCode codes.Code
		originURL  string
	}

	type wantMock struct {
		originURL string
		err       error
	}

	tests := []struct {
		name     string
		shortURL string
		want     want
		mock     wantMock
	}{
		{
			name:     "Simple test",
			shortURL: "WzYAhS",
			mock: wantMock{
				originURL: "https://practicum.yandex.ru/",
			},
			want: want{
				statusCode: codes.OK,
				originURL:  "https://practicum.yandex.ru/",
			},
		},
		{
			name:     "Deleted shortURL",
			shortURL: "WzYAhSs",
			mock: wantMock{
				err: constants.ErrIsDeleted,
			},
			want: want{
				statusCode: codes.DeadlineExceeded,
			},
		},
		{
			name:     "Bad request test",
			shortURL: "WzYAhSsss",
			mock: wantMock{
				err: errors.New("Not Found"),
			},
			want: want{
				statusCode: codes.NotFound,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			service.On("GetURLByID", mock.Anything, tt.shortURL).
				Return(tt.mock.originURL, tt.mock.err).Once()

			resp, err := server.GetURL(context.Background(), &pb.GetURLRequest{
				ShortId: tt.shortURL,
			})

			if err != nil {
				if e, ok := status.FromError(err); ok {
					assert.Equal(t, tt.want.statusCode, e.Code())
				}
			} else {
				assert.Equal(t, tt.want.originURL, resp.OriginUrl)
			}
		})
	}
}

func TestBatch(t *testing.T) {
	t.Parallel()

	type mockRequest struct {
		input  []model.ShortenerURLMapping
		result []model.ShortenerURLResponse
		err    error
	}

	tests := []struct {
		name       string
		req        *pb.BatchRequest
		want       *pb.BatchResponse
		mockResult mockRequest
		statusCode codes.Code
	}{
		{
			name: "Batch OK",
			req: &pb.BatchRequest{
				Urls: []*pb.URLMappingRequest{
					{CorrelationId: "1", OriginalUrl: "https://ya.ru/1"},
					{CorrelationId: "2", OriginalUrl: "https://ya.ru/2"},
				},
			},
			mockResult: mockRequest{
				input: []model.ShortenerURLMapping{
					{CorrelationID: "1", OriginalURL: "https://ya.ru/1"},
					{CorrelationID: "2", OriginalURL: "https://ya.ru/2"},
				},
				result: []model.ShortenerURLResponse{
					{CorrelationID: "1", ShortURL: "http://sht/abc1"},
					{CorrelationID: "2", ShortURL: "http://sht/abc2"},
				},
			},
			want: &pb.BatchResponse{
				Urls: []*pb.URLMappingResponse{
					{CorrelationId: "1", ShortUrl: "http://sht/abc1"},
					{CorrelationId: "2", ShortUrl: "http://sht/abc2"},
				},
			},
			statusCode: codes.OK,
		},
		{
			name: "Internal error",
			req: &pb.BatchRequest{
				Urls: []*pb.URLMappingRequest{
					{CorrelationId: "x", OriginalUrl: "invalid-url"},
				},
			},
			mockResult: mockRequest{
				input: []model.ShortenerURLMapping{
					{CorrelationID: "x", OriginalURL: "invalid-url"},
				},
				err: errors.New("db error"),
			},
			want:       nil,
			statusCode: codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewMockShortenerService(t)

			server := &Server{
				service: service,
			}

			service.On("InsertURLs", mock.Anything, tt.mockResult.input).
				Return(tt.mockResult.result, tt.mockResult.err).Once()

			resp, err := server.Batch(context.Background(), tt.req)

			if err != nil {
				st, ok := status.FromError(err)
				require.True(t, ok)
				assert.Equal(t, tt.statusCode, st.Code())
			} else {
				assert.Equal(t, tt.want, resp)
			}
		})
	}
}

func TestAuthInterceptor(t *testing.T) {
	inter := AuthInterceptor()

	t.Run("valid user_id", func(t *testing.T) {
		encoded, _ := crypto.EncodeUserID(crypto.GenerateUserID())
		md := metadata.New(map[string]string{
			"user_id": encoded,
		})
		ctx := metadata.NewIncomingContext(context.Background(), md)

		_, err := inter(ctx, nil, &grpc.UnaryServerInfo{FullMethod: "/test.Test/Test"}, func(ctx context.Context, req interface{}) (interface{}, error) {
			val := ctx.Value(crypto.KeyUserID)
			assert.Equal(t, encoded, val)
			return "ok", nil
		})

		assert.NoError(t, err)
	})

	t.Run("missing user_id in metadata", func(t *testing.T) {
		md := metadata.New(map[string]string{})
		ctx := metadata.NewIncomingContext(context.Background(), md)

		_, err := inter(ctx, nil, &grpc.UnaryServerInfo{FullMethod: "/test.Test/Test"}, func(ctx context.Context, req interface{}) (interface{}, error) {
			val := ctx.Value(crypto.KeyUserID)
			assert.NotEmpty(t, val)
			return "ok", nil
		})

		assert.NoError(t, err)
	})

	t.Run("no metadata at all", func(t *testing.T) {
		ctx := context.Background()

		_, err := inter(ctx, nil, &grpc.UnaryServerInfo{FullMethod: "/test.Test/Test"}, func(ctx context.Context, req interface{}) (interface{}, error) {
			t.Error("handler should not be called")
			return nil, nil
		})

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Unauthenticated, st.Code())
	})
}

func TestGetUserURLs(t *testing.T) {
	t.Parallel()

	type want struct {
		code codes.Code
		urls []*pb.URLPair
	}

	type mockData struct {
		userID string
		result []model.ShortenerURLSForUserResponse
		err    error
	}

	tests := []struct {
		name string
		ctx  context.Context
		mock mockData
		want want
	}{
		{
			name: "Success",
			ctx:  context.WithValue(context.Background(), "userID", "encodedUser123"),
			mock: mockData{
				userID: "encodedUser123",
				result: []model.ShortenerURLSForUserResponse{
					{ShortURL: "https://yan.dex/a1b2c3", OriginalURL: "https://example.com"},
				},
			},
			want: want{
				code: codes.OK,
				urls: []*pb.URLPair{
					{ShortUrl: "https://yan.dex/a1b2c3", OriginalUrl: "https://example.com"},
				},
			},
		},
		{
			name: "No userID in context",
			ctx:  context.Background(),
			want: want{
				code: codes.Unauthenticated,
			},
		},
		{
			name: "Repository error",
			ctx:  context.WithValue(context.Background(), "userID", "errorUser"),
			mock: mockData{
				userID: "errorUser",
				err:    errors.New("db error"),
			},
			want: want{
				code: codes.Internal,
			},
		},
		{
			name: "Empty URLs",
			ctx:  context.WithValue(context.Background(), "userID", "emptyUser"),
			mock: mockData{
				userID: "emptyUser",
				result: []model.ShortenerURLSForUserResponse{},
			},
			want: want{
				code: codes.NotFound,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewMockShortenerService(t)
			server := &Server{service: service}

			if tt.mock.userID != "" {
				service.On("GetURLSByUserID", mock.Anything, tt.mock.userID).
					Return(tt.mock.result, tt.mock.err).
					Once()
			}

			resp, err := server.GetUserURLs(tt.ctx, &pb.GetUserURLsRequest{})
			if err != nil {
				s, ok := status.FromError(err)
				require.True(t, ok)
				assert.Equal(t, tt.want.code, s.Code())
			} else {
				assert.Equal(t, tt.want.urls, resp.Urls)
			}
		})
	}
}
