# reltest

[![GoDoc](https://godoc.org/github.com/go-rel/reltest?status.svg)](https://pkg.go.dev/github.com/go-rel/reltest)
[![Test](https://github.com/go-rel/reltest/actions/workflows/test.yml/badge.svg)](https://github.com/go-rel/reltest/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-rel/reltest)](https://goreportcard.com/report/github.com/go-rel/reltest)
[![codecov](https://codecov.io/gh/go-rel/reltest/branch/main/graph/badge.svg?token=vxG9e5nJ3R)](https://codecov.io/gh/go-rel/reltest)
[![Gitter chat](https://badges.gitter.im/go-rel/rel.png)](https://gitter.im/go-rel/rel)

Database unit testing for Golang.

## Example 

```go
package main

import (
	"context"
	"fmt"

	"github.com/go-rel/rel/where"
	"github.com/go-rel/reltest"
)

type Movie struct {
	ID    int
	Title string
}

func main() {
	var (
		repo = reltest.New()
	)

	// Mock query
	repo.ExpectFind(where.Eq("id", 1)).Result(Movie{ID: 1, Title: "Golang"})

	// Application code
	var movie Movie
	repo.MustFind(context.Background(), &movie, where.Eq("id", 1))
	fmt.Println(movie.Title)
	// Output: Golang
}
```

**More Examples:**

- [gin-example](https://github.com/go-rel/gin-example) - Todo Backend using Gin and REL
- [go-todo-backend](https://github.com/Fs02/go-todo-backend) - Todo Backend using Chi and REL

## License

Released under the [MIT License](https://github.com/go-rel/reltest/blob/master/LICENSE)
