package validator

func ValidateFilename(value string) error {
	return ValidateString(value, 1, 128)
}
