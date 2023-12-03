package validator

func ValidateChatMessageId(value int64) error {
	return ValidateInt64(value, 1)
}
