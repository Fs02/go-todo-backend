package builder

import (
	"github.com/go-rel/rel"
)

// Insert builder.
type Insert struct {
	BufferFactory         BufferFactory
	ReturningPrimaryValue bool
	InsertDefaultValues   bool
}

// Build sql query and its arguments.
func (i Insert) Build(table string, primaryField string, mutates map[string]rel.Mutate) (string, []interface{}) {
	var (
		buffer = i.BufferFactory.Create()
		count  = len(mutates)
	)

	buffer.WriteString("INSERT INTO ")
	buffer.WriteEscape(table)

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

	if i.ReturningPrimaryValue && primaryField != "" {
		buffer.WriteString(" RETURNING ")
		buffer.WriteEscape(primaryField)
	}

	buffer.WriteString(";")

	return buffer.String(), buffer.Arguments()
}
