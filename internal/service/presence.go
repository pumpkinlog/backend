package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/rabbitmq/amqp091-go"

	"github.com/pumpkinlog/backend/internal/domain"
	"github.com/pumpkinlog/backend/internal/repository"
)

type PresenceService struct {
	logger *slog.Logger
	ch     *amqp091.Channel

	evaluationRepo domain.EvaluationRepository
	presenceRepo   domain.PresenceRepository
}

func NewPresenceService(logger *slog.Logger, conn repository.Connection, ch *amqp091.Channel) domain.PresenceService {
	return &PresenceService{
		logger: logger,
		ch:     ch,

		evaluationRepo: repository.NewPostgresEvaluationRepository(conn),
		presenceRepo:   repository.NewPostgresPresenceRepository(conn),
	}
}

func (s *PresenceService) GetByID(ctx context.Context, userID int64, regionID domain.RegionID, date time.Time) (*domain.Presence, error) {
	if userID < 0 {
		return nil, fmt.Errorf("%w: user ID cannot be negative", domain.ErrValidation)
	}

	if regionID == "" {
		return nil, fmt.Errorf("%w: region ID cannot be empty", domain.ErrValidation)
	}

	if date.IsZero() {
		return nil, fmt.Errorf("%w: date cannot be empty", domain.ErrValidation)
	}

	return s.presenceRepo.GetByID(ctx, userID, regionID, date)
}

func (s *PresenceService) List(ctx context.Context, userID int64, filter *domain.PresenceFilter) ([]*domain.Presence, error) {
	if userID < 0 {
		return nil, fmt.Errorf("%w: user ID cannot be negative", domain.ErrValidation)
	}

	if filter == nil {
		filter = &domain.PresenceFilter{}
	}

	return s.presenceRepo.List(ctx, userID, filter)
}

func (s *PresenceService) Create(ctx context.Context, userID int64, regionID domain.RegionID, deviceID *int64, start, end time.Time) error {
	if userID < 0 {
		return fmt.Errorf("%w: user ID cannot be negative", domain.ErrValidation)
	}

	if regionID == "" {
		return fmt.Errorf("%w: region ID cannot be empty", domain.ErrValidation)
	}

	if start.IsZero() {
		return fmt.Errorf("%w: start date cannot be empty", domain.ErrValidation)
	}

	if end.IsZero() {
		return fmt.Errorf("%w: end date cannot be empty", domain.ErrValidation)
	}

	if start.After(end) {
		return fmt.Errorf("start cannot be before end: %w", domain.ErrValidation)
	}

	if end.Before(start) {
		return fmt.Errorf("end cannot be before start: %w", domain.ErrValidation)
	}

	if err := s.presenceRepo.CreateRange(ctx, userID, regionID, deviceID, start, end); err != nil {
		return fmt.Errorf("create presence range: %w", err)
	}

	// Delete existing evaluations for the region
	if err := s.evaluationRepo.DeleteByRegionID(ctx, regionID); err != nil {
		return err
	}

	s.logger.Debug("cleared stale evaluations", "regionId", regionID, "userId", userID)

	body := map[string]any{
		"userId":   userID,
		"regionId": regionID,
	}

	encoded, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshal presence: %w", err)
	}

	msg := amqp091.Publishing{
		ContentType: "application/json",
		Body:        encoded,
	}

	if err := s.ch.PublishWithContext(ctx, "presence.events", "presence.create", false, false, msg); err != nil {
		return fmt.Errorf("publish presence: %w", err)
	}

	return nil
}

func (s *PresenceService) Delete(ctx context.Context, userID int64, regionID domain.RegionID, start, end time.Time) error {
	if userID < 0 {
		return fmt.Errorf("%w: user ID cannot be negative", domain.ErrValidation)
	}

	if regionID == "" {
		return fmt.Errorf("%w: region ID cannot be empty", domain.ErrValidation)
	}

	if start.IsZero() {
		return fmt.Errorf("%w: start date cannot be empty", domain.ErrValidation)
	}

	if end.IsZero() {
		return fmt.Errorf("%w: end date cannot be empty", domain.ErrValidation)
	}

	if start.After(end) {
		return fmt.Errorf("start cannot be before end: %w", domain.ErrValidation)
	}

	if end.Before(start) {
		return fmt.Errorf("end cannot be before start: %w", domain.ErrValidation)
	}

	if err := s.presenceRepo.DeleteRange(ctx, userID, regionID, start, end); err != nil {
		return fmt.Errorf("delete presence range: %w", err)
	}

	// Delete existing evaluations for the region
	if err := s.evaluationRepo.DeleteByRegionID(ctx, regionID); err != nil {
		return err
	}

	s.logger.Debug("cleared stale evaluations", "regionId", regionID)

	body := map[string]any{
		"userId":   userID,
		"regionId": regionID,
	}

	encoded, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshal presence: %w", err)
	}

	msg := amqp091.Publishing{
		ContentType: "application/json",
		Body:        encoded,
	}

	if err := s.ch.PublishWithContext(ctx, "presence.events", "presence.create", false, false, msg); err != nil {
		return fmt.Errorf("publish presence: %w", err)
	}

	return nil
}
