package seed

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/pumpkinlog/backend/internal/domain"
	"github.com/pumpkinlog/backend/internal/service"
)

type Region struct {
	domain.Region
	Rules      []domain.Rule
	Conditions []domain.Condition
}

type SeedData []Region

type Seeder struct {
	logger *slog.Logger
	conn   *pgxpool.Pool
}

func NewSeeder(logger *slog.Logger, conn *pgxpool.Pool) *Seeder {
	return &Seeder{logger, conn}
}

func (s *Seeder) Seed(ctx context.Context, fileName string) error {
	timestamp := time.Now()

	tx, err := s.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("cannot begin tx: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			s.logger.Error("cannot rollback seed tx", "error", err)
		}
	}()

	regionSvc := service.NewRegionService(s.logger, tx)
	ruleSvc := service.NewRuleService(s.logger, tx)
	conditionSvc := service.NewConditionService(s.logger, tx)

	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("cannot get working dir: %w", err)
	}

	path := filepath.Join(wd, fileName)

	file, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("cannot read seed data file: %w", err)
	}

	var regions []Region
	if err := json.Unmarshal(file, &regions); err != nil {
		return fmt.Errorf("cannot unmarshal seed data: %w", err)
	}

	for _, region := range regions {
		if err := regionSvc.CreateOrUpdate(ctx, &region.Region); err != nil {
			return fmt.Errorf("cannot upsert region: %w", err)
		}

		for _, rule := range region.Rules {
			rule.RegionID = region.ID
			rule.ID = domain.Code(strings.ToUpper(string(rule.ID)))

			if err := validID(region.ID, rule.ID); err != nil {
				return fmt.Errorf("cannot validate rule ID: %w", err)
			}

			if err := ruleSvc.CreateOrUpdate(ctx, &rule); err != nil {
				s.logger.Error("cannot upsert rule", "region", region.ID, "id", rule.ID)
				return fmt.Errorf("cannot upsert rule: %w", err)
			}
		}

		for _, condition := range region.Conditions {
			condition.RegionID = region.ID
			condition.ID = domain.Code(strings.ToUpper(string(condition.ID)))

			if err := validID(region.ID, condition.ID); err != nil {
				return fmt.Errorf("cannot validate condition ID: %w", err)
			}

			if err := conditionSvc.CreateOrUpdate(ctx, &condition); err != nil {
				s.logger.Error("cannot upsert condition", "region", region.ID, "id", condition.ID)
				return fmt.Errorf("cannot upsert condition: %w", err)
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("cannot commit tx: %w", err)
	}

	s.logger.Info("seeding complete", "duration", time.Since(timestamp).String())

	return nil
}

func validID(regionID domain.RegionID, code domain.Code) error {
	prefix := fmt.Sprintf("%s_", regionID)

	if !strings.HasPrefix(string(code), prefix) {
		return fmt.Errorf("ID %s for region %s must be prefixed with %s", code, regionID, prefix)
	}

	return nil
}
