package validator

func ValidateEmailId(value int64) error {
	return ValidateInt64(value, 1)
}
