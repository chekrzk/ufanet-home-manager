package notifications_client

import (
	"context"
	"time"

	"github.com/chekrzk/ufanet-home-manager/api-gateway/internal/models/domain"
	commonv1 "github.com/chekrzk/ufanet-home-manager/contracts/gen/go/common/v1"
	notificationsv1 "github.com/chekrzk/ufanet-home-manager/contracts/gen/go/notifications/v1"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Client struct {
	client notificationsv1.NotificationsServiceClient
	log    zerolog.Logger
}

func New(conn *grpc.ClientConn, log zerolog.Logger) *Client {
	return &Client{client: notificationsv1.NewNotificationsServiceClient(conn), log: log}
}

func (c *Client) RegisterDevice(ctx context.Context, actor domain.AuthContext, device domain.RegisterDevice) error {
	c.log.Debug().Str("user_id", actor.UserID).Msg("call notifications grpc register device")
	_, err := c.client.RegisterDevice(ctx, &notificationsv1.RegisterDeviceRequest{
		User:     userContext(actor),
		Token:    device.Token,
		Platform: device.Platform,
	})
	return err
}

func (c *Client) UnregisterDevice(ctx context.Context, actor domain.AuthContext, device domain.UnregisterDevice) error {
	c.log.Debug().Str("user_id", actor.UserID).Msg("call notifications grpc unregister device")
	_, err := c.client.UnregisterDevice(ctx, &notificationsv1.UnregisterDeviceRequest{
		User:  userContext(actor),
		Token: device.Token,
	})
	return err
}

func (c *Client) List(ctx context.Context, actor domain.AuthContext, page domain.Pagination) (domain.Page[domain.Notification], error) {
	c.log.Debug().Str("user_id", actor.UserID).Msg("call notifications grpc list")
	resp, err := c.client.ListNotifications(ctx, &notificationsv1.ListNotificationsRequest{
		User:       userContext(actor),
		Pagination: paginationToProto(page),
	})
	if err != nil {
		return domain.Page[domain.Notification]{}, err
	}

	items := make([]domain.Notification, 0, len(resp.GetItems()))
	for _, item := range resp.GetItems() {
		items = append(items, notificationFromProto(item))
	}

	return domain.Page[domain.Notification]{
		Items: items,
		Page:  int(resp.GetPage()),
		Limit: int(resp.GetLimit()),
		Total: int(resp.GetTotal()),
	}, nil
}

func (c *Client) MarkRead(ctx context.Context, actor domain.AuthContext, notificationID string) error {
	c.log.Debug().Str("user_id", actor.UserID).Str("notification_id", notificationID).Msg("call notifications grpc mark read")
	_, err := c.client.MarkRead(ctx, &notificationsv1.MarkReadRequest{
		User:           userContext(actor),
		NotificationId: notificationID,
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

func notificationFromProto(item *commonv1.Notification) domain.Notification {
	if item == nil {
		return domain.Notification{}
	}
	return domain.Notification{
		ID:        item.GetId(),
		UserID:    item.GetUserId(),
		HouseID:   item.GetHouseId(),
		Type:      item.GetType(),
		Title:     item.GetTitle(),
		Body:      item.GetBody(),
		EntityID:  item.GetEntityId(),
		Read:      item.GetRead(),
		CreatedAt: timeFromProto(item.GetCreatedAt()),
	}
}

func timeFromProto(ts *timestamppb.Timestamp) time.Time {
	if ts == nil {
		return time.Time{}
	}
	return ts.AsTime()
}
