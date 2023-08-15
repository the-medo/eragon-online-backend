package validator

func ValidateChatMessageText(value string) error {
	return ValidateString(value, 1, 1024)
}
