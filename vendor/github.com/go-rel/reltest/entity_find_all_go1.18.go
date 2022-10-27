//go:build go1.18
// +build go1.18

package reltest

// EntityMockFindAll mock wrapper
type EntityMockFindAll[T any] struct {
	*MockFindAll
}

// Result sets the result of this query.
func (emfa *EntityMockFindAll[T]) Result(result []T) *Assert {
	return emfa.MockFindAll.Result(result)
}
