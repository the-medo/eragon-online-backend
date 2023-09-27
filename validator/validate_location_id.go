package validator

func ValidateLocationId(value int32) error {
	return ValidateInt(value, 1)
}
