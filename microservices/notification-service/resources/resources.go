package resources

import (
	"context"

	"github.com/chekrzk/ufanet-home-manager/notification-service/resources/config"
	"github.com/chekrzk/ufanet-home-manager/notification-service/resources/logger"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Resources struct {
	Config *config.Config
	Log    zerolog.Logger
	DB     *gorm.DB
	Redis  *redis.Client
}

func New(ctx context.Context) (*Resources, error) {
	_ = ctx

	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}
	log := logger.New(cfg)

	db, err := gorm.Open(postgres.Open(cfg.DB.DSN), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	redisClient := redis.NewClient(&redis.Options{Addr: cfg.Redis.Addr})

	return &Resources{
		Config: cfg,
		Log:    log,
		DB:     db,
		Redis:  redisClient,
	}, nil
}

func (r *Resources) Close() error {
	if r == nil {
		return nil
	}
	if r.Redis != nil {
		_ = r.Redis.Close()
	}
	if r.DB == nil {
		return nil
	}
	db, err := r.DB.DB()
	if err != nil {
		return err
	}
	return db.Close()
}
