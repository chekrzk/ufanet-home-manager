package router

import (
	"github.com/gofiber/fiber/v2"
)

type HealthHandler interface {
	Check(c *fiber.Ctx) error
}

type AuthHandler interface {
	Login(c *fiber.Ctx) error
	Register(c *fiber.Ctx) error
	Refresh(c *fiber.Ctx) error
}

type NewsHandler interface {
	List(c *fiber.Ctx) error
	Create(c *fiber.Ctx) error
}

type NotificationsHandler interface {
	Register(c *fiber.Ctx) error
	Unregister(c *fiber.Ctx) error
}

type ProfileHandler interface {
	Me(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
}
