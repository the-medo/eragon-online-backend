package validator

func ValidateMenuId(value int32) error {
	return ValidateInt(value, 1)
}
