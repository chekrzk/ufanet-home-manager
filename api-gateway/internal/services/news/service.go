package news

import (
	"context"
	"time"

	"github.com/chekrzk/ufanet-home-manager/api-gateway/internal/models/domain"
	"github.com/chekrzk/ufanet-home-manager/api-gateway/internal/models/dto"
	commonv1 "github.com/chekrzk/ufanet-home-manager/contracts/gen/go/common/v1"
	newsv1 "github.com/chekrzk/ufanet-home-manager/contracts/gen/go/news/v1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type service struct {
	client newsv1.NewsServiceClient
}

func New(conn *grpc.ClientConn) Service {
	return service{client: newsv1.NewNewsServiceClient(conn)}
}

func (s service) List(ctx context.Context, actor domain.AuthContext, req dto.ListNewsRequest) (dto.Page[domain.News], error) {
	resp, err := s.client.ListNews(ctx, &newsv1.ListNewsRequest{
		User:       userContext(actor),
		Pagination: paginationToProto(req.Pagination),
		DateFrom:   req.DateFrom,
		DateTo:     req.DateTo,
	})
	if err != nil {
		return dto.Page[domain.News]{}, err
	}

	items := make([]domain.News, 0, len(resp.GetItems()))
	for _, item := range resp.GetItems() {
		items = append(items, newsFromProto(item))
	}

	return dto.Page[domain.News]{
		Items: items,
		Page:  int(resp.GetPage()),
		Limit: int(resp.GetLimit()),
		Total: int(resp.GetTotal()),
	}, nil
}

func (s service) Create(ctx context.Context, author domain.AuthContext, req dto.CreateNewsRequest) (domain.News, error) {
	resp, err := s.client.CreateNews(ctx, &newsv1.CreateNewsRequest{
		Author:  userContext(author),
		Title:   req.Title,
		Body:    req.Body,
		HouseId: req.HouseID,
	})
	if err != nil {
		return domain.News{}, err
	}
	return newsFromProto(resp), nil
}

func userContext(actor domain.AuthContext) *commonv1.UserContext {
	return &commonv1.UserContext{UserId: actor.UserID, Role: actor.Role}
}

func paginationToProto(page dto.Pagination) *commonv1.Pagination {
	page.Normalize()
	return &commonv1.Pagination{Page: int32(page.Page), Limit: int32(page.Limit)}
}

func newsFromProto(news *commonv1.News) domain.News {
	if news == nil {
		return domain.News{}
	}
	return domain.News{
		ID:        news.GetId(),
		Title:     news.GetTitle(),
		Body:      news.GetBody(),
		HouseID:   news.GetHouseId(),
		CreatedAt: timeFromProto(news.GetCreatedAt()),
	}
}

func timeFromProto(ts *timestamppb.Timestamp) time.Time {
	if ts == nil {
		return time.Time{}
	}
	return ts.AsTime()
}
