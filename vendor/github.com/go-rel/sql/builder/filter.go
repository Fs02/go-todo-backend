package builder

import (
	"github.com/go-rel/rel"
)

// Filter builder.
type Filter struct{}

// Write SQL to buffer.
func (f Filter) Write(buffer *Buffer, filter rel.FilterQuery, queryWriter QueryWriter) {
	switch filter.Type {
	case rel.FilterAndOp:
		f.BuildLogical(buffer, "AND", filter.Inner, queryWriter)
	case rel.FilterOrOp:
		f.BuildLogical(buffer, "OR", filter.Inner, queryWriter)
	case rel.FilterNotOp:
		buffer.WriteString("NOT ")
		f.BuildLogical(buffer, "AND", filter.Inner, queryWriter)
	case rel.FilterEqOp,
		rel.FilterNeOp,
		rel.FilterLtOp,
		rel.FilterLteOp,
		rel.FilterGtOp,
		rel.FilterGteOp:
		f.BuildComparison(buffer, filter, queryWriter)
	case rel.FilterNilOp:
		buffer.WriteEscape(filter.Field)
		buffer.WriteString(" IS NULL")
	case rel.FilterNotNilOp:
		buffer.WriteEscape(filter.Field)
		buffer.WriteString(" IS NOT NULL")
	case rel.FilterInOp,
		rel.FilterNinOp:
		f.BuildInclusion(buffer, filter, queryWriter)
	case rel.FilterLikeOp:
		buffer.WriteEscape(filter.Field)
		buffer.WriteString(" LIKE ")
		buffer.WriteValue(filter.Value)
	case rel.FilterNotLikeOp:
		buffer.WriteEscape(filter.Field)
		buffer.WriteString(" NOT LIKE ")
		buffer.WriteValue(filter.Value)
	case rel.FilterFragmentOp:
		buffer.WriteString(filter.Field)
		buffer.AddArguments(filter.Value.([]interface{})...)
	}
}

// BuildLogical SQL to buffer.
func (f Filter) BuildLogical(buffer *Buffer, op string, inner []rel.FilterQuery, queryWriter QueryWriter) {
	var (
		length = len(inner)
	)

	if length > 1 {
		buffer.WriteByte('(')
	}

	for i, c := range inner {
		f.Write(buffer, c, queryWriter)

		if i < length-1 {
			buffer.WriteByte(' ')
			buffer.WriteString(op)
			buffer.WriteByte(' ')
		}
	}

	if length > 1 {
		buffer.WriteByte(')')
	}
}

// BuildComparison SQL to buffer.
func (f Filter) BuildComparison(buffer *Buffer, filter rel.FilterQuery, queryWriter QueryWriter) {
	buffer.WriteEscape(filter.Field)

	switch filter.Type {
	case rel.FilterEqOp:
		buffer.WriteByte('=')
	case rel.FilterNeOp:
		buffer.WriteString("<>")
	case rel.FilterLtOp:
		buffer.WriteByte('<')
	case rel.FilterLteOp:
		buffer.WriteString("<=")
	case rel.FilterGtOp:
		buffer.WriteByte('>')
	case rel.FilterGteOp:
		buffer.WriteString(">=")
	}

	switch v := filter.Value.(type) {
	case rel.SubQuery:
		// For warped sub-queries
		f.buildSubQuery(buffer, v, queryWriter)
	case rel.Query:
		// For sub-queries without warp
		f.buildSubQuery(buffer, rel.SubQuery{Query: v}, queryWriter)
	default:
		// For simple values
		buffer.WriteValue(filter.Value)
	}
}

// BuildInclusion SQL to buffer.
func (f Filter) BuildInclusion(buffer *Buffer, filter rel.FilterQuery, queryWriter QueryWriter) {
	var (
		values = filter.Value.([]interface{})
	)

	buffer.WriteEscape(filter.Field)

	if filter.Type == rel.FilterInOp {
		buffer.WriteString(" IN ")
	} else {
		buffer.WriteString(" NOT IN ")
	}

	f.buildInclusionValues(buffer, values, queryWriter)
}

func (f Filter) buildInclusionValues(buffer *Buffer, values []interface{}, queryWriter QueryWriter) {
	if len(values) == 1 {
		if value, ok := values[0].(rel.Query); ok {
			f.buildSubQuery(buffer, rel.SubQuery{Query: value}, queryWriter)
			return
		}
	}

	buffer.WriteByte('(')
	for i := 0; i < len(values); i++ {
		if i > 0 {
			buffer.WriteByte(',')
		}
		buffer.WriteValue(values[i])
	}
	buffer.WriteByte(')')
}

func (f Filter) buildSubQuery(buffer *Buffer, sub rel.SubQuery, queryWriter QueryWriter) {
	buffer.WriteString(sub.Prefix)
	buffer.WriteByte('(')
	queryWriter.Write(buffer, sub.Query)
	buffer.WriteByte(')')
}
