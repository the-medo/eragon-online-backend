package validator

func ValidateMenuItemParentItemId(value int32) error {
	return ValidateInt(value, 1)
}
