package util

// Fine, I'll do it myself
func Ternary[T any](condition bool, trueValue T, falseValue T) T {
	if condition {
		return trueValue
	} else {
		return falseValue
	}
}
