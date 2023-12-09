package validator

func ValidateNullableImageId(value int32) error {
	return ValidateInt(value, 0)
}
