//go:build go1.18
// +build go1.18

package reltest

// EntityMockDeleteAll mock wrapper
type EntityMockDeleteAll[T any] struct {
	*MockDeleteAll
}

// For assert calls for given entity.
func (emda *EntityMockDeleteAll[T]) For(result *[]T) *EntityMockDeleteAll[T] {
	emda.MockDeleteAll.For(result)
	return emda
}
