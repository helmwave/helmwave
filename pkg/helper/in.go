package helper

// EqualChecker is an interface that allows to check interface for equality.
type EqualChecker[T any] interface {
	Equal(T) bool
}

// In returns whether `search` appears in `target` slice.
func In[T any, C EqualChecker[T]](search C, target []T) bool {
	for i := range target {
		if search.Equal(target[i]) {
			return true
		}
	}

	return false
}
