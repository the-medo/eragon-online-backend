package validator

func ValidateUniversalId(value int32) error {
	return ValidateInt(value, 1)
}
