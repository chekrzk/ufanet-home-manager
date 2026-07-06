package resources

import (
	"context"

	"github.com/chekrzk/ufanet-home-manager/auth-service/resources/config"
	"github.com/chekrzk/ufanet-home-manager/auth-service/resources/logger"
	"github.com/rs/zerolog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Resources struct {
	Config *config.Config
	Log    zerolog.Logger
	DB     *gorm.DB
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

	return &Resources{
		Config: cfg,
		Log:    log,
		DB:     db,
	}, nil
}

func (r *Resources) Close() error {
	if r == nil || r.DB == nil {
		return nil
	}
	db, err := r.DB.DB()
	if err != nil {
		return err
	}
	return db.Close()
}
