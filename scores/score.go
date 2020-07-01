package scores

import (
	"time"
)

// Score stores total points.
type Score struct {
	ID         int       `json:"id"`
	TotalPoint int       `json:"total_point"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
