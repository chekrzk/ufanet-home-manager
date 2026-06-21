package notifications

import (
	"github.com/gofiber/fiber/v2"

	gwerrors "github.com/chekrzk/ufanet-home-manager/api-gateway/internal/errors"
	"github.com/chekrzk/ufanet-home-manager/api-gateway/internal/models/constant"
	"github.com/chekrzk/ufanet-home-manager/api-gateway/internal/models/domain"
	"github.com/chekrzk/ufanet-home-manager/api-gateway/internal/models/dto"
	notificationsservice "github.com/chekrzk/ufanet-home-manager/api-gateway/internal/services/notifications"
)

type Handler struct {
	service notificationsservice.Service
}

func NewHandler(service notificationsservice.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Register(c *fiber.Ctx) error {
	req, err := gwerrors.ParseBody[dto.RegisterDeviceRequest](c)
	if err != nil {
		return err
	}

	if err := h.service.RegisterDevice(c.Context(), authContext(c), req); err != nil {
		return gwerrors.FromGRPC(err)
	}

	return gwerrors.NoContent(c)
}

func (h *Handler) Unregister(c *fiber.Ctx) error {
	req, err := gwerrors.ParseBody[dto.UnregisterDeviceRequest](c)
	if err != nil {
		return err
	}

	if err := h.service.UnregisterDevice(c.Context(), authContext(c), req); err != nil {
		return gwerrors.FromGRPC(err)
	}

	return gwerrors.NoContent(c)
}

func authContext(c *fiber.Ctx) domain.AuthContext {
	userID, _ := c.Locals(constant.CtxUserID).(string)
	role, _ := c.Locals(constant.CtxRole).(string)
	return domain.AuthContext{UserID: userID, Role: role}
}
