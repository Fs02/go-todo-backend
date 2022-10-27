//go:build go1.18
// +build go1.18

package reltest

// EntityMockIterate mock wrapper
type EntityMockIterate[T any] struct {
	*MockIterate
}

// Result sets the result of preload.
func (emi *EntityMockIterate[T]) Result(result []T) *Assert {
	return emi.MockIterate.Result(result)
}
