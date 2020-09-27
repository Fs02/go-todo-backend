package migrations

import (
	"github.com/go-rel/rel"
)

// MigrateCreateScores definition
func MigrateCreateScores(schema *rel.Schema) {
	schema.CreateTable("scores", func(t *rel.Table) {
		t.ID("id")
		t.DateTime("created_at")
		t.DateTime("updated_at")
		t.Int("total_point")
	})
}

// RollbackCreateScores definition
func RollbackCreateScores(schema *rel.Schema) {
	schema.DropTable("scores")
}
