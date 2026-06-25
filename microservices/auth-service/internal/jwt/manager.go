package jwt

import (
	"time"

	"github.com/chekrzk/ufanet-home-manager/auth-service/internal/models"
	"github.com/chekrzk/ufanet-home-manager/auth-service/resources/config"
	"github.com/golang-jwt/jwt/v5"
)

const (
	TokenTypeAccess  = "access"
	TokenTypeRefresh = "refresh"
)

type Manager struct {
	secret     []byte
	issuer     string
	accessTTL  time.Duration
	refreshTTL time.Duration
}

type Claims struct {
	Role      string `json:"role"`
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}

type Pair struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64
}

func NewManager(cfg config.JWTConfig) *Manager {
	return &Manager{
		secret:     []byte(cfg.Secret),
		issuer:     cfg.Issuer,
		accessTTL:  cfg.AccessTTL,
		refreshTTL: cfg.RefreshTTL,
	}
}

func (m *Manager) NewPair(user models.User) (Pair, error) {
	access, accessExp, err := m.newToken(user, TokenTypeAccess, m.accessTTL)
	if err != nil {
		return Pair{}, err
	}
	refresh, _, err := m.newToken(user, TokenTypeRefresh, m.refreshTTL)
	if err != nil {
		return Pair{}, err
	}
	return Pair{
		AccessToken:  access,
		RefreshToken: refresh,
		ExpiresIn:    int64(time.Until(accessExp).Seconds()),
	}, nil
}

func (m *Manager) ParseRefresh(token string) (string, error) {
	claims := &Claims{}
	parsed, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, jwt.ErrTokenSignatureInvalid
		}
		return m.secret, nil
	})
	if err != nil || !parsed.Valid || claims.TokenType != TokenTypeRefresh {
		return "", jwt.ErrTokenInvalidClaims
	}
	return claims.Subject, nil
}

func (m *Manager) newToken(user models.User, tokenType string, ttl time.Duration) (string, time.Time, error) {
	now := time.Now()
	expiresAt := now.Add(ttl)

	claims := Claims{
		Role:      string(user.Role),
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID,
			Issuer:    m.issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(m.secret)
	return signed, expiresAt, err
}
