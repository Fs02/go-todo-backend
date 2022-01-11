package builder

import (
	"strconv"
	"strings"

	"github.com/go-rel/rel"
)

type QueryWriter interface {
	Write(buffer *Buffer, query rel.Query)
}

// Query builder.
type Query struct {
	BufferFactory BufferFactory
	Filter        Filter
}

// Build SQL string and it arguments.
func (q Query) Build(query rel.Query) (string, []interface{}) {
	var (
		buffer = q.BufferFactory.Create()
	)

	q.Write(&buffer, query)

	return buffer.String(), buffer.Arguments()
}

// Write SQL to buffer.
func (q Query) Write(buffer *Buffer, query rel.Query) {
	if query.SQLQuery.Statement != "" {
		buffer.WriteString(query.SQLQuery.Statement)
		buffer.AddArguments(query.SQLQuery.Values...)
		return
	}

	rootQuery := buffer.Len() == 0

	q.WriteSelect(buffer, query.Table, query.SelectQuery)
	q.WriteQuery(buffer, query)

	if rootQuery {
		buffer.WriteByte(';')
	}
}

// WriteSelect SQL to buffer.
func (q Query) WriteSelect(buffer *Buffer, table string, selectQuery rel.SelectQuery) {
	if len(selectQuery.Fields) == 0 {
		buffer.WriteString("SELECT ")
		if selectQuery.OnlyDistinct {
			buffer.WriteString("DISTINCT ")
		}
		buffer.WriteField(table, "*")
		return
	}

	buffer.WriteString("SELECT ")

	if selectQuery.OnlyDistinct {
		buffer.WriteString("DISTINCT ")
	}

	l := len(selectQuery.Fields) - 1
	for i, f := range selectQuery.Fields {
		buffer.WriteField(table, f)

		if i < l {
			buffer.WriteByte(',')
		}
	}
}

// WriteQuery SQL to buffer.
func (q Query) WriteQuery(buffer *Buffer, query rel.Query) {
	q.WriteFrom(buffer, query.Table)
	q.WriteJoin(buffer, query.Table, query.JoinQuery)
	q.WriteWhere(buffer, query.Table, query.WhereQuery)

	if len(query.GroupQuery.Fields) > 0 {
		q.WriteGroupBy(buffer, query.Table, query.GroupQuery.Fields)
		q.WriteHaving(buffer, query.Table, query.GroupQuery.Filter)
	}

	q.WriteOrderBy(buffer, query.Table, query.SortQuery)
	q.WriteLimitOffet(buffer, query.LimitQuery, query.OffsetQuery)

	if query.LockQuery != "" {
		buffer.WriteByte(' ')
		buffer.WriteString(string(query.LockQuery))
	}
}

// WriteFrom SQL to buffer.
func (q Query) WriteFrom(buffer *Buffer, table string) {
	buffer.WriteString(" FROM ")
	buffer.WriteEscape(table)
}

// WriteJoin SQL to buffer.
func (q Query) WriteJoin(buffer *Buffer, table string, joins []rel.JoinQuery) {
	if len(joins) == 0 {
		return
	}

	for _, join := range joins {
		var (
			from = join.From
			to   = join.To
		)

		// TODO: move this to core functionality, and infer join condition using assoc data.
		if join.Arguments == nil && (join.From == "" || join.To == "") {
			from = table + "." + strings.TrimSuffix(join.Table, "s") + "_id"
			to = join.Table + ".id"
		}

		buffer.WriteByte(' ')
		buffer.WriteString(join.Mode)
		buffer.WriteByte(' ')

		if join.Table != "" {
			buffer.WriteEscape(join.Table)
			buffer.WriteString(" ON ")
			buffer.WriteEscape(from)
			buffer.WriteString("=")
			buffer.WriteEscape(to)
			if !join.Filter.None() {
				buffer.WriteString(" AND ")
				q.Filter.Write(buffer, join.Table, join.Filter, q)
			}
		}

		buffer.AddArguments(join.Arguments...)
	}
}

// WriteWhere SQL to buffer.
func (q Query) WriteWhere(buffer *Buffer, table string, filter rel.FilterQuery) {
	if filter.None() {
		return
	}

	buffer.WriteString(" WHERE ")
	q.Filter.Write(buffer, table, filter, q)
}

// WriteGroupBy SQL to buffer.
func (q Query) WriteGroupBy(buffer *Buffer, table string, fields []string) {
	buffer.WriteString(" GROUP BY ")

	l := len(fields) - 1
	for i, f := range fields {
		buffer.WriteField(table, f)

		if i < l {
			buffer.WriteByte(',')
		}
	}
}

// WriteHaving SQL to buffer.
func (q Query) WriteHaving(buffer *Buffer, table string, filter rel.FilterQuery) {
	if filter.None() {
		return
	}

	buffer.WriteString(" HAVING ")
	q.Filter.Write(buffer, table, filter, q)
}

// WriteOrderBy SQL to buffer.
func (q Query) WriteOrderBy(buffer *Buffer, table string, orders []rel.SortQuery) {
	var (
		length = len(orders)
	)

	if length == 0 {
		return
	}

	buffer.WriteString(" ORDER BY ")
	for i, order := range orders {
		if i > 0 {
			buffer.WriteString(", ")
		}

		buffer.WriteField(table, order.Field)

		if order.Asc() {
			buffer.WriteString(" ASC")
		} else {
			buffer.WriteString(" DESC")
		}
	}
}

// WriteLimitOffet SQL to buffer.
func (q Query) WriteLimitOffet(buffer *Buffer, limit rel.Limit, offset rel.Offset) {
	if limit > 0 {
		buffer.WriteString(" LIMIT ")
		buffer.WriteString(strconv.Itoa(int(limit)))

		if offset > 0 {
			buffer.WriteString(" OFFSET ")
			buffer.WriteString(strconv.Itoa(int(offset)))
		}
	}
}
