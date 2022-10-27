//go:build go1.18
// +build go1.18

package reltest

// EntityMockMutate mock wrapper
type EntityMockMutate[T any] struct {
	*MockMutate
}

// Result sets the result of this query.
func (emm *EntityMockMutate[T]) For(result *T) *EntityMockMutate[T] {
	emm.MockMutate.For(result)
	return emm
}
