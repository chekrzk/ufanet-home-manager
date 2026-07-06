package service

import (
	"context"
	"strings"
	"time"

	apperrors "github.com/chekrzk/ufanet-home-manager/requests-service/internal/errors"
	"github.com/chekrzk/ufanet-home-manager/requests-service/internal/models"
	"github.com/rs/zerolog"
)

type Service struct {
	repo      RequestRepository
	publisher NotificationPublisher
	log       zerolog.Logger
}

// New принимает repository и publisher через интерфейсы, чтобы заявки могли
// храниться и создавать уведомления без прямой зависимости от transport-кода.
func New(repo RequestRepository, publisher NotificationPublisher, log zerolog.Logger) *Service {
	return &Service{repo: repo, publisher: publisher, log: log}
}

// Create фиксирует заявку от текущего пользователя и сразу создает событие,
// чтобы заявитель и назначенный работник получили консистентное уведомление.
func (s *Service) Create(ctx context.Context, cmd models.CreateRequestCommand) (models.MaintenanceRequest, error) {
	if strings.TrimSpace(cmd.User.UserID) == "" || strings.TrimSpace(cmd.Category) == "" || strings.TrimSpace(cmd.Description) == "" {
		return models.MaintenanceRequest{}, apperrors.ErrInvalidArgument
	}
	request := models.MaintenanceRequest{
		UserID:        cmd.User.UserID,
		Category:      strings.TrimSpace(cmd.Category),
		Description:   strings.TrimSpace(cmd.Description),
		Status:        models.RequestStatusNew,
		PreferredDate: strings.TrimSpace(cmd.PreferredDate),
		Address:       strings.TrimSpace(cmd.Address),
		Apartment:     strings.TrimSpace(cmd.Apartment),
		Phone:         strings.TrimSpace(cmd.Phone),
	}
	assignRequest(&request, cmd.AssignedWorkerID)
	if err := s.repo.Create(ctx, &request); err != nil {
		return models.MaintenanceRequest{}, err
	}
	s.publishTo(ctx, request.UserID, request, "request.created", "New request created", "Request status: "+request.Status)
	if request.AssignedTo != nil {
		s.publishTo(ctx, *request.AssignedTo, request, "request.assigned", "New request assigned", request.Description)
	}
	return request, nil
}

// List нормализует пагинацию, чтобы список заявок оставался предсказуемым для UI.
func (s *Service) List(ctx context.Context, filter models.ListRequestsFilter) (models.RequestsPage, error) {
	items, total, err := s.repo.List(ctx, filter)
	if err != nil {
		return models.RequestsPage{}, err
	}
	page, limit := filter.Pagination.Page, filter.Pagination.Limit
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	return models.RequestsPage{Items: items, Page: page, Limit: limit, Total: total}, nil
}

// Get не раскрывает чужие заявки обычному жителю, но оставляет доступ ролям управления.
func (s *Service) Get(ctx context.Context, cmd models.GetRequestCommand) (models.MaintenanceRequest, error) {
	request, err := s.repo.FindByID(ctx, cmd.RequestID)
	if err != nil {
		return models.MaintenanceRequest{}, err
	}
	if cmd.Actor.Role == "resident" && request.UserID != cmd.Actor.UserID {
		return models.MaintenanceRequest{}, apperrors.ErrForbidden
	}
	return request, nil
}

// UpdateStatus ведет жизненный цикл заявки и публикует событие, чтобы статус
// в базе, UI и уведомлениях менялся согласованно.
func (s *Service) UpdateStatus(ctx context.Context, cmd models.UpdateRequestStatusCommand) (models.MaintenanceRequest, error) {
	if !canManageRequests(cmd.Actor.Role) || strings.TrimSpace(cmd.Status) == "" {
		return models.MaintenanceRequest{}, apperrors.ErrForbidden
	}
	request, err := s.repo.FindByID(ctx, cmd.RequestID)
	if err != nil {
		return models.MaintenanceRequest{}, err
	}
	status := strings.TrimSpace(cmd.Status)
	if !validStatus(status) {
		return models.MaintenanceRequest{}, apperrors.ErrInvalidArgument
	}
	now := time.Now()
	request.Status = status
	assignRequest(&request, cmd.AssignedTo)
	switch status {
	case models.RequestStatusInProgress:
		request.AcceptedAt = &now
	case models.RequestStatusCanceled:
		request.DeclinedAt = &now
	case models.RequestStatusDone:
		request.CompletedAt = &now
	}
	if err := s.repo.Save(ctx, &request); err != nil {
		return models.MaintenanceRequest{}, err
	}
	s.publish(ctx, request, "request.status_changed", "Статус заявки изменен", "Новый статус: "+request.Status)
	return request, nil
}

// AddComment сохраняет историю общения по заявке с проверкой доступа владельца.
func (s *Service) AddComment(ctx context.Context, cmd models.AddRequestCommentCommand) error {
	if strings.TrimSpace(cmd.Actor.UserID) == "" || strings.TrimSpace(cmd.Text) == "" {
		return apperrors.ErrInvalidArgument
	}
	request, err := s.repo.FindByID(ctx, cmd.RequestID)
	if err != nil {
		return err
	}
	if cmd.Actor.Role == "resident" && request.UserID != cmd.Actor.UserID {
		return apperrors.ErrForbidden
	}
	return s.repo.AddComment(ctx, &models.RequestComment{
		RequestID: request.ID,
		UserID:    cmd.Actor.UserID,
		Text:      strings.TrimSpace(cmd.Text),
	})
}

// publish оставляет уведомление о заявке рядом с изменением самой заявки.
func (s *Service) publish(ctx context.Context, request models.MaintenanceRequest, eventType string, title string, body string) {
	s.publishTo(ctx, request.UserID, request, eventType, title, body)
}

// publishTo изолирует best-effort уведомления, чтобы временный сбой notification
// service не откатывал уже сохраненное изменение заявки.
func (s *Service) publishTo(ctx context.Context, userID string, request models.MaintenanceRequest, eventType string, title string, body string) {
	if s.publisher == nil {
		return
	}
	if err := s.publisher.Publish(ctx, models.NotificationEvent{
		UserID:   userID,
		Type:     eventType,
		Title:    title,
		Body:     body,
		EntityID: request.ID,
	}); err != nil {
		s.log.Error().Err(err).Str("request_id", request.ID).Msg("failed to publish request notification")
	}
}

// canManageRequests фиксирует роли, которым можно менять рабочее состояние заявки.
func canManageRequests(role string) bool {
	return role == "admin" || role == "manager" || role == "employee"
}

// validStatus ограничивает жизненный цикл заявки известными состояниями.
func validStatus(status string) bool {
	switch status {
	case models.RequestStatusNew, models.RequestStatusInProgress, models.RequestStatusDone, models.RequestStatusCanceled:
		return true
	default:
		return false
	}
}

// assignRequest хранит пустое назначение как NULL, чтобы PostgreSQL uuid-поле
// не получало некорректную пустую строку.
func assignRequest(request *models.MaintenanceRequest, assignedTo string) {
	assignedTo = strings.TrimSpace(assignedTo)
	if assignedTo == "" {
		request.AssignedTo = nil
		return
	}
	request.AssignedTo = &assignedTo
}
