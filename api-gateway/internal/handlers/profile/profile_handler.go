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

// NewHandler позволяет тестировать HTTP-слой профиля через service interface.
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// Me берет пользователя из JWT context, чтобы profile endpoint не принимал userID.
func (h *Handler) Me(c *fiber.Ctx) error {
	user, err := h.service.Me(c.Context(), authContext(c))
	if err != nil {
		return gwerrors.FromGRPC(err)
	}

	return gwerrors.OK(c, user)
}

// Update принимает только изменяемые поля профиля, а владельца определяет JWT.
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

// AddWorker оставляет проверку роли manager/admin в service layer.
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

// ListWorkers использует house_id как фильтр, а не как источник прав.
func (h *Handler) ListWorkers(c *fiber.Ctx) error {
	workers, err := h.service.ListWorkers(c.Context(), authContext(c), c.Query("house_id"))
	if err != nil {
		return gwerrors.FromGRPC(err)
	}

	return gwerrors.OK(c, workers)
}

// SetWorkerAvailability публикует расписание текущего работника через service layer.
func (h *Handler) SetWorkerAvailability(c *fiber.Ctx) error {
	req, err := gwerrors.ParseBody[dto.SetWorkerAvailabilityRequest](c)
	if err != nil {
		return err
	}

	availability, err := h.service.SetWorkerAvailability(c.Context(), authContext(c), domain.SetWorkerAvailability{
		Specialization: req.Specialization,
		HouseID:        req.HouseID,
		AvailableDate:  req.AvailableDate,
		AvailableTime:  req.AvailableTime,
	})
	if err != nil {
		return gwerrors.FromGRPC(err)
	}

	return gwerrors.Created(c, availability)
}

// ListWorkerAvailability собирает фильтр подбора работников из query params.
func (h *Handler) ListWorkerAvailability(c *fiber.Ctx) error {
	items, err := h.service.ListWorkerAvailability(c.Context(), authContext(c), domain.WorkerAvailabilityFilter{
		HouseID:        c.Query("house_id"),
		Specialization: c.Query("specialization"),
		AvailableDate:  c.Query("available_date"),
	})
	if err != nil {
		return gwerrors.FromGRPC(err)
	}

	return gwerrors.OK(c, items)
}

// authContext переносит результат JWT middleware в domain-модель.
func authContext(c *fiber.Ctx) domain.AuthContext {
	userID, _ := c.Locals(constant.CtxUserID).(string)
	role, _ := c.Locals(constant.CtxRole).(string)
	return domain.AuthContext{UserID: userID, Role: role}
}
