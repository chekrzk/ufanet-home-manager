package requests

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

// NewHandler держит HTTP-заявки зависимыми от service interface, а не от клиента БД/gRPC.
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// Create переводит JSON в команду заявки, а владельца берет из JWT context.
func (h *Handler) Create(c *fiber.Ctx) error {
	req, err := gwerrors.ParseBody[dto.CreateRequestRequest](c)
	if err != nil {
		return err
	}

	request, err := h.service.Create(c.Context(), authContext(c), domain.CreateRequest{
		Category:         req.Category,
		Description:      req.Description,
		PreferredDate:    req.PreferredDate,
		AssignedWorkerID: req.AssignedWorkerID,
		Address:          req.Address,
		Apartment:        req.Apartment,
		Phone:            req.Phone,
	})
	if err != nil {
		return gwerrors.FromGRPC(err)
	}

	return gwerrors.Created(c, request)
}

// List нормализует pagination на HTTP-границе перед передачей в service layer.
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

// Get использует path id и actor context, чтобы service мог проверить доступ.
func (h *Handler) Get(c *fiber.Ctx) error {
	request, err := h.service.Get(c.Context(), authContext(c), c.Params("id"))
	if err != nil {
		return gwerrors.FromGRPC(err)
	}

	return gwerrors.OK(c, request)
}

// UpdateStatus оставляет правила перехода статусов requests-service.
func (h *Handler) UpdateStatus(c *fiber.Ctx) error {
	req, err := gwerrors.ParseBody[dto.UpdateRequestStatusRequest](c)
	if err != nil {
		return err
	}

	request, err := h.service.UpdateStatus(c.Context(), authContext(c), c.Params("id"), domain.UpdateRequestStatus{
		Status:     req.Status,
		AssignedTo: req.AssignedTo,
	})
	if err != nil {
		return gwerrors.FromGRPC(err)
	}

	return gwerrors.OK(c, request)
}

// AddComment связывает текст комментария с текущим actor на уровне service layer.
func (h *Handler) AddComment(c *fiber.Ctx) error {
	req, err := gwerrors.ParseBody[dto.AddRequestCommentRequest](c)
	if err != nil {
		return err
	}

	if err := h.service.AddComment(c.Context(), authContext(c), c.Params("id"), domain.AddRequestComment{Text: req.Text}); err != nil {
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
