package health

import "github.com/gofiber/fiber/v2"

type Handler struct {
	serviceName string
}

func NewHandler(serviceName string) *Handler {
	return &Handler{
		serviceName: serviceName,
	}
}

func (h *Handler) Check(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "ok",
		"service": h.serviceName,
	})
}
