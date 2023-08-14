package validator

func ValidatePassword(value string) error {
	return ValidateString(value, 6, 100)
}
