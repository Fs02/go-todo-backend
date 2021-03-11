package migrations

import (
	"github.com/go-rel/rel"
)

// MigrateCreatePoints definition
func MigrateCreatePoints(schema *rel.Schema) {
	schema.CreateTable("points", func(t *rel.Table) {
		t.ID("id")
		t.DateTime("created_at")
		t.DateTime("updated_at")
		t.String("name")
		t.Int("count")
		t.Int("score_id", rel.Unsigned(true))

		t.ForeignKey("score_id", "scores", "id")
	})
}

// RollbackCreatePoints definition
func RollbackCreatePoints(schema *rel.Schema) {
	schema.DropTable("points")
}
