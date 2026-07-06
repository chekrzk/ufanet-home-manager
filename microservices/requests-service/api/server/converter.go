package server

import (
	"time"

	commonv1 "github.com/chekrzk/ufanet-home-manager/contracts/gen/go/common/v1"
	requestsv1 "github.com/chekrzk/ufanet-home-manager/contracts/gen/go/requests/v1"
	"github.com/chekrzk/ufanet-home-manager/requests-service/internal/models"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// userContextFromProto переносит actor context в domain без proto-зависимости.
func userContextFromProto(user *commonv1.UserContext) models.UserContext {
	if user == nil {
		return models.UserContext{}
	}
	return models.UserContext{UserID: user.GetUserId(), Role: user.GetRole()}
}

// createRequestCommandFromProto собирает команду создания заявки на transport-границе.
func createRequestCommandFromProto(req *requestsv1.CreateRequestRequest) models.CreateRequestCommand {
	return models.CreateRequestCommand{
		User:             userContextFromProto(req.GetUser()),
		Category:         req.GetCategory(),
		Description:      req.GetDescription(),
		PreferredDate:    req.GetPreferredDate(),
		AssignedWorkerID: req.GetAssignedWorkerId(),
		Address:          req.GetAddress(),
		Apartment:        req.GetApartment(),
		Phone:            req.GetPhone(),
	}
}

// listRequestsFilterFromProto переносит actor и pagination в сценарий списка заявок.
func listRequestsFilterFromProto(req *requestsv1.ListRequestsRequest) models.ListRequestsFilter {
	return models.ListRequestsFilter{
		Actor: userContextFromProto(req.GetUser()),
		Pagination: models.Pagination{
			Page:  int(req.GetPagination().GetPage()),
			Limit: int(req.GetPagination().GetLimit()),
		},
	}
}

// getRequestCommandFromProto связывает request id с actor для проверки доступа.
func getRequestCommandFromProto(req *requestsv1.GetRequestRequest) models.GetRequestCommand {
	return models.GetRequestCommand{Actor: userContextFromProto(req.GetUser()), RequestID: req.GetRequestId()}
}

// updateRequestStatusCommandFromProto переводит изменение статуса в domain command.
func updateRequestStatusCommandFromProto(req *requestsv1.UpdateRequestStatusRequest) models.UpdateRequestStatusCommand {
	return models.UpdateRequestStatusCommand{
		Actor:      userContextFromProto(req.GetActor()),
		RequestID:  req.GetRequestId(),
		Status:     req.GetStatus(),
		AssignedTo: req.GetAssignedTo(),
	}
}

// addRequestCommentCommandFromProto строит команду комментария с actor context.
func addRequestCommentCommandFromProto(req *requestsv1.AddRequestCommentRequest) models.AddRequestCommentCommand {
	return models.AddRequestCommentCommand{
		Actor:     userContextFromProto(req.GetUser()),
		RequestID: req.GetRequestId(),
		Text:      req.GetText(),
	}
}

// requestToProto скрывает storage-модель заявки за общим gRPC-контрактом.
func requestToProto(request models.MaintenanceRequest) *commonv1.MaintenanceRequest {
	assignedTo := ""
	if request.AssignedTo != nil {
		assignedTo = *request.AssignedTo
	}
	return &commonv1.MaintenanceRequest{
		Id:            request.ID,
		UserId:        request.UserID,
		Category:      request.Category,
		Description:   request.Description,
		Status:        request.Status,
		CreatedAt:     timestamppb.New(request.CreatedAt),
		UpdatedAt:     timestamppb.New(request.UpdatedAt),
		AssignedTo:    assignedTo,
		PreferredDate: request.PreferredDate,
		Address:       request.Address,
		Apartment:     request.Apartment,
		Phone:         request.Phone,
		AcceptedAt:    timeToProto(request.AcceptedAt),
		DeclinedAt:    timeToProto(request.DeclinedAt),
		CompletedAt:   timeToProto(request.CompletedAt),
	}
}

// timeToProto сохраняет nullable timestamps при переходе из domain в proto.
func timeToProto(value *time.Time) *timestamppb.Timestamp {
	if value == nil {
		return nil
	}
	return timestamppb.New(*value)
}
