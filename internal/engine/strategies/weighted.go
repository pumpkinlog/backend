package strategies

import (
	"encoding/json"
	"fmt"
	"time"
)

// WeightedStrategy is weighted count across a range of years.
type WeightedStrategy struct {
	Threshold int       `json:"threshold"`
	Weights   []float32 `json:"weights"` // 0 = current year, 1 = previous year etc
}

func (s *WeightedStrategy) Evaluate(data []byte, presences map[time.Time]struct{}) (StrategyEvaluation, error) {
	if err := json.Unmarshal(data, &s); err != nil {
		return StrategyEvaluation{}, fmt.Errorf("invalid weighted strategy config: %w", err)
	}

	// Determine the base year (latest year found in presences)
	var baseYear int
	for date := range presences {
		if date.Year() > baseYear {
			baseYear = date.Year()
		}
	}

	// Initialize daysByOffset with zeros for all offsets from weights
	daysByOffset := make(map[int]int)
	for i := range s.Weights {
		daysByOffset[i] = 0
	}

	// Count presence days by offset
	for date := range presences {
		offset := baseYear - date.Year()
		if _, ok := daysByOffset[offset]; ok { // only count if offset is in weights
			daysByOffset[offset]++
		}
	}

	// Calculate weighted total
	var weightedTotal float32
	for i, w := range s.Weights {
		count := daysByOffset[i]
		weightedTotal += float32(count) * w
	}

	// Map to daysByYear including zero days
	daysByYear := make(map[int]int, len(daysByOffset))
	for offset := range daysByOffset {
		year := baseYear - offset
		daysByYear[year] = daysByOffset[offset]
	}
	return StrategyEvaluation{
		Passed:    weightedTotal >= float32(s.Threshold),
		Count:     int(weightedTotal),
		Remaining: int(float32(s.Threshold) - weightedTotal),
		Metadata: map[string]any{
			"baseYear":      baseYear,
			"daysByYear":    daysByYear,
			"weightedTotal": weightedTotal,
		},
	}, nil
}
