package util

func ArrayMap[T any, U any](array []T, mapper func(T) U) []U {
	var mapped = make([]U, 0, len(array))
	for _, element := range array {
		mapped = append(mapped, mapper(element))
	}

	return mapped
}

func ArrayFilter[T any](array []T, filter func(T) bool) []T {
	var filtered []T
	for _, element := range array {
		if filter(element) {
			filtered = append(filtered, element)
		}
	}

	return filtered
}

func ArrayDifference[T comparable](arrayA, arrayB []T) []T {
	return ArrayFilter(arrayA, func(e T) bool { return !ArrayContains(arrayB, e) })
}

// Retains only elements in array A that exist in array B
func ArrayUnion[T comparable](arrayA, arrayB []T) []T {
	return ArrayFilter(arrayA, func(e T) bool { return ArrayContains(arrayB, e) })
}

func ArrayContains[T comparable](array []T, element T) bool {
	return ArrayIndexOf(array, element) != -1
}

func ArrayIndexOf[T comparable](array []T, element T) int {
	for index, item := range array {
		if item == element {
			return index
		}
	}

	return -1
}
