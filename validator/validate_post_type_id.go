package validator

func ValidatePostTypeId(value int32) error {
	return ValidateInt(value, 1)
}
