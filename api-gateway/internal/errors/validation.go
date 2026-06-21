package errors

import "github.com/gofiber/fiber/v2"

type Validatable interface {
	Validate() map[string]string
}

func ParseBody[T any](c *fiber.Ctx) (T, error) {
	var req T
	if err := c.BodyParser(&req); err != nil {
		return req, New(fiber.StatusBadRequest, "invalid_body", "invalid request body")
	}
	if validatable, ok := any(req).(Validatable); ok {
		if fields := validatable.Validate(); len(fields) > 0 {
			return req, Validation(fields)
		}
	}
	return req, nil
}
