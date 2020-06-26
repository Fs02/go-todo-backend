package handler_test

import (
	"encoding/json"
)

// jo stands for json object.
// it's an alias that allows writing map for json response easier.
type jo map[string]interface{}

func (j jo) String() string {
	bytes, _ := json.Marshal(j)
	return string(bytes)
}
