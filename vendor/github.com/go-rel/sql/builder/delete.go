package builder

import (
	"github.com/go-rel/rel"
)

// Delete builder.
type Delete struct {
	BufferFactory BufferFactory
	Query         QueryWriter
	Filter        Filter
}

// Build SQL query and its arguments.
func (ds Delete) Build(table string, filter rel.FilterQuery) (string, []interface{}) {
	var (
		buffer = ds.BufferFactory.Create()
	)

	buffer.WriteString("DELETE FROM ")
	buffer.WriteEscape(table)

	if !filter.None() {
		buffer.WriteString(" WHERE ")
		ds.Filter.Write(&buffer, table, filter, ds.Query)
	}

	buffer.WriteString(";")

	return buffer.String(), buffer.Arguments()
}
