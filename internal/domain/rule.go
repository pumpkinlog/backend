package domain

import (
	"context"
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
	ID             string     `json:"id"`
	RegionID       string     `json:"regionId"`
	Name           string     `json:"name"`
	Description    string     `json:"description"`
	RuleType       RuleType   `json:"ruleType"`
	PeriodType     PeriodType `json:"periodType"`
	Threshold      int        `json:"threshold"`
	YearStartMonth time.Month `json:"-"`
	YearStartDay   int        `json:"-"`
	OffsetYears    int        `json:"-"`
	Years          int        `json:"-"`
	RollingDays    int        `json:"-"`
	RollingMonths  int        `json:"-"`
	RollingYears   int        `json:"-"`
}

type RuleFilter struct {
	RegionIDs []string
}

type RuleRepository interface {
	GetByID(ctx context.Context, id string) (*Rule, error)
	List(ctx context.Context, filter *RuleFilter) ([]*Rule, error)
	CreateOrUpdate(ctx context.Context, rule *Rule) error
}
