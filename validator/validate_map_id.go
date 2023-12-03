package validator

func ValidateMapId(value int32) error {
	return ValidateInt(value, 1)
}
