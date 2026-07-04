package middlewares

import (
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog"

	"github.com/chekrzk/ufanet-home-manager/api-gateway/config"
	gwerrors "github.com/chekrzk/ufanet-home-manager/api-gateway/internal/errors"
	"github.com/chekrzk/ufanet-home-manager/api-gateway/internal/models/constant"
)

type Middlewares struct {
	cfg       *config.Config
	log       zerolog.Logger
	blacklist *Blacklist
}

func New(cfg *config.Config, log zerolog.Logger) *Middlewares {
	return &Middlewares{
		cfg:       cfg,
		log:       log,
		blacklist: NewBlacklist(),
	}
}

func (m *Middlewares) Logger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()

		m.log.Info().
			Str("method", c.Method()).
			Str("path", c.Path()).
			Int("status", c.Response().StatusCode()).
			Dur("duration", time.Since(start)).
			Str("ip", c.IP()).
			Msg("http request")

		return err
	}
}

func (m *Middlewares) CORS() fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins: m.cfg.CORS.AllowedOrigins,
		AllowMethods: m.cfg.CORS.AllowedMethods,
		AllowHeaders: m.cfg.CORS.AllowedHeaders,
	})
}

func (m *Middlewares) RateLimit() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        m.cfg.RateLimit.Max,
		Expiration: m.cfg.RateLimit.Expiration,
		LimitReached: func(c *fiber.Ctx) error {
			return gwerrors.Fail(c, gwerrors.New(fiber.StatusTooManyRequests, "rate_limited", "too many requests"))
		},
	})
}

func (m *Middlewares) Blacklist() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := bearerToken(c)
		if token != "" && m.blacklist.Contains(token) {
			return gwerrors.ErrUnauthorized
		}
		c.Locals(constant.CtxToken, token)
		return c.Next()
	}
}

func (m *Middlewares) JWT() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := bearerToken(c)
		if token == "" {
			return gwerrors.ErrUnauthorized
		}

		claims, err := parseJWT(token, m.cfg.JWT.Secret)
		if err != nil {
			return gwerrors.ErrUnauthorized
		}

		c.Locals(constant.CtxToken, token)
		c.Locals(constant.CtxUserID, claims.Subject)
		c.Locals(constant.CtxRole, claims.Role)

		return c.Next()
	}
}

func (m *Middlewares) Role(roles ...string) fiber.Handler {
	allowed := make(map[string]struct{}, len(roles))
	for _, role := range roles {
		allowed[role] = struct{}{}
	}

	return func(c *fiber.Ctx) error {
		role, _ := c.Locals(constant.CtxRole).(string)
		if _, ok := allowed[role]; !ok {
			return gwerrors.ErrForbidden
		}
		return c.Next()
	}
}

type Blacklist struct {
	mu     sync.RWMutex
	tokens map[string]time.Time
}

func NewBlacklist() *Blacklist {
	return &Blacklist{tokens: make(map[string]time.Time)}
}

func (b *Blacklist) Add(token string, expiresAt time.Time) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.tokens[token] = expiresAt
}

func (b *Blacklist) Contains(token string) bool {
	b.mu.RLock()
	expiresAt, ok := b.tokens[token]
	b.mu.RUnlock()
	if !ok {
		return false
	}
	if !expiresAt.IsZero() && time.Now().After(expiresAt) {
		b.mu.Lock()
		delete(b.tokens, token)
		b.mu.Unlock()
		return false
	}
	return true
}

func bearerToken(c *fiber.Ctx) string {
	header := c.Get(fiber.HeaderAuthorization)
	if header == "" {
		return ""
	}
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}
	return strings.TrimSpace(parts[1])
}

type jwtClaims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func parseJWT(token string, secret string) (jwtClaims, error) {
	var claims jwtClaims

	parsed, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, gwerrors.ErrUnauthorized
		}
		return []byte(secret), nil
	})
	if err != nil || !parsed.Valid || claims.Subject == "" || claims.Role == "" {
		return claims, gwerrors.ErrUnauthorized
	}

	return claims, nil
}
