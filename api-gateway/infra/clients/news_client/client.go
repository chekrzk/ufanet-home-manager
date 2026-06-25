package news_client

import (
	"context"
	"time"

	"github.com/chekrzk/ufanet-home-manager/api-gateway/internal/models/domain"
	commonv1 "github.com/chekrzk/ufanet-home-manager/contracts/gen/go/common/v1"
	newsv1 "github.com/chekrzk/ufanet-home-manager/contracts/gen/go/news/v1"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Client struct {
	client newsv1.NewsServiceClient
	log    zerolog.Logger
}

func New(conn *grpc.ClientConn, log zerolog.Logger) *Client {
	return &Client{client: newsv1.NewNewsServiceClient(conn), log: log}
}

func (c *Client) List(ctx context.Context, actor domain.AuthContext, filter domain.NewsFilter) (domain.Page[domain.News], error) {
	c.log.Debug().Str("user_id", actor.UserID).Msg("call news grpc list")
	resp, err := c.client.ListNews(ctx, &newsv1.ListNewsRequest{
		User:       userContext(actor),
		Pagination: paginationToProto(filter.Pagination),
		DateFrom:   filter.DateFrom,
		DateTo:     filter.DateTo,
	})
	if err != nil {
		return domain.Page[domain.News]{}, err
	}

	items := make([]domain.News, 0, len(resp.GetItems()))
	for _, item := range resp.GetItems() {
		items = append(items, newsFromProto(item))
	}

	return domain.Page[domain.News]{
		Items: items,
		Page:  int(resp.GetPage()),
		Limit: int(resp.GetLimit()),
		Total: int(resp.GetTotal()),
	}, nil
}

func (c *Client) Create(ctx context.Context, author domain.AuthContext, command domain.CreateNews) (domain.News, error) {
	c.log.Debug().Str("user_id", author.UserID).Msg("call news grpc create")
	resp, err := c.client.CreateNews(ctx, &newsv1.CreateNewsRequest{
		Author:  userContext(author),
		Title:   command.Title,
		Body:    command.Body,
		HouseId: command.HouseID,
	})
	if err != nil {
		return domain.News{}, err
	}
	return newsFromProto(resp), nil
}

func userContext(actor domain.AuthContext) *commonv1.UserContext {
	return &commonv1.UserContext{UserId: actor.UserID, Role: actor.Role}
}

func paginationToProto(page domain.Pagination) *commonv1.Pagination {
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
