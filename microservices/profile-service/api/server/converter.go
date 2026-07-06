package server

import (
	commonv1 "github.com/chekrzk/ufanet-home-manager/contracts/gen/go/common/v1"
	profilev1 "github.com/chekrzk/ufanet-home-manager/contracts/gen/go/profile/v1"
	"github.com/chekrzk/ufanet-home-manager/profile-service/internal/models"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// userContextFromProto переносит auth context в domain без зависимости от proto.
func userContextFromProto(user *commonv1.UserContext) models.UserContext {
	if user == nil {
		return models.UserContext{}
	}
	return models.UserContext{UserID: user.GetUserId(), Role: user.GetRole()}
}

// updateProfileCommandFromProto собирает команду изменения профиля из контракта.
func updateProfileCommandFromProto(req *profilev1.UpdateProfileRequest) models.UpdateProfileCommand {
	return models.UpdateProfileCommand{
		Actor:     userContextFromProto(req.GetUser()),
		FullName:  req.GetFullName(),
		HouseID:   req.GetHouseId(),
		Apartment: req.GetApartment(),
	}
}

// addWorkerCommandFromProto отделяет transport payload от worker domain rules.
func addWorkerCommandFromProto(req *profilev1.AddWorkerRequest) models.AddWorkerCommand {
	return models.AddWorkerCommand{
		Actor:          userContextFromProto(req.GetActor()),
		UserID:         req.GetUserId(),
		FullName:       req.GetFullName(),
		Specialization: req.GetSpecialization(),
		Phone:          req.GetPhone(),
		HouseID:        req.GetHouseId(),
	}
}

// listWorkersFilterFromProto переносит фильтр работников вместе с actor context.
func listWorkersFilterFromProto(req *profilev1.ListWorkersRequest) models.ListWorkersFilter {
	return models.ListWorkersFilter{Actor: userContextFromProto(req.GetActor()), HouseID: req.GetHouseId()}
}

// setWorkerAvailabilityCommandFromProto строит команду публикации расписания.
func setWorkerAvailabilityCommandFromProto(req *profilev1.SetWorkerAvailabilityRequest) models.SetWorkerAvailabilityCommand {
	return models.SetWorkerAvailabilityCommand{
		Worker:         userContextFromProto(req.GetWorker()),
		Specialization: req.GetSpecialization(),
		HouseID:        req.GetHouseId(),
		AvailableDate:  req.GetAvailableDate(),
		AvailableTime:  req.GetAvailableTime(),
	}
}

// listWorkerAvailabilityFilterFromProto собирает доменный фильтр подбора работников.
func listWorkerAvailabilityFilterFromProto(req *profilev1.ListWorkerAvailabilityRequest) models.ListWorkerAvailabilityFilter {
	return models.ListWorkerAvailabilityFilter{
		Actor:          userContextFromProto(req.GetActor()),
		Specialization: req.GetSpecialization(),
		HouseID:        req.GetHouseId(),
		AvailableDate:  req.GetAvailableDate(),
	}
}

// profileToProto возвращает профиль в общем user-контракте gateway.
func profileToProto(profile models.Profile, role string) *commonv1.User {
	return &commonv1.User{Id: profile.UserID, FullName: profile.FullName, Role: role, HouseId: profile.HouseID, Apartment: profile.Apartment}
}

// workerToProto не раскрывает storage-модель worker за пределы сервиса.
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

// availabilityToProto стабилизирует ответ расписания для gateway и frontend.
func availabilityToProto(item models.WorkerAvailability) *commonv1.WorkerAvailability {
	return &commonv1.WorkerAvailability{
		Id:             item.ID,
		WorkerId:       item.WorkerID,
		UserId:         item.UserID,
		Specialization: item.Specialization,
		HouseId:        item.HouseID,
		AvailableDate:  item.AvailableDate,
		AvailableTime:  item.AvailableTime,
		CreatedAt:      timestamppb.New(item.CreatedAt),
	}
}
