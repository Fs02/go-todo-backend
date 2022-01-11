package builder

import (
	"database/sql/driver"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-rel/sql"
)

// UnescapeCharacter disable field escaping when it starts with this character.
var UnescapeCharacter byte = '^'

var escapeCache sync.Map

type escapeCacheKey struct {
	table  string
	value  string
	quoter Quoter
}

// Buffer is used to build query string.
type Buffer struct {
	strings.Builder
	Quoter              Quoter
	ValueConverter      driver.ValueConverter
	ArgumentPlaceholder string
	ArgumentOrdinal     bool
	InlineValues        bool
	BoolTrueValue       string
	BoolFalseValue      string
	valueCount          int
	arguments           []interface{}
}

// WriteValue query placeholder and append value to argument.
func (b *Buffer) WriteValue(value interface{}) {
	if !b.InlineValues {
		b.WritePlaceholder()
		b.arguments = append(b.arguments, value)
		return
	}

	// Detect float bits to not lose precision after converting to float64
	var floatBits = 64
	if value != nil && reflect.TypeOf(value).Kind() == reflect.Float32 {
		floatBits = 32
	}

	if v, err := b.ValueConverter.ConvertValue(value); err != nil {
		log.Printf("[WARN] unsupported inline value %v: %v", value, err)
	} else {
		value = v
	}

	if value == nil {
		b.WriteString("NULL")
		return
	}

	switch v := value.(type) {
	case string:
		b.WriteString(b.Quoter.Value(v))
		return
	case []byte:
		b.WriteString(b.Quoter.Value(string(v)))
		return
	case time.Time:
		b.WriteString(b.Quoter.Value(v.Format(sql.DefaultTimeLayout)))
		return
	}

	rv := reflect.ValueOf(value)
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		b.WriteString(strconv.FormatInt(rv.Int(), 10))
		return
	case reflect.Float32, reflect.Float64:
		b.WriteString(strconv.FormatFloat(rv.Float(), 'g', -1, floatBits))
		return
	case reflect.Bool:
		if rv.Bool() {
			b.WriteString(b.BoolTrueValue)
		} else {
			b.WriteString(b.BoolFalseValue)
		}
		return
	}
	b.WriteString(fmt.Sprintf("%v", value))
}

// WritePlaceholder without adding argument.
// argument can be added later using AddArguments function.
func (b *Buffer) WritePlaceholder() {
	b.valueCount++
	b.WriteString(b.ArgumentPlaceholder)
	if b.ArgumentOrdinal {
		b.WriteString(strconv.Itoa(b.valueCount))
	}
}

// WriteField writes table and field name.
func (b *Buffer) WriteField(table, field string) {
	b.WriteString(b.escape(table, field))
}

// WriteEscape string.
func (b *Buffer) WriteEscape(value string) {
	b.WriteString(b.escape("", value))
}

func (b Buffer) escape(table, value string) string {
	if value == "*" {
		if table == "" {
			return value
		}
		return b.Quoter.ID(table) + ".*"
	}

	key := escapeCacheKey{table: table, value: value, quoter: b.Quoter}
	escapedValue, ok := escapeCache.Load(key)
	if ok {
		return escapedValue.(string)
	}

	if len(value) > 0 && value[0] == UnescapeCharacter {
		escapedValue = value[1:]
	} else if _, err := strconv.Atoi(value); err == nil {
		escapedValue = value
	} else if i := strings.Index(strings.ToLower(value), " as "); i > -1 {
		escapedValue = b.escape(table, value[:i]) + " AS " + b.escape("", value[i+4:])
	} else if start, end := strings.IndexRune(value, '('), strings.IndexRune(value, ')'); start >= 0 && end >= 0 && end > start {
		escapedValue = value[:start+1] + b.escape(table, value[start+1:end]) + value[end:]
	} else {
		parts := strings.Split(value, ".")
		if len(parts) == 1 && table != "" {
			parts = []string{table, parts[0]}
		}
		for i, part := range parts {
			part = strings.TrimSpace(part)
			if part == "*" && i == len(parts)-1 {
				break
			}
			parts[i] = b.Quoter.ID(part)
		}
		escapedValue = strings.Join(parts, ".")
	}

	escapeCache.Store(key, escapedValue)
	return escapedValue.(string)
}

// AddArguments appends multiple arguments without writing placeholder query..
func (b *Buffer) AddArguments(args ...interface{}) {
	if b.arguments == nil {
		b.arguments = args
	} else {
		b.arguments = append(b.arguments, args...)
	}
}

func (b Buffer) Arguments() []interface{} {
	return b.arguments
}

// Reset buffer.
func (b *Buffer) Reset() {
	b.Builder.Reset()
	b.valueCount = 0
	b.arguments = nil
}

// BufferFactory is used to create buffer based on shared settings.
type BufferFactory struct {
	Quoter              Quoter
	ValueConverter      driver.ValueConverter
	ArgumentPlaceholder string
	ArgumentOrdinal     bool
	InlineValues        bool
	BoolTrueValue       string
	BoolFalseValue      string
}

func (bf BufferFactory) Create() Buffer {
	conv := bf.ValueConverter
	if conv == nil {
		conv = driver.DefaultParameterConverter
	}
	return Buffer{
		Quoter:              bf.Quoter,
		ValueConverter:      conv,
		ArgumentPlaceholder: bf.ArgumentPlaceholder,
		ArgumentOrdinal:     bf.ArgumentOrdinal,
		InlineValues:        bf.InlineValues,
		BoolTrueValue:       bf.BoolTrueValue,
		BoolFalseValue:      bf.BoolFalseValue,
	}
}
