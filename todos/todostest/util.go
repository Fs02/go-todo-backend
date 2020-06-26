package todostest

import (
	context "context"

	todos "github.com/Fs02/go-todo-backend/todos"
	rel "github.com/Fs02/rel"
	mock "github.com/stretchr/testify/mock"
)

// MockTodosFunc function.
type MockTodosFunc func(todo *Todos)

// MockTodos apply mock todo functions.
func MockTodos(todos *Todos, funcs ...MockTodosFunc) {
	for i := range funcs {
		if funcs[i] != nil {
			funcs[i](todos)
		}
	}
}

// MockTodosSearch util.
func MockTodosSearch(result []todos.Todo, filter todos.Filter, err error) MockTodosFunc {
	return func(svc *Todos) {
		svc.On("Search", mock.Anything, mock.Anything, filter).
			Return(func(ctx context.Context, out *[]todos.Todo, filter todos.Filter) error {
				*out = result
				return err
			})
	}
}

// MockTodosCreate util.
func MockTodosCreate(result todos.Todo, err error) MockTodosFunc {
	return func(svc *Todos) {
		svc.On("Create", mock.Anything, mock.Anything).
			Return(func(ctx context.Context, out *todos.Todo) error {
				*out = result
				return err
			})
	}
}

// MockTodosUpdate util.
func MockTodosUpdate(result todos.Todo, err error) MockTodosFunc {
	return func(svc *Todos) {
		svc.On("Update", mock.Anything, mock.Anything, mock.Anything).
			Return(func(ctx context.Context, out *todos.Todo, changeset rel.Changeset) error {
				if result.ID != out.ID {
					panic("inconsistent id")
				}

				*out = result
				return err
			})
	}
}

// MockTodosClear util.
func MockTodosClear() MockTodosFunc {
	return func(svc *Todos) {
		svc.On("Clear", mock.Anything)
	}
}

// MockTodosDelete util.
func MockTodosDelete() MockTodosFunc {
	return func(svc *Todos) {
		svc.On("Delete", mock.Anything, mock.Anything)
	}
}
