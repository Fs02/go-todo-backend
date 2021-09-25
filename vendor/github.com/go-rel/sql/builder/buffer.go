package builder

import (
	"strconv"
	"strings"
	"sync"
)

// UnescapeCharacter disable field escaping when it starts with this character.
var UnescapeCharacter byte = '^'

var escapeCache sync.Map

type escapeCacheKey struct {
	value  string
	prefix string
	suffix string
}

// Buffer is used to build query string.
type Buffer struct {
	strings.Builder
	ArgumentPlaceholder string
	ArgumentOrdinal     bool
	EscapePrefix        string
	EscapeSuffix        string
	valueCount          int
	arguments           []interface{}
}

// WriteValue query placeholder and append value to argument.
func (b *Buffer) WriteValue(value interface{}) {
	b.WritePlaceholder()
	b.arguments = append(b.arguments, value)
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

// WriteEscape string.
func (b *Buffer) WriteEscape(value string) {
	b.WriteString(b.escape(value))
}

func (b Buffer) escape(value string) string {
	if b.EscapePrefix == "" && b.EscapeSuffix == "" || value == "*" {
		return value
	}

	key := escapeCacheKey{value: value, prefix: b.EscapePrefix, suffix: b.EscapeSuffix}
	escapedValue, ok := escapeCache.Load(key)
	if ok {
		return escapedValue.(string)
	}

	if len(value) > 0 && value[0] == UnescapeCharacter {
		escapedValue = value[1:]
	} else if i := strings.Index(strings.ToLower(value), " as "); i > -1 {
		escapedValue = b.escape(value[:i]) + " AS " + b.escape(value[i+4:])
	} else if start, end := strings.IndexRune(value, '('), strings.IndexRune(value, ')'); start >= 0 && end >= 0 && end > start {
		escapedValue = value[:start+1] + b.escape(value[start+1:end]) + value[end:]
	} else if strings.HasSuffix(value, "*") {
		escapedValue = b.EscapePrefix + strings.Replace(value, ".", b.EscapeSuffix+".", 1)
	} else {
		escapedValue = b.EscapePrefix +
			strings.Replace(value, ".", b.EscapeSuffix+"."+b.EscapePrefix, 1) +
			b.EscapeSuffix
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
	ArgumentPlaceholder string
	ArgumentOrdinal     bool
	EscapePrefix        string
	EscapeSuffix        string
}

func (bf BufferFactory) Create() Buffer {
	return Buffer{
		ArgumentPlaceholder: bf.ArgumentPlaceholder,
		ArgumentOrdinal:     bf.ArgumentOrdinal,
		EscapePrefix:        bf.EscapePrefix,
		EscapeSuffix:        bf.EscapeSuffix,
	}
}
