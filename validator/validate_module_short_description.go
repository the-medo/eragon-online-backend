package validator

func ValidateModuleShortDescription(value string) error {
	if err := ValidateString(value, 0, 1000); err != nil {
		return err
	}
	return nil
}
