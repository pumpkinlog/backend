package sqlbuild

import (
	"fmt"
	"strings"
)

type Builder struct {
	baseSQL  string
	where    strings.Builder
	args     []any
	argIndex int
}

// NewBuilder initializes the builder with the base SQL query (e.g., SELECT ... FROM ...)
func NewBuilder(baseSQL string) *Builder {
	return &Builder{
		baseSQL:  baseSQL,
		argIndex: 1,
	}
}

// AddCondition appends a condition in the form of "field op $n" with the associated value.
// Example: AddCondition("name", "=", "foo") → WHERE name = $1
func (b *Builder) AddCondition(field, op string, value any) *Builder {
	if b.argIndex == 1 {
		b.where.WriteString(" WHERE ")
	} else {
		b.where.WriteString(" AND ")
	}
	b.where.WriteString(fmt.Sprintf("%s %s $%d", field, op, b.argIndex))
	b.args = append(b.args, value)
	b.argIndex++
	return b
}

// AddInCondition appends a `field IN (...)` clause with appropriate placeholders and args.
// It skips if `values` is empty.
func (b *Builder) AddInCondition(field string, values []any) *Builder {
	if len(values) == 0 {
		return b
	}
	if b.argIndex == 1 {
		b.where.WriteString(" WHERE ")
	} else {
		b.where.WriteString(" AND ")
	}
	b.where.WriteString(fmt.Sprintf("%s IN (", field))
	for i := range values {
		if i > 0 {
			b.where.WriteString(", ")
		}
		b.where.WriteString(fmt.Sprintf("$%d", b.argIndex))
		b.args = append(b.args, values[i])
		b.argIndex++
	}
	b.where.WriteString(")")
	return b
}

// SQL returns the final SQL query string with base SQL + WHERE clause.
func (b *Builder) SQL() string {
	return b.baseSQL + b.where.String()
}

// Args returns the collected parameters for use in db.Query/Exec.
func (b *Builder) Args() []any {
	return b.args
}

// AddConditionInt64 is a convenience function to add a condition for any value.
func AddInCondition[T any](b *Builder, field string, values []T) {
	vals := make([]any, len(values))
	for i, v := range values {
		vals[i] = v
	}
	b.AddInCondition(field, vals)
}
