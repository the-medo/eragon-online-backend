package util

// Upsert appends a value to a slice if it does not already exist in the slice.
// Works with any comparable type.
func Upsert[T comparable](slice []T, value T) []T {
	for _, v := range slice {
		if v == value {
			return slice // Value already exists, return the original slice
		}
	}
	return append(slice, value) // Value not found, append it
}
