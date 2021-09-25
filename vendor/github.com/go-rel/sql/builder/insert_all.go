package builder

import (
	"github.com/go-rel/rel"
)

// InsertAll builder.
type InsertAll struct {
	BufferFactory         BufferFactory
	ReturningPrimaryValue bool
}

// Build SQL string and its arguments.
func (ia InsertAll) Build(table string, primaryField string, fields []string, bulkMutates []map[string]rel.Mutate) (string, []interface{}) {
	var (
		buffer       = ia.BufferFactory.Create()
		fieldsCount  = len(fields)
		mutatesCount = len(bulkMutates)
	)

	buffer.WriteString("INSERT INTO ")
	buffer.WriteEscape(table)
	buffer.WriteString(" (")

	for i := range fields {
		buffer.WriteEscape(fields[i])

		if i < fieldsCount-1 {
			buffer.WriteByte(',')
		}
	}

	buffer.WriteString(") VALUES ")

	for i, mutates := range bulkMutates {
		buffer.WriteByte('(')

		for j, field := range fields {
			if mut, ok := mutates[field]; ok && mut.Type == rel.ChangeSetOp {
				buffer.WriteValue(mut.Value)
			} else {
				buffer.WriteString("DEFAULT")
			}

			if j < fieldsCount-1 {
				buffer.WriteByte(',')
			}
		}

		if i < mutatesCount-1 {
			buffer.WriteString("),")
		} else {
			buffer.WriteByte(')')
		}
	}

	if ia.ReturningPrimaryValue && primaryField != "" {
		buffer.WriteString(" RETURNING ")
		buffer.WriteEscape(primaryField)
	}

	buffer.WriteString(";")

	return buffer.String(), buffer.Arguments()

}
