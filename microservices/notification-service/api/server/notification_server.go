package server

import (
	"context"
	stderrors "errors"

	commonv1 "github.com/chekrzk/ufanet-home-manager/contracts/gen/go/common/v1"
	notificationsv1 "github.com/chekrzk/ufanet-home-manager/contracts/gen/go/notifications/v1"
	apperrors "github.com/chekrzk/ufanet-home-manager/notification-service/internal/errors"
	"github.com/chekrzk/ufanet-home-manager/notification-service/internal/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	notificationsv1.UnimplementedNotificationsServiceServer
	service NotificationService
}

func New(service NotificationService) *Server {
	return &Server{service: service}
}

func (s *Server) RegisterDevice(ctx context.Context, req *notificationsv1.RegisterDeviceRequest) (*commonv1.Empty, error) {
	err := s.service.RegisterDevice(ctx, models.RegisterDeviceCommand{
		User:     userContext(req.GetUser()),
		Token:    req.GetToken(),
		Platform: req.GetPlatform(),
	})
	if err != nil {
		return nil, grpcError(err)
	}
	return &commonv1.Empty{}, nil
}

func (s *Server) UnregisterDevice(ctx context.Context, req *notificationsv1.UnregisterDeviceRequest) (*commonv1.Empty, error) {
	err := s.service.UnregisterDevice(ctx, models.UnregisterDeviceCommand{
		User:  userContext(req.GetUser()),
		Token: req.GetToken(),
	})
	if err != nil {
		return nil, grpcError(err)
	}
	return &commonv1.Empty{}, nil
}

func (s *Server) Publish(ctx context.Context, req *notificationsv1.PublishNotificationRequest) (*commonv1.Empty, error) {
	err := s.service.Publish(ctx, models.PublishNotificationCommand{
		UserID:   req.GetUserId(),
		HouseID:  req.GetHouseId(),
		Type:     req.GetType(),
		Title:    req.GetTitle(),
		Body:     req.GetBody(),
		EntityID: req.GetEntityId(),
	})
	if err != nil {
		return nil, grpcError(err)
	}
	return &commonv1.Empty{}, nil
}

func (s *Server) ListNotifications(ctx context.Context, req *notificationsv1.ListNotificationsRequest) (*notificationsv1.ListNotificationsResponse, error) {
	page, err := s.service.List(ctx, models.ListNotificationsCommand{
		User: userContext(req.GetUser()),
		Pagination: models.Pagination{
			Page:  int(req.GetPagination().GetPage()),
			Limit: int(req.GetPagination().GetLimit()),
		},
	})
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

func (s *Server) MarkRead(ctx context.Context, req *notificationsv1.MarkReadRequest) (*commonv1.Empty, error) {
	if err := s.service.MarkRead(ctx, userContext(req.GetUser()), req.GetNotificationId()); err != nil {
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

func notificationToProto(notification models.Notification) *commonv1.Notification {
	return &commonv1.Notification{
		Id:        notification.ID,
		UserId:    notification.UserID,
		HouseId:   notification.HouseID,
		Type:      notification.Type,
		Title:     notification.Title,
		Body:      notification.Body,
		EntityId:  notification.EntityID,
		Read:      notification.Read,
		CreatedAt: timestamppb.New(notification.CreatedAt),
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
