//go:build go1.18
// +build go1.18

package reltest

// EntityMockPreload mock wrapper
type EntityMockPreload[T any] struct {
	*MockPreload
}

// For assert calls for given entity.
func (emp *EntityMockPreload[T]) For(result *T) *EntityMockPreload[T] {
	emp.MockPreload.For(result)
	return emp
}
