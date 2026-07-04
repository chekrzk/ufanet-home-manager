package profile

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

func NewHandler(service Service) *Handler {
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

	user, err := h.service.Update(c.Context(), authContext(c), domain.UpdateProfile{
		FullName:  req.FullName,
		HouseID:   req.HouseID,
		Apartment: req.Apartment,
	})
	if err != nil {
		return gwerrors.FromGRPC(err)
	}

	return gwerrors.OK(c, user)
}

func (h *Handler) AddWorker(c *fiber.Ctx) error {
	req, err := gwerrors.ParseBody[dto.AddWorkerRequest](c)
	if err != nil {
		return err
	}

	worker, err := h.service.AddWorker(c.Context(), authContext(c), domain.AddWorker{
		UserID:         req.UserID,
		FullName:       req.FullName,
		Specialization: req.Specialization,
		Phone:          req.Phone,
		HouseID:        req.HouseID,
	})
	if err != nil {
		return gwerrors.FromGRPC(err)
	}

	return gwerrors.Created(c, worker)
}

func (h *Handler) ListWorkers(c *fiber.Ctx) error {
	workers, err := h.service.ListWorkers(c.Context(), authContext(c), c.Query("house_id"))
	if err != nil {
		return gwerrors.FromGRPC(err)
	}

	return gwerrors.OK(c, workers)
}

func authContext(c *fiber.Ctx) domain.AuthContext {
	userID, _ := c.Locals(constant.CtxUserID).(string)
	role, _ := c.Locals(constant.CtxRole).(string)
	return domain.AuthContext{UserID: userID, Role: role}
}
