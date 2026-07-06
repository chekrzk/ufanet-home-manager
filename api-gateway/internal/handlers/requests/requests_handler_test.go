package requests

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	gwerrors "github.com/chekrzk/ufanet-home-manager/api-gateway/internal/errors"
	handlermocks "github.com/chekrzk/ufanet-home-manager/api-gateway/internal/handlers/requests/mocks"
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
	app.Post("/requests", handler.Create)
	app.Get("/requests", handler.List)
	app.Get("/requests/:id", handler.Get)
	app.Patch("/requests/:id/status", handler.UpdateStatus)
	app.Post("/requests/:id/comments", handler.AddComment)
	return app
}

func jsonRequest(method, path, body string) *http.Request {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	return req
}

func TestCreate(t *testing.T) {
	service := handlermocks.NewService(t)
	app := newTestApp(NewHandler(service))
	actor := domain.AuthContext{UserID: "user-1", Role: "resident"}
	command := domain.CreateRequest{
		Category:         "electrician",
		Description:      "No light",
		PreferredDate:    "2026-07-13",
		AssignedWorkerID: "worker-1",
		Address:          "Prospekt Oktyabrya, 107",
		Apartment:        "12",
		Phone:            "79990000001",
	}

	service.On("Create", mock.Anything, actor, command).Return(domain.Request{ID: "request-1"}, nil).Once()

	resp, err := app.Test(jsonRequest(http.MethodPost, "/requests", `{"category":"electrician","description":"No light","preferred_date":"2026-07-13","assigned_worker_id":"worker-1","address":"Prospekt Oktyabrya, 107","apartment":"12","phone":"79990000001"}`))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusCreated {
		t.Fatalf("status = %d, want %d", resp.StatusCode, fiber.StatusCreated)
	}
}

func TestList(t *testing.T) {
	service := handlermocks.NewService(t)
	app := newTestApp(NewHandler(service))
	actor := domain.AuthContext{UserID: "user-1", Role: "resident"}
	page := domain.Pagination{Page: 2, Limit: 10}

	service.On("List", mock.Anything, actor, page).Return(domain.Page[domain.Request]{Items: []domain.Request{{ID: "request-1"}}, Page: 2, Limit: 10, Total: 1}, nil).Once()

	resp, err := app.Test(httptest.NewRequest(http.MethodGet, "/requests?page=2&limit=10", nil))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("status = %d, want %d", resp.StatusCode, fiber.StatusOK)
	}
}

func TestGet(t *testing.T) {
	service := handlermocks.NewService(t)
	app := newTestApp(NewHandler(service))
	actor := domain.AuthContext{UserID: "user-1", Role: "resident"}

	service.On("Get", mock.Anything, actor, "request-1").Return(domain.Request{ID: "request-1"}, nil).Once()

	resp, err := app.Test(httptest.NewRequest(http.MethodGet, "/requests/request-1", nil))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("status = %d, want %d", resp.StatusCode, fiber.StatusOK)
	}
}

func TestUpdateStatus(t *testing.T) {
	service := handlermocks.NewService(t)
	app := newTestApp(NewHandler(service))
	actor := domain.AuthContext{UserID: "user-1", Role: "resident"}
	command := domain.UpdateRequestStatus{Status: "in_progress", AssignedTo: "worker-1"}

	service.On("UpdateStatus", mock.Anything, actor, "request-1", command).Return(domain.Request{ID: "request-1", Status: "in_progress"}, nil).Once()

	resp, err := app.Test(jsonRequest(http.MethodPatch, "/requests/request-1/status", `{"status":"in_progress","assigned_to":"worker-1"}`))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("status = %d, want %d", resp.StatusCode, fiber.StatusOK)
	}
}

func TestAddComment(t *testing.T) {
	service := handlermocks.NewService(t)
	app := newTestApp(NewHandler(service))
	actor := domain.AuthContext{UserID: "user-1", Role: "resident"}
	command := domain.AddRequestComment{Text: "Need help"}

	service.On("AddComment", mock.Anything, actor, "request-1", command).Return(nil).Once()

	resp, err := app.Test(jsonRequest(http.MethodPost, "/requests/request-1/comments", `{"text":"Need help"}`))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusNoContent {
		t.Fatalf("status = %d, want %d", resp.StatusCode, fiber.StatusNoContent)
	}
}
