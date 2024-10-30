package validator

func ValidateModuleId(value int32) error {
	return ValidateInt(value, 1)
}
