package server

import (
	"context"
	stderrors "errors"

	apperrors "github.com/chekrzk/ufanet-home-manager/auth-service/internal/errors"
	authv1 "github.com/chekrzk/ufanet-home-manager/contracts/gen/go/auth/v1"
	commonv1 "github.com/chekrzk/ufanet-home-manager/contracts/gen/go/common/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	authv1.UnimplementedAuthServiceServer
	service AuthService
}

// New связывает gRPC procedures с auth service interface, чтобы transport не
// создавал бизнес-зависимости самостоятельно.
func New(service AuthService) *Server {
	return &Server{service: service}
}

// Login конвертирует proto-запрос в доменную команду, чтобы service layer не
// зависел от формата gRPC-контракта.
func (s *Server) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.AuthTokens, error) {
	tokens, err := s.service.Login(ctx, loginCommandFromProto(req))
	if err != nil {
		return nil, grpcError(err)
	}
	return tokensToProto(tokens), nil
}

// Register оставляет создание учетной записи auth service, а server отвечает
// только за transport mapping и gRPC error.
func (s *Server) Register(ctx context.Context, req *authv1.RegisterRequest) (*commonv1.User, error) {
	user, err := s.service.Register(ctx, registerCommandFromProto(req))
	if err != nil {
		return nil, grpcError(err)
	}
	return userToProto(user), nil
}

// Refresh держит обновление токенов в auth service и возвращает proto-ответ gateway.
func (s *Server) Refresh(ctx context.Context, req *authv1.RefreshRequest) (*authv1.AuthTokens, error) {
	tokens, err := s.service.Refresh(ctx, req.GetRefreshToken())
	if err != nil {
		return nil, grpcError(err)
	}
	return tokensToProto(tokens), nil
}

// grpcError переводит доменные ошибки в стабильные gRPC status codes.
func grpcError(err error) error {
	switch {
	case stderrors.Is(err, apperrors.ErrInvalidArgument):
		return status.Error(codes.InvalidArgument, err.Error())
	case stderrors.Is(err, apperrors.ErrInvalidCredentials), stderrors.Is(err, apperrors.ErrUnauthorized):
		return status.Error(codes.Unauthenticated, err.Error())
	case stderrors.Is(err, apperrors.ErrUserAlreadyExists):
		return status.Error(codes.AlreadyExists, err.Error())
	case stderrors.Is(err, apperrors.ErrUserNotFound):
		return status.Error(codes.NotFound, err.Error())
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}
