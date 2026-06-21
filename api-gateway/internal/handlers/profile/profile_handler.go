package profile

import (
	"github.com/gofiber/fiber/v2"

	gwerrors "github.com/chekrzk/ufanet-home-manager/api-gateway/internal/errors"
	"github.com/chekrzk/ufanet-home-manager/api-gateway/internal/models/constant"
	"github.com/chekrzk/ufanet-home-manager/api-gateway/internal/models/domain"
	"github.com/chekrzk/ufanet-home-manager/api-gateway/internal/models/dto"
	profileservice "github.com/chekrzk/ufanet-home-manager/api-gateway/internal/services/profile"
)

type Handler struct {
	service profileservice.Service
}

func NewHandler(service profileservice.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Me(c *fiber.Ctx) error {
	user, err := h.service.Me(c.Context(), authContext(c))
	if err != nil {
		return gwerrors.FromGRPC(err)
	}

	return gwerrors.OK(c, user)
}

func (h *Handler) Update(c *fiber.Ctx) error {
	req, err := gwerrors.ParseBody[dto.UpdateProfileRequest](c)
	if err != nil {
		return err
	}

	user, err := h.service.Update(c.Context(), authContext(c), req)
	if err != nil {
		return gwerrors.FromGRPC(err)
	}

	return gwerrors.OK(c, user)
}

func authContext(c *fiber.Ctx) domain.AuthContext {
	userID, _ := c.Locals(constant.CtxUserID).(string)
	role, _ := c.Locals(constant.CtxRole).(string)
	return domain.AuthContext{UserID: userID, Role: role}
}
