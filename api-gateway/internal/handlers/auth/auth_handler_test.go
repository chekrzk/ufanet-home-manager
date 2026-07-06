package auth

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	gwerrors "github.com/chekrzk/ufanet-home-manager/api-gateway/internal/errors"
	handlermocks "github.com/chekrzk/ufanet-home-manager/api-gateway/internal/handlers/auth/mocks"
	"github.com/chekrzk/ufanet-home-manager/api-gateway/internal/models/domain"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/mock"
)

func newTestApp(handler *Handler) *fiber.App {
	app := fiber.New(fiber.Config{ErrorHandler: func(c *fiber.Ctx, err error) error {
		return gwerrors.Fail(c, err)
	}})
	app.Post("/login", handler.Login)
	app.Post("/register", handler.Register)
	app.Post("/refresh", handler.Refresh)
	return app
}

func jsonRequest(method, path, body string) *http.Request {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	return req
}

func TestLogin(t *testing.T) {
	service := handlermocks.NewService(t)
	app := newTestApp(NewHandler(service))
	credentials := domain.LoginCredentials{Phone: "79990000001", Password: "123456"}

	service.On("Login", mock.Anything, credentials).Return(domain.AuthTokens{AccessToken: "access"}, nil).Once()

	resp, err := app.Test(jsonRequest(http.MethodPost, "/login", `{"phone":"79990000001","password":"123456"}`))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("status = %d, want %d", resp.StatusCode, fiber.StatusOK)
	}
}

func TestRegister(t *testing.T) {
	service := handlermocks.NewService(t)
	app := newTestApp(NewHandler(service))
	user := domain.RegisterUser{Phone: "79990000001", Password: "123456", FullName: "Ivan", Role: "resident"}

	service.On("Register", mock.Anything, user).Return(domain.User{ID: "user-1", Phone: user.Phone}, nil).Once()

	resp, err := app.Test(jsonRequest(http.MethodPost, "/register", `{"phone":"79990000001","password":"123456","full_name":"Ivan","role":"resident"}`))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusCreated {
		t.Fatalf("status = %d, want %d", resp.StatusCode, fiber.StatusCreated)
	}
}

func TestRefresh(t *testing.T) {
	service := handlermocks.NewService(t)
	app := newTestApp(NewHandler(service))

	service.On("Refresh", mock.Anything, "refresh").Return(domain.AuthTokens{AccessToken: "access"}, nil).Once()

	resp, err := app.Test(jsonRequest(http.MethodPost, "/refresh", `{"refresh_token":"refresh"}`))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("status = %d, want %d", resp.StatusCode, fiber.StatusOK)
	}
}

func TestLoginReturnsValidationError(t *testing.T) {
	service := handlermocks.NewService(t)
	app := newTestApp(NewHandler(service))

	resp, err := app.Test(jsonRequest(http.MethodPost, "/login", `{"phone":""}`))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("status = %d, want %d", resp.StatusCode, fiber.StatusBadRequest)
	}
}
