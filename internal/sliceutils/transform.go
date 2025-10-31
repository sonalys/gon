package sliceutils

func Map[T, V any](slice []T, handler func(T) V) []V {
	output := make([]V, len(slice))
	for i := range slice {
		output[i] = handler(slice[i])
	}
	return output
}
