package util

// CleanEmptyStrings removes empty strings from a slice of strings in place.
// It modifies the existing slice, rather than returning a new one.
//
// Parameters:
//   - tags: A pointer to a slice of strings that you want to clean.
//
// Example:
//
//	tags := []string{"apple", "", "banana"}
//	CleanEmptyStrings(&tags)
//	// tags now contains: ["apple", "banana"]
func CleanEmptyStrings(strings *[]string) {
	cleanStrings := (*strings)[:0]
	for _, tag := range *strings {
		if tag != "" {
			cleanStrings = append(cleanStrings, tag)
		}
	}
	*strings = cleanStrings
}
