package server

import (
	"context"
	stderrors "errors"

	commonv1 "github.com/chekrzk/ufanet-home-manager/contracts/gen/go/common/v1"
	profilev1 "github.com/chekrzk/ufanet-home-manager/contracts/gen/go/profile/v1"
	apperrors "github.com/chekrzk/ufanet-home-manager/profile-service/internal/errors"
	"github.com/chekrzk/ufanet-home-manager/profile-service/internal/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	profilev1.UnimplementedProfileServiceServer
	service ProfileService
}

func New(service ProfileService) *Server {
	return &Server{service: service}
}

func (s *Server) Me(ctx context.Context, req *profilev1.MeRequest) (*commonv1.User, error) {
	profile, err := s.service.Me(ctx, userContext(req.GetUser()))
	if err != nil {
		return nil, grpcError(err)
	}
	return profileToProto(profile, req.GetUser().GetRole()), nil
}

func (s *Server) Update(ctx context.Context, req *profilev1.UpdateProfileRequest) (*commonv1.User, error) {
	profile, err := s.service.Update(ctx, models.UpdateProfileCommand{
		Actor:     userContext(req.GetUser()),
		FullName:  req.GetFullName(),
		HouseID:   req.GetHouseId(),
		Apartment: req.GetApartment(),
	})
	if err != nil {
		return nil, grpcError(err)
	}
	return profileToProto(profile, req.GetUser().GetRole()), nil
}

func (s *Server) AddWorker(ctx context.Context, req *profilev1.AddWorkerRequest) (*commonv1.Worker, error) {
	worker, err := s.service.AddWorker(ctx, models.AddWorkerCommand{
		Actor:          userContext(req.GetActor()),
		UserID:         req.GetUserId(),
		FullName:       req.GetFullName(),
		Specialization: req.GetSpecialization(),
		Phone:          req.GetPhone(),
		HouseID:        req.GetHouseId(),
	})
	if err != nil {
		return nil, grpcError(err)
	}
	return workerToProto(worker), nil
}

func (s *Server) ListWorkers(ctx context.Context, req *profilev1.ListWorkersRequest) (*profilev1.ListWorkersResponse, error) {
	workers, err := s.service.ListWorkers(ctx, models.ListWorkersFilter{Actor: userContext(req.GetActor()), HouseID: req.GetHouseId()})
	if err != nil {
		return nil, grpcError(err)
	}
	items := make([]*commonv1.Worker, 0, len(workers))
	for _, worker := range workers {
		items = append(items, workerToProto(worker))
	}
	return &profilev1.ListWorkersResponse{Items: items}, nil
}

func userContext(user *commonv1.UserContext) models.UserContext {
	if user == nil {
		return models.UserContext{}
	}
	return models.UserContext{UserID: user.GetUserId(), Role: user.GetRole()}
}

func profileToProto(profile models.Profile, role string) *commonv1.User {
	return &commonv1.User{Id: profile.UserID, FullName: profile.FullName, Role: role, HouseId: profile.HouseID, Apartment: profile.Apartment}
}

func workerToProto(worker models.Worker) *commonv1.Worker {
	return &commonv1.Worker{
		Id:             worker.ID,
		UserId:         worker.UserID,
		FullName:       worker.FullName,
		Specialization: worker.Specialization,
		Phone:          worker.Phone,
		HouseId:        worker.HouseID,
		CreatedAt:      timestamppb.New(worker.CreatedAt),
	}
}

func grpcError(err error) error {
	switch {
	case stderrors.Is(err, apperrors.ErrInvalidArgument):
		return status.Error(codes.InvalidArgument, err.Error())
	case stderrors.Is(err, apperrors.ErrForbidden):
		return status.Error(codes.PermissionDenied, err.Error())
	case stderrors.Is(err, apperrors.ErrNotFound):
		return status.Error(codes.NotFound, err.Error())
	case stderrors.Is(err, apperrors.ErrAlreadyExists):
		return status.Error(codes.AlreadyExists, err.Error())
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}
