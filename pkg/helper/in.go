package helper

import "golang.org/x/exp/slices"

// EqualChecker is an interface that allows to check interface for equality.
type EqualChecker[T any] interface {
	Equal(T) bool
}

// In returns whether `search` appears in `target` slice.
func In[T any, C EqualChecker[T]](search C, target []T) bool {
	i := slices.IndexFunc(target, func(t T) bool {
		return search.Equal(t)
	})

	return i != -1
}
