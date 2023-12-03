package validator

func ValidateImageId(value int32) error {
	return ValidateInt(value, 1)
}
