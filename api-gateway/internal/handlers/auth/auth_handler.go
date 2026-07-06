package auth

import (
	"github.com/gofiber/fiber/v2"

	gwerrors "github.com/chekrzk/ufanet-home-manager/api-gateway/internal/errors"
	"github.com/chekrzk/ufanet-home-manager/api-gateway/internal/models/domain"
	"github.com/chekrzk/ufanet-home-manager/api-gateway/internal/models/dto"
)

type Handler struct {
	service Service
}

// NewHandler принимает service через интерфейс, чтобы HTTP-слой зависел от
// нужного поведения, а не от конкретной реализации auth-сервиса.
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// Login переводит HTTP credentials в domain-команду и оставляет проверку
// пароля auth-сценарию.
func (h *Handler) Login(c *fiber.Ctx) error {
	req, err := gwerrors.ParseBody[dto.LoginRequest](c)
	if err != nil {
		return err
	}

	tokens, err := h.service.Login(c.Context(), domain.LoginCredentials{
		Phone:    req.Phone,
		Password: req.Password,
	})
	if err != nil {
		return gwerrors.FromGRPC(err)
	}

	return gwerrors.OK(c, tokens)
}

// Register валидирует входной JSON до service layer, чтобы бизнес-сценарий
// работал уже с полной командой регистрации.
func (h *Handler) Register(c *fiber.Ctx) error {
	req, err := gwerrors.ParseBody[dto.RegisterRequest](c)
	if err != nil {
		return err
	}

	user, err := h.service.Register(c.Context(), domain.RegisterUser{
		Phone:    req.Phone,
		Password: req.Password,
		FullName: req.FullName,
		Role:     req.Role,
	})
	if err != nil {
		return gwerrors.FromGRPC(err)
	}

	return gwerrors.Created(c, user)
}

// Refresh не читает JWT из header, потому что обновление основано на отдельном
// refresh token из тела запроса.
func (h *Handler) Refresh(c *fiber.Ctx) error {
	req, err := gwerrors.ParseBody[dto.RefreshRequest](c)
	if err != nil {
		return err
	}

	tokens, err := h.service.Refresh(c.Context(), req.RefreshToken)
	if err != nil {
		return gwerrors.FromGRPC(err)
	}

	return gwerrors.OK(c, tokens)
}
