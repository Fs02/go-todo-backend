package postgres

import (
	"database/sql/driver"
	"strings"
	"time"
)

// Quote PostgreSQL identifiers and literals.
type Quote struct{}

func (q Quote) ID(name string) string {
	end := strings.IndexRune(name, 0)
	if end > -1 {
		name = name[:end]
	}
	return `"` + strings.ReplaceAll(name, `"`, `""`) + `"`
}

func (q Quote) Value(v interface{}) string {
	switch v := v.(type) {
	default:
		panic("unsupported value")
	case string:
		v = strings.ReplaceAll(v, `'`, `''`)
		if strings.Contains(v, `\`) {
			v = strings.ReplaceAll(v, `\`, `\\`)
			v = ` E'` + v + `'`
		} else {
			v = `'` + v + `'`
		}
		return v
	}
}

// ValueConvert converts values to PostgreSQL literals.
type ValueConvert struct{}

func (c ValueConvert) ConvertValue(v interface{}) (driver.Value, error) {
	v, err := driver.DefaultParameterConverter.ConvertValue(v)
	if err != nil {
		return nil, err
	}
	switch v := v.(type) {
	default:
		return v, nil
	case time.Time:
		return FormatTime(v), nil
	}
}
