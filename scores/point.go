package scores

import (
	"time"
)

// Point component for score.
type Point struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Count     int       `json:"count"`
	ScoreID   int       `json:"score_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
