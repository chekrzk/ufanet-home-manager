package config

import (
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	App       AppConfig
	HTTP      HTTPConfig
	Log       LogConfig
	Timeout   TimeoutConfig
	JWT       JWTConfig
	RateLimit RateLimitConfig
	Service   ServiceConfig
}

type AppConfig struct {
	Name string `envconfig:"APP_NAME" default:"api-gateway"`
	Env  string `envconfig:"APP_ENV" default:"development"`
}

type HTTPConfig struct {
	Host string `envconfig:"HTTP_HOST" default:"0.0.0.0"`
	Port string `envconfig:"HTTP_PORT" default:"8080"`
}

func (c HTTPConfig) Addr() string {
	return c.Host + ":" + c.Port
}

type LogConfig struct {
	Level  string `envconfig:"LOG_LEVEL" default:"debug"`
	Pretty bool   `envconfig:"LOG_PRETTY" default:"true"`
}

type TimeoutConfig struct {
	Request time.Duration `envconfig:"REQUEST_TIMEOUT" default:"5s"`
}

type JWTConfig struct {
	Secret    string        `envconfig:"JWT_SECRET" default:"dev-secret-change-me"`
	AccessTTL time.Duration `envconfig:"JWT_ACCESS_TTL" default:"15m"`
}

type RateLimitConfig struct {
	Max        int           `envconfig:"RATE_LIMIT_MAX" default:"100"`
	Expiration time.Duration `envconfig:"RATE_LIMIT_EXPIRATION" default:"1m"`
}

type ServiceConfig struct {
	AuthAddr          string `envconfig:"AUTH_SERVICE_ADDR" default:"localhost:50051"`
	NewsAddr          string `envconfig:"NEWS_SERVICE_ADDR" default:"localhost:50052"`
	RequestsAddr      string `envconfig:"REQUESTS_SERVICE_ADDR" default:"localhost:50053"`
	NotificationsAddr string `envconfig:"NOTIFICATIONS_SERVICE_ADDR" default:"localhost:50054"`
	ProfileAddr       string `envconfig:"PROFILE_SERVICE_ADDR" default:"localhost:50055"`
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	var cfg Config

	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
