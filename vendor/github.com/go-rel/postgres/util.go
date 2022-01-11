package postgres

import (
	"time"
)

// TimeLayout used by PostgreSQL adapter.
const TimeLayout = "2006-01-02 15:04:05.999999999Z07:00:00"

// FormatTime formats time to PostgreSQL format.
func FormatTime(t time.Time) string {
	return t.Truncate(time.Microsecond).Format(TimeLayout)
}
