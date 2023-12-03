package validator

func ValidateModuleAdminApproved(value int32) error {
	return ValidateInt(value, 0, 2)
}
