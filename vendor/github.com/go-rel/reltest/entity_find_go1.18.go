//go:build go1.18
// +build go1.18

package reltest

// EntityMockFind mock wrapper
type EntityMockFind[T any] struct {
	*MockFind
}

// Result sets the result of this query.
func (emf *EntityMockFind[T]) Result(result T) *Assert {
	return emf.MockFind.Result(result)
}
