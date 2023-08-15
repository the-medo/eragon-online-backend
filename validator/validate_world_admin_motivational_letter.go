package validator

func ValidateWorldAdminMotivationalLetter(value string) error {
	if err := ValidateString(value, 0, 2000); err != nil {
		return err
	}
	return nil
}
