# postgres

[![GoDoc](https://godoc.org/github.com/go-rel/postgres?status.svg)](https://pkg.go.dev/github.com/go-rel/postgres)
[![Tesst](https://github.com/go-rel/postgres/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/go-rel/postgres/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-rel/postgres)](https://goreportcard.com/report/github.com/go-rel/postgres)
[![codecov](https://codecov.io/gh/go-rel/postgres/branch/main/graph/badge.svg?token=yxBdKVPXip)](https://codecov.io/gh/go-rel/postgres)
[![Gitter chat](https://badges.gitter.im/go-rel/rel.png)](https://gitter.im/go-rel/rel)

Postgres adapter for REL.

## Example

```go
package main

import (
	"context"

	_ "github.com/lib/pq"
	"github.com/go-rel/postgres"
	"github.com/go-rel/rel"
)

func main() {
	// open postgres connection.
	adapter, err := postgres.Open("postgres://postgres@localhost/rel_test?sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer adapter.Close()

	// initialize REL's repo.
	repo := rel.New(adapter)
	repo.Ping(context.TODO())
}
```

## Example Replication (Master/Standby)

```go
package main

import (
	"context"

	"github.com/go-rel/primaryreplica"
	_ "github.com/lib/pq"
	"github.com/go-rel/postgres"
	"github.com/go-rel/rel"
)

func main() {
	// open postgres connections.
	adapter := primaryreplica.New(
		postgres.MustOpen("postgres://postgres@master/rel_test?sslmode=disable"),
		postgres.MustOpen("postgres://postgres@standby/rel_test?sslmode=disable"),
	)
	defer adapter.Close()

	// initialize REL's repo.
	repo := rel.New(adapter)
	repo.Ping(context.TODO())
}
```

## Supported Driver

- github.com/lib/pq
- github.com/jackc/pgx/v4/stdlib

## Supported Database

- PostgreSQL 9.6, 10, 11, 12, 13 and 14

## Testing

### Start PostgreSQL server in Docker

```console
docker run -it --rm -p 5433:5432 -e "POSTGRES_USER=rel" -e "POSTGRES_PASSWORD=test" -e "POSTGRES_DB=rel_test" postgres:14-alpine
```

### Run tests

```console
POSTGRESQL_DATABASE="postgres://rel:test@localhost:5433/rel_test" go test ./...
```
