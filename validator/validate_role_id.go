package validator

func ValidateRoleId(value int32) error {
	return ValidateInt(value, 1, 2)
}
