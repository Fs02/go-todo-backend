//go:build go1.18
// +build go1.18

package reltest

// EntityMockPreloadAll mock wrapper
type EntityMockPreloadAll[T any] struct {
	*MockPreload
}

// For assert calls for given entity.
func (empa *EntityMockPreloadAll[T]) For(result *[]T) *EntityMockPreloadAll[T] {
	empa.MockPreload.For(result)
	return empa
}
