package service

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/pumpkinlog/backend/internal/domain"
	"github.com/pumpkinlog/backend/internal/repository"
)

type RegionService struct {
	logger *slog.Logger

	regionRepo     domain.RegionRepository
	evaluationRepo domain.EvaluationRepository
}

func NewRegionService(logger *slog.Logger, conn *pgxpool.Pool) domain.RegionService {
	return &RegionService{
		regionRepo:     repository.NewPostgresRegionRepository(conn),
		evaluationRepo: repository.NewPostgresEvaluationRepository(conn),
	}
}

func (s *RegionService) CreateOrUpdate(ctx context.Context, region *domain.Region) error {

	if err := region.Validate(); err != nil {
		return err
	}

	// @TODO: Tx?

	if err := s.regionRepo.CreateOrUpdate(ctx, region); err != nil {
		return err
	}

	// Delete existing evaluations for the region
	if err := s.evaluationRepo.DeleteByRegionID(ctx, region.ID); err != nil {
		return err
	}

	s.logger.Debug("cleared stale evaluations", "regionId", region.ID)

	return nil
}
