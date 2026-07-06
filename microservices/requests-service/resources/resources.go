package resources

import (
	"context"

	"github.com/chekrzk/ufanet-home-manager/requests-service/resources/config"
	"github.com/chekrzk/ufanet-home-manager/requests-service/resources/logger"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Resources struct {
	Config            *config.Config
	Log               zerolog.Logger
	DB                *gorm.DB
	NotificationsConn *grpc.ClientConn
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

	notificationsConn, err := grpc.NewClient(cfg.Service.NotificationsAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &Resources{
		Config:            cfg,
		Log:               log,
		DB:                db,
		NotificationsConn: notificationsConn,
	}, nil
}

func (r *Resources) Close() error {
	if r == nil {
		return nil
	}
	if r.NotificationsConn != nil {
		_ = r.NotificationsConn.Close()
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
