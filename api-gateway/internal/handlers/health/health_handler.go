package health

import "github.com/gofiber/fiber/v2"

type Handler struct {
	serviceName string
}

// NewHandler хранит имя сервиса в health response для быстрой диагностики окружения.
func NewHandler(serviceName string) *Handler {
	return &Handler{
		serviceName: serviceName,
	}
}

// Check остается легким endpoint без зависимостей, чтобы readiness gateway
// проверялась даже при недоступных downstream-сервисах.
func (h *Handler) Check(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "ok",
		"service": h.serviceName,
	})
}
