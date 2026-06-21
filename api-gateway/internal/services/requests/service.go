package requests

import (
	"context"
	"time"

	"github.com/chekrzk/ufanet-home-manager/api-gateway/internal/models/domain"
	"github.com/chekrzk/ufanet-home-manager/api-gateway/internal/models/dto"
	commonv1 "github.com/chekrzk/ufanet-home-manager/contracts/gen/go/common/v1"
	requestsv1 "github.com/chekrzk/ufanet-home-manager/contracts/gen/go/requests/v1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type service struct {
	client requestsv1.RequestsServiceClient
}

func New(conn *grpc.ClientConn) Service {
	return service{client: requestsv1.NewRequestsServiceClient(conn)}
}

func (s service) Create(ctx context.Context, actor domain.AuthContext, req dto.CreateRequestRequest) (domain.Request, error) {
	resp, err := s.client.CreateRequest(ctx, &requestsv1.CreateRequestRequest{
		User:        userContext(actor),
		Category:    req.Category,
		Description: req.Description,
	})
	if err != nil {
		return domain.Request{}, err
	}
	return requestFromProto(resp), nil
}

func (s service) List(ctx context.Context, actor domain.AuthContext, page dto.Pagination) (dto.Page[domain.Request], error) {
	resp, err := s.client.ListRequests(ctx, &requestsv1.ListRequestsRequest{
		User:       userContext(actor),
		Pagination: paginationToProto(page),
	})
	if err != nil {
		return dto.Page[domain.Request]{}, err
	}

	items := make([]domain.Request, 0, len(resp.GetItems()))
	for _, item := range resp.GetItems() {
		items = append(items, requestFromProto(item))
	}

	return dto.Page[domain.Request]{
		Items: items,
		Page:  int(resp.GetPage()),
		Limit: int(resp.GetLimit()),
		Total: int(resp.GetTotal()),
	}, nil
}

func (s service) Get(ctx context.Context, actor domain.AuthContext, requestID string) (domain.Request, error) {
	resp, err := s.client.GetRequest(ctx, &requestsv1.GetRequestRequest{
		User:      userContext(actor),
		RequestId: requestID,
	})
	if err != nil {
		return domain.Request{}, err
	}
	return requestFromProto(resp), nil
}

func (s service) UpdateStatus(ctx context.Context, actor domain.AuthContext, requestID string, req dto.UpdateRequestStatusRequest) (domain.Request, error) {
	resp, err := s.client.UpdateRequestStatus(ctx, &requestsv1.UpdateRequestStatusRequest{
		Actor:     userContext(actor),
		RequestId: requestID,
		Status:    req.Status,
	})
	if err != nil {
		return domain.Request{}, err
	}
	return requestFromProto(resp), nil
}

func (s service) AddComment(ctx context.Context, actor domain.AuthContext, requestID string, req dto.AddRequestCommentRequest) error {
	_, err := s.client.AddRequestComment(ctx, &requestsv1.AddRequestCommentRequest{
		User:      userContext(actor),
		RequestId: requestID,
		Text:      req.Text,
	})
	return err
}

func userContext(actor domain.AuthContext) *commonv1.UserContext {
	return &commonv1.UserContext{UserId: actor.UserID, Role: actor.Role}
}

func paginationToProto(page dto.Pagination) *commonv1.Pagination {
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
