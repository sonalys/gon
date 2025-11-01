package sliceutils

func Map[T, V any](slice []T, handler func(T) V) []V {
	output := make([]V, len(slice))
	for i := range slice {
		output[i] = handler(slice[i])
	}
	return output
}

func ToMap[T any, K comparable, V any](slice []T, handler func(T) (K, V)) map[K]V {
	output := make(map[K]V, len(slice))

	for i := range slice {
		key, value := handler(slice[i])
		output[key] = value
	}

	return output
}

func Filter[T any](slice []T, filter func(T) bool) []T {
	output := make([]T, 0, len(slice))

	for i := range slice {
		if filter(slice[i]) {
			output = append(output, slice[i])
		}
	}

	return output
}
