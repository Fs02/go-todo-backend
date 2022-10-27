//go:build go1.18
// +build go1.18

package reltest

// EntityMockInsertAll mock wrapper
type EntityMockInsertAll[T any] struct {
	*MockInsertAll
}

// For assert calls for given entity.
func (emia *EntityMockInsertAll[T]) For(result *[]T) *EntityMockInsertAll[T] {
	emia.MockInsertAll.For(result)
	return emia
}
