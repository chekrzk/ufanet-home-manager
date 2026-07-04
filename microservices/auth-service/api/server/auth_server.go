package server

import (
	"context"
	stderrors "errors"

	apperrors "github.com/chekrzk/ufanet-home-manager/auth-service/internal/errors"
	"github.com/chekrzk/ufanet-home-manager/auth-service/internal/models"
	authv1 "github.com/chekrzk/ufanet-home-manager/contracts/gen/go/auth/v1"
	commonv1 "github.com/chekrzk/ufanet-home-manager/contracts/gen/go/common/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	authv1.UnimplementedAuthServiceServer
	service AuthService
}

func New(service AuthService) *Server {
	return &Server{service: service}
}

func (s *Server) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.AuthTokens, error) {
	tokens, err := s.service.Login(ctx, models.LoginCommand{
		Phone:    req.GetPhone(),
		Password: req.GetPassword(),
	})
	if err != nil {
		return nil, grpcError(err)
	}
	return tokensToProto(tokens.AccessToken, tokens.RefreshToken, tokens.ExpiresIn), nil
}

func (s *Server) Register(ctx context.Context, req *authv1.RegisterRequest) (*commonv1.User, error) {
	user, err := s.service.Register(ctx, models.RegisterCommand{
		Phone:    req.GetPhone(),
		Password: req.GetPassword(),
		Role:     req.GetRole(),
	})
	if err != nil {
		return nil, grpcError(err)
	}
	return userToProto(user), nil
}

func (s *Server) Refresh(ctx context.Context, req *authv1.RefreshRequest) (*authv1.AuthTokens, error) {
	tokens, err := s.service.Refresh(ctx, req.GetRefreshToken())
	if err != nil {
		return nil, grpcError(err)
	}
	return tokensToProto(tokens.AccessToken, tokens.RefreshToken, tokens.ExpiresIn), nil
}

func tokensToProto(accessToken string, refreshToken string, expiresIn int64) *authv1.AuthTokens {
	return &authv1.AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
	}
}

func userToProto(user models.User) *commonv1.User {
	return &commonv1.User{
		Id:    user.ID,
		Phone: user.Phone,
		Role:  string(user.Role),
	}
}

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
