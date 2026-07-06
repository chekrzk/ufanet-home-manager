package server

import (
	"context"
	stderrors "errors"

	commonv1 "github.com/chekrzk/ufanet-home-manager/contracts/gen/go/common/v1"
	profilev1 "github.com/chekrzk/ufanet-home-manager/contracts/gen/go/profile/v1"
	apperrors "github.com/chekrzk/ufanet-home-manager/profile-service/internal/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	profilev1.UnimplementedProfileServiceServer
	service ProfileService
}

func New(service ProfileService) *Server {
	return &Server{service: service}
}

func (s *Server) Me(ctx context.Context, req *profilev1.MeRequest) (*commonv1.User, error) {
	profile, err := s.service.Me(ctx, userContextFromProto(req.GetUser()))
	if err != nil {
		return nil, grpcError(err)
	}
	return profileToProto(profile, req.GetUser().GetRole()), nil
}

func (s *Server) Update(ctx context.Context, req *profilev1.UpdateProfileRequest) (*commonv1.User, error) {
	profile, err := s.service.Update(ctx, updateProfileCommandFromProto(req))
	if err != nil {
		return nil, grpcError(err)
	}
	return profileToProto(profile, req.GetUser().GetRole()), nil
}

func (s *Server) AddWorker(ctx context.Context, req *profilev1.AddWorkerRequest) (*commonv1.Worker, error) {
	worker, err := s.service.AddWorker(ctx, addWorkerCommandFromProto(req))
	if err != nil {
		return nil, grpcError(err)
	}
	return workerToProto(worker), nil
}

func (s *Server) ListWorkers(ctx context.Context, req *profilev1.ListWorkersRequest) (*profilev1.ListWorkersResponse, error) {
	workers, err := s.service.ListWorkers(ctx, listWorkersFilterFromProto(req))
	if err != nil {
		return nil, grpcError(err)
	}
	items := make([]*commonv1.Worker, 0, len(workers))
	for _, worker := range workers {
		items = append(items, workerToProto(worker))
	}
	return &profilev1.ListWorkersResponse{Items: items}, nil
}

func (s *Server) SetWorkerAvailability(ctx context.Context, req *profilev1.SetWorkerAvailabilityRequest) (*commonv1.WorkerAvailability, error) {
	availability, err := s.service.SetWorkerAvailability(ctx, setWorkerAvailabilityCommandFromProto(req))
	if err != nil {
		return nil, grpcError(err)
	}
	return availabilityToProto(availability), nil
}

func (s *Server) ListWorkerAvailability(ctx context.Context, req *profilev1.ListWorkerAvailabilityRequest) (*profilev1.ListWorkerAvailabilityResponse, error) {
	items, err := s.service.ListWorkerAvailability(ctx, listWorkerAvailabilityFilterFromProto(req))
	if err != nil {
		return nil, grpcError(err)
	}
	respItems := make([]*commonv1.WorkerAvailability, 0, len(items))
	for _, item := range items {
		respItems = append(respItems, availabilityToProto(item))
	}
	return &profilev1.ListWorkerAvailabilityResponse{Items: respItems}, nil
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
