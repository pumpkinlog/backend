package service

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/pumpkinlog/backend/internal/domain"
	"github.com/pumpkinlog/backend/internal/repository"
)

type ConditionService struct {
	conn repository.Executor

	conditionRepo     domain.ConditionRepository
	ruleConditionRepo domain.RuleConditionRepository
}

func NewConditionService(conn repository.Executor) domain.ConditionService {
	return &ConditionService{
		conditionRepo:     repository.NewPostgresConditionRepository(conn),
		ruleConditionRepo: repository.NewPostgresRuleConditionRepository(conn),
	}
}

func (s *ConditionService) Create(ctx context.Context, condition *domain.Condition, ruleIDs []string) error {

	tx, err := s.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("failed to begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := s.conditionRepo.CreateOrUpdate(ctx, condition); err != nil {
		return fmt.Errorf("failed to create condition: %w", err)
	}

	for _, ruleID := range ruleIDs {
		if err := s.ruleConditionRepo.Link(ctx, ruleID, condition.ID); err != nil {
			return fmt.Errorf("failed to link rule to condition: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit tx: %w", err)
	}

	return nil
}

func (s *ConditionService) Delete(ctx context.Context, conditionID, ruleID string) error {

	tx, err := s.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("failed to begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := s.ruleConditionRepo.Unlink(ctx, ruleID, conditionID); err != nil {
		return fmt.Errorf("failed to unlink rule from condition: %w", err)
	}

	if err := s.conditionRepo.Delete(ctx, conditionID); err != nil {
		return fmt.Errorf("failed to delete condition: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit tx: %w", err)
	}

	return nil
}
