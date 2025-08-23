package util

func Map[T, U any](arr []T, mapFunc func(T, int) U) []U {
	output := make([]U, len(arr))
	for i, val := range arr {
		output[i] = mapFunc(val, i)
	}
	return output
}
