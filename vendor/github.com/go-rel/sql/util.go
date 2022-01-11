package sql

import (
	"strings"
	"time"

	"github.com/go-rel/rel"
)

// DefaultTimeLayout default time layout.
const DefaultTimeLayout = "2006-01-02 15:04:05"

// ColumnMapper function.
func ColumnMapper(column *rel.Column) (string, int, int) {
	var (
		typ        string
		m, n       int
		timeLayout = DefaultTimeLayout
	)

	switch column.Type {
	case rel.ID:
		typ = "INT UNSIGNED AUTO_INCREMENT"
	case rel.BigID:
		typ = "BIGINT UNSIGNED AUTO_INCREMENT"
	case rel.Bool:
		typ = "BOOL"
	case rel.Int:
		typ = "INT"
		m = column.Limit
	case rel.BigInt:
		typ = "BIGINT"
		m = column.Limit
	case rel.Float:
		typ = "FLOAT"
		m = column.Precision
	case rel.Decimal:
		typ = "DECIMAL"
		m = column.Precision
		n = column.Scale
	case rel.String:
		typ = "VARCHAR"
		m = column.Limit
		if m == 0 {
			m = 255
		}
	case rel.Text:
		typ = "TEXT"
		m = column.Limit
	case rel.JSON:
		typ = "TEXT"
	case rel.Date:
		typ = "DATE"
		timeLayout = "2006-01-02"
	case rel.DateTime:
		typ = "DATETIME"
	case rel.Time:
		typ = "TIME"
		timeLayout = "15:04:05"
	default:
		typ = string(column.Type)
	}

	if t, ok := column.Default.(time.Time); ok {
		column.Default = t.Format(timeLayout)
	}

	return typ, m, n
}

// ExtractString between two string.
func ExtractString(s, left, right string) string {
	var (
		start = strings.Index(s, left)
		end   = strings.LastIndex(s, right)
	)

	if start < 0 || end < 0 || start+len(left) >= end {
		return s
	}

	return s[start+len(left) : end]
}

func toInt64(i interface{}) int64 {
	var result int64

	switch s := i.(type) {
	case int:
		result = int64(s)
	case int64:
		result = s
	case int32:
		result = int64(s)
	case int16:
		result = int64(s)
	case int8:
		result = int64(s)
	case uint:
		result = int64(s)
	case uint64:
		result = int64(s)
	case uint32:
		result = int64(s)
	case uint16:
		result = int64(s)
	case uint8:
		result = int64(s)
	}

	return result
}
