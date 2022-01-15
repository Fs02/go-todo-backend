package builder

import (
	"strings"
)

// Quoter returns safe and valid SQL strings to use when building a SQL text.
type Quoter interface {
	// ID quotes identifiers such as schema, table, or column names.
	// ID does not operate on multipart identifiers such as "public.Table",
	// it only operates on single identifiers such as "public" and "Table".
	ID(name string) string

	// Value quotes database values such as string or []byte types as strings
	// that are suitable and safe to embed in SQL text. The returned value
	// of a string will include all surrounding quotes.
	//
	// If a value type is not supported it must panic.
	Value(v interface{}) string
}

// Quote is default implementation of Quoter interface.
type Quote struct {
	IDPrefix             string
	IDSuffix             string
	IDSuffixEscapeChar   string
	ValueQuote           string
	ValueQuoteEscapeChar string
}

func (q Quote) ID(name string) string {
	return q.IDPrefix + strings.ReplaceAll(name, q.IDSuffix, q.IDSuffixEscapeChar+q.IDSuffix) + q.IDSuffix
}

func (q Quote) Value(v interface{}) string {
	switch v := v.(type) {
	default:
		panic("unsupported value")
	case string:
		return q.ValueQuote + strings.ReplaceAll(v, q.ValueQuote, q.ValueQuoteEscapeChar+q.ValueQuote) + q.ValueQuote
	}
}
