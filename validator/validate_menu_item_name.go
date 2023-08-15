package validator

func ValidateMenuItemName(value string) error {
	return ValidateString(value, 1, 64)
}
