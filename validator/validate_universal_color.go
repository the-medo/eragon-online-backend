package validator

// ValidateUniversalName basic name validation - between 1 and 128 characters
func ValidateUniversalColor(value string) error {
	if err := ValidateString(value, 3, 64); err != nil {
		return err
	}
	return nil
}
