package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/pumpkinlog/backend/internal/domain"
	"github.com/pumpkinlog/backend/internal/repository"
)

type RegionService struct {
	logger *slog.Logger

	regionRepo     domain.RegionRepository
	evaluationRepo domain.EvaluationRepository
}

func NewRegionService(logger *slog.Logger, conn repository.Connection) domain.RegionService {
	return &RegionService{
		logger: logger,

		regionRepo:     repository.NewPostgresRegionRepository(conn),
		evaluationRepo: repository.NewPostgresEvaluationRepository(conn),
	}
}

func (s *RegionService) GetByID(ctx context.Context, regionID domain.RegionID) (*domain.Region, error) {
	if regionID == "" {
		return nil, fmt.Errorf("%w: region ID cannot be empty", domain.ErrValidation)
	}

	return s.regionRepo.GetByID(ctx, regionID)
}

func (s *RegionService) List(ctx context.Context, filter *domain.RegionFilter) ([]*domain.Region, error) {
	if filter == nil {
		filter = &domain.RegionFilter{}
	}

	return s.regionRepo.List(ctx, filter)
}

func (s *RegionService) CreateOrUpdate(ctx context.Context, region *domain.Region) error {
	if region.YearStartDay == 0 {
		region.YearStartDay = 1
	}

	if region.YearStartMonth == 0 {
		region.YearStartMonth = 1
	}

	if region.Sources == nil {
		region.Sources = make([]domain.Source, 0)
	}

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
