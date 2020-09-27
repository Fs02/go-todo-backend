package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Fs02/go-todo-backend/db/migrations"
	"github.com/go-rel/rel"
	"github.com/go-rel/rel/adapter/postgres"
	"github.com/go-rel/rel/migrator"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

var (
	logger, _ = zap.NewProduction(zap.Fields(zap.String("type", "main")))
	shutdowns []func() error
)

func main() {
	var (
		ctx = context.Background()
		dsn = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			os.Getenv("POSTGRESQL_USERNAME"),
			os.Getenv("POSTGRESQL_PASSWORD"),
			os.Getenv("POSTGRESQL_HOST"),
			os.Getenv("POSTGRESQL_PORT"),
			os.Getenv("POSTGRESQL_DATABASE"))
	)

	adapter, err := postgres.Open(dsn)
	if err != nil {
		logger.Fatal(err.Error(), zap.Error(err))
	}

	var (
		op   string
		repo = rel.New(adapter)
		m    = migrator.New(repo)
	)

	// There will be a command line like go test for this in the future.
	m.Register(20202806225100, migrations.MigrateCreateTodos, migrations.RollbackCreateTodos)
	m.Register(20203006230600, migrations.MigrateCreateScores, migrations.RollbackCreateScores)
	m.Register(20203006230700, migrations.MigrateCreatePoints, migrations.RollbackCreatePoints)

	if len(os.Args) > 1 {
		op = os.Args[1]
	}

	switch op {
	case "migrate", "up":
		m.Migrate(ctx)
	case "rollback", "down":
		m.Rollback(ctx)
	default:
		logger.Fatal("command not recognized", zap.String("command", op))
	}
}
