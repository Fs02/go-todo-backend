//go:build go1.18
// +build go1.18

package reltest

// EntityMockFindAndCountAll mock wrapper
type EntityMockFindAndCountAll[T any] struct {
	*MockFindAndCountAll
}

// Result sets the result of this query.
func (emfaca *EntityMockFindAndCountAll[T]) Result(result []T, count int) *Assert {
	return emfaca.MockFindAndCountAll.Result(result, count)
}
