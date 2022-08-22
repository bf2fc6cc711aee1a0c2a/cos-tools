package collections

func Contains[T comparable](elems []T, v T) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}

func Filter[T comparable](elems []T, predicate func(T) bool) []T {
	answer := make([]T, len(elems))

	for i := range elems {
		if predicate(elems[i]) {
			answer = append(answer, elems[i])
		}
	}

	return answer
}
