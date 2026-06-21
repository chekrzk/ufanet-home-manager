package router

import (
	"github.com/gofiber/fiber/v2"

	"github.com/chekrzk/ufanet-home-manager/api-gateway/internal/middlewares"
)

type Handlers struct {
	Health HealthHandler
}

type Router struct {
	app *fiber.App
	mw  *middlewares.Middlewares
	h   Handlers
}

func New(app *fiber.App, mw *middlewares.Middlewares, h Handlers) *Router {
	return &Router{
		app: app,
		mw:  mw,
		h:   h,
	}
}

func (r *Router) Register() {
	r.health()
}

func (r *Router) health() {
	r.app.Get("/health", r.h.Health.Check)
}
