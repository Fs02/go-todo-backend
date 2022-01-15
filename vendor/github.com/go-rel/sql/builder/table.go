package builder

import (
	"strconv"

	"github.com/go-rel/rel"
)

type ColumnMapper func(*rel.Column) (string, int, int)

// Table builder.
type Table struct {
	BufferFactory BufferFactory
	ColumnMapper  ColumnMapper
}

// Build SQL query for table creation and modification.
func (t Table) Build(table rel.Table) string {
	var (
		buffer = t.BufferFactory.Create()
	)

	switch table.Op {
	case rel.SchemaCreate:
		t.WriteCreateTable(&buffer, table)
	case rel.SchemaAlter:
		t.WriteAlterTable(&buffer, table)
	case rel.SchemaRename:
		t.WriteRenameTable(&buffer, table)
	case rel.SchemaDrop:
		t.WriteDropTable(&buffer, table)
	}

	return buffer.String()
}

// WriteCreateTable query to buffer.
func (t Table) WriteCreateTable(buffer *Buffer, table rel.Table) {
	buffer.WriteString("CREATE TABLE ")

	if table.Optional {
		buffer.WriteString("IF NOT EXISTS ")
	}

	buffer.WriteEscape(table.Name)
	if len(table.Definitions) > 0 {
		buffer.WriteString(" (")

		for i, def := range table.Definitions {
			if i > 0 {
				buffer.WriteString(", ")
			}
			switch v := def.(type) {
			case rel.Column:
				t.WriteColumn(buffer, v)
			case rel.Key:
				t.WriteKey(buffer, v)
			case rel.Raw:
				buffer.WriteString(string(v))
			}
		}

		buffer.WriteByte(')')
	}
	t.WriteOptions(buffer, table.Options)
	buffer.WriteByte(';')
}

// WriteAlterTable query to buffer.
func (t Table) WriteAlterTable(buffer *Buffer, table rel.Table) {
	for _, def := range table.Definitions {
		buffer.WriteString("ALTER TABLE ")
		buffer.WriteEscape(table.Name)
		buffer.WriteByte(' ')

		switch v := def.(type) {
		case rel.Column:
			switch v.Op {
			case rel.SchemaCreate:
				buffer.WriteString("ADD COLUMN ")
				t.WriteColumn(buffer, v)
			case rel.SchemaRename:
				// Add Change
				buffer.WriteString("RENAME COLUMN ")
				buffer.WriteEscape(v.Name)
				buffer.WriteString(" TO ")
				buffer.WriteEscape(v.Rename)
			case rel.SchemaDrop:
				buffer.WriteString("DROP COLUMN ")
				buffer.WriteEscape(v.Name)
			}
		case rel.Key:
			// TODO: Rename and Drop, PR welcomed.
			switch v.Op {
			case rel.SchemaCreate:
				buffer.WriteString("ADD ")
				t.WriteKey(buffer, v)
			}
		}

		t.WriteOptions(buffer, table.Options)
		buffer.WriteByte(';')
	}
}

// WriteRenameTable query to buffer.
func (t Table) WriteRenameTable(buffer *Buffer, table rel.Table) {
	buffer.WriteString("ALTER TABLE ")
	buffer.WriteEscape(table.Name)
	buffer.WriteString(" RENAME TO ")
	buffer.WriteEscape(table.Rename)
	buffer.WriteByte(';')
}

// WriteDropTable query to buffer.
func (t Table) WriteDropTable(buffer *Buffer, table rel.Table) {
	buffer.WriteString("DROP TABLE ")

	if table.Optional {
		buffer.WriteString("IF EXISTS ")
	}

	buffer.WriteEscape(table.Name)
	buffer.WriteByte(';')
}

// WriteColumn definition to buffer.
func (t Table) WriteColumn(buffer *Buffer, column rel.Column) {
	var (
		typ, m, n = t.ColumnMapper(&column)
	)

	buffer.WriteEscape(column.Name)
	buffer.WriteByte(' ')
	buffer.WriteString(typ)

	if m != 0 {
		buffer.WriteByte('(')
		buffer.WriteString(strconv.Itoa(m))

		if n != 0 {
			buffer.WriteByte(',')
			buffer.WriteString(strconv.Itoa(n))
		}

		buffer.WriteByte(')')
	}

	if column.Unsigned {
		buffer.WriteString(" UNSIGNED")
	}

	if column.Unique {
		buffer.WriteString(" UNIQUE")
	}

	if column.Required {
		buffer.WriteString(" NOT NULL")
	}

	if column.Primary {
		buffer.WriteString(" PRIMARY KEY")
	}

	if column.Default != nil {
		buffer.WriteString(" DEFAULT ")
		buffer.WriteValue(column.Default)
	}

	t.WriteOptions(buffer, column.Options)
}

// WriteKey definition to buffer.
func (t Table) WriteKey(buffer *Buffer, key rel.Key) {
	var (
		typ = string(key.Type)
	)

	buffer.WriteString(typ)

	if key.Name != "" {
		buffer.WriteByte(' ')
		buffer.WriteEscape(key.Name)
	}

	buffer.WriteString(" (")
	for i, col := range key.Columns {
		if i > 0 {
			buffer.WriteString(", ")
		}
		buffer.WriteEscape(col)
	}
	buffer.WriteString(")")

	if key.Type == rel.ForeignKey {
		buffer.WriteString(" REFERENCES ")
		buffer.WriteEscape(key.Reference.Table)

		buffer.WriteString(" (")
		for i, col := range key.Reference.Columns {
			if i > 0 {
				buffer.WriteString(", ")
			}
			buffer.WriteEscape(col)
		}
		buffer.WriteString(")")

		if onDelete := key.Reference.OnDelete; onDelete != "" {
			buffer.WriteString(" ON DELETE ")
			buffer.WriteString(onDelete)
		}

		if onUpdate := key.Reference.OnUpdate; onUpdate != "" {
			buffer.WriteString(" ON UPDATE ")
			buffer.WriteString(onUpdate)
		}
	}

	t.WriteOptions(buffer, key.Options)
}

// WriteOptions sql to buffer.
func (t Table) WriteOptions(buffer *Buffer, options string) {
	if options == "" {
		return
	}

	buffer.WriteByte(' ')
	buffer.WriteString(options)
}
