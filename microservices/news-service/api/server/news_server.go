package server

import (
	"context"
	stderrors "errors"

	commonv1 "github.com/chekrzk/ufanet-home-manager/contracts/gen/go/common/v1"
	newsv1 "github.com/chekrzk/ufanet-home-manager/contracts/gen/go/news/v1"
	apperrors "github.com/chekrzk/ufanet-home-manager/news-service/internal/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	newsv1.UnimplementedNewsServiceServer
	service NewsService
}

// New связывает gRPC endpoint с интерфейсом сценариев новостей, чтобы transport
// не зависел от repository и notification client.
func New(service NewsService) *Server {
	return &Server{service: service}
}

// ListNews конвертирует proto-фильтр в доменную модель, чтобы service layer
// отвечал за правила ленты, а не за формат контракта.
func (s *Server) ListNews(ctx context.Context, req *newsv1.ListNewsRequest) (*newsv1.ListNewsResponse, error) {
	page, err := s.service.List(ctx, newsFilterFromProto(req))
	if err != nil {
		return nil, grpcError(err)
	}
	items := make([]*commonv1.News, 0, len(page.Items))
	for _, item := range page.Items {
		items = append(items, newsToProto(item))
	}
	return &newsv1.ListNewsResponse{Items: items, Page: int32(page.Page), Limit: int32(page.Limit), Total: int32(page.Total)}, nil
}

// CreateNews оставляет публикацию и side effects сервисному слою, а server
// только связывает gRPC request/response с доменной командой.
func (s *Server) CreateNews(ctx context.Context, req *newsv1.CreateNewsRequest) (*commonv1.News, error) {
	item, err := s.service.Create(ctx, createNewsCommandFromProto(req))
	if err != nil {
		return nil, grpcError(err)
	}
	return newsToProto(item), nil
}

// grpcError сохраняет единый контракт ошибок между news-service и gateway.
func grpcError(err error) error {
	switch {
	case stderrors.Is(err, apperrors.ErrInvalidArgument):
		return status.Error(codes.InvalidArgument, err.Error())
	case stderrors.Is(err, apperrors.ErrForbidden):
		return status.Error(codes.PermissionDenied, err.Error())
	case stderrors.Is(err, apperrors.ErrNotFound):
		return status.Error(codes.NotFound, err.Error())
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}
