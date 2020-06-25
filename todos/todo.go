package todos

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

var (
	// TodoURLPrefix to be returned when encoding todo.
	TodoURLPrefix = os.Getenv("URL") + "todos/"
	// ErrTodoTitleBlank validation error.
	ErrTodoTitleBlank = errors.New("Title can't be blank")
)

// Todo respresent a record stored in todos table.
type Todo struct {
	ID        uint   `json:"id"`
	Title     string `json:"title"`
	Order     int    `json:"order"`
	Completed bool   `json:"completed"`
}

// Validate todo.
func (t Todo) Validate() error {
	var err error
	switch {
	case len(t.Title) == 0:
		err = ErrTodoTitleBlank
	}

	return err
}

// MarshalJSON implement custom marshaller to marshal url.
func (t Todo) MarshalJSON() ([]byte, error) {
	type Alias Todo

	return json.Marshal(struct {
		Alias
		URL string `json:"url"`
	}{
		Alias: Alias(t),
		URL:   fmt.Sprint(TodoURLPrefix, t.ID),
	})
}
