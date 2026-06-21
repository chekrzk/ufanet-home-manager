package errors

import "github.com/gofiber/fiber/v2"

type response struct {
	Success bool      `json:"success"`
	Data    any       `json:"data,omitempty"`
	Error   *AppError `json:"error,omitempty"`
}

func OK(c *fiber.Ctx, data any) error {
	return c.JSON(response{Success: true, Data: data})
}

func Created(c *fiber.Ctx, data any) error {
	return c.Status(fiber.StatusCreated).JSON(response{Success: true, Data: data})
}

func NoContent(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNoContent)
}

func Fail(c *fiber.Ctx, err error) error {
	appErr := Wrap(err, New(fiber.StatusInternalServerError, "internal_error", "internal server error"))
	return c.Status(appErr.Status).JSON(response{Success: false, Error: appErr})
}
