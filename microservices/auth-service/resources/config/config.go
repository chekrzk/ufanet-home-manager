package config

import (
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	App  AppConfig
	GRPC GRPCConfig
	DB   DBConfig
	JWT  JWTConfig
	Log  LogConfig
}

type AppConfig struct {
	Name string `envconfig:"APP_NAME" default:"auth-service"`
	Env  string `envconfig:"APP_ENV" default:"local"`
}

type GRPCConfig struct {
	Host string `envconfig:"GRPC_HOST" default:"0.0.0.0"`
	Port string `envconfig:"GRPC_PORT" default:"50051"`
}

func (c GRPCConfig) Addr() string {
	return c.Host + ":" + c.Port
}

type DBConfig struct {
	DSN string `envconfig:"DATABASE_DSN" default:"host=localhost user=postgres password=postgres dbname=ufanet_auth port=5432 sslmode=disable TimeZone=Asia/Yekaterinburg"`
}

type JWTConfig struct {
	Secret     string        `envconfig:"JWT_SECRET" default:"dev-secret-change-me"`
	AccessTTL  time.Duration `envconfig:"JWT_ACCESS_TTL" default:"15m"`
	RefreshTTL time.Duration `envconfig:"JWT_REFRESH_TTL" default:"168h"`
	Issuer     string        `envconfig:"JWT_ISSUER" default:"auth-service"`
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
