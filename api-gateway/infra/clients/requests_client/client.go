package requests_client

import (
	"context"
	"time"

	"github.com/chekrzk/ufanet-home-manager/api-gateway/internal/models/domain"
	commonv1 "github.com/chekrzk/ufanet-home-manager/contracts/gen/go/common/v1"
	requestsv1 "github.com/chekrzk/ufanet-home-manager/contracts/gen/go/requests/v1"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Client struct {
	client requestsv1.RequestsServiceClient
	log    zerolog.Logger
}

func New(conn *grpc.ClientConn, log zerolog.Logger) *Client {
	return &Client{client: requestsv1.NewRequestsServiceClient(conn), log: log}
}

func (c *Client) Create(ctx context.Context, actor domain.AuthContext, command domain.CreateRequest) (domain.Request, error) {
	c.log.Debug().Str("user_id", actor.UserID).Msg("call requests grpc create")
	resp, err := c.client.CreateRequest(ctx, &requestsv1.CreateRequestRequest{
		User:        userContext(actor),
		Category:    command.Category,
		Description: command.Description,
	})
	if err != nil {
		return domain.Request{}, err
	}
	return requestFromProto(resp), nil
}

func (c *Client) List(ctx context.Context, actor domain.AuthContext, page domain.Pagination) (domain.Page[domain.Request], error) {
	c.log.Debug().Str("user_id", actor.UserID).Msg("call requests grpc list")
	resp, err := c.client.ListRequests(ctx, &requestsv1.ListRequestsRequest{
		User:       userContext(actor),
		Pagination: paginationToProto(page),
	})
	if err != nil {
		return domain.Page[domain.Request]{}, err
	}

	items := make([]domain.Request, 0, len(resp.GetItems()))
	for _, item := range resp.GetItems() {
		items = append(items, requestFromProto(item))
	}

	return domain.Page[domain.Request]{
		Items: items,
		Page:  int(resp.GetPage()),
		Limit: int(resp.GetLimit()),
		Total: int(resp.GetTotal()),
	}, nil
}

func (c *Client) Get(ctx context.Context, actor domain.AuthContext, requestID string) (domain.Request, error) {
	c.log.Debug().Str("user_id", actor.UserID).Str("request_id", requestID).Msg("call requests grpc get")
	resp, err := c.client.GetRequest(ctx, &requestsv1.GetRequestRequest{
		User:      userContext(actor),
		RequestId: requestID,
	})
	if err != nil {
		return domain.Request{}, err
	}
	return requestFromProto(resp), nil
}

func (c *Client) UpdateStatus(ctx context.Context, actor domain.AuthContext, requestID string, command domain.UpdateRequestStatus) (domain.Request, error) {
	c.log.Debug().Str("user_id", actor.UserID).Str("request_id", requestID).Msg("call requests grpc update status")
	resp, err := c.client.UpdateRequestStatus(ctx, &requestsv1.UpdateRequestStatusRequest{
		Actor:      userContext(actor),
		RequestId:  requestID,
		Status:     command.Status,
		AssignedTo: command.AssignedTo,
	})
	if err != nil {
		return domain.Request{}, err
	}
	return requestFromProto(resp), nil
}

func (c *Client) AddComment(ctx context.Context, actor domain.AuthContext, requestID string, command domain.AddRequestComment) error {
	c.log.Debug().Str("user_id", actor.UserID).Str("request_id", requestID).Msg("call requests grpc add comment")
	_, err := c.client.AddRequestComment(ctx, &requestsv1.AddRequestCommentRequest{
		User:      userContext(actor),
		RequestId: requestID,
		Text:      command.Text,
	})
	return err
}

func userContext(actor domain.AuthContext) *commonv1.UserContext {
	return &commonv1.UserContext{UserId: actor.UserID, Role: actor.Role}
}

func paginationToProto(page domain.Pagination) *commonv1.Pagination {
	page.Normalize()
	return &commonv1.Pagination{Page: int32(page.Page), Limit: int32(page.Limit)}
}

func requestFromProto(request *commonv1.MaintenanceRequest) domain.Request {
	if request == nil {
		return domain.Request{}
	}
	return domain.Request{
		ID:          request.GetId(),
		UserID:      request.GetUserId(),
		Category:    request.GetCategory(),
		Description: request.GetDescription(),
		Status:      request.GetStatus(),
		AssignedTo:  request.GetAssignedTo(),
		CreatedAt:   timeFromProto(request.GetCreatedAt()),
		UpdatedAt:   timeFromProto(request.GetUpdatedAt()),
	}
}

func timeFromProto(ts *timestamppb.Timestamp) time.Time {
	if ts == nil {
		return time.Time{}
	}
	return ts.AsTime()
}
