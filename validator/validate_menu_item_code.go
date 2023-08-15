package validator

func ValidateMenuItemCode(value string) error {
	return ValidateString(value, 1, 64)
}
