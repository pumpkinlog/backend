package domain

import "time"

type TimeWindow struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}
