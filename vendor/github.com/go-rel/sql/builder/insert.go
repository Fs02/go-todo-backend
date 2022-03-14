package builder

import (
	"github.com/go-rel/rel"
)

// Insert builder.
type Insert struct {
	BufferFactory         BufferFactory
	ReturningPrimaryValue bool
	InsertDefaultValues   bool
	OnConflict            OnConflict
}

// Build sql query and its arguments.
func (i Insert) Build(table string, primaryField string, mutates map[string]rel.Mutate, onConflict rel.OnConflict) (string, []interface{}) {
	var (
		buffer = i.BufferFactory.Create()
	)

	i.WriteInsertInto(&buffer, table)
	i.WriteValues(&buffer, mutates)
	i.OnConflict.WriteMutates(&buffer, mutates, onConflict)
	i.WriteReturning(&buffer, primaryField)

	buffer.WriteString(";")

	return buffer.String(), buffer.Arguments()
}

func (i Insert) WriteInsertInto(buffer *Buffer, table string) {
	buffer.WriteString("INSERT INTO ")
	buffer.WriteEscape(table)
}

func (i Insert) WriteValues(buffer *Buffer, mutates map[string]rel.Mutate) {
	var (
		count = len(mutates)
	)

	if count == 0 && i.InsertDefaultValues {
		buffer.WriteString(" DEFAULT VALUES")
	} else {
		buffer.WriteString(" (")

		var (
			n         = 0
			arguments = make([]interface{}, 0, count)
		)

		for field, mut := range mutates {
			if mut.Type == rel.ChangeSetOp {
				if n > 0 {
					buffer.WriteByte(',')
				}

				buffer.WriteEscape(field)
				arguments = append(arguments, mut.Value)
				n++
			}
		}

		buffer.WriteString(") VALUES (")

		for i := range arguments {
			if i > 0 {
				buffer.WriteByte(',')
			}

			buffer.WritePlaceholder()
		}
		buffer.AddArguments(arguments...)
		buffer.WriteByte(')')
	}
}

func (i Insert) WriteReturning(buffer *Buffer, primaryField string) {
	if i.ReturningPrimaryValue && primaryField != "" {
		buffer.WriteString(" RETURNING ")
		buffer.WriteEscape(primaryField)
	}
}
