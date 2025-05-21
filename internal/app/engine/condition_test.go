package engine

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pumpkinlog/backend/internal/domain"
)

func TestCompare(t *testing.T) {
	tests := []struct {
		name     string
		expected any
		actual   any
		cmp      domain.Comparator
		want     bool
	}{
		{
			name:     "equal int",
			expected: 5,
			actual:   5,
			cmp:      domain.ComparatorEquals,
			want:     true,
		},
		{
			name:     "not equal int",
			expected: 5,
			actual:   6,
			cmp:      domain.ComparatorEquals,
			want:     false,
		},
		{
			name:     "equal float32",
			expected: float32(5.0),
			actual:   float32(5.0),
			cmp:      domain.ComparatorEquals,
			want:     true,
		},
		{
			name:     "not equal float32",
			expected: float32(5.0),
			actual:   float32(6.0),
			cmp:      domain.ComparatorEquals,
			want:     false,
		},
		{
			name:     "equal string",
			expected: "hello",
			actual:   "hello",
			cmp:      domain.ComparatorEquals,
			want:     true,
		},
		{
			name:     "not equal string",
			expected: "hello",
			actual:   "world",
			cmp:      domain.ComparatorEquals,
			want:     false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := compare(tc.cmp, tc.expected, tc.actual)
			if got != tc.want {
				require.Equal(t, tc.want, got)
			}
		})
	}
}
