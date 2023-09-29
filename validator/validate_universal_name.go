package validator

// ValidateUniversalName basic name validation - between 1 and 128 characters
func ValidateUniversalName(value string) error {
	if err := ValidateString(value, 1, 128); err != nil {
		return err
	}
	return nil
}
