package server

import (
	"context"
	stderrors "errors"

	commonv1 "github.com/chekrzk/ufanet-home-manager/contracts/gen/go/common/v1"
	requestsv1 "github.com/chekrzk/ufanet-home-manager/contracts/gen/go/requests/v1"
	apperrors "github.com/chekrzk/ufanet-home-manager/requests-service/internal/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	requestsv1.UnimplementedRequestsServiceServer
	service RequestsService
}

// New связывает gRPC transport с request scenarios, не раскрывая server слою DB
// и notification dependencies.
func New(service RequestsService) *Server {
	return &Server{service: service}
}

// CreateRequest переводит proto в доменную команду, чтобы владелец и назначение
// заявки обрабатывались внутри service layer.
func (s *Server) CreateRequest(ctx context.Context, req *requestsv1.CreateRequestRequest) (*commonv1.MaintenanceRequest, error) {
	request, err := s.service.Create(ctx, createRequestCommandFromProto(req))
	if err != nil {
		return nil, grpcError(err)
	}
	return requestToProto(request), nil
}

// ListRequests передает actor и pagination в service, чтобы видимость заявок
// определялась ролью, а не transport-кодом.
func (s *Server) ListRequests(ctx context.Context, req *requestsv1.ListRequestsRequest) (*requestsv1.ListRequestsResponse, error) {
	page, err := s.service.List(ctx, listRequestsFilterFromProto(req))
	if err != nil {
		return nil, grpcError(err)
	}
	items := make([]*commonv1.MaintenanceRequest, 0, len(page.Items))
	for _, item := range page.Items {
		items = append(items, requestToProto(item))
	}
	return &requestsv1.ListRequestsResponse{Items: items, Page: int32(page.Page), Limit: int32(page.Limit), Total: int32(page.Total)}, nil
}

// GetRequest оставляет проверку доступа к конкретной заявке доменному сценарию.
func (s *Server) GetRequest(ctx context.Context, req *requestsv1.GetRequestRequest) (*commonv1.MaintenanceRequest, error) {
	request, err := s.service.Get(ctx, getRequestCommandFromProto(req))
	if err != nil {
		return nil, grpcError(err)
	}
	return requestToProto(request), nil
}

// UpdateRequestStatus проводит изменение статуса через service, чтобы timestamps
// и уведомления создавались вместе с переходом состояния.
func (s *Server) UpdateRequestStatus(ctx context.Context, req *requestsv1.UpdateRequestStatusRequest) (*commonv1.MaintenanceRequest, error) {
	request, err := s.service.UpdateStatus(ctx, updateRequestStatusCommandFromProto(req))
	if err != nil {
		return nil, grpcError(err)
	}
	return requestToProto(request), nil
}

// AddRequestComment передает actor в service, чтобы комментарии не обходили
// правила доступа к заявке.
func (s *Server) AddRequestComment(ctx context.Context, req *requestsv1.AddRequestCommentRequest) (*commonv1.Empty, error) {
	err := s.service.AddComment(ctx, addRequestCommentCommandFromProto(req))
	if err != nil {
		return nil, grpcError(err)
	}
	return &commonv1.Empty{}, nil
}

// grpcError переводит domain errors заявок в gRPC status для gateway.
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
