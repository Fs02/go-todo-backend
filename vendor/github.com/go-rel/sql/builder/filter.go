package builder

import (
	"github.com/go-rel/rel"
)

// Filter builder.
type Filter struct{}

// Write SQL to buffer.
func (f Filter) Write(buffer *Buffer, table string, filter rel.FilterQuery, queryWriter QueryWriter) {
	switch filter.Type {
	case rel.FilterAndOp:
		f.WriteLogical(buffer, table, "AND", filter.Inner, queryWriter)
	case rel.FilterOrOp:
		f.WriteLogical(buffer, table, "OR", filter.Inner, queryWriter)
	case rel.FilterNotOp:
		buffer.WriteString("NOT ")
		f.WriteLogical(buffer, table, "AND", filter.Inner, queryWriter)
	case rel.FilterEqOp,
		rel.FilterNeOp,
		rel.FilterLtOp,
		rel.FilterLteOp,
		rel.FilterGtOp,
		rel.FilterGteOp:
		f.WriteComparison(buffer, table, filter, queryWriter)
	case rel.FilterNilOp:
		buffer.WriteField(table, filter.Field)
		buffer.WriteString(" IS NULL")
	case rel.FilterNotNilOp:
		buffer.WriteField(table, filter.Field)
		buffer.WriteString(" IS NOT NULL")
	case rel.FilterInOp,
		rel.FilterNinOp:
		f.WriteInclusion(buffer, table, filter, queryWriter)
	case rel.FilterLikeOp:
		buffer.WriteField(table, filter.Field)
		buffer.WriteString(" LIKE ")
		buffer.WriteValue(filter.Value)
	case rel.FilterNotLikeOp:
		buffer.WriteField(table, filter.Field)
		buffer.WriteString(" NOT LIKE ")
		buffer.WriteValue(filter.Value)
	case rel.FilterFragmentOp:
		buffer.WriteString(filter.Field)
		if !buffer.InlineValues {
			buffer.AddArguments(filter.Value.([]interface{})...)
		}
	}
}

// WriteLogical SQL to buffer.
func (f Filter) WriteLogical(buffer *Buffer, table, op string, inner []rel.FilterQuery, queryWriter QueryWriter) {
	var (
		length = len(inner)
	)

	if length > 1 {
		buffer.WriteByte('(')
	}

	for i, c := range inner {
		f.Write(buffer, table, c, queryWriter)

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

// WriteComparison SQL to buffer.
func (f Filter) WriteComparison(buffer *Buffer, table string, filter rel.FilterQuery, queryWriter QueryWriter) {
	buffer.WriteField(table, filter.Field)

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
		f.WriteSubQuery(buffer, v, queryWriter)
	case rel.Query:
		// For sub-queries without warp
		f.WriteSubQuery(buffer, rel.SubQuery{Query: v}, queryWriter)
	default:
		// For simple values
		buffer.WriteValue(filter.Value)
	}
}

// WriteInclusion SQL to buffer.
func (f Filter) WriteInclusion(buffer *Buffer, table string, filter rel.FilterQuery, queryWriter QueryWriter) {
	var (
		values = filter.Value.([]interface{})
	)

	if len(values) == 0 {
		if filter.Type == rel.FilterInOp {
			buffer.WriteString("1=0")
		} else {
			buffer.WriteString("1=1")
		}
	} else {
		buffer.WriteField(table, filter.Field)

		if filter.Type == rel.FilterInOp {
			buffer.WriteString(" IN ")
		} else {
			buffer.WriteString(" NOT IN ")
		}

		f.WriteInclusionValues(buffer, values, queryWriter)
	}
}

func (f Filter) WriteInclusionValues(buffer *Buffer, values []interface{}, queryWriter QueryWriter) {
	if len(values) == 1 {
		if value, ok := values[0].(rel.Query); ok {
			f.WriteSubQuery(buffer, rel.SubQuery{Query: value}, queryWriter)
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

func (f Filter) WriteSubQuery(buffer *Buffer, sub rel.SubQuery, queryWriter QueryWriter) {
	buffer.WriteString(sub.Prefix)
	buffer.WriteByte('(')
	queryWriter.Write(buffer, sub.Query)
	buffer.WriteByte(')')
}
