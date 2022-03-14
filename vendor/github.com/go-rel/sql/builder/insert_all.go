package builder

import (
	"github.com/go-rel/rel"
)

// InsertAll builder.
type InsertAll struct {
	BufferFactory         BufferFactory
	ReturningPrimaryValue bool
	OnConflict            OnConflict
}

// Build SQL string and its arguments.
func (ia InsertAll) Build(table string, primaryField string, fields []string, bulkMutates []map[string]rel.Mutate, onConflict rel.OnConflict) (string, []interface{}) {
	var (
		buffer = ia.BufferFactory.Create()
	)

	ia.WriteInsertInto(&buffer, table)
	ia.WriteValues(&buffer, fields, bulkMutates)
	ia.OnConflict.Write(&buffer, fields, onConflict)
	ia.WriteReturning(&buffer, primaryField)
	buffer.WriteString(";")

	return buffer.String(), buffer.Arguments()
}

func (ia InsertAll) WriteInsertInto(buffer *Buffer, table string) {
	buffer.WriteString("INSERT INTO ")
	buffer.WriteEscape(table)
}

func (ia InsertAll) WriteValues(buffer *Buffer, fields []string, bulkMutates []map[string]rel.Mutate) {
	var (
		fieldsCount  = len(fields)
		mutatesCount = len(bulkMutates)
	)

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
}

func (ia InsertAll) WriteReturning(buffer *Buffer, primaryField string) {
	if ia.ReturningPrimaryValue && primaryField != "" {
		buffer.WriteString(" RETURNING ")
		buffer.WriteEscape(primaryField)
	}
}
