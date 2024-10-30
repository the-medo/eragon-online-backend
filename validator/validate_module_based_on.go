package validator

func ValidateModuleBasedOn(value string) error {
	if err := ValidateString(value, 0, 100); err != nil {
		return err
	}
	return nil
}
