package extensions

// Unique removes duplicates from the slice.
func Unique[T comparable](xt []T) []T {
	m := make(map[T]bool)

	var result []T
	for _, t := range xt {
		if _, isInSlice := m[t]; !isInSlice {
			m[t] = true
			result = append(result, t)
		}
	}
	return result
}
