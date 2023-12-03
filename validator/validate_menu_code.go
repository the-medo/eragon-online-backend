package validator

func ValidateMenuCode(value string) error {
	return ValidateString(value, 1, 64)
}
