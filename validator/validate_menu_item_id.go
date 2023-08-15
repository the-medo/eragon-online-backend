package validator

func ValidateMenuItemId(value int32) error {
	return ValidateInt(value, 1)
}
