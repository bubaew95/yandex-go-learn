package grpc

import (
	"context"
	"errors"
	"github.com/bubaew95/yandex-go-learn/internal/adapters/constants"
	"github.com/bubaew95/yandex-go-learn/internal/core/model"
	pb "github.com/bubaew95/yandex-go-learn/internal/proto"
	"github.com/bubaew95/yandex-go-learn/pkg/crypto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"net/http"
)

// ShortenerService определяет бизнес-логику сервиса сокращения ссылок.
// Включает в себя генерацию ссылок, работу с пользователями и отложенное удаление.
//
//go:generate go run github.com/vektra/mockery/v2@v2.52.2 --name=ShortenerService --filename=servicemock_test.go --inpackage
type ShortenerService interface {
	// GenerateURL генерирует короткий URL на основе оригинального.
	GenerateURL(ctx context.Context, url string, randomStringLength int) (string, error)

	// GetURLByID возвращает оригинальный URL по его сокращённому ID.
	GetURLByID(ctx context.Context, id string) (string, error)

	// GetURLByOriginalURL возвращает ID, соответствующий оригинальному URL.
	GetURLByOriginalURL(ctx context.Context, originalURL string) (string, bool)

	// InsertURLs добавляет множество URL и возвращает их короткие представления.
	InsertURLs(ctx context.Context, urls []model.ShortenerURLMapping) ([]model.ShortenerURLResponse, error)

	// GetURLSByUserID возвращает список сокращённых URL, принадлежащих пользователю.
	GetURLSByUserID(ctx context.Context, userID string) ([]model.ShortenerURLSForUserResponse, error)

	// DeleteUserURLS помечает ссылки как удалённые.
	DeleteUserURLS(ctx context.Context, items []model.URLToDelete) error

	// ScheduleURLDeletion планирует асинхронное удаление ссылок (например, через очередь).
	ScheduleURLDeletion(ctx context.Context, items []model.URLToDelete)

	// RandStringBytes генерирует случайную строку заданной длины (обычно для ID короткой ссылки).
	RandStringBytes(n int) string

	// Ping проверяет доступность сервиса (например, для liveness-проб).
	Ping(ctx context.Context) error

	// Stats возвращает статистику
	Stats(ctx context.Context) (model.StatsRespose, error)
}

type Server struct {
	pb.UnimplementedURLShortenerServer
	service ShortenerService
}

func NewServer(s ShortenerService) *Server {
	return &Server{service: s}
}

func AuthInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "metadata not found")
		}

		userID := ""
		if vals := md.Get("user_id"); len(vals) > 0 {
			userID = vals[0]
		}

		if userID == "" || !crypto.IsInvalidUserID(&http.Cookie{Name: "user_id", Value: userID}) {
			rawID := crypto.GenerateUserID()
			encodedID, err := crypto.EncodeUserID(rawID)
			if err != nil {
				return nil, status.Error(codes.Internal, "user ID encoding failed")
			}
			userID = encodedID
		}

		ctx = context.WithValue(ctx, crypto.KeyUserID, userID)
		return handler(ctx, req)
	}
}

func (s *Server) CreateURL(ctx context.Context, req *pb.ShortenRequest) (*pb.ShortenResponse, error) {
	url, err := s.service.GenerateURL(ctx, req.OriginalUrl, 8)

	if err != nil {
		if errors.Is(err, constants.ErrParamsIsEmpty) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		if errors.Is(err, constants.ErrUniqueIndex) {
			originURL, ok := s.service.GetURLByOriginalURL(ctx, req.OriginalUrl)
			if ok {
				return &pb.ShortenResponse{
					ShortUrl: originURL,
				}, nil
			}
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.ShortenResponse{
		ShortUrl: url,
	}, nil
}

func (s *Server) GetURL(ctx context.Context, req *pb.GetURLRequest) (*pb.GetURLResponse, error) {
	url, err := s.service.GetURLByID(ctx, req.ShortId)

	if err != nil || url == "" {
		if errors.Is(err, constants.ErrIsDeleted) {
			return nil, status.Error(codes.DeadlineExceeded, err.Error())
		}

		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &pb.GetURLResponse{
		OriginUrl: url,
	}, nil
}

func (s *Server) Batch(ctx context.Context, req *pb.BatchRequest) (*pb.BatchResponse, error) {
	mappings := make([]model.ShortenerURLMapping, 0, len(req.Urls))
	for _, u := range req.Urls {
		mappings = append(mappings, model.ShortenerURLMapping{
			CorrelationID: u.CorrelationId,
			OriginalURL:   u.OriginalUrl,
		})
	}

	results, err := s.service.InsertURLs(ctx, mappings)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to insert URLs")
	}

	// Преобразуем обратно в proto-модель
	resp := &pb.BatchResponse{}
	for _, r := range results {
		resp.Urls = append(resp.Urls, &pb.URLMappingResponse{
			CorrelationId: r.CorrelationID,
			ShortUrl:      r.ShortURL,
		})
	}

	return resp, nil
}

func (s *Server) GetUserURLs(ctx context.Context, _ *pb.GetUserURLsRequest) (*pb.GetUserURLsResponse, error) {
	userID, ok := ctx.Value(crypto.KeyUserID).(string)
	if !ok || userID == "" {
		return nil, status.Error(codes.Unauthenticated, "user id is required")
	}

	urls, err := s.service.GetURLSByUserID(ctx, userID)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get user urls")
	}

	if len(urls) == 0 {
		return nil, status.Error(codes.NotFound, "no urls")
	}

	resp := &pb.GetUserURLsResponse{}
	for _, u := range urls {
		resp.Urls = append(resp.Urls, &pb.URLPair{
			ShortUrl:    u.ShortURL,
			OriginalUrl: u.OriginalURL,
		})
	}

	return resp, nil
}
