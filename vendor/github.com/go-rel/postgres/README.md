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

## Supported Driver

- github.com/lib/pq

## Supported Database

- PostgreSQL 9, 10, 11, 12 and 13
