package validator

func ValidateUserId(value int32) error {
	return ValidateInt(value, 1)
}
