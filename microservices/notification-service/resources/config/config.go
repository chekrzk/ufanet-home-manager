package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	App   AppConfig
	GRPC  GRPCConfig
	DB    DBConfig
	Redis RedisConfig
	Log   LogConfig
}

type AppConfig struct {
	Name string `envconfig:"APP_NAME" default:"notification-service"`
	Env  string `envconfig:"APP_ENV" default:"local"`
}

type GRPCConfig struct {
	Host string `envconfig:"GRPC_HOST" default:"0.0.0.0"`
	Port string `envconfig:"GRPC_PORT" default:"50054"`
}

func (c GRPCConfig) Addr() string { return c.Host + ":" + c.Port }

type DBConfig struct {
	DSN string `envconfig:"DATABASE_DSN" default:"host=localhost user=postgres password=postgres dbname=ufanet_notifications port=5432 sslmode=disable TimeZone=Asia/Yekaterinburg"`
}

type RedisConfig struct {
	Addr   string `envconfig:"REDIS_ADDR" default:"localhost:6379"`
	Stream string `envconfig:"REDIS_NOTIFICATIONS_STREAM" default:"notifications"`
}

type LogConfig struct {
	Level  string `envconfig:"LOG_LEVEL" default:"debug"`
	Pretty bool   `envconfig:"LOG_PRETTY" default:"true"`
}

func Load() (*Config, error) {
	_ = godotenv.Load()
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
