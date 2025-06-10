package engine

import (
	"fmt"
	"time"

	"github.com/pumpkinlog/backend/internal/domain"
	"github.com/pumpkinlog/backend/internal/engine/strategies"
)

type Engine struct {
	strategies *strategies.Strategies
}

func NewEngine() *Engine {
	return &Engine{
		strategies: strategies.NewStrategies(),
	}
}

// EvaluateRegion evaluates all rules and returns their details plus overall pass status.
func (e *Engine) EvaluateRegion(ctx *domain.EvaluationContext) ([]domain.EvaluationComponent, bool, error) {
	passed := true
	evaluations := make([]domain.EvaluationComponent, len(ctx.Rules))

	for i, rule := range ctx.Rules {
		evaluation, err := e.evaluateRuleNode(rule.Node, ctx)
		if err != nil {
			return nil, false, fmt.Errorf("evaluate rule %s: %w", rule.ID, err)
		}

		evaluations[i] = evaluation

		if passed && !evaluation.IsPassed() {
			passed = false
		}
	}

	return evaluations, passed, nil
}

// evaluateRuleNode dispatches evaluation based on node operator.
func (e *Engine) evaluateRuleNode(node domain.RuleNode, ctx *domain.EvaluationContext) (domain.EvaluationComponent, error) {
	switch node.Type {
	case domain.NodeTypeCompositeAnd, domain.NodeTypeCompositeAny:
		return e.evaluateCompositeNode(node, ctx)
	case domain.NodeTypeStrategy:
		return e.evaluateStrategyNode(node, ctx)
	case domain.NodeTypeCondition:
		return e.evaluateConditionNode(node, ctx)
	default:
		return nil, fmt.Errorf("unsupported node type: %s", node.Type)
	}
}

func (e *Engine) evaluateCompositeNode(node domain.RuleNode, ctx *domain.EvaluationContext) (domain.EvaluationComponent, error) {
	var nodes []domain.RuleNode
	if err := node.Unmarshal(&nodes); err != nil {
		return nil, fmt.Errorf("cannot unmarshal composite node: %w", err)
	}

	evaluations := make([]domain.EvaluationComponent, len(nodes))

	var passed bool
	switch node.Type {
	case domain.NodeTypeCompositeAnd:
		passed = true
	case domain.NodeTypeCompositeAny:
	default:
		return nil, fmt.Errorf("unknown composite node type: %s", node.Type)
	}

	for i, n := range nodes {
		var cn domain.CompositeNode
		if err := n.Unmarshal(&cn); err != nil {
			return nil, fmt.Errorf("cannot unmarshal composite node: %w", err)
		}

		evaluation, err := e.evaluateRuleNode(n, ctx)
		if err != nil {
			return nil, err
		}

		evaluations[i] = evaluation

		switch node.Type {
		case domain.NodeTypeCompositeAnd:
			passed = true
			for _, evaluation := range evaluations {
				if !evaluation.IsPassed() {
					passed = false
					break
				}
			}
		case domain.NodeTypeCompositeAny:
			passed = false
			for _, evaluation := range evaluations {
				if evaluation.IsPassed() {
					passed = true
					break
				}
			}
		default:
			return nil, fmt.Errorf("unknown composite node type: %s", node.Type)
		}
	}

	return &domain.CompositeEvaluation{
		NodeType:   node.Type,
		Status:     domain.EvaluationStatusEvaluated,
		Passed:     passed,
		Components: evaluations,
	}, nil
}

func (e *Engine) evaluateStrategyNode(node domain.RuleNode, ctx *domain.EvaluationContext) (domain.EvaluationComponent, error) {
	var sn domain.EvaluatorNode
	if err := node.Unmarshal(&sn); err != nil {
		return nil, fmt.Errorf("cannot unmarshal strategy node: %w", err)
	}

	start, end, err := ComputePeriod(ctx.At, ctx.Region, sn.Period)
	if err != nil {
		return nil, fmt.Errorf("compute period: %w", err)
	}

	presences := make(map[time.Time]struct{})
	for _, p := range ctx.Presences {
		if !p.Date.Before(start) && !p.Date.After(end) {
			presences[p.Date] = struct{}{}
		}
	}

	se, err := e.strategies.Evaluate(sn.Type, sn.Props, presences)
	if err != nil {
		return nil, fmt.Errorf("evaluate strategy %s: %w", node.Type, err)
	}

	return &domain.StrategyEvaluation{
		Type:      domain.ComponentTypeStrategy,
		Strategy:  sn.Type,
		Passed:    se.Passed,
		Status:    domain.EvaluationStatusEvaluated,
		Reason:    fmt.Sprintf("strategy %s evaluated", node.Type),
		Start:     start,
		End:       end,
		Count:     se.Count,
		Remaining: se.Remaining,
	}, nil
}

func (e *Engine) evaluateConditionNode(node domain.RuleNode, ctx *domain.EvaluationContext) (domain.EvaluationComponent, error) {
	var cn domain.ConditionNode
	if err := node.Unmarshal(&cn); err != nil {
		return nil, fmt.Errorf("cannot unmarshal condition node: %w", err)
	}

	answer, ok := ctx.Answers[cn.ConditionID]
	if !ok || answer.Value == nil {
		return &domain.ConditionEvaluation{
			Type:        domain.ComponentTypeCondition,
			ConditionID: cn.ConditionID,
			Expected:    cn.Equals,
			Comparator:  cn.Comparator,
			Status:      domain.EvaluationStatusUnanswered,
			Reason:      fmt.Sprintf("condition %s not answered", cn.ConditionID),
		}, nil
	}

	passed, err := compareGeneric(cn.Comparator, cn.Equals, answer.Value)
	if err != nil {
		return nil, fmt.Errorf("compare condition %s: %w", cn.ConditionID, err)
	}

	return &domain.ConditionEvaluation{
		Type:        domain.ComponentTypeCondition,
		ConditionID: cn.ConditionID,
		Expected:    cn.Equals,
		Actual:      answer.Value,
		Comparator:  cn.Comparator,
		Status:      domain.EvaluationStatusEvaluated,
		Passed:      passed,
		Reason:      fmt.Sprintf("condition %s evaluated", cn.ConditionID),
	}, nil
}

func compareGeneric(comparator domain.Comparator, expected, actual any) (bool, error) {
	switch comparator {
	case domain.ComparatorEquals:
		return expected == actual, nil
	case domain.ComparatorNotEquals:
		return expected != actual, nil
	default:
		return false, fmt.Errorf("unsupported operator: %s", comparator)
	}
}
