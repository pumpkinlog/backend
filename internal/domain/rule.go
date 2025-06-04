package domain

import (
	"context"
	"fmt"
	"time"
)

type (
	RuleType   string
	PeriodType string
)

const (
	RuleTypeAggregate   RuleType = "aggregate"
	RuleTypeAverage     RuleType = "average"
	RuleTypeWeighted    RuleType = "weighted"
	RuleTypeConsecutive RuleType = "consecutive"

	PeriodTypeYear    PeriodType = "fiscal_year"
	PeriodTypeRolling PeriodType = "rolling"
)

type Rule struct {
	ID             int64      `json:"id"`
	RegionID       string     `json:"regionId"`
	Name           string     `json:"name"`
	Description    string     `json:"description"`
	RuleType       RuleType   `json:"ruleType"`
	PeriodType     PeriodType `json:"periodType"`
	Weight         int        `json:"weight"`
	Threshold      int        `json:"threshold"`
	YearStartMonth time.Month `json:"-"`
	YearStartDay   int        `json:"-"`
	OffsetYears    int        `json:"-"`
	Years          int        `json:"-"`
	RollingDays    int        `json:"-"`
	RollingMonths  int        `json:"-"`
	RollingYears   int        `json:"-"`
}

func (rt RuleType) Valid() bool {
	switch rt {
	case RuleTypeAggregate, RuleTypeAverage, RuleTypeWeighted, RuleTypeConsecutive:
		return true
	default:
		return false
	}
}

func (pt PeriodType) Valid() bool {
	switch pt {
	case PeriodTypeYear, PeriodTypeRolling:
		return true
	default:
		return false
	}
}

func (r *Rule) Validate() error {

	if r.RegionID == "" {
		return fmt.Errorf("region ID cannot be empty: %w", ErrValidation)
	}

	if r.Name == "" {
		return fmt.Errorf("name cannot be empty: %w", ErrValidation)
	}

	if r.RuleType == "" {
		return fmt.Errorf("rule type cannot be empty: %w", ErrValidation)
	}

	if r.PeriodType == "" {
		return fmt.Errorf("period type cannot be empty: %w", ErrValidation)
	}

	if !r.RuleType.Valid() {
		return fmt.Errorf("invalid rule type: %s", r.RuleType)
	}

	if !r.PeriodType.Valid() {
		return fmt.Errorf("invalid period type: %s", r.PeriodType)
	}

	if r.Weight < 0 {
		return fmt.Errorf("weight cannot be negative: %w", ErrValidation)
	}

	if r.Threshold < 0 {
		return fmt.Errorf("threshold cannot be negative: %w", ErrValidation)
	}

	return nil
}

func (r *Rule) Period(asOf time.Time) (time.Time, time.Time, error) {
	switch r.PeriodType {
	case PeriodTypeYear:
		year := asOf.Year()
		boundary := time.Date(year, r.YearStartMonth, r.YearStartDay, 0, 0, 0, 0, asOf.Location())

		if asOf.Before(boundary) {
			year--
		}

		years := r.Years
		if years <= 0 {
			years = 1
		}

		finalYear := year - r.OffsetYears
		startYear := finalYear - (years - 1)

		start := time.Date(startYear, r.YearStartMonth, r.YearStartDay, 0, 0, 0, 0, asOf.Location())
		end := time.Date(finalYear, r.YearStartMonth, r.YearStartDay, 0, 0, 0, 0, asOf.Location()).AddDate(1, 0, 0).Add(-time.Second)

		return start, end, nil
	case PeriodTypeRolling:
		start := asOf.AddDate(-r.RollingYears, -r.RollingMonths, -r.RollingDays)
		end := asOf.Truncate(24*time.Hour).AddDate(0, 0, 1).Add(-time.Second)

		return start, end, nil
	default:
		return time.Time{}, time.Time{}, fmt.Errorf("unsupported period type: %s", r.PeriodType)
	}
}

type RuleService interface {
	CreateOrUpdate(ctx context.Context, rule *Rule) error
}

type RuleFilter struct {
	RegionIDs []string
	Page      *int
	Limit     *int
}

type RuleRepository interface {
	GetByID(ctx context.Context, ruleID int64) (*Rule, error)
	List(ctx context.Context, filter *RuleFilter) ([]*Rule, error)
	ListByRegionID(ctx context.Context, regionID string) ([]*Rule, error)
	CreateOrUpdate(ctx context.Context, rule *Rule) error
}
