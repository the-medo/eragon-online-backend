package validator

func ValidateSecretCode(value string) error {
	if err := ValidateString(value, 32, 128); err != nil {
		return err
	}
	return nil
}
