package redis_stream

import (
	"context"

	"github.com/chekrzk/ufanet-home-manager/notification-service/internal/models"
	"github.com/redis/go-redis/v9"
)

type Publisher struct {
	client *redis.Client
	stream string
}

func New(client *redis.Client, stream string) *Publisher {
	return &Publisher{client: client, stream: stream}
}

func (p *Publisher) Publish(ctx context.Context, event models.PublishNotificationCommand) error {
	return p.client.XAdd(ctx, &redis.XAddArgs{
		Stream: p.stream,
		Values: map[string]any{
			"user_id":   event.UserID,
			"house_id":  event.HouseID,
			"type":      event.Type,
			"title":     event.Title,
			"body":      event.Body,
			"entity_id": event.EntityID,
		},
	}).Err()
}
