package service

import (
	"context"
	"strings"

	apperrors "github.com/chekrzk/ufanet-home-manager/requests-service/internal/errors"
	"github.com/chekrzk/ufanet-home-manager/requests-service/internal/models"
	"github.com/rs/zerolog"
)

type Service struct {
	repo      RequestRepository
	publisher NotificationPublisher
	log       zerolog.Logger
}

func New(repo RequestRepository, publisher NotificationPublisher, log zerolog.Logger) *Service {
	return &Service{repo: repo, publisher: publisher, log: log}
}

func (s *Service) Create(ctx context.Context, cmd models.CreateRequestCommand) (models.MaintenanceRequest, error) {
	if strings.TrimSpace(cmd.User.UserID) == "" || strings.TrimSpace(cmd.Category) == "" || strings.TrimSpace(cmd.Description) == "" {
		return models.MaintenanceRequest{}, apperrors.ErrInvalidArgument
	}
	request := models.MaintenanceRequest{
		UserID:      cmd.User.UserID,
		Category:    strings.TrimSpace(cmd.Category),
		Description: strings.TrimSpace(cmd.Description),
		Status:      models.RequestStatusNew,
	}
	if err := s.repo.Create(ctx, &request); err != nil {
		return models.MaintenanceRequest{}, err
	}
	return request, nil
}

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

func (s *Service) UpdateStatus(ctx context.Context, cmd models.UpdateRequestStatusCommand) (models.MaintenanceRequest, error) {
	if !canManageRequests(cmd.Actor.Role) || strings.TrimSpace(cmd.Status) == "" {
		return models.MaintenanceRequest{}, apperrors.ErrForbidden
	}
	request, err := s.repo.FindByID(ctx, cmd.RequestID)
	if err != nil {
		return models.MaintenanceRequest{}, err
	}
	request.Status = strings.TrimSpace(cmd.Status)
	request.AssignedTo = strings.TrimSpace(cmd.AssignedTo)
	if err := s.repo.Save(ctx, &request); err != nil {
		return models.MaintenanceRequest{}, err
	}
	s.publish(ctx, request, "request.status_changed", "Статус заявки изменен", "Новый статус: "+request.Status)
	return request, nil
}

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

func (s *Service) publish(ctx context.Context, request models.MaintenanceRequest, eventType string, title string, body string) {
	if s.publisher == nil {
		return
	}
	if err := s.publisher.Publish(ctx, models.NotificationEvent{
		UserID:   request.UserID,
		Type:     eventType,
		Title:    title,
		Body:     body,
		EntityID: request.ID,
	}); err != nil {
		s.log.Error().Err(err).Str("request_id", request.ID).Msg("failed to publish request notification")
	}
}

func canManageRequests(role string) bool {
	return role == "admin" || role == "manager" || role == "employee"
}
