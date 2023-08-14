package validator

func ValidatePostId(value int32) error {
	return ValidateInt(value, 1)
}
