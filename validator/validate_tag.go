package validator

func ValidateTag(value string) error {
	return ValidateString(value, 1, 32)
}
