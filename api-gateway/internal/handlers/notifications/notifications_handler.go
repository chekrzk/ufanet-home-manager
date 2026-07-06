package notifications

import (
	"github.com/gofiber/fiber/v2"

	gwerrors "github.com/chekrzk/ufanet-home-manager/api-gateway/internal/errors"
	"github.com/chekrzk/ufanet-home-manager/api-gateway/internal/models/constant"
	"github.com/chekrzk/ufanet-home-manager/api-gateway/internal/models/domain"
	"github.com/chekrzk/ufanet-home-manager/api-gateway/internal/models/dto"
)

type Handler struct {
	service Service
}

// NewHandler отделяет HTTP transport уведомлений от service-сценариев.
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// Register связывает device token с текущим пользователем из JWT context.
func (h *Handler) Register(c *fiber.Ctx) error {
	req, err := gwerrors.ParseBody[dto.RegisterDeviceRequest](c)
	if err != nil {
		return err
	}

	if err := h.service.RegisterDevice(c.Context(), authContext(c), domain.RegisterDevice{
		Token:    req.Token,
		Platform: req.Platform,
	}); err != nil {
		return gwerrors.FromGRPC(err)
	}

	return gwerrors.NoContent(c)
}

// Unregister удаляет device token без раскрытия внутренних id устройств наружу.
func (h *Handler) Unregister(c *fiber.Ctx) error {
	req, err := gwerrors.ParseBody[dto.UnregisterDeviceRequest](c)
	if err != nil {
		return err
	}

	if err := h.service.UnregisterDevice(c.Context(), authContext(c), domain.UnregisterDevice{
		Token: req.Token,
	}); err != nil {
		return gwerrors.FromGRPC(err)
	}

	return gwerrors.NoContent(c)
}

// List строит pagination из query params и делегирует доступ service layer.
func (h *Handler) List(c *fiber.Ctx) error {
	var req dto.Pagination
	if err := c.QueryParser(&req); err != nil {
		return gwerrors.New(fiber.StatusBadRequest, "invalid_query", "invalid query params")
	}
	req.Normalize()

	page, err := h.service.List(c.Context(), authContext(c), domain.Pagination{Page: req.Page, Limit: req.Limit})
	if err != nil {
		return gwerrors.FromGRPC(err)
	}

	return gwerrors.OK(c, page)
}

// MarkRead использует path id и actor context, чтобы service проверил владельца.
func (h *Handler) MarkRead(c *fiber.Ctx) error {
	if err := h.service.MarkRead(c.Context(), authContext(c), c.Params("id")); err != nil {
		return gwerrors.FromGRPC(err)
	}

	return gwerrors.NoContent(c)
}

// authContext переносит результат JWT middleware в domain-модель.
func authContext(c *fiber.Ctx) domain.AuthContext {
	userID, _ := c.Locals(constant.CtxUserID).(string)
	role, _ := c.Locals(constant.CtxRole).(string)
	return domain.AuthContext{UserID: userID, Role: role}
}
