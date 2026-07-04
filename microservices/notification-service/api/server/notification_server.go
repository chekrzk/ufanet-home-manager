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

func userContext(user *commonv1.UserContext) models.UserContext {
	if user == nil {
		return models.UserContext{}
	}
	return models.UserContext{UserID: user.GetUserId(), Role: user.GetRole()}
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
