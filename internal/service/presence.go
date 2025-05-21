package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/rabbitmq/amqp091-go"

	"github.com/pumpkinlog/backend/internal/domain"
	"github.com/pumpkinlog/backend/internal/repository"
)

type PresenceService struct {
	ch *amqp091.Channel

	presenceRepo domain.PresenceRepository
}

func NewPresenceService(conn repository.Connection, ch *amqp091.Channel) domain.PresenceService {
	return &PresenceService{
		ch: ch,

		presenceRepo: repository.NewPostgresPresenceRepository(conn),
	}
}

func (s *PresenceService) Create(ctx context.Context, userID, regionID string, deviceID *string, start, end time.Time) error {

	if start.After(end) {
		return fmt.Errorf("start cannot be before end: %w", domain.ErrValidation)
	}

	if end.Before(start) {
		return fmt.Errorf("end cannot be before start: %w", domain.ErrValidation)
	}

	if err := s.presenceRepo.CreateRange(ctx, userID, regionID, deviceID, start, end); err != nil {
		return fmt.Errorf("create presence range: %w", err)
	}

	body := map[string]any{
		"userId":   userID,
		"regionId": regionID,
		"start":    start,
		"end":      end,
	}

	encoded, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshal presence: %w", err)
	}

	msg := amqp091.Publishing{
		ContentType: "application/json",
		Body:        encoded,
	}

	if err := s.ch.PublishWithContext(ctx, "presence", "presence.create", false, false, msg); err != nil {
		return fmt.Errorf("publish presence: %w", err)
	}

	return nil
}

func (s *PresenceService) Delete(ctx context.Context, userID, regionID string, start, end time.Time) error {

	if start.After(end) {
		return fmt.Errorf("start cannot be before end: %w", domain.ErrValidation)
	}

	if end.Before(start) {
		return fmt.Errorf("end cannot be before start: %w", domain.ErrValidation)
	}

	if err := s.presenceRepo.DeleteRange(ctx, userID, regionID, start, end); err != nil {
		return fmt.Errorf("delete presence range: %w", err)
	}

	body := map[string]any{
		"userId":   userID,
		"regionId": regionID,
		"start":    start,
		"end":      end,
	}

	encoded, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshal presence: %w", err)
	}

	msg := amqp091.Publishing{
		ContentType: "application/json",
		Body:        encoded,
	}

	if err := s.ch.PublishWithContext(ctx, "presence", "presence.create", false, false, msg); err != nil {
		return fmt.Errorf("publish presence: %w", err)
	}

	return nil
}
