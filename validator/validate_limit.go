package validator

func ValidateLimit(value int32) error {
	return ValidateInt(value, 0, 1000)
}
