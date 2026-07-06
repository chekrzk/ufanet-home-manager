package server

import (
	commonv1 "github.com/chekrzk/ufanet-home-manager/contracts/gen/go/common/v1"
	notificationsv1 "github.com/chekrzk/ufanet-home-manager/contracts/gen/go/notifications/v1"
	"github.com/chekrzk/ufanet-home-manager/notification-service/internal/models"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// userContextFromProto переносит trusted actor context из gateway в domain.
func userContextFromProto(user *commonv1.UserContext) models.UserContext {
	if user == nil {
		return models.UserContext{}
	}
	return models.UserContext{UserID: user.GetUserId(), Role: user.GetRole()}
}

// registerDeviceCommandFromProto изолирует proto device payload от service layer.
func registerDeviceCommandFromProto(req *notificationsv1.RegisterDeviceRequest) models.RegisterDeviceCommand {
	return models.RegisterDeviceCommand{
		User:     userContextFromProto(req.GetUser()),
		Token:    req.GetToken(),
		Platform: req.GetPlatform(),
	}
}

// unregisterDeviceCommandFromProto строит команду отвязки устройства из контракта.
func unregisterDeviceCommandFromProto(req *notificationsv1.UnregisterDeviceRequest) models.UnregisterDeviceCommand {
	return models.UnregisterDeviceCommand{
		User:  userContextFromProto(req.GetUser()),
		Token: req.GetToken(),
	}
}

// publishNotificationCommandFromProto делает внутреннюю публикацию transport-agnostic.
func publishNotificationCommandFromProto(req *notificationsv1.PublishNotificationRequest) models.PublishNotificationCommand {
	return models.PublishNotificationCommand{
		UserID:   req.GetUserId(),
		HouseID:  req.GetHouseId(),
		Type:     req.GetType(),
		Title:    req.GetTitle(),
		Body:     req.GetBody(),
		EntityID: req.GetEntityId(),
	}
}

// listNotificationsCommandFromProto переносит actor и pagination в domain command.
func listNotificationsCommandFromProto(req *notificationsv1.ListNotificationsRequest) models.ListNotificationsCommand {
	return models.ListNotificationsCommand{
		User: userContextFromProto(req.GetUser()),
		Pagination: models.Pagination{
			Page:  int(req.GetPagination().GetPage()),
			Limit: int(req.GetPagination().GetLimit()),
		},
	}
}

// notificationToProto скрывает storage-модель уведомления за общим контрактом.
func notificationToProto(notification models.Notification) *commonv1.Notification {
	userID := ""
	if notification.UserID != nil {
		userID = *notification.UserID
	}
	return &commonv1.Notification{
		Id:        notification.ID,
		UserId:    userID,
		HouseId:   notification.HouseID,
		Type:      notification.Type,
		Title:     notification.Title,
		Body:      notification.Body,
		EntityId:  notification.EntityID,
		Read:      notification.Read,
		CreatedAt: timestamppb.New(notification.CreatedAt),
	}
}
