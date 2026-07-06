package notifications

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	gwerrors "github.com/chekrzk/ufanet-home-manager/api-gateway/internal/errors"
	handlermocks "github.com/chekrzk/ufanet-home-manager/api-gateway/internal/handlers/notifications/mocks"
	"github.com/chekrzk/ufanet-home-manager/api-gateway/internal/models/constant"
	"github.com/chekrzk/ufanet-home-manager/api-gateway/internal/models/domain"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/mock"
)

func newTestApp(handler *Handler) *fiber.App {
	app := fiber.New(fiber.Config{ErrorHandler: func(c *fiber.Ctx, err error) error {
		return gwerrors.Fail(c, err)
	}})
	app.Use(func(c *fiber.Ctx) error {
		c.Locals(constant.CtxUserID, "user-1")
		c.Locals(constant.CtxRole, "resident")
		return c.Next()
	})
	app.Post("/notifications/register", handler.Register)
	app.Post("/notifications/unregister", handler.Unregister)
	app.Get("/notifications", handler.List)
	app.Post("/notifications/:id/read", handler.MarkRead)
	return app
}

func jsonRequest(method, path, body string) *http.Request {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	return req
}

func TestRegister(t *testing.T) {
	service := handlermocks.NewService(t)
	app := newTestApp(NewHandler(service))
	actor := domain.AuthContext{UserID: "user-1", Role: "resident"}
	device := domain.RegisterDevice{Token: "token-1", Platform: "web"}

	service.On("RegisterDevice", mock.Anything, actor, device).Return(nil).Once()

	resp, err := app.Test(jsonRequest(http.MethodPost, "/notifications/register", `{"token":"token-1","platform":"web"}`))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusNoContent {
		t.Fatalf("status = %d, want %d", resp.StatusCode, fiber.StatusNoContent)
	}
}

func TestUnregister(t *testing.T) {
	service := handlermocks.NewService(t)
	app := newTestApp(NewHandler(service))
	actor := domain.AuthContext{UserID: "user-1", Role: "resident"}
	device := domain.UnregisterDevice{Token: "token-1"}

	service.On("UnregisterDevice", mock.Anything, actor, device).Return(nil).Once()

	resp, err := app.Test(jsonRequest(http.MethodPost, "/notifications/unregister", `{"token":"token-1"}`))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusNoContent {
		t.Fatalf("status = %d, want %d", resp.StatusCode, fiber.StatusNoContent)
	}
}

func TestList(t *testing.T) {
	service := handlermocks.NewService(t)
	app := newTestApp(NewHandler(service))
	actor := domain.AuthContext{UserID: "user-1", Role: "resident"}
	page := domain.Pagination{Page: 1, Limit: 20}

	service.On("List", mock.Anything, actor, page).Return(domain.Page[domain.Notification]{Items: []domain.Notification{{ID: "n-1"}}, Page: 1, Limit: 20, Total: 1}, nil).Once()

	resp, err := app.Test(httptest.NewRequest(http.MethodGet, "/notifications", nil))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("status = %d, want %d", resp.StatusCode, fiber.StatusOK)
	}
}

func TestMarkRead(t *testing.T) {
	service := handlermocks.NewService(t)
	app := newTestApp(NewHandler(service))
	actor := domain.AuthContext{UserID: "user-1", Role: "resident"}

	service.On("MarkRead", mock.Anything, actor, "notification-1").Return(nil).Once()

	resp, err := app.Test(httptest.NewRequest(http.MethodPost, "/notifications/notification-1/read", nil))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusNoContent {
		t.Fatalf("status = %d, want %d", resp.StatusCode, fiber.StatusNoContent)
	}
}
