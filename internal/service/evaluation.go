package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"golang.org/x/sync/errgroup"

	"github.com/pumpkinlog/backend/internal/domain"
	"github.com/pumpkinlog/backend/internal/engine"
	"github.com/pumpkinlog/backend/internal/repository"
)

type EvaluationService struct {
	logger *slog.Logger
	ch     *amqp091.Channel
	engine *engine.Engine

	regionRepo     domain.RegionRepository
	ruleRepo       domain.RuleRepository
	answerRepo     domain.AnswerRepository
	evaluationRepo domain.EvaluationRepository
	presenceRepo   domain.PresenceRepository
}

func NewEvaluationService(logger *slog.Logger, conn repository.Connection, ch *amqp091.Channel) domain.EvaluationService {
	return &EvaluationService{
		logger: logger,
		ch:     ch,
		engine: engine.NewEngine(),

		regionRepo:     repository.NewPostgresRegionRepository(conn),
		ruleRepo:       repository.NewPostgresRuleRepository(conn),
		answerRepo:     repository.NewPostgresAnswerRepository(conn),
		evaluationRepo: repository.NewPostgresEvaluationRepository(conn),
		presenceRepo:   repository.NewPostgresPresenceRepository(conn),
	}
}

// EvaluateRegion evaluates a specific region and generates an evaluation profile for the user.
// If an evaluation already exists, it returns the existing evaluation.
//
// A point-in-time (PIT) date can be provided to evaluate the region as of that specific time.
func (s *EvaluationService) EvaluateRegion(ctx context.Context, userID int64, regionID domain.RegionID, opts *domain.EvaluateOpts) (*domain.RegionEvaluation, error) {

	if !opts.Recompute {
		evaluation, err := s.evaluationRepo.GetByID(ctx, userID, regionID)
		if err != nil && !errors.Is(err, domain.ErrNotFound) {
			return nil, fmt.Errorf("get evaluation by ID: %w", err)
		}

		if evaluation != nil {
			return evaluation, nil
		}
	}

	s.logger.Debug("running region evaluation", "userId", userID, "regionId", regionID)

	timestamp := time.Now().UTC()

	if opts.PointInTime.IsZero() {
		opts.PointInTime = timestamp
	}

	evalCtx, err := s.buildEvaluationContext(ctx, userID, regionID, opts.PointInTime)
	if err != nil {
		return nil, fmt.Errorf("load aggregate: %w", err)
	}

	details, passed, err := s.engine.EvaluateRegion(evalCtx)
	if err != nil {
		return nil, fmt.Errorf("evaluation service: evaluate region: %w", err)
	}

	evaluation := &domain.RegionEvaluation{
		RegionID:    regionID,
		UserID:      userID,
		Passed:      passed,
		Nodes:       details,
		PointInTime: opts.PointInTime,
		EvaluatedAt: timestamp,
	}

	if opts.Cache {
		if err := s.evaluationRepo.CreateOrUpdate(ctx, evaluation); err != nil {
			return nil, fmt.Errorf("create or update evaluation: %w", err)
		}
	}

	if opts.Publish {
		body := map[string]any{
			"userId":   userID,
			"regionId": regionID,
		}

		encoded, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal presence: %w", err)
		}

		msg := amqp091.Publishing{
			ContentType: "application/json",
			Body:        encoded,
		}

		if err := s.ch.Publish("evaluation.events", "evaluation.created", false, false, msg); err != nil {
			return nil, fmt.Errorf("publish evaluation created: %w", err)
		}
	}

	return evaluation, nil
}

func (s *EvaluationService) buildEvaluationContext(ctx context.Context, userID int64, regionID domain.RegionID, pit time.Time) (*domain.EvaluationContext, error) {
	g, groupCtx := errgroup.WithContext(ctx)

	var (
		region  *domain.Region
		rules   []*domain.Rule
		answers []*domain.Answer
	)

	g.Go(func() error {
		r, err := s.regionRepo.GetByID(groupCtx, regionID)
		if err != nil {
			return fmt.Errorf("get region: %w", err)
		}
		region = r
		return nil
	})

	g.Go(func() error {
		r, err := s.ruleRepo.ListByRegionID(groupCtx, regionID)
		if err != nil {
			return fmt.Errorf("list rules by region ID: %w", err)
		}
		rules = r
		return nil
	})

	g.Go(func() error {
		a, err := s.answerRepo.ListByRegionID(groupCtx, userID, regionID)
		if err != nil {
			return fmt.Errorf("list answers by region ID: %w", err)
		}
		answers = a
		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, fmt.Errorf("load evaluation context data: %w", err)
	}

	start, end, err := engine.ComputeMaxPeriod(pit, region, rules)
	if err != nil {
		return nil, fmt.Errorf("compute max period: %w", err)
	}

	presences, err := s.presenceRepo.ListByRegionPeriod(ctx, userID, region.ID, start, end)
	if err != nil {
		return nil, fmt.Errorf("list presences by region period: %w", err)
	}

	answerMap := make(map[domain.Code]*domain.Answer, len(answers))
	for _, a := range answers {
		answerMap[a.ConditionID] = a
	}

	ec := &domain.EvaluationContext{
		At:        pit,
		Region:    region,
		Presences: presences,
		Rules:     rules,
		Answers:   answerMap,
	}

	return ec, nil
}
