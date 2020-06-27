## snake-case

Convert given string to snake-case. ([Original Code](https://gist.github.com/elwinar/14e1e897fdbe4d3432e1)) by [@elwinar](http://github.com/elwinar)

## Install

```
go get github.com/azer/snakecase
```

## Usage Example

```go
import (
  "github.com/azer/snakecase"
)

snakecase.SnakeCase("APIResponse")
// => api_response

snakecase.SnakeCase("Hello World")
// => hello_world

snakecase.SnakeCase("CreateDB")
// => create_db
```
