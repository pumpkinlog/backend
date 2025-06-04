package sqlbuild

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBuilder(t *testing.T) {
	tests := []struct {
		name         string
		baseSQL      string
		build        func(b *Builder)
		expectedSQL  string
		expectedArgs []any
	}{
		{
			name:         "No conditions",
			baseSQL:      "SELECT * FROM rules",
			build:        func(b *Builder) {},
			expectedSQL:  "SELECT * FROM rules",
			expectedArgs: nil,
		},
		{
			name:    "Single condition",
			baseSQL: "SELECT * FROM rules",
			build: func(b *Builder) {
				b.AddCondition("name", "=", "foo")
			},
			expectedSQL:  "SELECT * FROM rules WHERE name = $1",
			expectedArgs: []any{"foo"},
		},
		{
			name:    "IN condition",
			baseSQL: "SELECT * FROM rules",
			build: func(b *Builder) {
				b.AddInCondition("region_id", []any{1, 2, 3})
			},
			expectedSQL:  "SELECT * FROM rules WHERE region_id IN ($1, $2, $3)",
			expectedArgs: []any{1, 2, 3},
		},
		{
			name:    "Combined conditions",
			baseSQL: "SELECT * FROM rules",
			build: func(b *Builder) {
				b.AddCondition("name", "=", "foo")
				b.AddInCondition("region_id", []any{10, 20})
				b.AddCondition("enabled", "=", true)
			},
			expectedSQL:  "SELECT * FROM rules WHERE name = $1 AND region_id IN ($2, $3) AND enabled = $4",
			expectedArgs: []any{"foo", 10, 20, true},
		},
		{
			name:    "Empty IN condition skipped",
			baseSQL: "SELECT * FROM rules",
			build: func(b *Builder) {
				b.AddInCondition("region_id", []any{})
			},
			expectedSQL:  "SELECT * FROM rules",
			expectedArgs: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			builder := NewBuilder(tc.baseSQL)
			tc.build(builder)

			require.Equal(t, tc.expectedSQL, builder.SQL())
			require.Equal(t, tc.expectedArgs, builder.Args())
		})
	}
}
