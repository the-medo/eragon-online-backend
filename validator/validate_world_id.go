package validator

func ValidateWorldId(value int32) error {
	return ValidateInt(value, 1)
}
