package validator

func ValidateImageTypeId(value int32) error {
	return ValidateInt(value, 1)
}
