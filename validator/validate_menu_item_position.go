package validator

func ValidateMenuItemPosition(value int32) error {
	return ValidateInt(value, 1)
}
