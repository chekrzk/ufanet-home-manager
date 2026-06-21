package news

import (
	"github.com/gofiber/fiber/v2"

	gwerrors "github.com/chekrzk/ufanet-home-manager/api-gateway/internal/errors"
	"github.com/chekrzk/ufanet-home-manager/api-gateway/internal/models/constant"
	"github.com/chekrzk/ufanet-home-manager/api-gateway/internal/models/domain"
	"github.com/chekrzk/ufanet-home-manager/api-gateway/internal/models/dto"
	newsservice "github.com/chekrzk/ufanet-home-manager/api-gateway/internal/services/news"
)

type Handler struct {
	service newsservice.Service
}

func NewHandler(service newsservice.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) List(c *fiber.Ctx) error {
	var req dto.ListNewsRequest
	if err := c.QueryParser(&req); err != nil {
		return gwerrors.New(fiber.StatusBadRequest, "invalid_query", "invalid query params")
	}
	req.Normalize()

	page, err := h.service.List(c.Context(), authContext(c), req)
	if err != nil {
		return gwerrors.FromGRPC(err)
	}

	return gwerrors.OK(c, page)
}

func (h *Handler) Create(c *fiber.Ctx) error {
	req, err := gwerrors.ParseBody[dto.CreateNewsRequest](c)
	if err != nil {
		return err
	}

	item, err := h.service.Create(c.Context(), authContext(c), req)
	if err != nil {
		return gwerrors.FromGRPC(err)
	}

	return gwerrors.Created(c, item)
}

func authContext(c *fiber.Ctx) domain.AuthContext {
	userID, _ := c.Locals(constant.CtxUserID).(string)
	role, _ := c.Locals(constant.CtxRole).(string)
	return domain.AuthContext{UserID: userID, Role: role}
}
