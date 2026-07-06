package router

import (
	"github.com/gofiber/fiber/v2"

	"github.com/chekrzk/ufanet-home-manager/api-gateway/internal/middlewares"
	"github.com/chekrzk/ufanet-home-manager/api-gateway/internal/models/constant"
)

type Handlers struct {
	Health        HealthHandler
	Auth          AuthHandler
	News          NewsHandler
	Notifications NotificationsHandler
	Profile       ProfileHandler
	Requests      RequestsHandler
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
	r.app.Use(r.mw.Logger())
	r.app.Use(r.mw.CORS())
	r.app.Use(r.mw.RateLimit())

	r.health()
	r.auth()
	r.protected()
}

func (r *Router) health() {
	r.app.Get("/health", r.h.Health.Check)
}

func (r *Router) auth() {
	auth := r.app.Group("/auth")
	auth.Post("/login", r.h.Auth.Login)
	auth.Post("/registr", r.h.Auth.Register)
	auth.Post("/register", r.h.Auth.Register)
	auth.Post("/refresh", r.h.Auth.Refresh)
}

func (r *Router) protected() {
	api := r.app.Group("", r.mw.Blacklist(), r.mw.JWT())

	api.Get("/profile", r.h.Profile.Me)
	api.Patch("/profile", r.h.Profile.Update)
	api.Get("/profile/workers", r.mw.Role(constant.RoleAdmin, constant.RoleManager), r.h.Profile.ListWorkers)
	api.Post("/profile/workers", r.mw.Role(constant.RoleAdmin, constant.RoleManager), r.h.Profile.AddWorker)
	api.Get("/profile/workers/availability", r.h.Profile.ListWorkerAvailability)
	api.Post("/profile/workers/availability", r.mw.Role(constant.RoleEmployee), r.h.Profile.SetWorkerAvailability)

	api.Get("/news", r.h.News.List)
	api.Post("/news", r.mw.Role(constant.RoleAdmin, constant.RoleManager), r.h.News.Create)

	api.Get("/requests", r.h.Requests.List)
	api.Post("/requests", r.h.Requests.Create)
	api.Get("/requests/:id", r.h.Requests.Get)
	api.Patch("/requests/:id/status", r.mw.Role(constant.RoleAdmin, constant.RoleManager, constant.RoleEmployee), r.h.Requests.UpdateStatus)
	api.Post("/requests/:id/comments", r.h.Requests.AddComment)

	api.Post("/notifications/register", r.h.Notifications.Register)
	api.Delete("/notifications/unregister", r.h.Notifications.Unregister)
	api.Get("/notifications", r.h.Notifications.List)
	api.Patch("/notifications/:id/read", r.h.Notifications.MarkRead)
}
