package util

func ArrayMap[T any, U any](array []T, mapper func(T) U) []U {
	var mapped = make([]U, 0, len(array))
	for _, element := range array {
		mapped = append(mapped, mapper(element))
	}

	return mapped
}

func ArrayFlatMap[T any, U any](array []T, mapper func(T) []U) []U {
	var mapped = make([]U, 0, len(array))
	for _, element := range array {
		mapped = append(mapped, mapper(element)...)
	}

	return mapped
}

// For all those multiple conversions where an error might be returned in any one of them
func ArrayMapKeepSecondIfPresent[T any, R any, S any](array []T, mapper func(T) (R, *S)) ([]R, *S) {
	var mapped = make([]R, 0, len(array))
	var retained *S
	for _, element := range array {
		mappedElement, maybeRetained := mapper(element)
		if maybeRetained != nil {
			retained = maybeRetained
		}
		mapped = append(mapped, mappedElement)
	}

	return mapped, retained
}

// Specifically for the error type, but otherwise the same as ArrayMapKeepSecondIfPresent
func ArrayMapTError[T any, R any](array []T, mapper func(T) (R, error)) ([]R, error) {
	var mapped = make([]R, 0, len(array))
	for _, element := range array {
		mappedElement, err := mapper(element)
		if err != nil {
			return nil, err
		}
		mapped = append(mapped, mappedElement)
	}

	return mapped, nil
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
