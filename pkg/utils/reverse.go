package utils

// ReverseSlice reverses the order of the elements in a slice of any type.
func ReverseSlice[T any](s []T) {
	// Initialize two pointers, left and right, at the beginning and end of the slice.
	left := 0
	right := len(s) - 1

	// Swap elements at the left and right indices until they meet in the middle.
	for left < right {
		s[left], s[right] = s[right], s[left]
		left++
		right--
	}
}