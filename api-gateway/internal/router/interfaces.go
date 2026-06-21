package router

import (
	"github.com/gofiber/fiber/v2"
)

type HealthHandler interface {
	Check(c *fiber.Ctx) error
}
