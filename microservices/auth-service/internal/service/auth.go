package service

import (
	"context"
	"strings"

	apperrors "github.com/chekrzk/ufanet-home-manager/auth-service/internal/errors"
	"github.com/chekrzk/ufanet-home-manager/auth-service/internal/hasher"
	jwtmanager "github.com/chekrzk/ufanet-home-manager/auth-service/internal/jwt"
	"github.com/chekrzk/ufanet-home-manager/auth-service/internal/models"
	"github.com/rs/zerolog"
)

type AuthService struct {
	users  UserRepository
	hasher *hasher.Hasher
	tokens *jwtmanager.Manager
	log    zerolog.Logger
}

// NewAuthService получает зависимости снаружи, чтобы авторизация не создавала
// DB, hasher или JWT manager внутри бизнес-слоя.
func NewAuthService(users UserRepository, hasher *hasher.Hasher, tokens *jwtmanager.Manager, log zerolog.Logger) *AuthService {
	return &AuthService{users: users, hasher: hasher, tokens: tokens, log: log}
}

// Register создает identity пользователя и хранит только password hash, чтобы
// остальные сервисы могли доверять auth-service как источнику учетных записей.
func (s *AuthService) Register(ctx context.Context, cmd models.RegisterCommand) (models.User, error) {
	if err := validateRegister(cmd); err != nil {
		return models.User{}, err
	}

	exists, err := s.users.ExistsByPhone(ctx, cmd.Phone)
	if err != nil {
		return models.User{}, err
	}
	if exists {
		return models.User{}, apperrors.ErrUserAlreadyExists
	}

	hash, err := s.hasher.Hash(cmd.Password)
	if err != nil {
		return models.User{}, err
	}

	user := models.User{
		Phone:        strings.TrimSpace(cmd.Phone),
		PasswordHash: hash,
		Role:         registerRole(cmd.Role),
	}
	if err := s.users.Create(ctx, &user); err != nil {
		return models.User{}, err
	}

	s.log.Info().Str("user_id", user.ID).Msg("user registered")
	return user, nil
}

// Login выдает пару токенов только после проверки hash, чтобы gateway и другие
// сервисы работали с JWT, а не с паролем пользователя.
func (s *AuthService) Login(ctx context.Context, cmd models.LoginCommand) (jwtmanager.Pair, error) {
	if strings.TrimSpace(cmd.Phone) == "" || cmd.Password == "" {
		return jwtmanager.Pair{}, apperrors.ErrInvalidArgument
	}

	user, err := s.users.FindByPhone(ctx, strings.TrimSpace(cmd.Phone))
	if err != nil {
		return jwtmanager.Pair{}, apperrors.ErrInvalidCredentials
	}
	if !s.hasher.Compare(user.PasswordHash, cmd.Password) {
		return jwtmanager.Pair{}, apperrors.ErrInvalidCredentials
	}

	return s.tokens.NewPair(user)
}

// Refresh перевыпускает токены через refresh token, чтобы короткий access token
// можно было безопасно обновлять без повторной передачи пароля.
func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (jwtmanager.Pair, error) {
	if strings.TrimSpace(refreshToken) == "" {
		return jwtmanager.Pair{}, apperrors.ErrUnauthorized
	}

	userID, err := s.tokens.ParseRefresh(refreshToken)
	if err != nil {
		return jwtmanager.Pair{}, apperrors.ErrUnauthorized
	}

	user, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return jwtmanager.Pair{}, apperrors.ErrUnauthorized
	}

	return s.tokens.NewPair(user)
}

// validateRegister отсекает некорректные учетные данные до обращения к БД и hashing.
func validateRegister(cmd models.RegisterCommand) error {
	if strings.TrimSpace(cmd.Phone) == "" || cmd.Password == "" {
		return apperrors.ErrInvalidArgument
	}
	if len(cmd.Password) < 6 {
		return apperrors.ErrInvalidArgument
	}
	return nil
}

// registerRole ограничивает self-registration безопасными ролями.
func registerRole(role string) models.Role {
	switch strings.TrimSpace(role) {
	case string(models.RoleEmployee):
		return models.RoleEmployee
	default:
		return models.RoleResident
	}
}
