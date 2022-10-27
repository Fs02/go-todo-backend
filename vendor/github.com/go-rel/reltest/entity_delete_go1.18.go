//go:build go1.18
// +build go1.18

package reltest

// EntityMockDelete mock wrapper
type EntityMockDelete[T any] struct {
	*MockDelete
}

// For assert calls for given entity.
func (emd *EntityMockDelete[T]) For(result *T) *EntityMockDelete[T] {
	emd.MockDelete.For(result)
	return emd
}
