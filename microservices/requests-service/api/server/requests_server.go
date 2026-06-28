package server

import (
	"context"
	stderrors "errors"

	commonv1 "github.com/chekrzk/ufanet-home-manager/contracts/gen/go/common/v1"
	requestsv1 "github.com/chekrzk/ufanet-home-manager/contracts/gen/go/requests/v1"
	apperrors "github.com/chekrzk/ufanet-home-manager/requests-service/internal/errors"
	"github.com/chekrzk/ufanet-home-manager/requests-service/internal/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	requestsv1.UnimplementedRequestsServiceServer
	service RequestsService
}

func New(service RequestsService) *Server {
	return &Server{service: service}
}

func (s *Server) CreateRequest(ctx context.Context, req *requestsv1.CreateRequestRequest) (*commonv1.MaintenanceRequest, error) {
	request, err := s.service.Create(ctx, models.CreateRequestCommand{
		User:        userContext(req.GetUser()),
		Category:    req.GetCategory(),
		Description: req.GetDescription(),
	})
	if err != nil {
		return nil, grpcError(err)
	}
	return requestToProto(request), nil
}

func (s *Server) ListRequests(ctx context.Context, req *requestsv1.ListRequestsRequest) (*requestsv1.ListRequestsResponse, error) {
	page, err := s.service.List(ctx, models.ListRequestsFilter{
		Actor: userContext(req.GetUser()),
		Pagination: models.Pagination{
			Page:  int(req.GetPagination().GetPage()),
			Limit: int(req.GetPagination().GetLimit()),
		},
	})
	if err != nil {
		return nil, grpcError(err)
	}
	items := make([]*commonv1.MaintenanceRequest, 0, len(page.Items))
	for _, item := range page.Items {
		items = append(items, requestToProto(item))
	}
	return &requestsv1.ListRequestsResponse{Items: items, Page: int32(page.Page), Limit: int32(page.Limit), Total: int32(page.Total)}, nil
}

func (s *Server) GetRequest(ctx context.Context, req *requestsv1.GetRequestRequest) (*commonv1.MaintenanceRequest, error) {
	request, err := s.service.Get(ctx, models.GetRequestCommand{Actor: userContext(req.GetUser()), RequestID: req.GetRequestId()})
	if err != nil {
		return nil, grpcError(err)
	}
	return requestToProto(request), nil
}

func (s *Server) UpdateRequestStatus(ctx context.Context, req *requestsv1.UpdateRequestStatusRequest) (*commonv1.MaintenanceRequest, error) {
	request, err := s.service.UpdateStatus(ctx, models.UpdateRequestStatusCommand{
		Actor:      userContext(req.GetActor()),
		RequestID:  req.GetRequestId(),
		Status:     req.GetStatus(),
		AssignedTo: req.GetAssignedTo(),
	})
	if err != nil {
		return nil, grpcError(err)
	}
	return requestToProto(request), nil
}

func (s *Server) AddRequestComment(ctx context.Context, req *requestsv1.AddRequestCommentRequest) (*commonv1.Empty, error) {
	err := s.service.AddComment(ctx, models.AddRequestCommentCommand{
		Actor:     userContext(req.GetUser()),
		RequestID: req.GetRequestId(),
		Text:      req.GetText(),
	})
	if err != nil {
		return nil, grpcError(err)
	}
	return &commonv1.Empty{}, nil
}

func userContext(user *commonv1.UserContext) models.UserContext {
	if user == nil {
		return models.UserContext{}
	}
	return models.UserContext{UserID: user.GetUserId(), Role: user.GetRole()}
}

func requestToProto(request models.MaintenanceRequest) *commonv1.MaintenanceRequest {
	return &commonv1.MaintenanceRequest{
		Id:          request.ID,
		UserId:      request.UserID,
		Category:    request.Category,
		Description: request.Description,
		Status:      request.Status,
		CreatedAt:   timestamppb.New(request.CreatedAt),
		UpdatedAt:   timestamppb.New(request.UpdatedAt),
		AssignedTo:  request.AssignedTo,
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
