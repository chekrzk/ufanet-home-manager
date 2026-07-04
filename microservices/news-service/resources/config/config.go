package config

import (
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	App     AppConfig
	GRPC    GRPCConfig
	DB      DBConfig
	Log     LogConfig
	Service ServiceConfig
}

type AppConfig struct {
	Name string `envconfig:"APP_NAME" default:"news-service"`
	Env  string `envconfig:"APP_ENV" default:"local"`
}

type GRPCConfig struct {
	Host string `envconfig:"GRPC_HOST" default:"0.0.0.0"`
	Port string `envconfig:"GRPC_PORT" default:"50052"`
}

func (c GRPCConfig) Addr() string { return c.Host + ":" + c.Port }

type DBConfig struct {
	DSN string `envconfig:"DATABASE_DSN" default:"host=localhost user=postgres password=postgres dbname=ufanet_news port=5432 sslmode=disable TimeZone=Asia/Yekaterinburg"`
}

type LogConfig struct {
	Level  string `envconfig:"LOG_LEVEL" default:"debug"`
	Pretty bool   `envconfig:"LOG_PRETTY" default:"true"`
}

type ServiceConfig struct {
	NotificationsAddr string        `envconfig:"NOTIFICATIONS_SERVICE_ADDR" default:"localhost:50054"`
	RequestTimeout    time.Duration `envconfig:"REQUEST_TIMEOUT" default:"5s"`
}

func Load() (*Config, error) {
	_ = godotenv.Load()
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
