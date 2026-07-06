package news

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	gwerrors "github.com/chekrzk/ufanet-home-manager/api-gateway/internal/errors"
	handlermocks "github.com/chekrzk/ufanet-home-manager/api-gateway/internal/handlers/news/mocks"
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
		c.Locals(constant.CtxUserID, "manager-1")
		c.Locals(constant.CtxRole, "manager")
		return c.Next()
	})
	app.Get("/news", handler.List)
	app.Post("/news", handler.Create)
	return app
}

func jsonRequest(method, path, body string) *http.Request {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	return req
}

func TestList(t *testing.T) {
	service := handlermocks.NewService(t)
	app := newTestApp(NewHandler(service))
	actor := domain.AuthContext{UserID: "manager-1", Role: "manager"}
	filter := domain.NewsFilter{Pagination: domain.Pagination{Page: 2, Limit: 10}, DateFrom: "2026-07-01"}

	service.On("List", mock.Anything, actor, filter).Return(domain.Page[domain.News]{Items: []domain.News{{ID: "news-1"}}, Page: 2, Limit: 10, Total: 1}, nil).Once()

	resp, err := app.Test(httptest.NewRequest(http.MethodGet, "/news?page=2&limit=10&date_from=2026-07-01", nil))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("status = %d, want %d", resp.StatusCode, fiber.StatusOK)
	}
}

func TestCreate(t *testing.T) {
	service := handlermocks.NewService(t)
	app := newTestApp(NewHandler(service))
	actor := domain.AuthContext{UserID: "manager-1", Role: "manager"}
	command := domain.CreateNews{Title: "Title", Body: "Body", HouseID: "house-1"}

	service.On("Create", mock.Anything, actor, command).Return(domain.News{ID: "news-1", Title: "Title"}, nil).Once()

	resp, err := app.Test(jsonRequest(http.MethodPost, "/news", `{"title":"Title","body":"Body","house_id":"house-1"}`))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusCreated {
		t.Fatalf("status = %d, want %d", resp.StatusCode, fiber.StatusCreated)
	}
}

func TestCreateReturnsValidationError(t *testing.T) {
	service := handlermocks.NewService(t)
	app := newTestApp(NewHandler(service))

	resp, err := app.Test(jsonRequest(http.MethodPost, "/news", `{"title":""}`))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("status = %d, want %d", resp.StatusCode, fiber.StatusBadRequest)
	}
}
