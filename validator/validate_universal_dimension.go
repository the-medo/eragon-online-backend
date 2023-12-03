package validator

func ValidateUniversalDimension(value int32) error {
	return ValidateInt(value, 1, 10000)
}
