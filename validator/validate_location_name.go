package validator

func ValidateLocationName(value string) error {
	if err := ValidateString(value, 2, 64); err != nil {
		return err
	}
	return nil
}
