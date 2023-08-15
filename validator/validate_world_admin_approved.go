package validator

func ValidateWorldAdminApproved(value int32) error {
	return ValidateInt(value, 0, 2)
}
