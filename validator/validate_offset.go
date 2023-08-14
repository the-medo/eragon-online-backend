package validator

func ValidateOffset(value int32) error {
	return ValidateInt(value, 0)
}
