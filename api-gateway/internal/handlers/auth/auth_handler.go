package auth

import (
	"github.com/gofiber/fiber/v2"

	gwerrors "github.com/chekrzk/ufanet-home-manager/api-gateway/internal/errors"
	"github.com/chekrzk/ufanet-home-manager/api-gateway/internal/models/domain"
	"github.com/chekrzk/ufanet-home-manager/api-gateway/internal/models/dto"
	authservice "github.com/chekrzk/ufanet-home-manager/api-gateway/internal/services/auth"
)

type Handler struct {
	service authservice.Service
}

func NewHandler(service authservice.Service) *Handler {
	return &Handler{service: service}
}

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

func (h *Handler) Register(c *fiber.Ctx) error {
	req, err := gwerrors.ParseBody[dto.RegisterRequest](c)
	if err != nil {
		return err
	}

	user, err := h.service.Register(c.Context(), domain.RegisterUser{
		Phone:     req.Phone,
		Password:  req.Password,
		FullName:  req.FullName,
		HouseID:   req.HouseID,
		Apartment: req.Apartment,
	})
	if err != nil {
		return gwerrors.FromGRPC(err)
	}

	return gwerrors.Created(c, user)
}

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
