package util

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
