package validator

func ValidateTagId(value int32) error {
	return ValidateInt(value, 1)
}
