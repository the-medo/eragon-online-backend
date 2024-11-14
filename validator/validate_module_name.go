package validator

func ValidateModuleName(value string) error {
	if err := ValidateString(value, 3, 64); err != nil {
		return err
	}
	return nil
}
