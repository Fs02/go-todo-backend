package builder

import (
	"log"

	"github.com/go-rel/rel"
)

// Index builder.
type Index struct {
	BufferFactory    BufferFactory
	Query            QueryWriter
	Filter           Filter
	DropIndexOnTable bool
	SupportFilter    bool
}

// Build sql query for index.
func (i Index) Build(index rel.Index) string {
	buffer := i.BufferFactory.Create()

	switch index.Op {
	case rel.SchemaCreate:
		i.WriteCreateIndex(&buffer, index)
	case rel.SchemaDrop:
		i.WriteDropIndex(&buffer, index)
	}

	i.WriteOptions(&buffer, index.Options)
	buffer.WriteByte(';')

	return buffer.String()
}

// WriteCreateIndex to buffer
func (i Index) WriteCreateIndex(buffer *Buffer, index rel.Index) {
	buffer.WriteString("CREATE ")
	if index.Unique {
		buffer.WriteString("UNIQUE ")
	}
	buffer.WriteString("INDEX ")

	if index.Optional {
		buffer.WriteString("IF NOT EXISTS ")
	}

	buffer.WriteEscape(index.Name)
	buffer.WriteString(" ON ")
	buffer.WriteEscape(index.Table)

	buffer.WriteString(" (")
	for n, col := range index.Columns {
		if n > 0 {
			buffer.WriteString(", ")
		}
		buffer.WriteEscape(col)
	}
	buffer.WriteString(")")
	if !index.Filter.None() {
		if !i.SupportFilter {
			log.Print("[WARN] Adapter does not support filtered/partial indexes")
			return
		}
		buffer.WriteString(" WHERE ")
		i.Filter.Write(buffer, "", index.Filter, i.Query)
	}
}

// WriteDropIndex to buffer
func (i Index) WriteDropIndex(buffer *Buffer, index rel.Index) {
	buffer.WriteString("DROP INDEX ")

	if index.Optional {
		buffer.WriteString("IF EXISTS ")
	}

	buffer.WriteEscape(index.Name)

	if i.DropIndexOnTable {
		buffer.WriteString(" ON ")
		buffer.WriteEscape(index.Table)
	}
}

// WriteOptions sql to buffer.
func (i Index) WriteOptions(buffer *Buffer, options string) {
	if options == "" {
		return
	}

	buffer.WriteByte(' ')
	buffer.WriteString(options)
}
