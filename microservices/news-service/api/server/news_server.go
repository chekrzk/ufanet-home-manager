package server

import (
	"context"
	stderrors "errors"

	commonv1 "github.com/chekrzk/ufanet-home-manager/contracts/gen/go/common/v1"
	newsv1 "github.com/chekrzk/ufanet-home-manager/contracts/gen/go/news/v1"
	apperrors "github.com/chekrzk/ufanet-home-manager/news-service/internal/errors"
	"github.com/chekrzk/ufanet-home-manager/news-service/internal/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	newsv1.UnimplementedNewsServiceServer
	service NewsService
}

func New(service NewsService) *Server {
	return &Server{service: service}
}

func (s *Server) ListNews(ctx context.Context, req *newsv1.ListNewsRequest) (*newsv1.ListNewsResponse, error) {
	page, err := s.service.List(ctx, models.NewsFilter{
		Actor: userContext(req.GetUser()),
		Pagination: models.Pagination{
			Page:  int(req.GetPagination().GetPage()),
			Limit: int(req.GetPagination().GetLimit()),
		},
		DateFrom: req.GetDateFrom(),
		DateTo:   req.GetDateTo(),
	})
	if err != nil {
		return nil, grpcError(err)
	}
	items := make([]*commonv1.News, 0, len(page.Items))
	for _, item := range page.Items {
		items = append(items, newsToProto(item))
	}
	return &newsv1.ListNewsResponse{Items: items, Page: int32(page.Page), Limit: int32(page.Limit), Total: int32(page.Total)}, nil
}

func (s *Server) CreateNews(ctx context.Context, req *newsv1.CreateNewsRequest) (*commonv1.News, error) {
	item, err := s.service.Create(ctx, models.CreateNewsCommand{
		Author:  userContext(req.GetAuthor()),
		Title:   req.GetTitle(),
		Body:    req.GetBody(),
		HouseID: req.GetHouseId(),
	})
	if err != nil {
		return nil, grpcError(err)
	}
	return newsToProto(item), nil
}

func userContext(user *commonv1.UserContext) models.UserContext {
	if user == nil {
		return models.UserContext{}
	}
	return models.UserContext{UserID: user.GetUserId(), Role: user.GetRole()}
}

func newsToProto(item models.News) *commonv1.News {
	return &commonv1.News{
		Id:        item.ID,
		Title:     item.Title,
		Body:      item.Body,
		HouseId:   item.HouseID,
		CreatedAt: timestamppb.New(item.CreatedAt),
	}
}

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
