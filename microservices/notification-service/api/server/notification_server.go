package server

import (
	"context"
	stderrors "errors"

	commonv1 "github.com/chekrzk/ufanet-home-manager/contracts/gen/go/common/v1"
	notificationsv1 "github.com/chekrzk/ufanet-home-manager/contracts/gen/go/notifications/v1"
	apperrors "github.com/chekrzk/ufanet-home-manager/notification-service/internal/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	notificationsv1.UnimplementedNotificationsServiceServer
	service NotificationService
}

// New подключает gRPC transport к notification scenarios без знания DB или Redis.
func New(service NotificationService) *Server {
	return &Server{service: service}
}

// RegisterDevice переводит transport request в команду привязки устройства.
func (s *Server) RegisterDevice(ctx context.Context, req *notificationsv1.RegisterDeviceRequest) (*commonv1.Empty, error) {
	err := s.service.RegisterDevice(ctx, registerDeviceCommandFromProto(req))
	if err != nil {
		return nil, grpcError(err)
	}
	return &commonv1.Empty{}, nil
}

// UnregisterDevice держит отвязку устройства в service layer, а не в gRPC handler.
func (s *Server) UnregisterDevice(ctx context.Context, req *notificationsv1.UnregisterDeviceRequest) (*commonv1.Empty, error) {
	err := s.service.UnregisterDevice(ctx, unregisterDeviceCommandFromProto(req))
	if err != nil {
		return nil, grpcError(err)
	}
	return &commonv1.Empty{}, nil
}

// Publish нужен для внутренних сервисов, чтобы они создавали уведомления через
// один контракт, не зная о Redis Stream и таблицах notification-service.
func (s *Server) Publish(ctx context.Context, req *notificationsv1.PublishNotificationRequest) (*commonv1.Empty, error) {
	err := s.service.Publish(ctx, publishNotificationCommandFromProto(req))
	if err != nil {
		return nil, grpcError(err)
	}
	return &commonv1.Empty{}, nil
}

// ListNotifications возвращает историю уведомлений через service rules владельца.
func (s *Server) ListNotifications(ctx context.Context, req *notificationsv1.ListNotificationsRequest) (*notificationsv1.ListNotificationsResponse, error) {
	page, err := s.service.List(ctx, listNotificationsCommandFromProto(req))
	if err != nil {
		return nil, grpcError(err)
	}
	items := make([]*commonv1.Notification, 0, len(page.Items))
	for _, item := range page.Items {
		items = append(items, notificationToProto(item))
	}
	return &notificationsv1.ListNotificationsResponse{
		Items: items,
		Page:  int32(page.Page),
		Limit: int32(page.Limit),
		Total: int32(page.Total),
	}, nil
}

// MarkRead оставляет проверку владельца уведомления в service layer.
func (s *Server) MarkRead(ctx context.Context, req *notificationsv1.MarkReadRequest) (*commonv1.Empty, error) {
	if err := s.service.MarkRead(ctx, userContextFromProto(req.GetUser()), req.GetNotificationId()); err != nil {
		return nil, grpcError(err)
	}
	return &commonv1.Empty{}, nil
}

// grpcError переводит доменные ошибки уведомлений в стабильные gRPC codes.
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
